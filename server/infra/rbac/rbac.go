package rbac

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/glog"
)

type RbacHolder interface {
	// 检查用户是否拥有访问接口权限
	CheckUserAccessApi(uid string, path string, method string) error
	// 查询接口权限信息
	FindApiPermission(path string, method string) (*ApiPermission, error)
	// 查询用户权限信息
	FindUserPermission(uid string) (*UserPermission, error)
	// 查询用户有哪些角色
	//FindUserHasRoles(uid string) (roles []int, err error)
	// 查询接口需要哪些角色
	//FindApiNeedRoles(path string, method string) (roles []int, err error)
	// 查询接口是否可用
	//IsApiOpen(path string, method string) (bool, error)
	// 查询接口是否需要记录操作日志 (所有接口都会记录访问日志，但不是所有接口都需要记录操作日志)
	//IsApiTraceable(path string, method string) (bool, error)
}

type PermissionHolder struct {
	DbEngin    *gorm.DB
	CacheEngin CacheStrategy
	logger     glog.Logger
	debug      bool
}

func (s *PermissionHolder) CheckUserAccessApi(uid string, path string, method string) error {
	ap, err := s.FindApiPermission(path, method)
	if err != nil {
		return err
	}

	if ap == nil {
		return nil
	}

	up, err := s.FindUserPermission(uid)
	if err != nil {
		return err
	}

	if len(ap.Roles) == 0 {
		return nil
	}
	// 遍历资源角色
	for _, ar := range ap.Roles {
		// 匹配用户角色
		for _, ur := range up.Roles {
			if ar == ur {
				return nil
			}
		}
	}

	return fmt.Errorf("permissions not match,user:%+v,api:%+v", up.Roles, ap.Roles)
}

func (s *PermissionHolder) FindUserPermission(uid string) (*UserPermission, error) {
	// 从缓存查找
	permission, err := s.CacheEngin.GetUserPermission(uid)
	if err != nil {
		s.info("load user from database:error with %v", err)
		// 加载api
		permission, err = s.LoadUser(uid)
		if err != nil {
			return nil, err
		}
		s.info("find user from database:%+v", permission)
		return permission, nil
	}

	s.info("load user from cache:%+v", permission)
	return permission, nil
}

func (s *PermissionHolder) FindApiPermission(path string, method string) (*ApiPermission, error) {
	api := fmt.Sprintf("%v-%v", path, method)
	// 从缓存查找
	permission, err := s.CacheEngin.GetApiPermission(api)
	if err != nil {
		s.info("load api from database:error with %v", err)
		// 加载api
		permission, err = s.LoadApi(path, method)
		if err != nil {
			return nil, err
		}
		s.info("find api from database:%+v", permission)
		return permission, nil
	}

	s.info("load api from cache:%+v", permission)
	return permission, nil
}

// 加载用户
func (s *PermissionHolder) LoadUser(uid string) (*UserPermission, error) {
	db := s.DbEngin
	// 查询用户角色
	var userApis []entity.UserRole
	err := db.Where("user_id = ?", uid).First(&userApis).Error
	if err != nil {
		return nil, err
	}

	var roleIds []int
	for _, item := range userApis {
		roleIds = append(roleIds, item.UserID)
	}

	var roles []*entity.Role
	err = db.Where("id in (?)", roleIds).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// 保存到缓存
	permission := &UserPermission{uid: uid, Roles: roles}
	err = s.CacheEngin.SetUserPermission(uid, permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// 加载api
func (s *PermissionHolder) LoadApi(path string, method string) (*ApiPermission, error) {
	db := s.DbEngin

	// 查询接口
	var api entity.Api
	err := db.Where("path = ? and method = ?", path, method).First(&api).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 查询接口分组
	var parent entity.Api
	err = db.Where("id", api.ParentID).First(&parent).Error
	if err != nil {
		return nil, err
	}

	// 查询接口角色
	var roleApis []*entity.RoleApi
	err = db.Where("api_id = ?", api.ID).Find(&roleApis).Error
	if err != nil {
		return nil, err
	}

	var roleIds []int
	for _, item := range roleApis {
		roleIds = append(roleIds, item.RoleID)
	}

	var roles []*entity.Role
	err = db.Where("id in (?)", roleIds).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// 保存到缓存
	key := fmt.Sprintf("%v-%v", path, method)
	permission := &ApiPermission{Api: api, Group: parent.Name, Roles: roles}
	err = s.CacheEngin.SetApiPermission(key, permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// 加载角色
func (s *PermissionHolder) LoadRole(rid string) (*RolePermission, error) {
	db := s.DbEngin

	// 查询角色菜单
	var role entity.Role
	err := db.Where("id = ?", rid).First(&role).Error
	if err != nil {
		return nil, err
	}

	// 查询接口
	var roleApis []*entity.RoleApi
	err = db.Where("role_id = ?", role.ID).First(&roleApis).Error
	if err != nil {
		return nil, err
	}

	var apiIds []int
	for _, item := range roleApis {
		apiIds = append(apiIds, item.ApiID)
	}

	var apis []*entity.Api
	err = db.Where("id in (?)", apiIds).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	// 查询接口
	var roleMenus []*entity.RoleMenu
	err = db.Where("role_id = ?", role.ID).Find(&roleMenus).Error
	if err != nil {
		return nil, err
	}

	var menuIds []int
	for _, item := range roleMenus {
		apiIds = append(apiIds, item.MenuID)
	}

	var menus []*entity.Menu
	err = db.Where("id in (?)", menuIds).Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 保存到缓存
	permission := &RolePermission{Role: role, Apis: apis, Menus: menus}
	err = s.CacheEngin.SetRolePermission(rid, permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *PermissionHolder) info(format string, args ...interface{}) {
	if !s.debug || s.logger == nil {
		return
	}
	s.logger.Infof(format, args...)
}

func NewPermissionHolder(db *gorm.DB, logger glog.Logger) RbacHolder {
	return &PermissionHolder{
		DbEngin:    db,
		CacheEngin: NewCacheStrategy(),
		logger:     logger,
		debug:      false,
	}
}
