package stomphook

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
	"github.com/ve-weiyi/stompws/server/client"
	"github.com/ve-weiyi/vkit/adapter/ipx"
	"github.com/ve-weiyi/vkit/x/jsonconv"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
)

type ChatRoomEventHook struct {
	client.DefaultEventHook
	connectTime sync.Map // key: clientId, value: connectTime
	onlineUser  sync.Map // key: clientId, value: *userservice.User

	UserService       userservice.UserService
	DiscussionService discussionservice.DiscussionService
}

func NewChatRoomEventHook(userSvc userservice.UserService, discussionSvc discussionservice.DiscussionService) *ChatRoomEventHook {
	return &ChatRoomEventHook{
		UserService:       userSvc,
		DiscussionService: discussionSvc,
	}
}

func (h *ChatRoomEventHook) OnConnect(server *client.StompHubServer, c *client.Client) {
	id, login, _ := c.GetClientInfo()
	ip := c.GetIpAddress()

	h.connectTime.Store(id, time.Now())

	if login != "" {
		out, err := h.UserService.GetUser(context.Background(), &userservice.GetUserRequest{
			UserId: login,
		})
		if err == nil && out.User != nil {
			h.onlineUser.Store(id, out.User)
		}
	}

	log.Printf("✅ User connected: (client: %s, user: %s, ip: %s) ", id, login, ip)
}

func (h *ChatRoomEventHook) OnDisconnect(server *client.StompHubServer, c *client.Client) {
	id, login, _ := c.GetClientInfo()
	ip := c.GetIpAddress()

	var duration time.Duration
	if value, exists := h.connectTime.LoadAndDelete(id); exists {
		connectTime := value.(time.Time)
		duration = time.Since(connectTime)
	}

	h.onlineUser.Delete(id)

	log.Printf("❌ User disconnected: (client: %s, user: %s, ip: %s), online: %v", id, login, ip, duration.Round(time.Second))
}

