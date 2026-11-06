package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ TLikeModel = (*defaultTLikeModel)(nil)

type (
	// 接口定义
	TLikeModel interface {
		Table[TLike]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TLikeModel)
		// add extra method in here
		FindOneByTargetTypeTargetIdUserId(ctx context.Context, targetType int64, targetId int64, userId string) (out *TLike, err error)
		// Like 点赞，返回本次调用是否改变了状态；重复点赞返回 false
		Like(ctx context.Context, targetType int64, targetId int64, userId string) (changed bool, err error)
		// Unlike 取消点赞，返回本次调用是否改变了状态；未点赞返回 false
		Unlike(ctx context.Context, targetType int64, targetId int64, userId string) (changed bool, err error)
		// Liked 返回目标 ID 到「当前用户是否已赞」的映射；未登录时全部为未赞
		Liked(ctx context.Context, targetType int64, userId string, ids []int64) (map[int64]bool, error)
		// Counts 返回目标 ID 到点赞数的映射；无人点赞的目标不出现在结果里
		Counts(ctx context.Context, targetType int64, ids []int64) (map[int64]int64, error)
	}

	// 表字段定义
	TLike struct {
		Id         int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                                                                                                                                   // 主键
		TargetType int64     `gorm:"column:target_type;type:tinyint;not null;uniqueIndex:uk_target_user,priority:1;index:idx_user_target,priority:2;index:idx_user_created,priority:2;default:0;comment:目标类型 1.文章 2.说说 3.评论" json:"target_type"` // 目标类型 1.文章 2.说说 3.评论
		TargetId   int64     `gorm:"column:target_id;type:bigint;not null;uniqueIndex:uk_target_user,priority:2;index:idx_user_target,priority:3;default:0;comment:目标id" json:"target_id"`                                                       // 目标id
		UserId     string    `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_target_user,priority:3;index:idx_user_target,priority:1;index:idx_user_created,priority:1;default:'';comment:点赞用户id" json:"user_id"`                 // 点赞用户id
		CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;index:idx_user_created,priority:3;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                                        // 创建时间
	}

	// 接口实现
	defaultTLikeModel struct {
		baseModel[TLike]
	}
)

func NewTLikeModel(db *gorm.DB) TLikeModel {
	return &defaultTLikeModel{newBaseModel[TLike](db, "`t_like`")}
}

// 在事务中操作
func (m *defaultTLikeModel) WithTx(tx *gorm.DB) (out TLikeModel) {
	return &defaultTLikeModel{newBaseModel[TLike](tx, m.table)}
}

// add extra method in here
func (m *defaultTLikeModel) FindOneByTargetTypeTargetIdUserId(ctx context.Context, targetType int64, targetId int64, userId string) (out *TLike, err error) {
	db := m.scope(ctx)

	err = db.Where("`target_type` = ? and `target_id` = ? and `user_id` = ?", targetType, targetId, userId).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

// 点赞以 t_like 明细为唯一状态：是否已赞与点赞数都在读取时派生，不维护计数列。

// Like 点赞，返回本次调用是否改变了状态；重复点赞返回 false。
//
// 去重由唯一键承担：冲突时插入不修改任何列，受影响行数为 0。
func (m *defaultTLikeModel) Like(ctx context.Context, targetType int64, targetId int64, userId string) (bool, error) {
	rows, err := m.InsertIgnore(ctx, &TLike{
		TargetType: targetType,
		TargetId:   targetId,
		UserId:     userId,
	})
	return rows == 1, err
}

// Unlike 取消点赞，返回本次调用是否改变了状态；未点赞返回 false。
func (m *defaultTLikeModel) Unlike(ctx context.Context, targetType int64, targetId int64, userId string) (bool, error) {
	rows, err := m.DeleteBatch(ctx, "target_type = ? AND target_id = ? AND user_id = ?", targetType, targetId, userId)
	return rows == 1, err
}

// Liked 返回目标 ID 到「当前用户是否已赞」的映射；未登录时全部为未赞。
func (m *defaultTLikeModel) Liked(ctx context.Context, targetType int64, userId string, ids []int64) (map[int64]bool, error) {
	liked := make(map[int64]bool, len(ids))
	if userId == "" || len(ids) == 0 {
		return liked, nil
	}

	list, err := m.FindAll(ctx, "user_id = ? AND target_type = ? AND target_id IN ?", userId, targetType, ids)
	if err != nil {
		return nil, err
	}
	for _, row := range list {
		liked[row.TargetId] = true
	}
	return liked, nil
}

// Counts 返回目标 ID 到点赞数的映射；无人点赞的目标不出现在结果里。
// 分组计数不属于通用 CRUD 契约，直接按本表下发。
func (m *defaultTLikeModel) Counts(ctx context.Context, targetType int64, ids []int64) (map[int64]int64, error) {
	counts := make(map[int64]int64, len(ids))
	if len(ids) == 0 {
		return counts, nil
	}

	var rows []struct {
		TargetId int64
		Total    int64
	}
	if err := m.scope(ctx).
		Select("target_id, COUNT(*) AS total").
		Where("target_type = ? AND target_id IN ?", targetType, ids).
		Group("target_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		counts[row.TargetId] = row.Total
	}
	return counts, nil
}
