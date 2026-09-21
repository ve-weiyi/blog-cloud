package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateNotifyMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotifyMessageLogic {
	return &CreateNotifyMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建通知消息
func (l *CreateNotifyMessageLogic) CreateNotifyMessage(in *notificationrpc.CreateNotifyMessageRequest) (*notificationrpc.CreateNotifyMessageResponse, error) {
	msg := convertProtoToTNotifyMessage(in)

	_, err := l.svcCtx.TNotifyMessageModel.Insert(l.ctx, msg)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.CreateNotifyMessageResponse{
		Id: msg.Id,
	}, nil
}
