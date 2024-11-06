package talkrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/talkrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTalkLogic {
	return &GetTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询说说
func (l *GetTalkLogic) GetTalk(in *talkrpc.IdReq) (*talkrpc.TalkDetails, error) {
	entity, err := l.svcCtx.TTalkModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return convertTalkOut(entity), nil
}
