package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetNotifyMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotifyMessageLogic {
	return &GetNotifyMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据ID查询通知消息
func (l *GetNotifyMessageLogic) GetNotifyMessage(in *notificationrpc.GetNotifyMessageRequest) (*notificationrpc.GetNotifyMessageResponse, error) {
	record, err := l.svcCtx.TNotifyMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.GetNotifyMessageResponse{
		Message: convertTNotifyMessageToProto(record),
	}, nil
}
