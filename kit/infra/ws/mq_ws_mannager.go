package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/mq"
)

type MqWebsocketManager struct {
	// 用于升级 HTTP 连接
	websocket.Upgrader
	// 用于存储连接
	connections map[string]*websocket.Conn
	mu          sync.Mutex // 保护clients

	// 消息发布者和订阅者
	Publisher  mq.MessagePublisher
	Subscriber mq.MessageSubscriber
}

// WebSocket 处理函数
func (m *MqWebsocketManager) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// 将新连接加入clients
	m.mu.Lock()
	clients[conn] = true
	m.mu.Unlock()

	// 持续读取来自客户端的消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			m.mu.Lock()
			delete(clients, conn)
			m.mu.Unlock()
			conn.Close()
			break
		}

		// 将消息发送到RabbitMQ
		m.publishToRabbitMQ(msg)
	}
}

// 将消息发布到RabbitMQ
func (m *MqWebsocketManager) publishToRabbitMQ(message []byte) {
	err := m.Publisher.PublishMessage(nil, message)

	if err != nil {
		log.Println("Failed to publish message:", err)
	}
}
