package rabbitmqx

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqProducer struct {
	conn *RabbitmqConn
	opts PublisherOptions
}

func NewRabbitmqProducer(conn *RabbitmqConn, opts ...PublisherOption) *RabbitmqProducer {
	opt := PublisherOptions{
		Exchange:   "",
		RoutingKey: "",
		Mandatory:  false,
		Immediate:  false,
	}

	for _, o := range opts {
		o(&opt)
	}

	return &RabbitmqProducer{
		conn: conn,
		opts: opt,
	}
}

func (r *RabbitmqProducer) PublishMessage(ctx context.Context, msg []byte) error {
	err := r.conn.Publish(amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	}, r.opts)

	if err != nil {
		r.conn.Infof("rabbitmq publish message failed: %v", err)
		return err
	}

	r.conn.Infof("rabbitmq publish message success")
	return nil
}
