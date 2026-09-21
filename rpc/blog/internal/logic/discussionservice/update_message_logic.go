package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMessageLogic {
	return &UpdateMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新留言
func (l *UpdateMessageLogic) UpdateMessage(in *discussionrpc.UpdateMessageRequest) (*discussionrpc.UpdateMessageResponse, error) {
	entity, err := l.svcCtx.TMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.MessageContent = in.MessageContent
	entity.Status = in.Status

	_, err = l.svcCtx.TMessageModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.UpdateMessageResponse{Success: true}, nil
}
