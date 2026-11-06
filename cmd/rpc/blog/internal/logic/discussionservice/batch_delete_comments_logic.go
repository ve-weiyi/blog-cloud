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

type BatchDeleteCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteCommentsLogic {
	return &BatchDeleteCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除评论
func (l *BatchDeleteCommentsLogic) BatchDeleteComments(in *discussionrpc.BatchDeleteCommentsRequest) (*discussionrpc.BatchDeleteCommentsResponse, error) {
	var rows int64
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		// 明细与目标同生共死：目标没了，点赞记录不再有归属
		if _, err := l.svcCtx.TLikeModel.WithTx(tx).DeleteBatch(l.ctx,
			"target_type = ? AND target_id in (?)", enums.LikeTargetComment, in.Ids); err != nil {
			return err
		}

		var err error
		rows, err = l.svcCtx.TCommentModel.WithTx(tx).DeleteBatch(l.ctx, "id in (?)", in.Ids)
		if err != nil {
			return err
		}

		// 对象没了，指向它的收件箱条目就是死链：先终止事件，再删条目与投递
		return l.svcCtx.Notify.TerminateObject(l.ctx, tx, notify.ObjectComment, in.Ids)
	})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchDeleteCommentsResponse{SuccessCount: rows}, nil
}
