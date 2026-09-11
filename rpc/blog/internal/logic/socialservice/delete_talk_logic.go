package socialservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/socialrpc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTalkLogic {
	return &DeleteTalkLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteTalkLogic) DeleteTalk(in *socialrpc.DeleteTalkRequest) (*socialrpc.DeleteTalkResponse, error) {
	rows, err := l.svcCtx.TTalkModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	return &socialrpc.DeleteTalkResponse{SuccessCount: rows}, nil
}
