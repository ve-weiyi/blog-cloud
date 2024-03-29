package logic

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository/svc"
)

type ArticleTagRepository struct {
	DbEngin *gorm.DB
	Cache   *redis.Client
}

func NewArticleTagRepository(svcCtx *svc.RepositoryContext) *ArticleTagRepository {
	return &ArticleTagRepository{
		DbEngin: svcCtx.DbEngin,
		Cache:   svcCtx.Cache,
	}
}

// 创建ArticleTag记录
func (s *ArticleTagRepository) Create(ctx context.Context, item *entity.ArticleTag) (out *entity.ArticleTag, err error) {
	db := s.DbEngin.WithContext(ctx)

	err = db.Create(&item).Error
	if err != nil {
		return nil, err
	}
	return item, err
}

// 更新ArticleTag记录
func (s *ArticleTagRepository) Update(ctx context.Context, item *entity.ArticleTag) (out *entity.ArticleTag, err error) {
	db := s.DbEngin.WithContext(ctx)

	err = db.Save(&item).Error
	if err != nil {
		return nil, err
	}
	return item, err
}

// 删除ArticleTag记录
func (s *ArticleTagRepository) Delete(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	query := db.Delete(&entity.ArticleTag{})
	err = query.Error
	rows = query.RowsAffected
	return rows, err
}

// 查询ArticleTag记录
func (s *ArticleTagRepository) First(ctx context.Context, conditions string, args ...interface{}) (out *entity.ArticleTag, err error) {
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

func (s *ArticleTagRepository) FindALL(ctx context.Context, conditions string, args ...interface{}) (out []*entity.ArticleTag, err error) {
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

// 分页查询ArticleTag记录
func (s *ArticleTagRepository) FindList(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*entity.ArticleTag, err error) {
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
func (s *ArticleTagRepository) Count(ctx context.Context, conditions string, args ...interface{}) (count int64, err error) {
	db := s.DbEngin.WithContext(ctx)

	// 如果有条件语句
	if len(conditions) != 0 {
		db = db.Where(conditions, args...)
	}

	err = db.Model(&entity.ArticleTag{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
