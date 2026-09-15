package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchPatchMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchPatchMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchPatchMessagesLogic {
	return &BatchPatchMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量更新留言状态
func (l *BatchPatchMessagesLogic) BatchPatchMessages(in *discussionrpc.BatchPatchMessagesRequest) (*discussionrpc.BatchPatchMessagesResponse, error) {
	rows, err := l.svcCtx.TMessageModel.UpdateFields(l.ctx, map[string]interface{}{
		"status": in.Status,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchPatchMessagesResponse{SuccessCount: rows}, nil
}
