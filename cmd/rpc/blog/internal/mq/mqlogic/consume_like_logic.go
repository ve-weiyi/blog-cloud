package mqlogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
)

// ConsumeLikeLogic 点赞落库逻辑
type ConsumeLikeLogic struct {
	svcCtx *svc.ServiceContext
}

func NewConsumeLikeLogic(svcCtx *svc.ServiceContext) *ConsumeLikeLogic {
	return &ConsumeLikeLogic{svcCtx: svcCtx}
}

// Consume 把点赞动作幂等落库：明细与计数在同一事务内收敛。
func (l *ConsumeLikeLogic) Consume(ctx context.Context, event *mq.LikeEvent) error {
	if event.Action != mq.LikeActionLike && event.Action != mq.LikeActionUnlike {
		// 未知动作重试也不会成功，记录后丢弃
		logx.WithContext(ctx).Errorf("未知点赞动作: %s", event.Action)
		return nil
	}

	wantLike := event.Action == mq.LikeActionLike
	return l.svcCtx.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		delta, err := likes.Apply(ctx, l.svcCtx.TLikeModel.WithTx(tx),
			event.TargetType, event.TargetId, event.UserId, wantLike)
		if err != nil || delta == 0 {
			return err
		}
		affected, err := l.applyCount(ctx, tx, event.TargetType, event.TargetId, delta)
		if err != nil {
			return err
		}
		// 目标已被删除：撤销刚写入的明细，避免排队事件把它写回
		if affected == 0 {
			_, err = l.svcCtx.TLikeModel.WithTx(tx).DeleteBatch(ctx,
				"target_type = ? AND target_id = ?", event.TargetType, event.TargetId)
		}
		return err
	})
}

// applyCount 按目标类型增减计数列，返回受影响行数。
func (l *ConsumeLikeLogic) applyCount(ctx context.Context, tx *gorm.DB, targetType, targetID, delta int64) (int64, error) {
	fields := map[string]interface{}{"like_count": gorm.Expr("like_count + ?", delta)}

	switch targetType {
	case enums.LikeTargetArticle:
		return l.svcCtx.TArticleModel.WithTx(tx).UpdateFields(ctx, fields, "id = ?", targetID)
	case enums.LikeTargetTalk:
		return l.svcCtx.TTalkModel.WithTx(tx).UpdateFields(ctx, fields, "id = ?", targetID)
	case enums.LikeTargetComment:
		return l.svcCtx.TCommentModel.WithTx(tx).UpdateFields(ctx, fields, "id = ?", targetID)
	}

	return 0, fmt.Errorf("未知点赞目标类型: %d", targetType)
}
