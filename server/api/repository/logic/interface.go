package logic

import (
	"context"
)

// 定义泛型接口
type Repository[T any] interface {
	// 增删改查
	Create(ctx context.Context, item T) (out T, err error)
	Update(ctx context.Context, item T) (out T, err error)
	Delete(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
	First(ctx context.Context, conditions string, args ...interface{}) (out T, err error)
	FindALL(ctx context.Context, conditions string, args ...interface{}) (list []T, err error)
	FindList(ctx context.Context, limit int, offset int, sorts string, conditions string, args ...interface{}) (list []T, err error)
	Count(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
}

// 定义泛型接口
type IModel[T any] interface {
	// 增删改查
	Create(ctx context.Context, in T) (out T, err error)
	Update(ctx context.Context, in T) (out T, err error)
	Delete(ctx context.Context, id int) (rows int64, err error)
	First(ctx context.Context, conditions string, args ...interface{}) (out T, err error)
	// 批量操作
	BatchCreate(ctx context.Context, in ...T) (rows int64, err error)
	BatchDelete(ctx context.Context, conditions string, args ...interface{}) (rows int64, err error)
	// 查询
	Count(ctx context.Context, conditions string, args ...interface{}) (count int64, err error)
	FindALL(ctx context.Context, conditions string, args ...interface{}) (list []T, err error)
	FindList(ctx context.Context, limit int, offset int, sorts string, conditions string, args ...interface{}) (list []T, err error)
}
