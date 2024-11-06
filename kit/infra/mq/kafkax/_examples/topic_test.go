package _examples

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/mq/kafkax"
)

var c = &kafkax.KafkaConf{
	Username: "",
	Password: "",
	Brokers:  []string{"localhost:19094", "localhost:29094", "localhost:39094"},
	GroupID:  "consumer-group-id",
	Topic:    "Topic-A",
}

func Test_Topic_Publish(t *testing.T) {

	mq1 := kafkax.NewKafkaProducer(c)

	for i := 0; i <= 100; i++ {
		fmt.Println(i)
		err := mq1.PublishMessage(nil, []byte("user online: "+strconv.Itoa(i)))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func Test_Topic_Subscribe1(t *testing.T) {

	mq := kafkax.NewKafkaConsumer(c)

	mq.SubscribeMessage(func(msg []byte) error {
		log.Printf("1 receive message: %s", string(msg))
		return nil
	})

	select {}
}

func Test_Topic_Subscribe2(t *testing.T) {
	mq := kafkax.NewKafkaConsumer(c)

	mq.SubscribeMessage(func(msg []byte) error {
		log.Printf("2 receive message: %s", string(msg))
		return nil
	})

	select {}
}
