package kafkax

import (
	"context"
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/logz"
)

type KafkaProducer struct {
	logz.Logger

	w *kafka.Writer
}

func NewKafkaProducer(c *KafkaConf) *KafkaProducer {
	mechanism, err := scram.Mechanism(scram.SHA512, "username", "password")
	if err != nil {
		panic(err)
	}

	sharedTransport := &kafka.Transport{
		SASL: mechanism,
		TLS:  &tls.Config{},
	}

	w := &kafka.Writer{
		Addr:      kafka.TCP(c.Brokers...),
		Topic:     c.Topic,
		Transport: sharedTransport,
	}

	return &KafkaProducer{
		w:      w,
		Logger: logz.S(),
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
