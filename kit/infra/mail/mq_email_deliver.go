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

func NewMqEmailDeliver(deliver *EmailDeliver, pb mq.MessagePublisher, sb mq.MessageSubscriber) *MqEmailDeliver {
	return &MqEmailDeliver{
		Deliver:    deliver,
		Publisher:  pb,
		Subscriber: sb,
	}
}

func (m *MqEmailDeliver) DeliveryEmail(msg *EmailMessage) error {
	ms := m.Publisher
	if ms == nil {
		return fmt.Errorf("mq publisher is nil")
	}

	// 序列化消息
	jb, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 发送邮件
	err = ms.PublishMessage(nil, jb)
	if err != nil {
		return err
	}
	return nil
}

func (m *MqEmailDeliver) SubscribeEmail() {
	ms := m.Subscriber
	deliver := m.Deliver

	if ms == nil {
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
		err = deliver.DeliveryEmail(&msg)
		if err != nil {
			return err
		}

		return nil
	}

	// 订阅消息队列，发送邮件
	ms.SubscribeMessage(handler)
	return
}
