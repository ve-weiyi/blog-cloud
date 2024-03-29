package rabbitmq

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// 连接信息amqp://guest:guest@127.0.0.1:5672/guest这个信息是固定不变的amqp://是固定参数 后面两个是用户名密码ip地址端口号Virtual Host
const rabbitMQURL = "amqp://guest:guest@localhost:5672/"

func Test_SimplePublish(t *testing.T) {
	var err error
	rabbitmq := NewRabbitmqConn(rabbitMQURL,
		Queue(QueueOptions{
			Name:    "code",
			Durable: true,
			Args:    nil,
		}),
		WithoutExchange(),
	)

	err = rabbitmq.Connect(nil)
	t.Log(err)

	for i := 0; i <= 10; i++ {
		err := rabbitmq.PublishMessage([]byte("Hello,RabbitMQ!" + strconv.Itoa(i)))
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}

func Test_SimpleConsumer(t *testing.T) {
	var err error
	rabbitmq := NewRabbitmqConn(rabbitMQURL,
		Queue(QueueOptions{
			Name:    "code",
			Durable: true,
			Args:    nil,
		}),
		WithoutExchange(),
	)

	err = rabbitmq.Connect(nil)
	t.Log(err)

	err = rabbitmq.SubscribeMessage(func(message []byte) error {
		fmt.Println(string(message))
		return nil
	})
	t.Log(err)
}
