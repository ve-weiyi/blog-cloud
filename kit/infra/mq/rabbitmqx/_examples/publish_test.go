package RabbitMQ

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/mq/rabbitmqx"
)

// 模拟一个fanout模式：邮件订阅

const ExchangeEmail = "email_exchange"

// 发布者
const TopicEmail = "blog.email" // fanout模式不需要routing key

// 订阅者
const QueueEmail = "email_queue" //相同队列名称会争抢消息

const QueueEmail2 = "email_queue2"

func Test_Fanout_Publish(t *testing.T) {
	conn, err := rabbitmqx.NewRabbitmqConn(
		"amqp://veweiyi:rabbitmq7914@localhost:5672",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	exchange := &rabbitmqx.ExchangeOptions{
		Name:    ExchangeEmail,
		Kind:    rabbitmqx.ExchangeTypeFanout,
		Durable: true, // 是否持久化
	}

	err = conn.Declare(nil, exchange, nil)
	if err != nil {
		log.Fatal(err)
	}

	// pub/sub 模式发布者只需要声明交换机
	mq1 := rabbitmqx.NewRabbitmqProducer(conn,
		rabbitmqx.WithPublisherExchange(exchange.Name),
		rabbitmqx.WithPublisherMandatory(true),
	)

	for i := 0; i <= 100; i++ {
		fmt.Println(i)
		err = mq1.PublishMessage(nil, []byte("user email: "+strconv.Itoa(i)))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func Test_Fanout_Subscribe1(t *testing.T) {
	conn, err := rabbitmqx.NewRabbitmqConn(
		"amqp://veweiyi:rabbitmq7914@localhost:5672",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queue := &rabbitmqx.QueueOptions{
		Name:    QueueEmail, //不填是随机队列名称
		Durable: true,       // 是否持久化
	}

	exchange := &rabbitmqx.ExchangeOptions{
		Name:    ExchangeEmail,
		Kind:    rabbitmqx.ExchangeTypeFanout,
		Durable: true, // 是否持久化
	}

	binding := &rabbitmqx.BindingOptions{
		RoutingKey: TopicEmail,
	}

	err = conn.Declare(queue, exchange, binding)
	if err != nil {
		log.Fatal(err)
	}

	mq := rabbitmqx.NewRabbitmqConsumer(conn,
		rabbitmqx.WithConsumerQueue(queue.Name),
		rabbitmqx.WithConsumerAutoAck(true),
	)

	mq.SubscribeMessage(func(ctx context.Context, msg []byte) error {
		log.Printf("receive message: %s", string(msg))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

func Test_Fanout_Subscribe2(t *testing.T) {
	conn, err := rabbitmqx.NewRabbitmqConn(
		"amqp://veweiyi:rabbitmq7914@localhost:5672",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queue := &rabbitmqx.QueueOptions{
		Name:    QueueEmail, //不填是随机队列名称
		Durable: true,       // 是否持久化
	}

	exchange := &rabbitmqx.ExchangeOptions{
		Name:    ExchangeEmail,
		Kind:    rabbitmqx.ExchangeTypeFanout,
		Durable: true, // 是否持久化
	}

	binding := &rabbitmqx.BindingOptions{
		RoutingKey: TopicEmail,
	}

	err = conn.Declare(queue, exchange, binding)
	if err != nil {
		log.Fatal(err)
	}

	mq := rabbitmqx.NewRabbitmqConsumer(conn,
		rabbitmqx.WithConsumerQueue(queue.Name),
		rabbitmqx.WithConsumerAutoAck(true),
	)

	mq.SubscribeMessage(func(ctx context.Context, msg []byte) error {
		log.Printf("receive message: %s", string(msg))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

func Test_Fanout_Subscribe3(t *testing.T) {
	conn, err := rabbitmqx.NewRabbitmqConn(
		"amqp://veweiyi:rabbitmq7914@localhost:5672",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queue := &rabbitmqx.QueueOptions{
		Name:    QueueEmail2, //不填是随机队列名称
		Durable: true,        // 是否持久化
	}

	exchange := &rabbitmqx.ExchangeOptions{
		Name:    ExchangeEmail,
		Kind:    rabbitmqx.ExchangeTypeFanout,
		Durable: true, // 是否持久化
	}

	binding := &rabbitmqx.BindingOptions{
		RoutingKey: TopicEmail,
	}

	err = conn.Declare(queue, exchange, binding)
	if err != nil {
		log.Fatal(err)
	}

	mq := rabbitmqx.NewRabbitmqConsumer(conn,
		rabbitmqx.WithConsumerQueue(queue.Name),
		rabbitmqx.WithConsumerAutoAck(true),
	)

	mq.SubscribeMessage(func(ctx context.Context, msg []byte) error {
		log.Printf("receive message: %s", string(msg))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
