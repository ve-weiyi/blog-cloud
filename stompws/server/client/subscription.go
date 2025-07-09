package client

import (
	"github.com/go-stomp/stomp/v3/frame"
)

type Subscription struct {
	client      *Client
	id          string
	destination string
	ack         string
	msgId       uint64
	subList     *SubscriptionList // am I in a list
	frame       *frame.Frame
}

func newSubscription(c *Client, dest string, id string, ack string) *Subscription {
	return &Subscription{
		client:      c,
		id:          id,
		destination: dest,
		ack:         ack,
	}
}

func (s *Subscription) Destination() string {
	return s.destination
}

func (s *Subscription) Ack() string {
	return s.ack
}

func (s *Subscription) Id() string {
	return s.id
}

func (s *Subscription) IsAckedBy(msgId uint64) bool {
	switch s.ack {
	case frame.AckAuto:
		return true
	case frame.AckClient:
		// any later message acknowledges an earlier message
		return msgId >= s.msgId
	case frame.AckClientIndividual:
		return msgId == s.msgId
	default:
		return false
	}
}

func (s *Subscription) IsNackedBy(msgId uint64) bool {
	return msgId == s.msgId
}

func (s *Subscription) SendQueueFrame(f *frame.Frame) {
	s.setSubscriptionHeader(f)
	s.frame = f

	// 设置frame并发送到subChan
	s.client.subChan <- s
}

func (s *Subscription) SendTopicFrame(f *frame.Frame) {
	s.setSubscriptionHeader(f)

	// let the connection deal with the subscription
	// acknowledgement
	s.client.writeChan <- f
}

func (s *Subscription) setSubscriptionHeader(f *frame.Frame) {
	if s.frame != nil {
		panic("subscription already has a frame pending")
	}
	f.Header.Set(frame.Subscription, s.id)
}
