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
