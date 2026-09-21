package chatservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteChatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteChatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteChatsLogic {
	return &BatchDeleteChatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除聊天记录
func (l *BatchDeleteChatsLogic) BatchDeleteChats(in *chatrpc.BatchDeleteChatsRequest) (*chatrpc.BatchDeleteChatsResponse, error) {
	rows, err := l.svcCtx.TChatModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &chatrpc.BatchDeleteChatsResponse{SuccessCount: rows}, nil
}
