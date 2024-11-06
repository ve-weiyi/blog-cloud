package kafkax

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/glog"
)

type KafkaConsumer struct {
	glog.Logger

	r *kafka.Reader
}

func NewKafkaConsumer(c *KafkaConf) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Brokers,
		GroupID:  c.GroupID,
		Topic:    c.Topic,
		MaxWait:  10 * time.Second,
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		Logger: glog.Default(),
		r:      r,
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
