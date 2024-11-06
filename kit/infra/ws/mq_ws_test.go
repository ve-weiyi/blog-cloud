package ws

import (
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var mu sync.Mutex // 保护clients

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

// WebSocket 处理函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// 将新连接加入clients
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	// 持续读取来自客户端的消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			conn.Close()
			break
		}

		// 将消息发送到RabbitMQ
		publishToRabbitMQ(msg)
	}
}

// 将消息发布到RabbitMQ
func publishToRabbitMQ(message []byte) {
	err := rabbitChannel.Publish(
		"",              // exchange
		"chat_messages", // routing key (队列名称)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)

	if err != nil {
		log.Println("Failed to publish message:", err)
	}
}

// 从RabbitMQ消费消息并广播给所有WebSocket客户端
func consumeFromRabbitMQ() {
	msgs, err := rabbitChannel.Consume(
		"chat_messages", // 队列名称
		"",              // 消费者名称
		true,            // autoAck
		false,           // exclusive
		false,           // noLocal
		false,           // noWait
		nil,             // 参数
	)

	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	for msg := range msgs {
		broadcast <- msg.Body
	}
}

// 广播消息给所有连接的WebSocket客户端
func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

// 初始化RabbitMQ连接和队列
func initRabbitMQ() {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}

	_, err = rabbitChannel.QueueDeclare(
		"chat_messages", // 队列名称
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}
}

func main() {
	// 初始化RabbitMQ
	initRabbitMQ()
	go consumeFromRabbitMQ()

	// 处理消息广播
	go handleMessages()

	// 启动WebSocket服务器
	http.HandleFunc("/ws", wsHandler)
	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Test_MQWS(t *testing.T) {

}
