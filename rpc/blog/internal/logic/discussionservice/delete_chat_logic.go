package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteChatLogic {
	return &DeleteChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除聊天记录
func (l *DeleteChatLogic) DeleteChat(in *discussionrpc.DeleteChatRequest) (*discussionrpc.DeleteChatResponse, error) {
	rows, err := l.svcCtx.TChatModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.DeleteChatResponse{SuccessCount: rows}, nil
}
