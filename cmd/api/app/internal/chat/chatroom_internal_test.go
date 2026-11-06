package chat

import (
	"testing"
	"time"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/chatservice"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
)

// 突发 3 条，之后按每秒 1 条补充。
func TestSendLimiter(t *testing.T) {
	limiter := newSendLimiter()
	now := time.Unix(100, 0)

	for i := 0; i < chatRateBurst; i++ {
		if !limiter.allow("client-1", now) {
			t.Fatalf("第 %d 条应放行", i+1)
		}
	}
	if limiter.allow("client-1", now) {
		t.Fatal("突发额度用完后应拒绝")
	}

	if !limiter.allow("client-1", now.Add(time.Second)) {
		t.Fatal("等待 1 秒后应补充 1 条额度")
	}
	if !limiter.allow("client-2", now) {
		t.Fatal("不同连接使用独立额度")
	}

	limiter.remove("client-1")
	if !limiter.allow("client-1", now) {
		t.Fatal("移除后应重新获得满额")
	}
}

// 私发帧必须带订阅 id，客户端才会把它交给订阅回调。
func TestReplyFrameCarriesSubscription(t *testing.T) {
	h := NewChatRoomEventHook(nil, nil)
	h.subscriptions.Store("c1", "sub-0")

	got := h.replyFrame("c1", "/topic/chatroom", nil).Header.Get(frame.Subscription)
	if got != "sub-0" {
		t.Fatalf("私发帧应带订阅 id, got %q", got)
	}

	if got := h.replyFrame("c2", "/topic/chatroom", nil).Header.Get(frame.Subscription); got != "" {
		t.Fatalf("没有订阅的连接不应被填上订阅 id, got %q", got)
	}
}

// 一页里混有已删除记录时，过滤不影响「还有更早」，游标取最早的原始行。
func TestNewHistoryEvent(t *testing.T) {
	out := &chatservice.ListChatsResponse{
		ListResult: &chatservice.ListResult{Page: 1, PageSize: 2, Total: 5},
		List: []*chatservice.Chat{
			{Id: 5, CreatedAt: 5000},
			{Id: 4, CreatedAt: 4000, Status: enums.ChatStatusDeleted},
		},
	}

	event := newHistoryEvent(out)

	if len(event.List) != 1 || event.List[0].Id != 5 {
		t.Fatalf("已删除记录应被过滤，got %+v", event.List)
	}
	if !event.HasMore {
		t.Fatal("原始行数达到一页时应还有更早的消息")
	}
	if event.NextBefore != 4000 {
		t.Fatalf("游标应取最早的原始行, got %d", event.NextBefore)
	}
}

func TestNewHistoryEventLastPage(t *testing.T) {
	out := &chatservice.ListChatsResponse{
		ListResult: &chatservice.ListResult{Page: 1, PageSize: 20, Total: 1},
		List:       []*chatservice.Chat{{Id: 1, CreatedAt: 1000}},
	}

	event := newHistoryEvent(out)

	if event.HasMore {
		t.Fatal("不足一页时应没有更早的消息")
	}
}
