package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// ErrEventNotRedeliverable 表示事件已走完或被终止，重放会让扇出期间新注册的用户收到旧通知。
var ErrEventNotRedeliverable = errors.New("notification: event is not redeliverable")

// ErrDeliveryNotRedeliverable 表示投递记录已发送或已取消，不应再发一次。
var ErrDeliveryNotRedeliverable = errors.New("notification: delivery is not redeliverable")

// ErrAudienceEmpty 表示定向公告没有受众：发出去不会有人收到，属于参数问题而非系统故障。
var ErrAudienceEmpty = errors.New("notification: audience is empty")

// DeleteAnnouncements 删除公告内容，并回收派生条目与投递。
//
// 先终止事件再删条目：只删条目的话，仍在扇出的事件随后会把死链写回来。
func (s *Service) DeleteAnnouncements(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	var affected int64
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.terminateObject(ctx, tx, notify.ObjectAnnouncement, ids); err != nil {
			return err
		}
		rows, err := s.announcements.WithTx(tx).DeleteBatch(ctx, "id IN ?", ids)
		if err != nil {
			return fmt.Errorf("notification: delete announcements: %w", err)
		}
		affected = rows
		return nil
	})
	if err != nil {
		return 0, err
	}
	return affected, nil
}

// PublishAnnouncement 把草稿置为已发布，并在同一事务里发出通知事件。
//
// 返回 false 表示当前状态不是草稿（已被并发发布或已撤回），调用方据此提示「已处理」。
// 提交后再广播公告状态，供管理端多实例刷新列表；广播失败只影响实时性。
func (s *Service) PublishAnnouncement(ctx context.Context, id int64, publishedBy string) (bool, error) {
	// 先取回实体：受众与展示字段要按管理员看到的那一份来
	row, err := s.announcements.FindById(ctx, id)
	if err != nil {
		return false, err
	}
	if row.Status != AnnouncementDraft {
		return false, nil
	}

	at := time.Now()
	var published bool

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// WHERE 带上原状态：先读再按 id 更新不是原子的
		rows, err := s.announcements.WithTx(tx).UpdateFields(ctx, map[string]interface{}{
			"status":       AnnouncementPublished,
			"published_at": at,
			"published_by": publishedBy,
		}, "id = ? AND status = ?", id, AnnouncementDraft)
		if err != nil {
			return fmt.Errorf("notification: publish announcement: %w", err)
		}
		if rows == 0 {
			return nil
		}
		published = true

		audienceIDs, err := decodeAudience(row.TargetIds)
		if err != nil {
			return err
		}
		if row.TargetType == AudienceUserIDs && len(audienceIDs) == 0 {
			return ErrAudienceEmpty
		}
		return s.Emit(ctx, tx, Request{
			Scene:       notify.SceneAnnouncement,
			ActorID:     publishedBy,
			ObjectType:  notify.ObjectAnnouncement,
			ObjectID:    id,
			Audience:    row.TargetType,
			AudienceIDs: audienceIDs,
			Payload:     notify.Payload{AnnouncementID: id},
		})
	})
	if err != nil {
		return false, err
	}

	if published && s.announcementNotifier != nil {
		s.announcementNotifier.NotifyAnnouncement(ctx, true, row.Id, row.Title, row.Category, row.Level, at)
	}
	return published, nil
}

// RevokeAnnouncement 把已发布的公告置为已撤回，并回收派生条目与投递。
func (s *Service) RevokeAnnouncement(ctx context.Context, id int64) (bool, error) {
	row, err := s.announcements.FindById(ctx, id)
	if err != nil {
		return false, err
	}
	if row.Status != AnnouncementPublished {
		return false, nil
	}

	var revoked bool
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rows, err := s.announcements.WithTx(tx).UpdateFields(ctx,
			map[string]interface{}{"status": AnnouncementRevoked},
			"id = ? AND status = ?", id, AnnouncementPublished)
		if err != nil {
			return fmt.Errorf("notification: revoke announcement: %w", err)
		}
		if rows == 0 {
			return nil
		}
		revoked = true
		return s.terminateObject(ctx, tx, notify.ObjectAnnouncement, []int64{id})
	})
	if err != nil {
		return false, err
	}

	if revoked && s.announcementNotifier != nil {
		s.announcementNotifier.NotifyAnnouncement(ctx, false, row.Id, row.Title, row.Category, row.Level, time.Now())
	}
	return revoked, nil
}

// RedeliverEvent 把失败或停滞的事件复位后重新扇出。
//
// 已完成的事件拒绝重放：重放会从头扫描受众，扇出期间新注册的用户会收到一条旧通知。
// 已终止的事件同样拒绝 —— 它的内容已经不可见。
func (s *Service) RedeliverEvent(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ev, err := lockEvent(ctx, tx, s.events.TableName(), id)
		if err != nil {
			return err
		}
		if ev == nil {
			return gorm.ErrRecordNotFound
		}
		if ev.Status == StatusDone || ev.Status == StatusTerminated {
			return ErrEventNotRedeliverable
		}

		// 游标与计数一起归零：受众可能已经变化，重放要重新算分母
		if _, err := s.events.WithTx(tx).UpdateFields(ctx, map[string]interface{}{
			"status":         StatusPending,
			"cursor":         0,
			"processed":      0,
			"audience_total": 0,
		}, "id = ?", id); err != nil {
			return fmt.Errorf("notification: reset event: %w", err)
		}

		msg, err := mq.FanoutMessage.Message(&mq.NotificationFanoutEvent{EventID: id})
		if err != nil {
			return fmt.Errorf("notification: build fanout wakeup: %w", err)
		}
		if err := s.outbox.Enqueue(ctx, tx, msg); err != nil {
			return fmt.Errorf("notification: enqueue fanout wakeup: %w", err)
		}
		return nil
	})
}

// RedeliverDelivery 把失败或待发送的投递记录复位为待发送，由投递消费者重新领取。
func (s *Service) RedeliverDelivery(ctx context.Context, id int64) error {
	row, err := s.deliveries.FindById(ctx, id)
	if err != nil {
		return err
	}
	if row.Status == DeliverySent || row.Status == DeliveryCancelled {
		return ErrDeliveryNotRedeliverable
	}

	// 重投递是一次新的人工尝试，重试次数归零才不会被旧计数卡在重试上限外
	_, err = s.deliveries.UpdateFields(ctx, map[string]interface{}{
		"status":        DeliveryPending,
		"attempts":      0,
		"error_msg":     "",
		"next_retry_at": gorm.Expr("NULL"),
	}, "id = ?", id)
	if err != nil {
		return fmt.Errorf("notification: reset delivery: %w", err)
	}
	return nil
}
