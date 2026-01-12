package kafkax

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/logz"
)

type KafkaConsumer struct {
	logz.Logger

	r *kafka.Reader
}

func NewKafkaConsumer(c *KafkaConf) *KafkaConsumer {

	mechanism, err := scram.Mechanism(scram.SHA512, c.Username, c.Password)
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Brokers,
		GroupID:  c.GroupID,
		Topic:    c.Topic,
		Dialer:   dialer,
		MaxWait:  10 * time.Second,
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		r:      r,
		Logger: logz.S(),
	}
}

func (ms *KafkaConsumer) SubscribeMessage(handler func(ctx context.Context, msg []byte) error) {
	ms.Info("KafkaConsumer [*] Waiting for messages. To exit press CTRL+C")
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					ms.Errorf("KafkaConsumer Recovered  err:%v", r)
				}
			}()

			ctx := context.Background()
			m, err := ms.r.ReadMessage(ctx)
			if err != nil {
				ms.Errorf("KafkaConsumer readMessage err:%v", err)
				return
			}

			err = handler(ctx, m.Value)
			if err != nil {
				ms.Errorf("KafkaConsumer handler err:%v", err)
				return
			}

			ms.Info("KafkaConsumer handler success", string(m.Key))
		}()
	}
}
