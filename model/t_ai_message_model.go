package model

import (
	"time"

	"gorm.io/gorm"
)

var _ TAiMessageModel = (*defaultTAiMessageModel)(nil)

type (
	// 接口定义
	TAiMessageModel interface {
		Table[TAiMessage]
		// 在事务中操作
		WithTx(tx *gorm.DB) (out TAiMessageModel)
	}

	// 表字段定义
	TAiMessage struct {
		Id               int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:主键" json:"id"`                                     // 主键
		SessionId        int64     `gorm:"column:session_id;type:bigint;not null;index:idx_session,priority:1;default:0;comment:会话id" json:"session_id"` // 会话id
		Role             string    `gorm:"column:role;type:varchar(16);not null;default:'';comment:角色：user/assistant" json:"role"`                       // 角色：user/assistant
		Content          string    `gorm:"column:content;type:longtext;not null;comment:正文" json:"content"`                                              // 正文
		PromptTokens     int64     `gorm:"column:prompt_tokens;type:int;not null;default:0;comment:提示词token数" json:"prompt_tokens"`                      // 提示词token数
		CompletionTokens int64     `gorm:"column:completion_tokens;type:int;not null;default:0;comment:生成token数" json:"completion_tokens"`               // 生成token数
		CreatedAt        time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`            // 创建时间
	}

	// 接口实现
	defaultTAiMessageModel struct {
		baseModel[TAiMessage]
	}
)

func NewTAiMessageModel(db *gorm.DB) TAiMessageModel {
	return &defaultTAiMessageModel{newBaseModel[TAiMessage](db, "`t_ai_message`")}
}

// 在事务中操作
func (m *defaultTAiMessageModel) WithTx(tx *gorm.DB) (out TAiMessageModel) {
	return &defaultTAiMessageModel{newBaseModel[TAiMessage](tx, m.table)}
}
