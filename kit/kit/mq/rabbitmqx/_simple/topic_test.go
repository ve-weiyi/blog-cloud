package RabbitMQ

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	kutengOne := NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	kutengTwo := NewRabbitMQTopic("exKutengTopic", "kuteng.topic.two")
	for i := 0; i <= 100; i++ {
		kutengOne.PublishTopic("Hello kuteng topic one!" + strconv.Itoa(i))
		kutengTwo.PublishTopic("Hello kuteng topic Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}

func TestConsume1(t *testing.T) {
	kutengOne := NewRabbitMQTopic("exKutengTopic", "#")
	kutengOne.ReceiveTopic()

}

func TestConsume2(t *testing.T) {
	kutengOne := NewRabbitMQTopic("exKutengTopic", "kuteng.*.two")
	kutengOne.ReceiveTopic()
}
