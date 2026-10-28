package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
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
		return err
	})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchDeleteTalksResponse{SuccessCount: rows}, nil
}
