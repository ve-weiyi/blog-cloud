package client

import (
	"github.com/go-stomp/stomp/v3/frame"
)

// EventHook defines callbacks for STOMP events
type EventHook interface {
	// OnConnect is called after a client successfully connects
	// server parameter allows sending messages back to client
	OnConnect(server *StompHubServer, c *Client)

	// OnDisconnect is called when a client disconnects
	OnDisconnect(server *StompHubServer, c *Client)

	// OnSubscribe is called after a client subscribes to a destination
	// server parameter allows sending welcome messages
	OnSubscribe(server *StompHubServer, c *Client, destination string, subscriptionId string)

	// OnUnsubscribe is called after a client unsubscribes from a destination
	OnUnsubscribe(server *StompHubServer, c *Client, destination string, subscriptionId string)

	// OnSend is called when a client sends a message
	// Return false to prevent the message from being sent
	OnSend(server *StompHubServer, c *Client, message *frame.Frame) bool

	// OnMessage is called before a message is delivered to a client
	// Return false to prevent the message from being delivered
	OnMessage(c *Client, destination string, message *frame.Frame) bool

	// OnAck is called when a client acknowledges a message
	OnAck(c *Client, messageId string)

	// OnNack is called when a client negatively acknowledges a message
	OnNack(c *Client, messageId string)

	// OnError is called when an error occurs
	OnError(c *Client, err error)
}

// DefaultEventHook provides default no-op implementations
type DefaultEventHook struct{}

func NewDefaultEventHook() *DefaultEventHook {
	return &DefaultEventHook{}
}

func (h *DefaultEventHook) OnConnect(server *StompHubServer, c *Client)    {}
func (h *DefaultEventHook) OnDisconnect(server *StompHubServer, c *Client) {}
func (h *DefaultEventHook) OnSubscribe(server *StompHubServer, c *Client, destination string, subscriptionId string) {
}
func (h *DefaultEventHook) OnUnsubscribe(server *StompHubServer, c *Client, destination string, subscriptionId string) {
}
func (h *DefaultEventHook) OnSend(server *StompHubServer, c *Client, message *frame.Frame) bool {
	return true
}
func (h *DefaultEventHook) OnMessage(c *Client, destination string, message *frame.Frame) bool {
	return true
}
func (h *DefaultEventHook) OnAck(c *Client, messageId string)  {}
func (h *DefaultEventHook) OnNack(c *Client, messageId string) {}
func (h *DefaultEventHook) OnError(c *Client, err error)       {}
