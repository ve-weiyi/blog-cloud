package logic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
)

// 获取Api记录
func (s *ApiRepository) FindApiRoles(ctx context.Context, apiId int) (list []*entity.Role, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)
	var roleApis []*entity.RoleApi

	err = db.Where("api_id = ?", apiId).Find(&roleApis).Error
	if err != nil {
		return nil, err
	}

	var roleIds []int
	for _, item := range roleApis {
		roleIds = append(roleIds, item.ApiID)
	}

	var roles []*entity.Role
	err = db.Where("id in (?)", roleIds).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// 清空菜单
func (s *ApiRepository) CleanApis(ctx context.Context) (data interface{}, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)

	err = ClearTable(db, entity.TableNameApi)
	if err != nil {
		return nil, err
	}

	err = ClearTable(db, entity.TableNameRoleApi)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
