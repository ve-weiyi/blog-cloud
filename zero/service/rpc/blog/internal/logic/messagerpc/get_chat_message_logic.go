package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatMessageLogic {
	return &GetChatMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询聊天记录
func (l *GetChatMessageLogic) GetChatMessage(in *messagerpc.IdReq) (*messagerpc.ChatMessageDetails, error) {
	entity, err := l.svcCtx.TChatMessageModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return convertChatMessageOut(entity), nil
}
