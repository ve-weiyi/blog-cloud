// Package likes 维护点赞状态：Redis 承载实时状态，数据库由消费端异步收敛。
package likes

import (
	"context"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/outbox"
	"github.com/ve-weiyi/blog-cloud/model"
)

// Deps 是发布端依赖：Redis 承载实时状态，outbox 保证落库事件不丢。
type Deps struct {
	Redis  *redis.Client
	DB     *gorm.DB
	Outbox outbox.Store
}

// Toggle 切换 Redis 中的点赞状态，并把落库动作投递到 outbox。
// 动作在发布端按实时状态判定，事件携带确定动作，重投不会翻转结果。
func Toggle(ctx context.Context, d Deps, targetType, targetID int64, userID, likeKey, countKey string) error {
	member := strconv.FormatInt(targetID, 10)

	liked, err := d.Redis.SIsMember(ctx, likeKey, member).Result()
	if err != nil {
		return err
	}

	action, delta := mq.LikeActionLike, int64(1)
	if liked {
		action, delta = mq.LikeActionUnlike, -1
	}

	if liked {
		err = d.Redis.SRem(ctx, likeKey, member).Err()
	} else {
		err = d.Redis.SAdd(ctx, likeKey, member).Err()
	}
	if err != nil {
		return err
	}

	// 计数排行缓存是派生数据，写失败不影响点赞结果
	if err := d.Redis.ZIncrBy(ctx, countKey, float64(delta), member).Err(); err != nil {
		logx.Errorf("更新点赞计数缓存失败: %v", err)
	}

	// 事件未落库时回滚实时状态，避免用户看到的结果与数据库分叉
	rollback := func() {
		if liked {
			_ = d.Redis.SAdd(ctx, likeKey, member).Err()
		} else {
			_ = d.Redis.SRem(ctx, likeKey, member).Err()
		}
		_ = d.Redis.ZIncrBy(ctx, countKey, float64(-delta), member).Err()
	}

	msg, err := mq.LikeMessage.Message(&mq.LikeEvent{
		UserId:     userID,
		TargetType: targetType,
		TargetId:   targetID,
		Action:     action,
	})
	if err != nil {
		rollback()
		return err
	}

	if err := d.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return d.Outbox.Enqueue(ctx, tx, msg)
	}); err != nil {
		rollback()
		return err
	}

	return nil
}

// Apply 把点赞状态收敛到 wantLike，返回计数应有的变化（0 表示无需变更）。
// 消费端必须幂等：事件重投时状态已一致，不再改动计数。
func Apply(ctx context.Context, m model.TLikeModel, targetType, targetID int64, userID string, wantLike bool) (int64, error) {
	_, err := m.FindOneByTargetTypeTargetIdUserId(ctx, targetType, targetID, userID)
	exists := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	switch {
	case wantLike && !exists:
		if _, err := m.Insert(ctx, &model.TLike{
			TargetType: targetType,
			TargetId:   targetID,
			UserId:     userID,
		}); err != nil {
			return 0, err
		}
		return 1, nil

	case !wantLike && exists:
		if _, err := m.DeleteBatch(ctx, "target_type = ? AND target_id = ? AND user_id = ?", targetType, targetID, userID); err != nil {
			return 0, err
		}
		return -1, nil
	}

	return 0, nil
}
