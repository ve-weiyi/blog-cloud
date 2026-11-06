package ssex

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EventSource 事件来源。只要求"能订阅一串事件"，由 Broker 实现。
type EventSource[T any] interface {
	Subscribe() (<-chan T, func())
}

// Bridge 把事件源桥接到推送通道，mapEvent 把单条事件转成推送帧，直到 ctx 结束或事件源关闭。
// 返回 nil 表示正常结束（客户端断开或服务关停）。
//
// 心跳不在这里：它以 SSE 注释行的形式由 handler 层写入，注释行不进入这条通道，
// 也就不参与帧的编解码。
func Bridge[S, T any](ctx context.Context, src EventSource[S], client chan<- T, mapEvent func(S) T) error {
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
			if !send(ctx, client, mapEvent(event)) {
				return nil
			}
		}
	}
}

// BridgeSignals 订阅指定 Redis 频道，每收到一条消息写一帧 frame()。
//
// 信号不带正文与收件人：连接本身已经绑定到登录用户，客户端收到后回查拉取接口。
// 返回 nil 表示正常结束（客户端断开或服务关停）。
func BridgeSignals[T any](ctx context.Context, rdb *redis.Client, channel string, client chan<- T, frame func() T) error {
	if rdb == nil || channel == "" {
		<-ctx.Done()
		return nil
	}

	sub := rdb.Subscribe(ctx, channel)
	defer func() { _ = sub.Close() }()

	return forward(ctx, sub.Channel(), client, frame)
}

// forward 把消息流转发成推送帧。
//
// ctx 必须与消息通道放在同一个 select 里：只写 `for range ch` 的话，客户端断开且
// 随后没有新信号时，这个 goroutine 会永远阻塞在接收上，订阅也一起泄漏。
func forward[T any](ctx context.Context, events <-chan *redis.Message, client chan<- T, frame func() T) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case _, ok := <-events:
			if !ok {
				// 订阅被关闭
				return nil
			}
			if !send(ctx, client, frame()) {
				return nil
			}
		}
	}
}

// send 写一帧。发送**必须**纳入 ctx 选择：客户端断开后 handler 不再排空该通道，
// 直接写会永久阻塞在这里，把桥接 goroutine 泄漏掉。
func send[T any](ctx context.Context, client chan<- T, frame T) bool {
	select {
	case client <- frame:
		return true
	case <-ctx.Done():
		return false
	}
}
