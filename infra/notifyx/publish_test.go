package notifyx

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
)

func newTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	return mr, rdb
}

// 订阅侧应当收到事件（锁定 channel 名与载荷编码）
func TestPublishReachesSubscriber(t *testing.T) {
	_, rdb := newTestRedis(t)

	sub := rdb.Subscribe(context.Background(), cachekey.NotifyEventChannel)
	defer func() { _ = sub.Close() }()

	received := make(chan string, 1)
	go func() {
		for msg := range sub.Channel() {
			received <- msg.Payload
			return
		}
	}()

	want := Event{Type: EventNotice, MessageId: 42, Title: "维护通知"}
	if err := Publish(context.Background(), rdb, want); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case payload := <-received:
		got, err := Decode(payload)
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
		if got != want {
			t.Fatalf("收到的事件不符:\n got=%+v\nwant=%+v", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("未收到事件")
	}
}

// Redis 短暂不可用时应当重试：推送是"尽力而为"的提醒，
// 但一次抖动就丢掉信号会让客户端长时间停在旧未读数。
//
// 这里只断言"恢复后 Publish 成功"——它就是重试存在的决定性证据
// （没有重试时首次失败即返回错误）。订阅可达性由上一个用例覆盖。
func TestPublishRetriesUntilRedisRecovers(t *testing.T) {
	mr, rdb := newTestRedis(t)

	mr.SetError("redis is down")
	go func() {
		time.Sleep(80 * time.Millisecond)
		mr.SetError("")
	}()

	if err := Publish(context.Background(), rdb, Event{Type: EventNotice, MessageId: 7}); err != nil {
		t.Fatalf("Redis 恢复后应当重试成功, got err=%v", err)
	}
}

// 彻底不可用时返回错误（由调用方决定是否记日志），且不应无限阻塞
func TestPublishGivesUpAndReturnsError(t *testing.T) {
	mr, rdb := newTestRedis(t)
	mr.SetError("redis is down")

	done := make(chan error, 1)
	go func() {
		done <- Publish(context.Background(), rdb, Event{Type: EventNotice})
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("应当返回错误")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Publish 重试未收敛，阻塞超过 3 秒")
	}
}
