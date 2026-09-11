package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type RevokeNotifyMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeNotifyMessageLogic {
	return &RevokeNotifyMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤回通知消息
func (l *RevokeNotifyMessageLogic) RevokeNotifyMessage(in *notificationrpc2.RevokeNotifyMessageRequest) (*notificationrpc2.RevokeNotifyMessageResponse, error) {
	msg, err := l.svcCtx.TNotifyMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if msg.Status != "published" {
		return &notificationrpc.RevokeNotifyMessageResponse{Success: false}, nil
	}

	fields := map[string]interface{}{
		"status": "revoked",
	}

	_, err = l.svcCtx.TNotifyMessageModel.UpdateFields(l.ctx, fields, "id = ?", in.Id)
	if err != nil {
		return nil, err
	}

	// 将未读的投递记录标记为已撤回，已读的保留原样
	_, err = l.svcCtx.TNotifyRecordModel.UpdateFields(l.ctx,
		map[string]interface{}{"status": "revoked"},
		"message_id = ? AND channel = ? AND status = ?",
		in.Id, "inbox", "unread")
	if err != nil {
		l.Errorf("RevokeNotifyMessage UpdateFields error: %v", err)
	}

	return &notificationrpc.RevokeNotifyMessageResponse{
		Success: true,
	}, nil
}