func (h *ChatRoomEventHook) OnSubscribe(server *client.StompHubServer, c *client.Client, destination string, subscriptionId string) {
	_, login, _ := c.GetClientInfo()
	ip := c.GetIpAddress()
	log.Printf("📢 User %s subscribed to %s", login, destination)
	count, _ := server.OnlineTracker().GetOnlineCount(context.Background())

	// 1. 私发问候消息
	greeting := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	greeting.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeGreeting,
		Data: jsonconv.AnyToJsonNE(GreetingMessageEvent{
			Content:   fmt.Sprintf("👋 welcome %s to the chat channel", login),
			IpAddress: ip,
			IpSource:  ipx.GetIpSourceByBaidu(ip),
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	c.SendFrame(greeting)

	// 2. 广播在线人数
	online := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	online.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeOnline,
		Data: jsonconv.AnyToJsonNE(OnlineMessageEvent{
			Online: true,
			Count:  count,
			Tips:   fmt.Sprintf("👋 %s joined the chat channel", login),
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	server.RouteMessage(nil, online)

	// 3. 私发历史消息
	out, err := h.DiscussionService.ListChats(context.Background(), &discussionservice.ListChatsRequest{
		After:  time.Now().Add(-365 * 24 * time.Hour).UnixMilli(),
		Before: time.Now().UnixMilli(),
		Limit:  0,
	})
	if err != nil {
		return
	}
	list := make([]*ChatMessageEvent, 0)
	for _, msg := range out.List {
		list = append(list, &ChatMessageEvent{
			Id:        msg.Id,
			UserId:    msg.UserId,
			DeviceId:  msg.DeviceId,
			Nickname:  msg.Nickname,
			Avatar:    msg.Avatar,
			IpAddress: msg.IpAddress,
			IpSource:  msg.IpSource,
			Type:      msg.Type,
			Content:   msg.Content,
			Status:    msg.Status,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		})
	}

	history := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	history.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeHistory,
		Data: jsonconv.AnyToJsonNE(HistoryMessageEvent{
			List:  list,
			Page:  out.PageResult.Page,
			Size:  out.PageResult.PageSize,
			Total: out.PageResult.Total,
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	c.SendFrame(history)
}

func (h *ChatRoomEventHook) OnUnsubscribe(server *client.StompHubServer, c *client.Client, destination string, subscriptionId string) {
	_, login, _ := c.GetClientInfo()
	log.Printf("📤 User %s unsubscribed from %s", login, destination)
	count, _ := server.OnlineTracker().GetOnlineCount(context.Background())

	online := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	online.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeOnline,
		Data: jsonconv.AnyToJsonNE(OnlineMessageEvent{
			Online: true,
			Count:  count,
			Tips:   fmt.Sprintf("👋 %s left the chat channel", login),
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	server.RouteMessage(nil, online)
}

func (h *ChatRoomEventHook) OnSend(server *client.StompHubServer, c *client.Client, message *frame.Frame) bool {
	clientId, login, _ := c.GetClientInfo()
	ip := c.GetIpAddress()

	destination := message.Header.Get(frame.Destination)
	body := string(message.Body)

	log.Printf("💬 Message from %s to %s: %s", login, destination, body)
	if len(body) > 1000 {
		log.Printf("⚠️ Long message from %s (%d chars)", login, len(body))
	}

	var event MessageEvent
	jsonconv.JsonToAny(body, &event)

	switch event.Type {
	case MessageTypeSend:
		var send SendMessageEvent
		jsonconv.JsonToAny(event.Data, &send)

		ipSource := ipx.GetIpSourceByBaidu(ip)
		var userId, nickname, avatar string
		if userInfo, ok := h.onlineUser.Load(clientId); ok {
			if user, ok := userInfo.(*userservice.User); ok {
				userId = user.UserId
				nickname = user.Nickname
				avatar = user.Avatar
			}
		}

		now := time.Now().UnixMilli()
		out, err := h.DiscussionService.CreateChat(context.Background(), &discussionservice.CreateChatRequest{
			UserId:    userId,
			DeviceId:  clientId,
			IpAddress: ip,
			IpSource:  ipSource,
			Nickname:  nickname,
			Avatar:    avatar,
			Type:      send.Type,
			Content:   send.Content,
			Status:    enums.ChatStatusNormal,
		})
		if err != nil {
			return false
		}

		msgFrame := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msgFrame.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
			Type: MessageTypeMessage,
			Data: jsonconv.AnyToJsonNE(ChatMessageEvent{
				Id:        out.Id,
				UserId:    userId,
				DeviceId:  clientId,
				Nickname:  nickname,
				Avatar:    avatar,
				IpAddress: ip,
				IpSource:  ipSource,
				Type:      send.Type,
				Content:   send.Content,
				Status:    enums.ChatStatusNormal,
				CreatedAt: now,
				UpdatedAt: now,
			}),
			TimeStamp: now,
		}))
		server.RouteMessage(c, msgFrame)

	case MessageTypeEdit:
		var edit EditMessageEvent
		jsonconv.JsonToAny(event.Data, &edit)

		now := time.Now().UnixMilli()
		_, err := h.DiscussionService.UpdateChat(context.Background(), &discussionservice.UpdateChatRequest{
			Id:      edit.Id,
			Type:    edit.Type,
			Content: edit.Content,
			Status:  edit.Status,
		})
		if err != nil {
			return false
		}

		msgFrame := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msgFrame.Body = []byte(jsonconv.AnyToJsonNE(MessageEvent{
			Type: MessageTypeEdit,
			Data: jsonconv.AnyToJsonNE(EditMessageEvent{
				Id:        edit.Id,
				Type:      edit.Type,
				Content:   edit.Content,
				Status:    edit.Status,
				UpdatedAt: now,
			}),
			TimeStamp: now,
		}))
		server.RouteMessage(c, msgFrame)

	default:
		return false
	}

	return true
}

func (h *ChatRoomEventHook) OnAck(c *client.Client, messageId string) {
	_, login, _ := c.GetClientInfo()
	log.Printf("✅ Message %s acknowledged by %s", messageId, login)
}

func (h *ChatRoomEventHook) OnNack(c *client.Client, messageId string) {
	_, login, _ := c.GetClientInfo()
	log.Printf("❌ Message %s rejected by %s", messageId, login)
}

func (h *ChatRoomEventHook) OnError(c *client.Client, err error) {
	_, login, _ := c.GetClientInfo()
	log.Printf("🚨 Error for user %s: %v", login, err)
}
