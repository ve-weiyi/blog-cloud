package mail

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/mq"
)

type MqEmailDeliver struct {
	Deliver *EmailDeliver

	Publisher  mq.MessagePublisher
	Subscriber mq.MessageSubscriber
}

func NewMqEmailDeliver(cfg *EmailConfig, pb mq.MessagePublisher, sb mq.MessageSubscriber) *MqEmailDeliver {
	deliver := &MqEmailDeliver{
		Deliver:    NewEmailDeliver(cfg),
		Publisher:  pb,
		Subscriber: sb,
	}

	return deliver
}

func (m *MqEmailDeliver) DeliveryEmail(msg *EmailMessage) error {
	if m.Publisher == nil {
		return fmt.Errorf("mq publisher is nil")
	}

	// 序列化消息
	jb, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 发送邮件
	err = m.Publisher.PublishMessage(nil, jb)
	if err != nil {
		return err
	}
	return nil
}

func (m *MqEmailDeliver) SubscribeEmail() {
	if m.Subscriber == nil {
		return
	}

	// 消息处理函数
	handler := func(ctx context.Context, message []byte) (err error) {
		var msg EmailMessage
		// 反序列化消息
		err = json.Unmarshal(message, &msg)
		if err != nil {
			return err
		}

		// 发送邮件
		err = m.Deliver.DeliveryEmail(&msg)
		if err != nil {
			return err
		}

		return nil
	}

	// 订阅消息队列，发送邮件
	m.Subscriber.SubscribeMessage(handler)
	return
}
