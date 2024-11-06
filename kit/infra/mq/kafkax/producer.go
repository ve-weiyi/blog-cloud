package kafkax

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/glog"
)

type KafkaProducer struct {
	glog.Logger

	w *kafka.Writer
}

func NewKafkaProducer(c *KafkaConf) *KafkaProducer {
	w := &kafka.Writer{
		Addr:  kafka.TCP(c.Brokers...),
		Topic: c.Topic,
	}

	return &KafkaProducer{
		w:      w,
		Logger: glog.Default(),
	}
}

func (mq *KafkaProducer) PublishMessage(ctx context.Context, msg []byte) error {
	message := kafka.Message{
		Value: msg,
	}

	err := mq.w.WriteMessages(ctx, message)

	if err != nil {
		mq.Errorf("KafkaProducer publish message failed: %v", err)
		return err
	}

	mq.Info("KafkaProducer publish message success", message.Key)
	return nil
}
