package stomphook

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/stompws/server/client"
)

type ChatRoomEventHook struct {
	client.DefaultEventHook
	onlineCount int64
	connectTime sync.Map // key: clientId, value: connectTime
	onlineUser  sync.Map // key: clientId, value: userInfo

	AccountRpc accountrpc.AccountRpc
	MessageRpc messagerpc.MessageRpc
}

func NewChatRoomEventHook(accountRpc accountrpc.AccountRpc, messageRpc messagerpc.MessageRpc) *ChatRoomEventHook {
	return &ChatRoomEventHook{
		AccountRpc: accountRpc,
		MessageRpc: messageRpc,
	}
}

func (h *ChatRoomEventHook) OnConnect(server *client.StompHubServer, c *client.Client) {
	id, login, ip, _ := c.GetClientInfo()

	h.connectTime.Store(id, time.Now())

	// 用户登录
	if login != "" {
		// 加载用户信息
		userInfo, err := h.AccountRpc.GetUserInfo(context.Background(), &accountrpc.GetUserInfoReq{
			UserId: login,
		})
		if err == nil {
			h.onlineUser.Store(id, userInfo)
		}
	}

	log.Printf("✅ User connected: (client: %s, user: %s, ip: %s) ", id, login, ip)
}

func (h *ChatRoomEventHook) OnDisconnect(server *client.StompHubServer, c *client.Client) {
	id, login, ip, _ := c.GetClientInfo()

	var duration time.Duration
	if value, exists := h.connectTime.LoadAndDelete(id); exists {
		connectTime := value.(time.Time)
		duration = time.Since(connectTime)
	}

	// 用户退出
	h.onlineUser.Delete(id)

	log.Printf("❌ User disconnected: (client: %s, user: %s, ip: %s), online: %v", id, login, ip, duration.Round(time.Second))
}

