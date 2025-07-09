package queue

import (
	"github.com/go-stomp/stomp/v3/frame"
)

// Subscription is the interface for queue subscribers.
type Subscription interface {
	// Send a message frame to the queue subscriber.
	SendQueueFrame(f *frame.Frame)

	// Destination returns the subscription destination.
	Destination() string
}
