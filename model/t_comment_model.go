package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TCommentModel = (*defaultTCommentModel)(nil)

type (
	// 接口定义
	TCommentModel interface {
		Table[TComment]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TCommentModel)
	}

	// 表字段定义
	TComment struct {
		Id             int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                                                              // 主键
		UserId         string    `gorm:"column:user_id;type:varchar(64);not null;index:fk_comment_user,priority:1;default:'';comment:评论用户id" json:"user_id"`                    // 评论用户id
		DeviceId       string    `gorm:"column:device_id;type:varchar(64);not null;default:'';comment:设备id" json:"device_id"`                                                   // 设备id
		TopicId        int64     `gorm:"column:topic_id;type:bigint;not null;index:idx_type_topic_created,priority:2;default:0;comment:主题id" json:"topic_id"`                   // 主题id
		ParentId       int64     `gorm:"column:parent_id;type:bigint;not null;index:fk_comment_parent,priority:1;default:0;comment:父评论id" json:"parent_id"`                     // 父评论id
		ReplyId        int64     `gorm:"column:reply_id;type:bigint;not null;default:0;comment:回复评论id" json:"reply_id"`                                                         // 回复评论id
		ReplyUserId    string    `gorm:"column:reply_user_id;type:varchar(64);not null;default:'';comment:评论回复用户id" json:"reply_user_id"`                                       // 评论回复用户id
		CommentContent string    `gorm:"column:comment_content;type:text;not null;default:'';comment:评论内容" json:"comment_content"`                                              // 评论内容
		Type           int64     `gorm:"column:type;type:tinyint;not null;index:idx_type_topic_created,priority:1;default:0;comment:评论类型 1.文章 2.友链 3.说说" json:"type"`           // 评论类型 1.文章 2.友链 3.说说
		Status         int64     `gorm:"column:status;type:tinyint;not null;default:0;comment:状态 0.正常 1.已编辑 2.已删除" json:"status"`                                               // 状态 0.正常 1.已编辑 2.已删除
		IsReview       int64     `gorm:"column:is_review;type:tinyint;not null;index:idx_review_created,priority:1;default:0;comment:审核状态：0-未审核 1-已通过" json:"is_review"`        // 审核状态：0-未审核 1-已通过
		CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;index:idx_review_created,priority:2;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
		UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                     // 更新时间
	}

	// 接口实现
	defaultTCommentModel struct {
		baseModel[TComment]
	}
)

func NewTCommentModel(db *gorm.DB) TCommentModel {
	return &defaultTCommentModel{newBaseModel[TComment](db, "`t_comment`")}
}

// 在事务中操作
func (m *defaultTCommentModel) WithTx(tx *gorm.DB) (out TCommentModel) {
	return &defaultTCommentModel{newBaseModel[TComment](tx, m.table)}
}
