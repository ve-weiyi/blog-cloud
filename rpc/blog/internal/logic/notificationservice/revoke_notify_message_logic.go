package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"

	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
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
func (l *RevokeNotifyMessageLogic) RevokeNotifyMessage(in *notificationrpc.RevokeNotifyMessageRequest) (*notificationrpc.RevokeNotifyMessageResponse, error) {
	msg, err := l.svcCtx.TNotifyMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if msg.Status != "published" {
		return &notificationrpc.RevokeNotifyMessageResponse{Success: false}, nil
	}

	// 消息状态与投递记录必须一起改：只改其一会让已撤回的通知仍留在读者的可见列表里
	// （列表按投递记录状态过滤），而且外面看不出失败。
	// WHERE 带上原状态，并发撤回由影响行数裁决——先 FindById 判状态再按 id 更新不是原子的。
	var affected int64
	err = l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		rows, err := l.svcCtx.TNotifyMessageModel.WithTx(tx).UpdateFields(l.ctx,
			map[string]interface{}{"status": "revoked"},
			"id = ? AND status = ?", in.Id, "published")
		if err != nil {
			return err
		}
		affected = rows

		// 将未读的投递记录标记为已撤回，已读的保留原样
		_, err = l.svcCtx.TNotifyRecordModel.WithTx(tx).UpdateFields(l.ctx,
			map[string]interface{}{"status": "revoked"},
			"message_id = ? AND channel = ? AND status = ?",
			in.Id, "inbox", "unread")
		return err
	})
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		// 已被并发撤回（或状态已不是 published），本次不产生副作用
		return &notificationrpc.RevokeNotifyMessageResponse{Success: false}, nil
	}

	// 与发布对称：撤回同样广播一次事件，客户端据此刷新未读。
	// 撤回只改状态、不删记录，所以信号发出后未读列表里那条会变成已撤回
	if err := notifyx.Publish(l.ctx, l.svcCtx.Redis, notifyx.Event{
		Type:      notifyx.EventNoticeRevoke,
		MessageId: msg.Id,
	}); err != nil {
		l.Errorf("广播通知撤回事件失败: %v", err)
	}

	return &notificationrpc.RevokeNotifyMessageResponse{
		Success: true,
	}, nil
}
