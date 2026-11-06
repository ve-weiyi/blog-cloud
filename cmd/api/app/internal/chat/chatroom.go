package chat

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/go-stomp/stomp/v3/frame"
	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/chatservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/userservice"
	chatdomain "github.com/ve-weiyi/blog-cloud/infra/domain/chat"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/stompwsx"
	"github.com/ve-weiyi/stompws/server"
	"github.com/ve-weiyi/stompws/server/client"
	"github.com/ve-weiyi/vkit/adapter/iplocx"
	"github.com/ve-weiyi/vkit/x/jsonv"
)

// rpcTimeout 限定回调内远程调用的最长等待时间，避免连接协程被悬挂的 RPC 占住。
const rpcTimeout = 3 * time.Second

type ChatRoomEventHook struct {
	server.DefaultEventHook
	connectTime sync.Map // key: clientId, value: connectTime
	onlineUser  sync.Map // key: clientId, value: *userservice.User
	// subscriptions key: clientId, value: subscriptionId。私发的 MESSAGE 必须带
	// subscription 头，否则客户端按订阅分发时会当作「未订阅的消息」丢掉。
	subscriptions sync.Map

	sendLimiter *sendLimiter

	UserService userservice.UserService
	ChatService chatservice.ChatService
}

func NewChatRoomEventHook(userSvc userservice.UserService, chatSvc chatservice.ChatService) *ChatRoomEventHook {
	return &ChatRoomEventHook{
		UserService: userSvc,
		ChatService: chatSvc,
		sendLimiter: newSendLimiter(),
	}
}

func (h *ChatRoomEventHook) OnConnect(server *server.StompHubServer, c *client.Conn) {
	id, login, _ := c.GetSessionInfo()
	ip := c.GetIpAddress()

	h.connectTime.Store(id, time.Now())

	if login != "" {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		out, err := h.UserService.GetUser(ctx, &userservice.GetUserRequest{
			UserId: login,
		})
		if err == nil && out.User != nil {
			h.onlineUser.Store(id, out.User)
		}
	}

	log.Printf("✅ User connected: (client: %s, user: %s, ip: %s) ", id, login, ip)
}

func (h *ChatRoomEventHook) OnDisconnect(server *server.StompHubServer, c *client.Conn) {
	id, login, _ := c.GetSessionInfo()
	ip := c.GetIpAddress()

	var duration time.Duration
	if value, exists := h.connectTime.LoadAndDelete(id); exists {
		connectTime := value.(time.Time)
		duration = time.Since(connectTime)
	}

	h.onlineUser.Delete(id)
	h.subscriptions.Delete(id)
	h.sendLimiter.remove(id)

	log.Printf("❌ User disconnected: (client: %s, user: %s, ip: %s), online: %v", id, login, ip, duration.Round(time.Second))
}

