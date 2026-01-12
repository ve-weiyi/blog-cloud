package rabbitmqx

import (
	"context"
)

type RabbitmqConsumer struct {
	conn *RabbitmqConn
	opts ConsumerOptions
}

func NewRabbitmqConsumer(conn *RabbitmqConn, opts ...ConsumerOption) *RabbitmqConsumer {
	opt := ConsumerOptions{
		Queue:     "",
		Name:      "",
		AutoAck:   false,
		Exclusive: false,
		NoWait:    false,
		NoLocal:   false,
		Args:      nil,
	}

	for _, o := range opts {
		o(&opt)
	}

	return &RabbitmqConsumer{
		conn: conn,
		opts: opt,
	}
}

func (r *RabbitmqConsumer) SubscribeMessage(handler func(ctx context.Context, message []byte) error) {
	//接收消息
	msgs, err := r.conn.Consume(r.opts)
	if err != nil {
		r.conn.Infof("RabbitmqConsumer rabbitmq consume error: %v", err)
		return
	}

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for msg := range msgs {
			// 处理消息
			err = handler(context.Background(), msg.Body)
			if err != nil {
				r.conn.Infof("RabbitmqConsumer rabbitmq subscribe handle error: %v", err)

				// 消费失败，根据配置决定是否重入队列
				if err = msg.Nack(false, true); err != nil {
					r.conn.Infof("Warning: rabbitmq nack msg error: %v", err)
				}
				continue
			}

			// AutoAck为true时，不进行手动ack确认
			if r.opts.AutoAck {
				continue
			}

			// 消费成功时确认消息
			if err = msg.Ack(false); err != nil {
				r.conn.Infof("RabbitmqConsumer rabbitmq ack msg error: %v", err)
				// 这里可以决定是否重新处理消息或报警
				continue
			}

			r.conn.Infof("RabbitmqConsumer handler message success: %s", msg.MessageId)
		}
	}()

	r.conn.Infof("RabbitmqConsumer [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return
}
