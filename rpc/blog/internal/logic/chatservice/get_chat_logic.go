package chatservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogic {
	return &GetChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询聊天记录
func (l *GetChatLogic) GetChat(in *chatrpc.GetChatRequest) (*chatrpc.GetChatResponse, error) {
	entity, err := l.svcCtx.TChatModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &chatrpc.GetChatResponse{Chat: convertChatOut(entity)}, nil
}
