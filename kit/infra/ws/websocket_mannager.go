package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	// 用于升级 HTTP 连接
	websocket.Upgrader
	// 用于存储连接
	connections map[string]*WebSocketClient
}

// NewWebSocketManager 创建一个 WebSocket 管理器
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// 允许所有请求来源
				return true
			},
		},
		connections: make(map[string]*WebSocketClient),
	}
}

// 开始处理 WebSocket 连接
func (m *WebSocketManager) RegisterWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocketClient, error) {
	conn, err := m.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	name := conn.RemoteAddr().String()
	// 将新的连接加入客户端列表
	client := NewWebSocketClient(name, conn)
	m.connections[name] = client

	return client, nil
}

// 关闭指定的客户端连接
func (m *WebSocketManager) CloseClient(name string) error {
	conn, ok := m.connections[name]
	if !ok {
		return nil
	}

	delete(m.connections, name)
	return conn.Close()
}

// 广播消息
func (m *WebSocketManager) BroadcastMsg(msg []byte) error {
	for _, conn := range m.connections {
		err := conn.PublishMessage(nil, msg)
		if err != nil {
			m.CloseClient(conn.Name)
		}
	}
	return nil
}
