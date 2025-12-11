package main

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/mq/redismqx"
)

func main() {
	ctx := context.Background()

	// 创建 RedisMessenger 实例
	messenger := redismqx.NewRedisMqConn("localhost:6379", "events")

	// 在一个 Goroutine 中启动订阅
	go messenger.SubscribeMessage(func(ctx context.Context, msg []byte) error {
		fmt.Printf("[收到消息] %s\n", msg)
		return nil
	})

	// 模拟发布消息
	for i := 1; i <= 3; i++ {
		messenger.PublishMessage(ctx, []byte(fmt.Sprintf("事件 %d", i)))
	}

	// 防止主程序退出
	select {}
}
