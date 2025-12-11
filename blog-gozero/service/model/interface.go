package model

import (
	"context"
)

type IModel[T any] interface {
	// 插入
	Insert(ctx context.Context, in *T) (rows int64, err error)
	Inserts(ctx context.Context, in ...*T) (rows int64, err error)
	// 删除
	Delete(ctx context.Context, id int64) (rows int64, err error)
	Deletes(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
	// 更新
	Update(ctx context.Context, in *T) (rows int64, err error)
	Updates(ctx context.Context, columns map[string]interface{}, conditions string, args ...interface{}) (rows int64, err error)
	// 保存
	Save(ctx context.Context, in *T) (rows int64, err error)
	// 查询
	FindById(ctx context.Context, id int64) (out *T, err error)
	FindALL(ctx context.Context, conditions string, args ...interface{}) (list []*T, err error)
	FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) (list []*T, total int64, err error)
	FindCount(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
}

var _ IModel[TArticle] = (*defaultTArticleModel)(nil)
