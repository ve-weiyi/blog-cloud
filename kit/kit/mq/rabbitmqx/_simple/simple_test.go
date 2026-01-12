package RabbitMQ

import (
	"fmt"
	"testing"
)

func Test_SimplePublish(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("code")
	rabbitmq.PublishSimple("Hello kuteng222!")
	fmt.Println("发送成功！")
}

func Test_SimpleConsume1(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("code")
	rabbitmq.ConsumeSimple()
}

func Test_SimpleConsume2(t *testing.T) {
	rabbitmq := NewRabbitMQSimple("code")
	rabbitmq.ConsumeSimple()
}
