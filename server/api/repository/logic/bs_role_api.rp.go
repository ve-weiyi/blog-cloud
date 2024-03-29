package logic

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository/svc"
)

type RoleApiRepository struct {
	DbEngin *gorm.DB
	Cache   *redis.Client
}

func NewRoleApiRepository(svcCtx *svc.RepositoryContext) *RoleApiRepository {
	return &RoleApiRepository{
		DbEngin: svcCtx.DbEngin,
		Cache:   svcCtx.Cache,
	}
}

// 创建RoleApi记录
func (s *RoleApiRepository) Create(ctx context.Context, item *entity.RoleApi) (out *entity.RoleApi, err error) {
	db := s.DbEngin.WithContext(ctx)

	err = db.Create(&item).Error
	if err != nil {
		return nil, err
	}
	return item, err
}

// 更新RoleApi记录
func (s *RoleApiRepository) Update(ctx context.Context, item *entity.RoleApi) (out *entity.RoleApi, err error) {
	db := s.DbEngin.WithContext(ctx)

	err = db.Save(&item).Error
	if err != nil {
		return nil, err
	}
	return item, err
}

// 删除RoleApi记录
func (s *RoleApiRepository) Delete(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	query := db.Delete(&entity.RoleApi{})
	err = query.Error
	rows = query.RowsAffected
	return rows, err
}

// 查询RoleApi记录
func (s *RoleApiRepository) First(ctx context.Context, conditions string, args ...interface{}) (out *entity.RoleApi, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.First(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

func (s *RoleApiRepository) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*entity.RoleApi, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 分页查询RoleApi记录
func (s *RoleApiRepository) FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*entity.RoleApi, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)

	// 如果有搜索条件
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	// 如果有排序参数
	if len(sorts) != 0 {
		db = db.Order(sorts)
	}

	// 如果有分页参数
	if page > 0 && size > 0 {
		limit := size
		offset := (page - 1) * limit
		db = db.Limit(limit).Offset(offset)
	}

	// 查询数据
	err = db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 查询总数
func (s *RoleApiRepository) Count(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&entity.RoleApi{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
