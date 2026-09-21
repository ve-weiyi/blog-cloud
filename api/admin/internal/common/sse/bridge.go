package sse

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
)

// EventSource 事件来源。只要求"能订阅一串事件"，由 Broker 实现。
type EventSource interface {
	Subscribe() (<-chan notifyx.Event, func())
}

// Bridge 把事件源桥接到推送通道，直到 ctx 结束或事件源关闭。
// 返回 nil 表示正常结束（客户端断开或服务关停）。
//
// 心跳不在这里：它以 SSE 注释行的形式由 handler 层写入（见 handler 的 beat），
// 注释行不进入这条通道，也就不参与帧的编解码。
func Bridge(ctx context.Context, src EventSource, client chan<- *types.NotifyStreamEvent) error {
	events, unsubscribe := src.Subscribe()
	defer unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-events:
			if !ok {
				// 事件源关停
				return nil
			}
			if !send(ctx, client, ToStreamEvent(event)) {
				return nil
			}
		}
	}
}

// ToStreamEvent 把跨进程事件映射成推送帧
func ToStreamEvent(event notifyx.Event) *types.NotifyStreamEvent {
	return &types.NotifyStreamEvent{
		Event:       event.Type,
		MessageId:   event.MessageId,
		Title:       event.Title,
		Category:    event.Category,
		Level:       event.Level,
		PublishedAt: event.PublishedAt,
	}
}

// send 写一帧。发送**必须**纳入 ctx 选择：客户端断开后 handler 不再排空该通道，
// 直接写会永久阻塞在这里，把桥接 goroutine 泄漏掉。
func send(ctx context.Context, client chan<- *types.NotifyStreamEvent, frame *types.NotifyStreamEvent) bool {
	select {
	case client <- frame:
		return true
	case <-ctx.Done():
		return false
	}
}
