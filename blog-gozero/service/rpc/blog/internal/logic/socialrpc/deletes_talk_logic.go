package socialrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesTalkLogic {
	return &DeletesTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除说说
func (l *DeletesTalkLogic) DeletesTalk(in *socialrpc.IdsReq) (*socialrpc.BatchResp, error) {
	rows, err := l.svcCtx.TTalkModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &socialrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
