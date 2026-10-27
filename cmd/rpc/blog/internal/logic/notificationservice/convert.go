package notificationservicelogic

import (
	"strings"
	"time"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/model"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/vkit/x/jsonv"
)

// ============================================
// NotifyTemplate 转换函数
// ============================================

func convertTNotifyTemplateToProto(m *model.TNotifyTemplate) *notificationrpc.NotifyTemplate {
	if m == nil {
		return nil
	}

	return &notificationrpc.NotifyTemplate{
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

func convertProtoToTNotifyTemplate(p *notificationrpc.CreateNotifyTemplateRequest) *model.TNotifyTemplate {
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
		TargetIds:   targetIDsFromJSON(m.TargetIds),
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
		TargetIds:  targetIDsToJSON(p.TargetIds),
		Status:     enums.NoticeStatusDraft,
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

// targetIDsToJSON 把逗号分隔的目标用户列表转成 JSON 数组；空内容写 NULL。
func targetIDsToJSON(raw string) *string {
	ids := splitTargetIDs(raw)
	if len(ids) == 0 {
		return nil
	}
	s, err := jsonv.AnyToJson(ids)
	if err != nil {
		return nil
	}
	return &s
}

// targetIDsFromJSON 把 JSON 数组还原成逗号分隔字符串。
func targetIDsFromJSON(raw *string) string {
	if raw == nil || *raw == "" {
		return ""
	}
	var ids []string
	if err := jsonv.JsonToAny(*raw, &ids); err != nil {
		return ""
	}
	return strings.Join(ids, ",")
}

func splitTargetIDs(raw string) []string {
	parts := strings.Split(raw, ",")
	ids := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			ids = append(ids, p)
		}
	}
	return ids
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
