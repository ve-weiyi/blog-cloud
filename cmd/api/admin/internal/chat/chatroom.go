package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/stompws/server"
	"github.com/ve-weiyi/stompws/server/client"
)

type ChatRoomEventHook struct {
	server.DefaultEventHook
	connectTime sync.Map
}

func NewChatRoomEventHook() *ChatRoomEventHook {
	return &ChatRoomEventHook{}
}

// marshalChatMessage 用 json.Marshal 生成消息体，避免用户内容里的引号破坏 JSON。
func marshalChatMessage(sender, content string) []byte {
	body, err := json.Marshal(struct {
		Sender  string `json:"sender"`
		Content string `json:"content"`
	}{Sender: sender, Content: content})
	if err != nil {
		return []byte(`{}`)
	}
	return body
}

func (h *ChatRoomEventHook) OnConnect(server *server.StompHubServer, c *client.Conn) {
	id, login, _ := c.GetSessionInfo()
	h.connectTime.Store(id, time.Now())

	log.Printf("✅ User connected: %s (id: %s)", login, id)
	msg := frame.New(frame.MESSAGE, frame.Destination, "/topic/system/broadcast", frame.MessageId, "0")
	msg.Body = []byte(fmt.Sprintf(`{"sender":"System","content":"🟢 %s is online"}`, login))
	server.RouteMessage(nil, msg)
}

func (h *ChatRoomEventHook) OnDisconnect(server *server.StompHubServer, c *client.Conn) {
	id, login, _ := c.GetSessionInfo()

	var duration time.Duration
	if value, exists := h.connectTime.LoadAndDelete(id); exists {
		connectTime := value.(time.Time)
		duration = time.Since(connectTime)
	}

	log.Printf("❌ User disconnected: %s (id: %s), online: %v", login, id, duration.Round(time.Second))
	msg := frame.New(frame.MESSAGE, frame.Destination, "/topic/system/broadcast", frame.MessageId, "0")
	msg.Body = []byte(fmt.Sprintf(`{"sender":"System","content":"🔴 %s went offline"}`, login))
	server.RouteMessage(nil, msg)
}

func (h *ChatRoomEventHook) OnSubscribe(server *server.StompHubServer, c *client.Conn, destination string, subscriptionId string) {
	_, login, _ := c.GetSessionInfo()
	log.Printf("📢 User %s subscribed to %s", login, destination)

	// 只在聊天频道发送加入消息
	if strings.HasPrefix(destination, "/topic/public/chat") {
		msg := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msg.Body = marshalChatMessage(login, fmt.Sprintf("👋 %s joined the chat channel", login))
		server.RouteMessage(c, msg)
	}
}

func (h *ChatRoomEventHook) OnUnsubscribe(server *server.StompHubServer, c *client.Conn, destination string, subscriptionId string) {
	_, login, _ := c.GetSessionInfo()
	log.Printf("📤 User %s unsubscribed from %s", login, destination)

	// 只在聊天频道发送离开消息
	if strings.HasPrefix(destination, "/topic/public/chat") {
		msg := frame.New(frame.MESSAGE, frame.Destination, destination, frame.MessageId, "1")
		msg.Body = marshalChatMessage(login, fmt.Sprintf("👋 %s left the chat channel", login))
		server.RouteMessage(c, msg)
	}
}

// /app/system/broadcast → /topic/system/broadcast
// /app/chat/{roomId} → /topic/chat/{roomId}
// /app/private/{userId} → /queue/user/{userId}
func (h *ChatRoomEventHook) OnSend(server *server.StompHubServer, c *client.Conn, message *frame.Frame) bool {
	_, login, _ := c.GetSessionInfo()
	destination := message.Header.Get(frame.Destination)
	msgContent := string(message.Body)

	log.Printf("💬 Message from %s to %s: %s", login, destination, msgContent)
	if len(msgContent) > 1000 {
		log.Printf("⚠️ Long message from %s (%d chars)", login, len(msgContent))
	}

	switch {
	// 系统广播消息
	case strings.HasPrefix(destination, "/app/system/broadcast"):
		dest := fmt.Sprintf("/topic/system/broadcast")
		msg := frame.New(frame.MESSAGE, frame.Destination, dest, frame.MessageId, "1")
		msg.Body = marshalChatMessage(login, msgContent)
		server.RouteMessage(c, msg)

	// 聊天室消息
	case strings.HasPrefix(destination, "/app/chat/"):
		roomId := strings.TrimPrefix(destination, "/app/chat/")
		dest := fmt.Sprintf("/topic/chat/%s", roomId)
		msg := frame.New(frame.MESSAGE, frame.Destination, dest, frame.MessageId, "1")
		msg.Body = marshalChatMessage(login, msgContent)
		server.RouteMessage(c, msg)

	// 私聊消息
	case strings.HasPrefix(destination, "/app/user/"):
		userId := strings.TrimPrefix(destination, "/app/user/")
		dest := fmt.Sprintf("/queue/user/%s", userId)
		msg := frame.New(frame.MESSAGE, frame.Destination, dest, frame.MessageId, "1")
		msg.Body = marshalChatMessage(login, msgContent)
		server.RouteMessage(c, msg)
	default:
		return false
	}

	return true
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
