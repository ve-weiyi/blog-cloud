package stomphook

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/ve-blog-golang/stompws/server/client"
)

type OnlineEventHook struct {
	client.DefaultEventHook
	onlineCount int64
}

func NewOnlineEventHook() *OnlineEventHook {
	return &OnlineEventHook{}
}

func (h *OnlineEventHook) OnConnect(server *client.StompHubServer, c *client.Client) {
	count := atomic.AddInt64(&h.onlineCount, 1)
	h.broadcastOnlineCount(server, count)
}

func (h *OnlineEventHook) OnDisconnect(server *client.StompHubServer, c *client.Client) {
	count := atomic.AddInt64(&h.onlineCount, -1)
	h.broadcastOnlineCount(server, count)
}

func (h *OnlineEventHook) OnSubscribe(server *client.StompHubServer, c *client.Client, destination string, subscriptionId string) {
	_, login, _ := c.GetClientInfo()
	log.Printf("📢 User %s subscribed to %s", login, destination)

	// 只在聊天频道发送加入消息
	if strings.HasPrefix(destination, "/topic/system/online") {
		msg := frame.New(frame.MESSAGE, frame.Destination, "/topic/system/online", frame.MessageId, "0", frame.Subscription, subscriptionId)
		msg.Body = []byte(fmt.Sprintf(`{"type":"online_count","count":%d}`, h.onlineCount))
		c.SendFrame(msg)
	}
}

func (h *OnlineEventHook) broadcastOnlineCount(server *client.StompHubServer, count int64) {
	msg := frame.New(frame.MESSAGE, frame.Destination, "/topic/system/online", frame.MessageId, "0")
	msg.Body = []byte(fmt.Sprintf(`{"type":"online_count","count":%d}`, count))
	server.RouteMessage(nil, msg)
}
