package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageLogic {
	return &GetMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询留言
func (l *GetMessageLogic) GetMessage(in *discussionrpc.GetMessageRequest) (*discussionrpc.GetMessageResponse, error) {
	entity, err := l.svcCtx.TMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.GetMessageResponse{Message: convertMessageOut(entity)}, nil
}
