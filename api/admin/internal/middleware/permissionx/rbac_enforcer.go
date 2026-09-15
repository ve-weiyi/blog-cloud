package permissionx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/mikespook/gorbac/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

// apiPermission 实现 gorbac.Permission，ID 格式为 "METHOD:PATH"
// Match 由策略侧调用：p（策略）.Match(req)，支持通配符
type apiPermission string

func newApiPermission(method, path string) gorbac.Permission {
	return apiPermission(method + ":" + path)
}

func (p apiPermission) ID() string { return string(p) }

func (p apiPermission) Match(other gorbac.Permission) bool {
	return matchPermission(string(p), other.ID())
}

var _ Enforcer = &RbacEnforcer{}

// RbacEnforcer 基于 gorbac 的 RBAC 权限执行器
type RbacEnforcer struct {
	mu  sync.RWMutex
	pr  accessservice.AccessService
	rds *redis.Client

	rbac         *gorbac.RBAC        // role -> permissions（全量加载，原子替换）
	allPerms     map[string]struct{} // 全量已注册的权限ID（METHOD:PATH）
	userRoles    map[string][]string // 内存缓存: userID -> []roleKey
	userRoleTTL  time.Duration
	policyLoaded bool
}

func NewRbacEnforcer(rds *redis.Client, pr accessservice.AccessService) *RbacEnforcer {
	h := &RbacEnforcer{
		pr:          pr,
		rds:         rds,
		rbac:        gorbac.New(),
		userRoles:   make(map[string][]string),
		userRoleTTL: 5 * time.Minute,
	}
	h.startSubscribe()
	return h
}

func (m *RbacEnforcer) ReloadPolicy() error {
	return m.LoadPolicy()
}

// LoadPolicy 从 RPC 全量拉取角色权限，原子替换 rbac 实例
func (m *RbacEnforcer) LoadPolicy() error {
	if m.pr == nil {
		return errors.New("permission service is nil")
	}

	logx.Info("Loading permissions...")

	ctx := context.Background()

	roleList, err := m.pr.ListRoles(ctx, &accessservice.ListRolesRequest{})
	if err != nil {
		return err
	}

	apiList, err := m.pr.ListApis(ctx, &accessservice.ListApisRequest{})
	if err != nil {
		return err
	}

	flatApis := flattenApiTree(apiList.List)
	apis := make(map[int64]*accessservice.Api, len(flatApis))
	for _, v := range flatApis {
		apis[v.Id] = v
	}

	newRbac := gorbac.New()
	allPerms := make(map[string]struct{})

	for _, role := range roleList.List {
		if role.Status == enums.RoleStatusDisabled {
			continue
		}
		rk := roleKey(role.RoleKey, role.Id)

		resource, err := m.pr.GetRoleResource(ctx, &accessservice.GetRoleResourceRequest{RoleId: role.Id})
		if err != nil {
			return err
		}

		r := gorbac.NewStdRole(rk)
		seen := make(map[string]struct{})
		for _, apiId := range resource.ApiIds {
			api, ok := apis[apiId]
			if !ok || api.Status == enums.APIStatusDisabled {
				continue
			}
			method := normalizeMethod(api.Method)
			path := normalizePath(api.Path)
			if method == "" || path == "" {
				continue
			}
			permID := method + ":" + path
			allPerms[permID] = struct{}{}
			if _, exists := seen[permID]; exists {
				continue
			}
			seen[permID] = struct{}{}
			_ = r.Assign(newApiPermission(method, path))
		}
		if err := newRbac.Add(r); err != nil && !errors.Is(err, gorbac.ErrRoleExist) {
			return err
		}
	}

	m.mu.Lock()
	m.rbac = newRbac
	m.allPerms = allPerms
	m.userRoles = make(map[string][]string)
	m.policyLoaded = true
	m.mu.Unlock()

	return nil
}

// Enforce 校验用户是否有权限访问指定资源
func (m *RbacEnforcer) Enforce(user string, resource string, action string) (bool, error) {
	if strings.TrimSpace(user) == "" {
		return false, errors.New("user is empty")
	}
	if strings.TrimSpace(resource) == "" {
		return false, errors.New("resource is empty")
	}
	if strings.TrimSpace(action) == "" {
		return false, errors.New("action is empty")
	}

	m.mu.RLock()
	loaded := m.policyLoaded
	m.mu.RUnlock()
	if !loaded {
		if err := m.LoadPolicy(); err != nil {
			return false, err
		}
	}

	roles, err := m.getUserRoles(user)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		if strings.EqualFold(r, "root") {
			logx.Infof("[Perm] user=%s role=root (super admin bypass)", user)
			return true, nil
		}
	}

	method := normalizeMethod(action)
	path := normalizePath(resource)
	if method == "" || path == "" {
		return false, errors.New("invalid action or resource")
	}

	permKey := method + ":" + path

	m.mu.RLock()
	rbac := m.rbac
	allPerms := m.allPerms
	m.mu.RUnlock()

	if _, registered := allPerms[permKey]; !registered {
		logx.Infof("[Perm] user=%s method=%s path=%s -> unregistered, allowed", user, method, path)
		return true, nil
	}

	req := newApiPermission(method, path)
	for _, roleKey := range roles {
		if rbac.IsGranted(roleKey, req, nil) {
			logx.Infof("[Perm] user=%s role=%s method=%s path=%s -> granted", user, roleKey, method, path)
			return true, nil
		}
	}

	logx.Infof("[Perm] user=%s roles=%v method=%s path=%s -> denied", user, roles, method, path)
	return false, fmt.Errorf("用户[%s]无权限访问资源[%s %s]", user, method, path)
}

