package permissionx

import (
	"context"
	"fmt"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
)

// 角色资源管理器
type MemoryHolder struct {
	rw sync.RWMutex

	pr permissionrpc.PermissionRpc

	// 角色-权限 key: roleId, value: apiIds
	policy map[string][]string

	// 用户-角色 key: (user), value: role_key[]
	user map[string][]string
}

func NewMemoryHolder(pr permissionrpc.PermissionRpc) *MemoryHolder {

	return &MemoryHolder{
		rw:     sync.RWMutex{},
		pr:     pr,
		user:   make(map[string][]string),
		policy: make(map[string][]string),
	}
}

func (m *MemoryHolder) ReloadPolicy() error {
	return m.LoadPolicy()
}

func (m *MemoryHolder) LoadPolicy() error {
	logx.Info("Reloading permissions...")
	m.rw.Lock()
	defer m.rw.Unlock()
	logx.Info("Reloading permissions... done")

	// 重置所有权限
	m.policy = make(map[string][]string)
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

			m.policy[role.RoleKey] = append(m.policy[role.RoleKey], fmt.Sprintf("%s %s", api.Path, api.Method))
		}
	}

	return nil
}

func (m *MemoryHolder) Enforce(user string, resource string, action string) error {
	m.rw.RUnlock()
	defer m.rw.RUnlock()

	err := m.dynamicLoadUserRoles(user)
	if err != nil {
		return err
	}

	ok, err := m.enforce(user, resource, action)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("permission denied, user: %s, resource: %s, action: %s", user, resource, action)
	}

	return nil
}

// 加载用户角色
func (m *MemoryHolder) dynamicLoadUserRoles(userId string) error {
	if _, ok := m.user[userId]; ok {
		return nil
	}

	// 获取用户角色
	userRoles, err := m.pr.FindUserRoles(context.Background(), &permissionrpc.FindUserRolesReq{UserId: userId})
	if err != nil {
		return err
	}

	for _, v := range userRoles.List {
		m.user[userId] = append(m.user[userId], v.RoleKey)
	}

	return nil
}

func (m *MemoryHolder) enforce(user string, resource string, action string) (bool, error) {
	roles, ok := m.user[user]
	if !ok {
		return false, nil
	}

	p := fmt.Sprintf("%s %s", resource, action)
	for _, role := range roles {
		policies, ok := m.policy[role]
		if !ok {
			continue
		}

		for _, policy := range policies {
			if policy == p {
				return true, nil
			}
		}
	}

	return false, nil
}
