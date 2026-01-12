package permissionx

import (
	"context"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	jsonadapter "github.com/casbin/json-adapter/v2"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
)

const SubjectObjectAction = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

type PermissionHolder interface {
	// 加载权限
	LoadPolicy() error
	// 清除权限
	ReloadPolicy() error
	// 验证用户是否有权限
	Enforce(user string, resource string, action string) error
}

// 角色资源管理器
type CasbinHolder struct {
	rw sync.RWMutex

	pr permissionrpc.PermissionRpc

	// casbin
	enforcer *casbin.SyncedCachedEnforcer

	// 用户角色缓存 key: (user), value: role_key[]
	user map[string][]string
}

func NewCasbinHolder(redisAddr string, pr permissionrpc.PermissionRpc) *CasbinHolder {
	// 载入模型
	m, err := model.NewModelFromString(SubjectObjectAction)
	if err != nil {
		panic(fmt.Errorf("字符串加载模型失败: %v", err))
	}

	// 载入适配器
	b := []byte{}                         // b stores Casbin policy in JSON bytes.
	adapter := jsonadapter.NewAdapter(&b) // Use b as the data source.

	// 初始化
	enforcer, _ := casbin.NewSyncedCachedEnforcer(m, adapter)
	enforcer.SetExpireTime(60 * 60)
	return &CasbinHolder{
		rw:       sync.RWMutex{},
		pr:       pr,
		enforcer: enforcer,
		user:     make(map[string][]string),
	}
}

func (m *CasbinHolder) ReloadPolicy() error {
	return m.LoadPolicy()
}

func (m *CasbinHolder) LoadPolicy() error {
	m.rw.Lock()
	defer m.rw.Unlock()

	// 重置所有权限
	m.enforcer.ClearPolicy()
	m.user = make(map[string][]string)

	var rs = make(map[int64][]int64)
	var roles = make(map[int64]*permissionrpc.Role)
	var apis = make(map[int64]*permissionrpc.Api)

	// 收集角色
	roleList, err := m.pr.FindAllRole(context.Background(), &permissionrpc.FindAllRoleReq{})
	if err != nil {
		return err
	}
	for _, v := range roleList.List {
		roles[v.Id] = v
	}

	// 收集资源
	apiList, err := m.pr.FindAllApi(context.Background(), &permissionrpc.FindAllApiReq{})
	if err != nil {
		return err
	}
	for _, v := range apiList.List {
		apis[v.Id] = v
	}

	// 收集角色-资源
	for _, v := range roleList.List {
		resource, err := m.pr.FindRoleResources(context.Background(), &permissionrpc.FindRoleResourcesReq{Id: v.Id})
		if err != nil {
			return err
		}

		rs[v.Id] = resource.ApiIds
	}

	for roleId, apiIds := range rs {
		role, ok := roles[roleId]
		if !ok {
			return fmt.Errorf("role %d not found", roleId)
		}

		if role.Status == 1 {
			continue
		}

		for _, apiId := range apiIds {
			api, ok := apis[apiId]
			if !ok {
				return fmt.Errorf("api %d not found", apiId)
			}

			if api.Status == 1 {
				continue
			}

			_, err = m.enforcer.AddPolicy(role.RoleKey, api.Path, api.Method)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *CasbinHolder) Enforce(user string, resource string, action string) error {
	m.rw.RUnlock()
	defer m.rw.RUnlock()

	err := m.dynamicLoadUserRoles(user)
	if err != nil {
		return err
	}

	ok, err := m.enforcer.Enforce(user, resource, action)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("permission denied, user: %s, resource: %s, action: %s", user, resource, action)
	}

	return nil
}

// 加载用户角色
func (m *CasbinHolder) dynamicLoadUserRoles(userId string) error {
	if _, ok := m.user[userId]; ok {
		return nil
	}

	// 获取用户角色
	userRoles, err := m.pr.FindUserRoles(context.Background(), &permissionrpc.FindUserRolesReq{UserId: userId})
	if err != nil {
		return err
	}

	for _, v := range userRoles.List {
		_, err = m.enforcer.AddRoleForUser(userId, v.RoleKey)
		if err != nil {
			return err
		}

		m.user[userId] = append(m.user[userId], v.RoleKey)
	}

	return nil
}
