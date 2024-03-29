package logic

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository/svc"
)

type PhotoAlbumRepository struct {
	DbEngin *gorm.DB
	Cache   *redis.Client
}

func NewPhotoAlbumRepository(svcCtx *svc.RepositoryContext) *PhotoAlbumRepository {
	return &PhotoAlbumRepository{
		DbEngin: svcCtx.DbEngin,
		Cache:   svcCtx.Cache,
	}
}

// 创建PhotoAlbum记录
func (s *PhotoAlbumRepository) Create(ctx context.Context, item *entity.PhotoAlbum) (out *entity.PhotoAlbum, err error) {
	db := s.DbEngin.WithContext(ctx)

	err = db.Create(&item).Error
	if err != nil {
		return nil, err
	}
	return item, err
}

// 更新PhotoAlbum记录
func (s *PhotoAlbumRepository) Update(ctx context.Context, item *entity.PhotoAlbum) (out *entity.PhotoAlbum, err error) {
	db := s.DbEngin.WithContext(ctx)

	err = db.Save(&item).Error
	if err != nil {
		return nil, err
	}
	return item, err
}

// 删除PhotoAlbum记录
func (s *PhotoAlbumRepository) Delete(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	query := db.Delete(&entity.PhotoAlbum{})
	err = query.Error
	rows = query.RowsAffected
	return rows, err
}

// 查询PhotoAlbum记录
func (s *PhotoAlbumRepository) First(ctx context.Context, conditions string, args ...interface{}) (out *entity.PhotoAlbum, err error) {
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

func (s *PhotoAlbumRepository) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*entity.PhotoAlbum, err error) {
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

// 分页查询PhotoAlbum记录
func (s *PhotoAlbumRepository) FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*entity.PhotoAlbum, err error) {
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
func (s *PhotoAlbumRepository) Count(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&entity.PhotoAlbum{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