func (h *ChatRoomEventHook) OnSubscribe(server *client.StompHubServer, c *client.Client, destination string, subscriptionId string) {
	_, login, ip, _ := c.GetClientInfo()
	log.Printf("📢 User %s subscribed to %s", login, destination)
	count := atomic.AddInt64(&h.onlineCount, 1)

	// 用户进入聊天室
	// 1. 私发问候消息
	greeting := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	greeting.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeGreeting,
		Data: jsonconv.AnyToJsonNE(
			GreetingMessageEvent{
				Content:   fmt.Sprintf("👋 welcome %s to the chat channel", login),
				IpAddress: ip,
				IpSource:  ipx.GetIpSourceByBaidu(ip),
			},
		),
		TimeStamp: time.Now().UnixMilli(),
	}))
	// 私发
	c.SendFrame(greeting)

	// 2. 发送在线人数
	online := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	online.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeOnline,
		Data: jsonconv.AnyToJsonNE(
			OnlineMessageEvent{
				Online: true,
				Count:  count,
				Tips:   fmt.Sprintf("👋 %s joined the chat channel", login),
			}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	// 广播
	server.RouteMessage(nil, online)

	// 3. 私发历史消息
	out, err := h.MessageRpc.FindChatList(context.Background(), &messagerpc.FindChatListReq{
		After:   time.Now().Add(-365 * 24 * time.Hour).UnixMilli(),
		Before:  time.Now().UnixMilli(),
		Limit:   0,
		UserId:  "",
		Type:    "",
		Content: "",
		Status:  0,
	})
	if err != nil {
		return
	}
	list := make([]*ChatMessageEvent, 0)
	for _, msg := range out.List {
		item := &ChatMessageEvent{
			Id:         msg.Id,
			UserId:     msg.UserId,
			TerminalId: msg.TerminalId,
			Nickname:   msg.Nickname,
			Avatar:     msg.Avatar,
			IpAddress:  msg.IpAddress,
			IpSource:   msg.IpSource,
			Type:       msg.Type,
			Content:    msg.Content,
			Status:     msg.Status,
			CreatedAt:  msg.CreatedAt,
			UpdatedAt:  msg.UpdatedAt,
		}

		list = append(list, item)
	}

	history := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	history.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeHistory,
		Data: jsonconv.AnyToJsonNE(HistoryMessageEvent{
			List:  list,
			Page:  out.Pagination.Page,
			Size:  out.Pagination.PageSize,
			Total: out.Pagination.Total,
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	// 私发
	c.SendFrame(history)
}

func (h *ChatRoomEventHook) OnUnsubscribe(server *client.StompHubServer, c *client.Client, destination string, subscriptionId string) {
	_, login, _, _ := c.GetClientInfo()
	log.Printf("📤 User %s unsubscribed from %s", login, destination)
	count := atomic.AddInt64(&h.onlineCount, -1)
	// 发送在线人数
	online := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	online.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeOnline,
		Data: jsonconv.AnyToJsonNE(
			OnlineMessageEvent{
				Online: true,
				Count:  count,
				Tips:   fmt.Sprintf("👋 %s left the chat channel", login),
			}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	// 广播
	server.RouteMessage(nil, online)
}

func (h *ChatRoomEventHook) OnSend(server *client.StompHubServer, c *client.Client, message *frame.Frame) bool {
	clientId, login, ip, _ := c.GetClientInfo()

	destination := message.Header.Get(frame.Destination)
	body := string(message.Body)

	log.Printf("💬 Message from %s to %s: %s", login, destination, body)
	if len(body) > 1000 {
		log.Printf("⚠️ Long message from %s (%d chars)", login, len(body))
	}

	// 解析发送消息
	var event MessageEvent
	jsonconv.JsonToAny(body, &event)

	switch event.Type {
	case MessageTypeSend:
		// 发送消息
		// 1. 解析内容
		var send SendMessageEvent
		jsonconv.JsonToAny(event.Data, &send)
		// 2. 存储到数据库
		ipSource := ipx.GetIpSourceByBaidu(ip)
		var userId, nickname, avatar string
		// 获取用户信息
		userInfo, ok := h.onlineUser.Load(clientId)
		if ok {
			user, okk := userInfo.(*accountrpc.UserInfo)
			if okk {
				userId = user.UserId
				nickname = user.Nickname
				avatar = user.Avatar
			}
		}

		out, err := h.MessageRpc.AddChat(context.Background(), &messagerpc.AddChatReq{
			UserId:     userId,
			TerminalId: clientId,
			IpAddress:  ip,
			IpSource:   ipSource,
			Nickname:   nickname,
			Avatar:     avatar,
			Type:       send.Type,
			Content:    send.Content,
			Status:     constant.ChatStatusNormal,
		})
		if err != nil {
			return false
		}

		// 3. 转发消息
		msgFrame := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msgFrame.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
			Type: MessageTypeMessage,
			Data: jsonconv.AnyToJsonNE(
				ChatMessageEvent{
					Id:         out.Chat.Id,
					UserId:     out.Chat.UserId,
					TerminalId: out.Chat.TerminalId,
					Nickname:   out.Chat.Nickname,
					Avatar:     out.Chat.Avatar,
					IpAddress:  out.Chat.IpAddress,
					IpSource:   out.Chat.IpSource,
					Type:       out.Chat.Type,
					Content:    out.Chat.Content,
					Status:     out.Chat.Status,
					CreatedAt:  out.Chat.CreatedAt,
					UpdatedAt:  out.Chat.UpdatedAt,
				}),
			TimeStamp: time.Now().UnixMilli(),
		}))
		server.RouteMessage(c, msgFrame)
	case MessageTypeEdit:
		// 编辑消息
		// 1. 解析内容
		var edit EditMessageEvent
		jsonconv.JsonToAny(event.Data, &edit)
		// 2. 更新数据库
		out, err := h.MessageRpc.UpdateChat(context.Background(), &messagerpc.UpdateChatReq{
			Id:      edit.Id,
			Type:    edit.Type,
			Content: edit.Content,
			Status:  edit.Status,
		})
		if err != nil {
			return false
		}
		// 3. 转发消息
		msgFrame := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msgFrame.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
			Type: MessageTypeEdit,
			Data: jsonconv.AnyToJsonNE(
				EditMessageEvent{
					Id:        out.Chat.Id,
					Type:      out.Chat.Type,
					Content:   out.Chat.Content,
					Status:    out.Chat.Status,
					UpdatedAt: out.Chat.UpdatedAt,
				}),
			TimeStamp: time.Now().UnixMilli(),
		}))
		server.RouteMessage(c, msgFrame)
	default:
		return false
	}

	return true
}

func (h *ChatRoomEventHook) OnAck(c *client.Client, messageId string) {
	_, login, _, _ := c.GetClientInfo()
	log.Printf("✅ Message %s acknowledged by %s", messageId, login)
}

func (h *ChatRoomEventHook) OnNack(c *client.Client, messageId string) {
	_, login, _, _ := c.GetClientInfo()
	log.Printf("❌ Message %s rejected by %s", messageId, login)
}

func (h *ChatRoomEventHook) OnError(c *client.Client, err error) {
	_, login, _, _ := c.GetClientInfo()
	log.Printf("🚨 Error for user %s: %v", login, err)
}
