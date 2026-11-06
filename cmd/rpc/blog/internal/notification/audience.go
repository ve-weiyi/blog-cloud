package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
)

// recipient 是本批的一个收件人及其投递地址。
type recipient struct {
	ID     int64 // 数字主键，仅 all 受众用于游标
	UserID string
	Email  string
	Mobile string
}

// countAudience 计算受众总数，作为进度分母。只在扇出开始时算一次，之后不再变。
func countAudience(ctx context.Context, tx *gorm.DB, table string, ev *model.TNotificationEvent) (int64, error) {
	if ev.AudienceType == AudienceUserIDs {
		ids, err := decodeAudience(ev.AudienceIds)
		if err != nil {
			return 0, err
		}
		// 与 fetchRecipients 用同一套清洗，否则分母把触发者与重复值也算进去
		return int64(len(NormalizeRecipients(ids, ev.ActorId))), nil
	}

	q := tx.WithContext(ctx).Table(table).Where("status = ?", enums.UserStatusNormal)
	if ev.ActorId != "" {
		// 与 fetchRecipients 一致：触发者不接收自己的通知
		q = q.Where("user_id <> ?", ev.ActorId)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// fetchRecipients 按游标取一批收件人，多取一条用于判断是否已穷尽。
func fetchRecipients(ctx context.Context, tx *gorm.DB, table string, ev *model.TNotificationEvent) (batch []recipient, nextCursor int64, exhausted bool, err error) {
	if ev.AudienceType == AudienceUserIDs {
		return fetchByIDs(ev)
	}

	rows := make([]recipient, 0, fanoutBatchSize+1)
	q := tx.WithContext(ctx).Table(table).
		Select("id, user_id, email, mobile").
		Where("status = ? AND id > ?", enums.UserStatusNormal, ev.Cursor).
		Order("id ASC").
		Limit(fanoutBatchSize + 1)
	if ev.ActorId != "" {
		q = q.Where("user_id <> ?", ev.ActorId)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, ev.Cursor, false, err
	}

	if len(rows) > fanoutBatchSize {
		rows = rows[:fanoutBatchSize]
		exhausted = false
	} else {
		exhausted = true
	}
	nextCursor = ev.Cursor
	if len(rows) > 0 {
		nextCursor = rows[len(rows)-1].ID
	}
	return rows, nextCursor, exhausted, nil
}

// fetchByIDs 处理定向受众：游标是已处理到的下标。
func fetchByIDs(ev *model.TNotificationEvent) ([]recipient, int64, bool, error) {
	ids, err := decodeAudience(ev.AudienceIds)
	if err != nil {
		return nil, ev.Cursor, false, err
	}
	ids = NormalizeRecipients(ids, ev.ActorId)

	start := int(ev.Cursor)
	if start > len(ids) {
		start = len(ids)
	}

	end := start + fanoutBatchSize
	exhausted := true
	if end < len(ids) {
		exhausted = false
	} else {
		end = len(ids)
	}

	rows := make([]recipient, 0, end-start)
	for _, id := range ids[start:end] {
		rows = append(rows, recipient{UserID: id})
	}
	return rows, int64(end), exhausted, nil
}

func decodeAudience(raw *string) ([]string, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	var ids []string
	if err := json.Unmarshal([]byte(*raw), &ids); err != nil {
		return nil, fmt.Errorf("notification: decode audience: %w", err)
	}
	return ids, nil
}

// fillAddress 为定向受众补齐邮箱与手机号。
func fillAddress(ctx context.Context, tx *gorm.DB, table string, address map[string]recipient) error {
	ids := make([]string, 0, len(address))
	for id := range address {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil
	}

	rows := make([]recipient, 0, len(ids))
	if err := tx.WithContext(ctx).Table(table).
		Select("id, user_id, email, mobile").
		Where("user_id IN ?", ids).
		Find(&rows).Error; err != nil {
		return fmt.Errorf("notification: load addresses: %w", err)
	}
	for _, r := range rows {
		address[r.UserID] = r
	}
	return nil
}
