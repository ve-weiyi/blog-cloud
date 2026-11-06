package notification

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MarkRead 把一组条目置为已读，返回真正被改动的条数。
// 收件人由调用方从会话主体取，方法不接受来自请求的用户标识。
func (s *Service) MarkRead(ctx context.Context, recipientID, groupKey string) (int64, error) {
	return s.notifs.UpdateFields(ctx,
		map[string]interface{}{"read_at": time.Now()},
		"recipient_id = ? AND group_key = ? AND read_at IS NULL", recipientID, groupKey)
}

// MarkAllRead 把该收件人的全部未读置为已读。
func (s *Service) MarkAllRead(ctx context.Context, recipientID string) (int64, error) {
	return s.notifs.UpdateFields(ctx,
		map[string]interface{}{"read_at": time.Now()},
		"recipient_id = ? AND read_at IS NULL", recipientID)
}

// DeleteGroup 删除一组条目及其投递记录。
//
// 条目与投递必须一起删：只删条目会留下孤儿投递行，消费者仍可能据它把内容发出去。
func (s *Service) DeleteGroup(ctx context.Context, recipientID, groupKey string) (int64, error) {
	var affected int64

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rows, err := s.notifs.WithTx(tx).FindAll(ctx,
			"recipient_id = ? AND group_key = ?", recipientID, groupKey)
		if err != nil {
			return fmt.Errorf("notification: query entries: %w", err)
		}
		if len(rows) == 0 {
			return nil
		}

		ids := make([]int64, 0, len(rows))
		for _, row := range rows {
			ids = append(ids, row.Id)
		}

		if _, err := s.deliveries.WithTx(tx).DeleteBatch(ctx, "notification_id IN ?", ids); err != nil {
			return fmt.Errorf("notification: delete deliveries: %w", err)
		}

		affected, err = s.notifs.WithTx(tx).DeleteBatch(ctx,
			"recipient_id = ? AND group_key = ?", recipientID, groupKey)
		return err
	})
	if err != nil {
		return 0, err
	}
	return affected, nil
}