// InvalidateUser 主动失效用户角色缓存
func (m *RbacEnforcer) InvalidateUser(userId string) {
	m.mu.Lock()
	delete(m.userRoles, userId)
	m.mu.Unlock()

	if m.rds != nil {
		_ = m.rds.Del(context.Background(), cachekey.UserRoleCacheKey(userId)).Err()
	}
}

// getUserRoles 获取用户角色：内存缓存 -> Redis -> RPC
func (m *RbacEnforcer) getUserRoles(userId string) ([]string, error) {
	// 1. 内存缓存
	m.mu.RLock()
	roles, ok := m.userRoles[userId]
	m.mu.RUnlock()
	if ok {
		return slices.Clone(roles), nil
	}

	// 2. Redis 缓存
	if m.rds != nil {
		val, err := m.rds.Get(context.Background(), cachekey.UserRoleCacheKey(userId)).Result()
		if err == nil {
			var roles []string
			if jsonErr := json.Unmarshal([]byte(val), &roles); jsonErr != nil {
				logx.Errorf("unmarshal user roles failed: %v", jsonErr)
			} else {
				m.mu.Lock()
				m.userRoles[userId] = roles
				m.mu.Unlock()
				return slices.Clone(roles), nil
			}
		} else if !errors.Is(err, redis.Nil) {
			logx.Errorf("load user roles from redis failed: %v", err)
		}
	}

	// 3. RPC 拉取
	if m.pr == nil {
		return nil, errors.New("permission service is nil")
	}
	resp, err := m.pr.ListUserRoles(context.Background(), &accessservice.ListUserRolesRequest{UserId: userId})
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	var roleKeys []string
	for _, v := range resp.List {
		if v.Status == enums.RoleStatusDisabled {
			continue
		}
		key := roleKey(v.RoleKey, v.Id)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		roleKeys = append(roleKeys, key)
	}

	if m.rds != nil {
		if payload, err := json.Marshal(roleKeys); err == nil {
			if err = m.rds.Set(context.Background(), cachekey.UserRoleCacheKey(userId), payload, m.userRoleTTL).Err(); err != nil {
				logx.Errorf("cache user roles failed: %v", err)
			}
		}
	}

	m.mu.Lock()
	m.userRoles[userId] = roleKeys
	m.mu.Unlock()

	return slices.Clone(roleKeys), nil
}

// startSubscribe 启动 Pub/Sub 订阅，监听策略变更和用户角色失效
func (m *RbacEnforcer) startSubscribe() {
	if m.rds == nil {
		return
	}
	go func() {
		for {
			if err := m.subscribe(); err != nil {
				logx.Errorf("pubsub subscription lost: %v, reconnecting...", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

func (m *RbacEnforcer) subscribe() error {
	ctx := context.Background()
	sub := m.rds.Subscribe(ctx, cachekey.PolicyInvalidateChannel, cachekey.UserRoleInvalidateChannel)
	defer sub.Close()

	const debounce = 500 * time.Millisecond
	var timer *time.Timer

	ch := sub.Channel()
	for msg := range ch {
		switch msg.Channel {
		case cachekey.PolicyInvalidateChannel:
			if timer == nil {
				timer = time.AfterFunc(debounce, func() {
					logx.Info("received policy invalidate, reloading...")
					if err := m.LoadPolicy(); err != nil {
						logx.Errorf("reload policy failed: %v", err)
					}
				})
			} else {
				timer.Reset(debounce)
			}
		case cachekey.UserRoleInvalidateChannel:
			if msg.Payload != "" {
				logx.Infof("received user role invalidate: userId=%s", msg.Payload)
				m.InvalidateUser(msg.Payload)
			}
		}
	}
	return errors.New("channel closed")
}

func flattenApiTree(nodes []*accessservice.Api) []*accessservice.Api {
	var result []*accessservice.Api
	for _, node := range nodes {
		result = append(result, node)
		if len(node.Children) > 0 {
			result = append(result, flattenApiTree(node.Children)...)
		}
	}
	return result
}

// roleKey 返回角色的唯一 key，空时回退为 "role:<id>"
func roleKey(key string, id int64) string {
	if strings.TrimSpace(key) == "" {
		return fmt.Sprintf("role:%d", id)
	}
	return key
}