func (h *ChatRoomEventHook) OnSubscribe(server *server.StompHubServer, c *client.Conn, destination string, subscriptionId string) {
	clientId, login, _ := c.GetSessionInfo()
	ip := c.GetIpAddress()
	h.subscriptions.Store(clientId, subscriptionId)
	log.Printf("📢 User %s subscribed to %s", login, destination)
	count := stompwsx.OnlineUserCount(server)

	// 1. 私发问候消息
	greeting := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	greeting.Body = []byte(jsonv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeGreeting,
		Data: jsonv.AnyToJsonNE(GreetingMessageEvent{
			Content:  fmt.Sprintf("👋 welcome %s to the chat channel", login),
			IpSource: iplocx.GetIpSourceByBaidu(ip),
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	c.SendFrame(greeting)

	// 2. 广播在线人数
	online := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	online.Body = []byte(jsonv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeOnline,
		Data: jsonv.AnyToJsonNE(OnlineMessageEvent{
			Online: true,
			Count:  count,
			Tips:   fmt.Sprintf("👋 %s joined the chat channel", login),
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	server.RouteMessage(nil, online)

	// 3. 私发历史消息
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	out, err := h.ChatService.ListChats(ctx, &chatservice.ListChatsRequest{
		After:  time.Now().Add(-365 * 24 * time.Hour).UnixMilli(),
		Before: time.Now().UnixMilli(),
		Limit:  0,
	})
	if err != nil {
		return
	}
	history := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	history.Body = []byte(jsonv.AnyToJsonNE(MessageEvent{
		Type:      MessageTypeHistory,
		Data:      jsonv.AnyToJsonNE(newHistoryEvent(out)),
		TimeStamp: time.Now().UnixMilli(),
	}))
	c.SendFrame(history)
}

func (h *ChatRoomEventHook) OnUnsubscribe(server *server.StompHubServer, c *client.Conn, destination string, subscriptionId string) {
	clientId, login, _ := c.GetSessionInfo()
	h.subscriptions.Delete(clientId)
	log.Printf("📤 User %s unsubscribed from %s", login, destination)
	count := stompwsx.OnlineUserCount(server)

	online := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0", frame.Subscription, subscriptionId)
	online.Body = []byte(jsonv.AnyToJsonNE(MessageEvent{
		Type: MessageTypeOnline,
		Data: jsonv.AnyToJsonNE(OnlineMessageEvent{
			Online: true,
			Count:  count,
			Tips:   fmt.Sprintf("👋 %s left the chat channel", login),
		}),
		TimeStamp: time.Now().UnixMilli(),
	}))
	server.RouteMessage(nil, online)
}

func (h *ChatRoomEventHook) OnSend(server *server.StompHubServer, c *client.Conn, message *frame.Frame) bool {
	clientId, login, _ := c.GetSessionInfo()
	ip := c.GetIpAddress()
	destination := message.Header.Get(frame.Destination)

	var event MessageEvent
	if err := jsonv.JsonToAny(string(message.Body), &event); err != nil {
		log.Printf("invalid chat event from %s: %v", login, err)
		return h.reject(c, "消息格式错误", "")
	}

	switch event.Type {
	case MessageTypeSend:
		var send SendMessageEvent
		if err := jsonv.JsonToAny(event.Data, &send); err != nil {
			return h.reject(c, "消息格式错误", "")
		}
		if !chatdomain.ValidType(send.Type) {
			return h.reject(c, "不支持的消息类型", send.ClientMsgId)
		}
		content, err := chatdomain.NormalizeContent(send.Content)
		if err != nil {
			return h.reject(c, fmt.Sprintf("消息长度必须在 1-%d 个字符之间", chatdomain.MaxContentRunes), send.ClientMsgId)
		}
		if !h.sendLimiter.allow(clientId, time.Now()) {
			return h.reject(c, "发送过于频繁，请稍后再试", send.ClientMsgId)
		}

		ipSource := iplocx.GetIpSourceByBaidu(ip)
		var userId, nickname, avatar string
		if userInfo, ok := h.onlineUser.Load(clientId); ok {
			if user, ok := userInfo.(*userservice.User); ok {
				userId = user.UserId
				nickname = user.Nickname
				avatar = user.Avatar
			}
		}

		now := time.Now().UnixMilli()
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		out, err := h.ChatService.CreateChat(ctx, &chatservice.CreateChatRequest{
			UserId:    userId,
			DeviceId:  clientId,
			IpAddress: ip,
			IpSource:  ipSource,
			Nickname:  nickname,
			Avatar:    avatar,
			Type:      send.Type,
			Content:   content,
			Status:    enums.ChatStatusNormal,
		})
		if err != nil {
			// 带上 client-msg-id，客户端才能立刻把本地气泡标记为失败而不是等超时
			return h.reject(c, "消息发送失败，请稍后重试", send.ClientMsgId)
		}

		msgFrame := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msgFrame.Body = []byte(jsonv.AnyToJsonNE(MessageEvent{
			Type: MessageTypeMessage,
			Data: jsonv.AnyToJsonNE(ChatMessageEvent{
				Id:          out.Id,
				UserId:      userId,
				DeviceId:    clientId,
				Nickname:    nickname,
				Avatar:      avatar,
				IpSource:    ipSource,
				Type:        send.Type,
				Content:     content,
				Status:      enums.ChatStatusNormal,
				CreatedAt:   now,
				UpdatedAt:   now,
				ClientMsgId: send.ClientMsgId,
			}),
			TimeStamp: now,
		}))
		server.RouteMessage(c, msgFrame)

	case MessageTypeEdit:
		var edit EditMessageEvent
		if err := jsonv.JsonToAny(event.Data, &edit); err != nil {
			return h.reject(c, "消息格式错误", "")
		}

		content := edit.Content
		messageType := edit.Type
		switch edit.Status {
		case enums.ChatStatusRecall:
			// 撤回的展示文案由服务端固定，避免客户端把任意内容塞进「已撤回」消息
			content = "[消息已撤回]"
			messageType = chatdomain.TypeText
		case enums.ChatStatusNormal, enums.ChatStatusEdited:
			if !chatdomain.ValidType(messageType) {
				return h.reject(c, "不支持的消息类型", "")
			}
			normalized, err := chatdomain.NormalizeContent(content)
			if err != nil {
				return h.reject(c, fmt.Sprintf("消息长度必须在 1-%d 个字符之间", chatdomain.MaxContentRunes), "")
			}
			content = normalized
		default:
			return h.reject(c, "消息状态不合法", "")
		}

		var userId string
		if userInfo, ok := h.onlineUser.Load(clientId); ok {
			if user, ok := userInfo.(*userservice.User); ok {
				userId = user.UserId
			}
		}

		now := time.Now().UnixMilli()
		rpcCtx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		ctx := metadata.AppendToOutgoingContext(rpcCtx,
			headerx.HeaderUid, userId,
			headerx.HeaderXDeviceId, clientId,
		)
		_, err := h.ChatService.UpdateChat(ctx, &chatservice.UpdateChatRequest{
			Id:      edit.Id,
			Type:    messageType,
			Content: content,
			Status:  edit.Status,
		})
		if err != nil {
			return h.reject(c, "操作失败，请稍后重试", "")
		}

		msgFrame := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msgFrame.Body = []byte(jsonv.AnyToJsonNE(MessageEvent{
			Type: MessageTypeEdit,
			Data: jsonv.AnyToJsonNE(EditMessageEvent{
				Id:        edit.Id,
				Type:      messageType,
				Content:   content,
				Status:    edit.Status,
				UpdatedAt: now,
			}),
			TimeStamp: now,
		}))
		server.RouteMessage(c, msgFrame)

	case MessageTypeHistoryMore:
		var req HistoryMoreMessageEvent
		if err := jsonv.JsonToAny(event.Data, &req); err != nil {
			return h.reject(c, "消息格式错误", "")
		}
		if req.Before <= 0 {
			req.Before = time.Now().UnixMilli()
		}
		if req.Limit <= 0 || req.Limit > 50 {
			req.Limit = 20
		}

		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		out, err := h.ChatService.ListChats(ctx, &chatservice.ListChatsRequest{
			Before: req.Before,
			Limit:  int64(req.Limit),
		})
		if err != nil {
			return h.reject(c, "加载历史消息失败", "")
		}

		c.SendFrame(h.replyFrame(clientId, destination, []byte(jsonv.AnyToJsonNE(MessageEvent{
			Type:      MessageTypeHistoryMore,
			Data:      jsonv.AnyToJsonNE(newHistoryEvent(out)),
			TimeStamp: time.Now().UnixMilli(),
		}))))

	default:
		return false
	}

	return true
}

// replyFrame 构造私发 MESSAGE。
//
// 私发帧必须带 subscription 头：客户端按订阅分发 MESSAGE，缺这个头会落到
// 「未订阅消息」分支被丢弃（广播帧由 hub 投递时自动补上，这里绕过了那一步）。
func (h *ChatRoomEventHook) replyFrame(clientId, destination string, body []byte) *frame.Frame {
	f := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "0")
	if subId, ok := h.subscriptions.Load(clientId); ok {
		if id, ok := subId.(string); ok && id != "" {
			f.Header.Set(frame.Subscription, id)
		}
	}
	f.Body = body
	return f
}

// newHistoryEvent 把一次查询结果转成广播事件，过滤已删除记录。
func newHistoryEvent(out *chatservice.ListChatsResponse) HistoryMessageEvent {
	list := make([]*ChatMessageEvent, 0, len(out.List))
	for _, msg := range out.List {
		if msg.Status == enums.ChatStatusDeleted {
			continue
		}
		list = append(list, &ChatMessageEvent{
			Id:        msg.Id,
			UserId:    msg.UserId,
			DeviceId:  msg.DeviceId,
			Nickname:  msg.Nickname,
			Avatar:    msg.Avatar,
			IpSource:  msg.IpSource,
			Type:      msg.Type,
			Content:   msg.Content,
			Status:    msg.Status,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		})
	}

	var page, size, total int64
	if out.ListResult != nil {
		page, size, total = out.ListResult.Page, out.ListResult.PageSize, out.ListResult.Total
	}

	event := HistoryMessageEvent{List: list, Page: page, Size: size, Total: total}
	// 是否有更早的页按**原始行数**判断：已删除记录也占一页名额，
	// 用过滤后的条数推断会让分页提前终止。
	event.HasMore = size > 0 && int64(len(out.List)) >= size
	// 游标取本页最早的原始行（含已删除），下一页不会重复拉取这一页的行；
	// 结果按 created_at 倒序，最后一条即最早。
	if n := len(out.List); n > 0 {
		event.NextBefore = out.List[n-1].CreatedAt
	}
	return event
}

// reject 回一帧 ERROR 并拒绝该 SEND。client-msg-id 供客户端定位失败的本地消息。
func (h *ChatRoomEventHook) reject(c *client.Conn, message, clientMsgID string) bool {
	errFrame := frame.New(frame.ERROR, frame.Message, message)
	if clientMsgID != "" {
		errFrame.Header.Set("client-msg-id", clientMsgID)
	}
	c.SendFrame(errFrame)
	return false
}

const (
	chatRatePerSecond = 1
	chatRateBurst     = 3
)

type sendBucket struct {
	tokens float64
	last   time.Time
}

// sendLimiter 按连接做令牌桶：突发 3 条，之后每秒 1 条。
type sendLimiter struct {
	mu      sync.Mutex
	buckets map[string]*sendBucket
}

func newSendLimiter() *sendLimiter {
	return &sendLimiter{buckets: make(map[string]*sendBucket)}
}

func (l *sendLimiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, ok := l.buckets[key]
	if !ok {
		bucket = &sendBucket{tokens: chatRateBurst, last: now}
		l.buckets[key] = bucket
	}
	if elapsed := now.Sub(bucket.last).Seconds(); elapsed > 0 {
		bucket.tokens = math.Min(chatRateBurst, bucket.tokens+elapsed*chatRatePerSecond)
	}
	bucket.last = now
	if bucket.tokens < 1 {
		return false
	}
	bucket.tokens--
	return true
}

func (l *sendLimiter) remove(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

func (h *ChatRoomEventHook) OnAck(c *client.Conn, messageId string) {
	_, login, _ := c.GetSessionInfo()
	log.Printf("✅ Message %s acknowledged by %s", messageId, login)
}

func (h *ChatRoomEventHook) OnNack(c *client.Conn, messageId string) {
	_, login, _ := c.GetSessionInfo()
	log.Printf("❌ Message %s rejected by %s", messageId, login)
}

func (h *ChatRoomEventHook) OnError(c *client.Conn, err error) {
	_, login, _ := c.GetSessionInfo()
	log.Printf("🚨 Error for user %s: %v", login, err)
}
