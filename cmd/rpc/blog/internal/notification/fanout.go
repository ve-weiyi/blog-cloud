package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// fanoutBatchSize 单次扇出处理的收件人上限。
// 每批独立事务并持久化游标，进程中途退出也能从游标继续。
const fanoutBatchSize = 200

// lockEvent 读取并锁住事件行；不存在时返回 nil，表示无需处理。
//
// 行锁覆盖整批写入，让「终止事件」与「扇出」串行：撤回先到就停止，扇出先到就写完这批。
func lockEvent(ctx context.Context, tx *gorm.DB, table string, eventID int64) (*model.TNotificationEvent, error) {
	var ev model.TNotificationEvent
	err := tx.WithContext(ctx).Table(table).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", eventID).
		First(&ev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

// Fanout 处理一批收件人。
//
// 一次调用只推进一批：事件行上持有行锁，写完条目、投递与游标后提交；
// 未穷尽时写一条唤醒消息继续下一批，丢了由兜底扫描补上。
func (s *Service) Fanout(ctx context.Context, eventID int64) error {
	var inserted []string

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ev, err := lockEvent(ctx, tx, s.events.TableName(), eventID)
		if err != nil || ev == nil {
			return err
		}
		// 已终止或已完成：撤回、内容删除跑在扇出前面时走这里
		if IsTerminal(ev.Status) {
			return nil
		}

		scene, ok := notify.Lookup(ev.Scene)
		if !ok {
			// 未知场景是编程错误，重试没有意义
			_, err := s.events.WithTx(tx).UpdateFields(ctx,
				map[string]interface{}{"status": StatusFailed}, "id = ?", ev.Id)
			return err
		}

		if ev.AudienceTotal == 0 {
			total, err := countAudience(ctx, tx, s.users.TableName(), ev)
			if err != nil {
				return err
			}
			if _, err := s.events.WithTx(tx).UpdateFields(ctx,
				map[string]interface{}{"audience_total": total}, "id = ?", ev.Id); err != nil {
				return err
			}
			ev.AudienceTotal = total
		}

		batch, nextCursor, exhausted, err := fetchRecipients(ctx, tx, s.users.TableName(), ev)
		if err != nil {
			return err
		}

		if len(batch) > 0 {
			inserted, err = s.writeBatch(ctx, tx, scene, ev, batch)
			if err != nil {
				return err
			}
		}

		fields := map[string]interface{}{"cursor": nextCursor, "status": StatusFanning}
		if exhausted {
			// 完成即对齐分母：all 受众的分母含触发者，期间也可能有用户被删
			fields["status"] = StatusDone
			fields["processed"] = ev.AudienceTotal
		} else {
			fields["processed"] = gorm.Expr("processed + ?", len(batch))
		}
		if _, err := s.events.WithTx(tx).UpdateFields(ctx, fields, "id = ?", ev.Id); err != nil {
			return err
		}

		// 续跑批次的唤醒与游标同事务：提交后再发消息，一旦发送失败事件就停在「扇出中」，
		// 只能等兜底扫描。本批可能一条新条目都没产生（全被去重键吞掉），但那不代表后面没有。
		if exhausted {
			return nil
		}
		msg, err := mq.FanoutMessage.Message(&mq.NotificationFanoutEvent{EventID: eventID})
		if err != nil {
			return fmt.Errorf("notification: build fanout wakeup: %w", err)
		}
		if err := s.outbox.Enqueue(ctx, tx, msg); err != nil {
			return fmt.Errorf("notification: enqueue fanout wakeup: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 信号只在提交之后发：客户端收到就会回查，早于写入的回查什么都拉不到
	if len(inserted) > 0 && s.notifier != nil {
		s.notifier.NotifyChanged(ctx, inserted)
	}
	return nil
}

// writeBatch 写入条目与投递，返回本次真正新增的收件人。
func (s *Service) writeBatch(ctx context.Context, tx *gorm.DB, scene notify.Scene, ev *model.TNotificationEvent, batch []recipient) ([]string, error) {
	recipients := make([]string, 0, len(batch))
	address := make(map[string]recipient, len(batch))
	for _, r := range batch {
		recipients = append(recipients, r.UserID)
		address[r.UserID] = r
	}

	prefs, err := s.loadPrefs(ctx, tx, recipients, scene.Group)
	if err != nil {
		return nil, err
	}

	// 定向受众的收件人还没查过用户记录，这里补上投递地址
	if ev.AudienceType == AudienceUserIDs {
		if err := fillAddress(ctx, tx, s.users.TableName(), address); err != nil {
			return nil, err
		}
	}

	plans := PlanEntries(scene, ev.ObjectType, ev.ObjectId, ev.Id, ev.ActorId, recipients, prefs)
	if len(plans) == 0 {
		return nil, nil
	}

	existing, err := s.entryIDs(ctx, tx, plans)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	recs := make([]*model.TNotification, 0, len(plans))
	for _, p := range plans {
		if existing[key(p.RecipientID, p.DedupeKey)] != 0 {
			continue
		}
		recs = append(recs, &model.TNotification{
			EventId:     ev.Id,
			RecipientId: p.RecipientID,
			Scene:       ev.Scene,
			ObjectType:  ev.ObjectType,
			ObjectId:    ev.ObjectId,
			GroupKey:    p.GroupKey,
			DedupeKey:   p.DedupeKey,
			CreatedAt:   now,
		})
	}
	if len(recs) > 0 {
		if _, err := s.notifs.WithTx(tx).InsertIgnoreBatch(ctx, recs...); err != nil {
			return nil, fmt.Errorf("notification: insert entries: %w", err)
		}
	}

	// 投递的依据是「本批收件人对应的条目」，不是「本次新插入的条目」：
	// 上批可能插了条目却没写投递，重放时插入会被去重键吞掉，按新行算就永远补不上。
	ids, err := s.entryIDs(ctx, tx, plans)
	if err != nil {
		return nil, err
	}
	if err := s.writeDeliveries(ctx, tx, plans, address, ids, now); err != nil {
		return nil, err
	}

	inserted := make([]string, 0, len(recs))
	for _, rec := range recs {
		inserted = append(inserted, rec.RecipientId)
	}
	return inserted, nil
}

// loadPrefs 读取这批收件人在本场景所属分组下的显式开关。
func (s *Service) loadPrefs(ctx context.Context, tx *gorm.DB, recipients []string, group notify.Group) (map[string]map[notify.Channel]bool, error) {
	out := make(map[string]map[notify.Channel]bool, len(recipients))

	rows, err := s.prefs.WithTx(tx).FindAll(ctx,
		"user_id IN ? AND pref_group = ?", recipients, string(group))
	if err != nil {
		return nil, fmt.Errorf("notification: load preferences: %w", err)
	}
	for _, row := range rows {
		if out[row.UserId] == nil {
			out[row.UserId] = make(map[notify.Channel]bool, 2)
		}
		out[row.UserId][notify.Channel(row.Channel)] = row.Enabled == 1
	}
	return out, nil
}

// entryIDs 回查本批收件人对应的条目 id；投递按它挂载。
func (s *Service) entryIDs(ctx context.Context, tx *gorm.DB, plans []Entry) (map[string]int64, error) {
	recipients, dedupeKeys := collect(plans)

	// 两个 IN 是笛卡尔条件，会多取到「非本组合」的行；按 (recipient, dedupe) 取值即可，
	// 去重键对不可分组场景含事件、对可分组场景含触发者，跨事件不会误配。
	rows, err := s.notifs.WithTx(tx).FindAll(ctx,
		"recipient_id IN ? AND dedupe_key IN ?", recipients, dedupeKeys)
	if err != nil {
		return nil, fmt.Errorf("notification: query entries: %w", err)
	}

	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[key(row.RecipientId, row.DedupeKey)] = row.Id
	}
	return out, nil
}

func collect(plans []Entry) (recipients []string, dedupeKeys []string) {
	recipients = make([]string, 0, len(plans))
	dedupeKeys = make([]string, 0, len(plans))
	for _, p := range plans {
		recipients = append(recipients, p.RecipientID)
		dedupeKeys = append(dedupeKeys, p.DedupeKey)
	}
	return recipients, dedupeKeys
}

func key(recipientID, dedupeKey string) string {
	return recipientID + "\x00" + dedupeKey
}
