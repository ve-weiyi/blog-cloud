package talkrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/talkrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTalkLogic {
	return &UpdateTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新说说
func (l *UpdateTalkLogic) UpdateTalk(in *talkrpc.TalkNewReq) (*talkrpc.TalkDetails, error) {
	entity := convertTalkIn(in)

	_, err := l.svcCtx.TTalkModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertTalkOut(entity), nil
}
