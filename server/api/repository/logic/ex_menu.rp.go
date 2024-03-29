package logic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
)

// 清空菜单
func (s *MenuRepository) CleanMenus(ctx context.Context) (data interface{}, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)

	err = ClearTable(db, entity.TableNameMenu)
	if err != nil {
		return nil, err
	}

	err = ClearTable(db, entity.TableNameRoleMenu)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
