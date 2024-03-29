package entity

import "time"

// TableNameComment return the table name of <comment>
const TableNameComment = "comment"

// Comment mapped from table <comment>
type Comment struct {
	ID             int       `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement:true;comment:主键" json:"id"`                       // 主键
	UserID         int       `gorm:"column:user_id;type:int;not null;index:fk_comment_user,priority:1;comment:评论用户Id" json:"user_id"`      // 评论用户Id
	TopicID        int       `gorm:"column:topic_id;type:int;not null;comment:评论主题id" json:"topic_id"`                                     // 评论主题id
	CommentContent string    `gorm:"column:comment_content;type:text;not null;comment:评论内容" json:"comment_content"`                        // 评论内容
	ReplyUserID    int       `gorm:"column:reply_user_id;type:int;not null;comment:回复用户id" json:"reply_user_id"`                           // 回复用户id
	ParentID       int       `gorm:"column:parent_id;type:int;not null;index:fk_comment_parent,priority:1;comment:父评论id" json:"parent_id"` // 父评论id
	Type           int       `gorm:"column:type;type:tinyint;not null;comment:评论类型 1.文章 2.友链 3.说说" json:"type"`                            // 评论类型 1.文章 2.友链 3.说说
	IsDelete       int       `gorm:"column:is_delete;type:tinyint;not null;comment:是否删除  0否 1是" json:"is_delete"`                          // 是否删除  0否 1是
	IsReview       int       `gorm:"column:is_review;type:tinyint(1);not null;default:1;comment:是否审核" json:"is_review"`                    // 是否审核
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:评论时间" json:"created_at"`    // 评论时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`    // 更新时间
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}
