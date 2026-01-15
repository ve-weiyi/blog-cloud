package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
func (l *GetChatLogic) GetChat(in *messagerpc.GetChatReq) (*messagerpc.GetChatResp, error) {
	entity, err := l.svcCtx.TChatModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &messagerpc.GetChatResp{
		Chat: convertChatOut(entity),
	}, nil
}
