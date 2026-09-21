package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type CreateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMessageLogic {
	return &CreateMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建留言
func (l *CreateMessageLogic) CreateMessage(in *discussionrpc.CreateMessageRequest) (*discussionrpc.CreateMessageResponse, error) {
	entity := &model.TMessage{
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		MessageContent: in.MessageContent,
		Status:         in.Status,
	}
	_, err := l.svcCtx.TMessageModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.CreateMessageResponse{Id: entity.Id}, nil
}
