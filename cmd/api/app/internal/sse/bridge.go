// Package sse 把 Redis 上的按收件人信号桥接到单个连接的推送流。
package sse

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
)

// Bridge 订阅该收件人的频道，把「收件箱变了」写成推送帧。
//
// 信号不带正文与收件人：连接本身已经绑定到登录用户，客户端收到后回查拉取接口。
// 返回 nil 表示正常结束（客户端断开或服务关停）。
func Bridge(ctx context.Context, rdb *redis.Client, userID string, client chan<- *types.NotificationStreamEvent) error {
	if rdb == nil || userID == "" {
		<-ctx.Done()
		return nil
	}

	sub := rdb.Subscribe(ctx, cachekey.NotifyUserChannel(userID))
	defer func() { _ = sub.Close() }()

	return forward(ctx, sub.Channel(), client)
}

// forward 把消息流转发成推送帧。
//
// ctx 必须与消息通道放在同一个 select 里：只写 `for range ch` 的话，客户端断开且
// 随后没有新信号时，这个 goroutine 会永远阻塞在接收上，订阅也一起泄漏。
func forward(ctx context.Context, events <-chan *redis.Message, client chan<- *types.NotificationStreamEvent) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case _, ok := <-events:
			if !ok {
				// 订阅被关闭
				return nil
			}
			select {
			case client <- &types.NotificationStreamEvent{Event: notifyx.EventInboxChanged}:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
