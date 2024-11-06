package ws

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Name string
	// 用于存储连接
	conn *websocket.Conn
}

func NewWebSocketClient(name string, conn *websocket.Conn) *WebSocketClient {
	return &WebSocketClient{
		Name: name,
		conn: conn,
	}
}

// 发送消息
func (c *WebSocketClient) PublishMessage(ctx context.Context, msg []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

// 接收消息
func (c *WebSocketClient) SubscribeMessage(handler Receive) error {
	for {
		//读取客户端发送的消息
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		err = handler(msg)
		if err != nil {
			response := fmt.Sprintf("Failed to handle message: %v", err)
			_ = c.conn.WriteMessage(websocket.TextMessage, []byte(response))
			continue
		}
	}
}

func (c *WebSocketClient) Close() error {
	return c.conn.Close()
}
