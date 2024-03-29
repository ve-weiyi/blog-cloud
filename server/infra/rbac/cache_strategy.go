package rbac

import (
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
)

type UserPermission struct {
	uid   string
	Roles []*entity.Role
}

// 判断资源需要哪些角色（Resource-Based Access Control，RBAC）
type ApiPermission struct {
	entity.Api
	Group string
	Roles []*entity.Role
}

// 判断角色拥有哪些资源（Role-Based Access Control，RBAC）
type RolePermission struct {
	entity.Role
	Apis  []*entity.Api
	Menus []*entity.Menu
}

type CacheStrategy interface {
	GetUserPermission(uid string) (*UserPermission, error)
	GetApiPermission(api string) (*ApiPermission, error)
	GetRolePermission(role string) (*RolePermission, error)
	SetUserPermission(uid string, permission *UserPermission) error
	SetApiPermission(api string, permission *ApiPermission) error
	SetRolePermission(role string, permission *RolePermission) error
}

type RuntimeStrategy struct {
	UserPermissions map[string]*UserPermission
	ApiPermissions  map[string]*ApiPermission
	RolePermissions map[string]*RolePermission
}

func (s *RuntimeStrategy) GetUserPermission(uid string) (*UserPermission, error) {
	value, ok := s.UserPermissions[uid]
	if ok {
		return value, nil
	}
	return nil, fmt.Errorf("record not found")
}

func (s *RuntimeStrategy) GetApiPermission(api string) (*ApiPermission, error) {
	value, ok := s.ApiPermissions[api]
	if ok {
		return value, nil
	}
	return nil, fmt.Errorf("record not found")
}

func (s *RuntimeStrategy) GetRolePermission(role string) (*RolePermission, error) {
	value, ok := s.RolePermissions[role]
	if ok {
		return value, nil
	}
	return nil, fmt.Errorf("record not found")
}

func (s *RuntimeStrategy) SetUserPermission(uid string, permission *UserPermission) error {
	s.UserPermissions[uid] = permission
	return nil
}

func (s *RuntimeStrategy) SetApiPermission(api string, permission *ApiPermission) error {
	s.ApiPermissions[api] = permission
	return nil
}

func (s *RuntimeStrategy) SetRolePermission(role string, permission *RolePermission) error {
	s.RolePermissions[role] = permission
	return nil
}

func NewCacheStrategy() CacheStrategy {
	return &RuntimeStrategy{
		UserPermissions: map[string]*UserPermission{},
		ApiPermissions:  map[string]*ApiPermission{},
	}
}
