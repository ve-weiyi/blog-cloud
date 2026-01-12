package _simple

import (
	"context"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func Test_Send(t *testing.T) {
	connect()
	createTopic()
	go producer()
	go consumer()
	select {}
}

func producer() {

	w := &kafka.Writer{
		Addr:  kafka.TCP("localhost:19094", "localhost:29094", "localhost:39094"),
		Topic: "topic-A",
	}

	for {
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("Key-A"),
				Value: []byte("Hello World! " + time.Now().String()),
			},
		)
		if err != nil {
			log.Println("failed to write messages:", err)
			continue
		}

		log.Println("send message success")
		time.Sleep(2 * time.Second)
	}
}

func consumer() {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:19094", "localhost:29094", "localhost:39094"},
		GroupID:  "consumer-group-id",
		Topic:    "topic-A",
		MaxWait:  10 * time.Second,
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("failed to read message:", err)
			continue
		}

		log.Printf("message at offset %d: %s = %s", m.Offset, string(m.Key), string(m.Value))
	}
}

// 列出主题
func connect() {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	conn, err := dialer.Dial("tcp", "localhost:19094")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		log.Println(k)
	}

	log.Println("connect success")
}

func createTopic() {
	// to create topics when auto.create.topics.enable='false'
	topic := "topic-A"

	conn, err := kafka.Dial("tcp", "localhost:19094")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
