package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// Request 是一次通知意图。业务方只描述「发生了什么、给谁看」，不关心扇出与投递。
type Request struct {
	Scene       string
	ActorID     string
	ObjectType  string
	ObjectID    int64
	Audience    string
	AudienceIDs []string
	Payload     notify.Payload
}

// Emit 在调用方的事务里写入一条通知事件，并写入唤醒消息。
//
// 事件行是通知的权威状态：载荷、受众、游标都在它上面；发件箱只负责把扇出消费者叫醒。
// 因此这条调用必须与业务写入同事务 —— 事件写失败，业务写入一并回滚。
func (s *Service) Emit(ctx context.Context, tx *gorm.DB, req Request) error {
	if tx == nil {
		return errors.New("notification: transaction is required")
	}
	if _, ok := notify.Lookup(req.Scene); !ok {
		return fmt.Errorf("notification: unknown scene %q", req.Scene)
	}
	if !ValidAudience(req.Audience) {
		return fmt.Errorf("notification: invalid audience %q", req.Audience)
	}
	if req.Audience == AudienceUserIDs && len(req.AudienceIDs) == 0 {
		return errors.New("notification: user_ids audience is empty")
	}

	payload, err := json.Marshal(req.Payload)
	if err != nil {
		return fmt.Errorf("notification: encode payload: %w", err)
	}
	payloadStr := string(payload)

	var audienceIDs *string
	if req.Audience == AudienceUserIDs {
		ids, err := json.Marshal(req.AudienceIDs)
		if err != nil {
			return fmt.Errorf("notification: encode audience: %w", err)
		}
		sv := string(ids)
		audienceIDs = &sv
	}

	rec := &model.TNotificationEvent{
		Scene:        req.Scene,
		ActorId:      req.ActorID,
		ObjectType:   req.ObjectType,
		ObjectId:     req.ObjectID,
		AudienceType: req.Audience,
		AudienceIds:  audienceIDs,
		Payload:      &payloadStr,
		Status:       StatusPending,
	}
	if _, err := s.events.WithTx(tx).Insert(ctx, rec); err != nil {
		return fmt.Errorf("notification: insert event: %w", err)
	}

	msg, err := mq.FanoutMessage.Message(&mq.NotificationFanoutEvent{EventID: rec.Id})
	if err != nil {
		return fmt.Errorf("notification: build wakeup message: %w", err)
	}
	if err := s.outbox.Enqueue(ctx, tx, msg); err != nil {
		return fmt.Errorf("notification: enqueue wakeup: %w", err)
	}
	return nil
}
