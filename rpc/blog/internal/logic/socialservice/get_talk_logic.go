package socialservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/socialrpc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTalkLogic {
	return &GetTalkLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetTalkLogic) GetTalk(in *socialrpc.GetTalkRequest) (*socialrpc.GetTalkResponse, error) {
	entity, err := l.svcCtx.TTalkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &socialrpc.GetTalkResponse{Talk: convertTalkOut(entity)}, nil
}
