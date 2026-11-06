package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

type BatchDeleteTalksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteTalksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteTalksLogic {
	return &BatchDeleteTalksLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// 批量删除说说
func (l *BatchDeleteTalksLogic) BatchDeleteTalks(in *discussionrpc.BatchDeleteTalksRequest) (*discussionrpc.BatchDeleteTalksResponse, error) {
	var rows int64
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		// 明细与目标同生共死：目标没了，点赞记录不再有归属
		if _, err := l.svcCtx.TLikeModel.WithTx(tx).DeleteBatch(l.ctx,
			"target_type = ? AND target_id in (?)", enums.LikeTargetTalk, in.Ids); err != nil {
			return err
		}

		var err error
		rows, err = l.svcCtx.TTalkModel.WithTx(tx).DeleteBatch(l.ctx, "id in (?)", in.Ids)
		if err != nil {
			return err
		}

		// 对象没了，指向它的收件箱条目就是死链：先终止事件，再删条目与投递。
		// 说说下的评论仍在库里，挂在它们上面的通知一并回收
		return l.svcCtx.Notify.TerminateTopic(l.ctx, tx, notify.ObjectTalk, in.Ids)
	})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchDeleteTalksResponse{SuccessCount: rows}, nil
}
