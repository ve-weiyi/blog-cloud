package notificationservicelogic

import (
	"time"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

// ============================================
// NotifyTemplate 转换函数
// ============================================

func convertTNotifyTemplateToProto(m *model.TNotifyTemplate) *notificationrpc2.NotifyTemplate {
	if m == nil {
		return nil
	}

	return &notificationrpc2.NotifyTemplate{
		Id:        m.Id,
		Code:      m.Code,
		Channel:   m.Channel,
		Scene:     m.Scene,
		Title:     m.Title,
		Content:   m.Content,
		Enabled:   m.Enabled,
		CreatedAt: m.CreatedAt.UnixMilli(),
		UpdatedAt: m.UpdatedAt.UnixMilli(),
	}
}

func convertProtoToTNotifyTemplate(p *notificationrpc2.CreateNotifyTemplateRequest) *model.TNotifyTemplate {
	if p == nil {
		return nil
	}

	return &model.TNotifyTemplate{
		Code:      p.Code,
		Channel:   p.Channel,
		Scene:     p.Scene,
		Title:     p.Title,
		Content:   p.Content,
		Enabled:   p.Enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ============================================
// NotifyMessage 转换函数
// ============================================

func convertTNotifyMessageToProto(m *model.TNotifyMessage) *notificationrpc.NotifyMessage {
	if m == nil {
		return nil
	}

	msg := &notificationrpc.NotifyMessage{
		Id:          m.Id,
		Title:       m.Title,
		Content:     stringPtr(m.Content),
		Category:    m.Category,
		Level:       m.Level,
		TargetType:  m.TargetType,
		TargetIds:   stringPtr(m.TargetIds),
		Status:      m.Status,
		PublishedBy: m.PublishedBy,
		CreatedAt:   m.CreatedAt.UnixMilli(),
		UpdatedAt:   m.UpdatedAt.UnixMilli(),
	}

	if m.PublishedAt != nil {
		msg.PublishedAt = m.PublishedAt.UnixMilli()
	}

	return msg
}

func convertProtoToTNotifyMessage(p *notificationrpc.CreateNotifyMessageRequest) *model.TNotifyMessage {
	if p == nil {
		return nil
	}

	return &model.TNotifyMessage{
		Title:      p.Title,
		Content:    stringToPtr(p.Content),
		Category:   p.Category,
		Level:      p.Level,
		TargetType: p.TargetType,
		TargetIds:  stringToPtr(p.TargetIds),
		Status:     "draft",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// ============================================
// NotifyRecord 转换函数
// ============================================

func convertTNotifyRecordToProto(m *model.TNotifyRecord) *notificationrpc.NotifyRecord {
	if m == nil {
		return nil
	}

	d := &notificationrpc.NotifyRecord{
		Id:           m.Id,
		MessageId:    m.MessageId,
		Channel:      m.Channel,
		Recipient:    m.Recipient,
		TemplateCode: m.TemplateCode,
		Content:      stringPtr(m.Content),
		Status:       m.Status,
		BizId:        m.BizId,
		ErrorMsg:     m.ErrorMsg,
		CreatedAt:    m.CreatedAt.UnixMilli(),
	}

	if m.ReadAt != nil {
		d.ReadAt = m.ReadAt.UnixMilli()
	}
	if m.SentAt != nil {
		d.SentAt = m.SentAt.UnixMilli()
	}

	return d
}

func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
