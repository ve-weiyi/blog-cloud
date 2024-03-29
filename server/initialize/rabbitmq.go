package initialize

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/mail"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/rabbitmq"
)

const (
	// email交换机
	EMAIL_EXCHANGE = "email_exchange"
	// email队列
	EMAIL_QUEUE = "email_queue"
)

// 订阅消息
func RabbitMq() {
	url := global.CONFIG.RabbitMQ.GetUrl()

	// 消息发布者只需要声明交换机
	mq := rabbitmq.NewRabbitmqConn(url,
		rabbitmq.Exchange(rabbitmq.ExchangeOptions{
			Name:    EMAIL_EXCHANGE,
			Type:    rabbitmq.ExchangeTypeFanout,
			Durable: true,
		}),
		rabbitmq.DisableAutoAck(),
		rabbitmq.Requeue(),
	)

	err := mq.Connect(nil)
	if err != nil {
		log.Fatal("rabbitmq 初始化失败!", err)
	}

	global.EmailMQ = mq

	go SubscribeMessage()
}

func SubscribeMessage() {
	url := global.CONFIG.RabbitMQ.GetUrl()
	// 消息订阅者需要声明交换机和队列
	mq := rabbitmq.NewRabbitmqConn(url,
		rabbitmq.Queue(rabbitmq.QueueOptions{
			Name:    EMAIL_QUEUE,
			Durable: true,
			Args:    nil,
		}),
		rabbitmq.Exchange(rabbitmq.ExchangeOptions{
			Name:    EMAIL_EXCHANGE,
			Type:    rabbitmq.ExchangeTypeFanout,
			Durable: true,
		}),
		rabbitmq.Key("email"),
	)
	err := mq.Connect(nil)
	if err != nil {
		log.Fatal("rabbitmq 初始化失败!", err)
	}

	cfg := global.CONFIG.Email
	emailSender := &mail.EmailSender{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		Nickname: cfg.Nickname,
		Deliver:  strings.Split(cfg.Deliver, ","),
		IsSSL:    cfg.IsSSL,
	}

	//订阅消息队列，发送邮件
	err = mq.SubscribeMessage(func(message []byte) (err error) {
		var msg mail.EmailMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			return err
		}

		err = emailSender.SendEmailMessage(msg)
		if err != nil {
			global.LOG.Error("邮件发送失败!", err)
		}
		return err
	})
	if err != nil {
		log.Fatal("订阅消息失败!", err)
	}
}
