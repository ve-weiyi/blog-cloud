package client

import (
	"bytes"
	"strings"
	"sync"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
	"github.com/gorilla/websocket"

	"github.com/ve-weiyi/ve-blog-golang/stompws/logws"
)

type Client struct {
	conn          *websocket.Conn
	id            string // id is usually client-id or terminal-id, and its default value is the ip.
	login         string // login is usually user-id or username, and its default value is the ip.
	ip            string // ip address.
	version       string
	subscriptions map[string]*Subscription
	subList       *SubscriptionList
	readChan      chan *frame.Frame
	writeChan     chan *frame.Frame
	subChan       chan *Subscription
	closeChan     chan struct{}
	closeOnce     sync.Once
	lastMsgId     uint64
	lastReadTime  time.Time
	readTimeout   time.Duration
	writeTimeout  time.Duration
	validator     stomp.Validator
	connected     bool
	txStore       *txStore
	log           logws.Logger
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:          conn,
		subscriptions: make(map[string]*Subscription),
		subList:       NewSubscriptionList(),
		readChan:      make(chan *frame.Frame, 16),
		writeChan:     make(chan *frame.Frame, 256),
		subChan:       make(chan *Subscription, 16),
		closeChan:     make(chan struct{}),
		lastReadTime:  time.Now(),
		txStore:       newTxStore(),
	}
}

func (c *Client) readLoop() {
	defer c.close()

	for {
		if c.readTimeout > 0 {
			c.conn.SetReadDeadline(time.Now().Add(c.readTimeout * 2))
		}

		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.log.Errorf("client=%s: read error: %v", c.id, err)
			return
		}

		c.lastReadTime = time.Now()

		reader := bytes.NewReader(data)
		f, err := frame.NewReader(reader).Read()
		if err != nil {
			c.log.Warningf("client=%s: frame parse error: %v", c.id, err)
			continue
		}
		if f == nil {
			continue
		}

		select {
		case c.readChan <- f:
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) writeLoop(s *StompHubServer) {
	var ticker *time.Ticker
	var tickerChan <-chan time.Time

	if c.writeTimeout > 0 {
		ticker = time.NewTicker(c.writeTimeout)
		tickerChan = ticker.C
		defer ticker.Stop()
	}

	for {
		select {
		case f := <-c.writeChan:
			if err := c.conn.WriteMessage(websocket.TextMessage, frameToBytes(f)); err != nil {
				c.close()
				return
			}
		case sub := <-c.subChan:
			c.allocateMessageId(sub.frame, sub)
			if err := c.conn.WriteMessage(websocket.TextMessage, frameToBytes(sub.frame)); err != nil {
				c.close()
				return
			}
			if sub.ack == frame.AckAuto {
				sub.frame = nil
				if strings.HasPrefix(sub.destination, "/queue/") {
					s.queueManager.Find(sub.destination).Subscribe(sub)
				}
			} else {
				c.subList.Add(sub)
			}
		case <-tickerChan:
			if c.writeTimeout > 0 {
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					c.close()
					return
				}
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) processLoop(s *StompHubServer) {
	defer s.disconnect(c)

	var ticker *time.Ticker
	var tickerChan <-chan time.Time

	if c.readTimeout > 0 {
		ticker = time.NewTicker(c.readTimeout)
		tickerChan = ticker.C
		defer ticker.Stop()
	}

	for {
		select {
		case f := <-c.readChan:
			c.handleFrame(s, f)
		case <-tickerChan:
			if c.readTimeout > 0 && time.Since(c.lastReadTime) > c.readTimeout*3 {
				c.close()
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) handleFrame(s *StompHubServer, f *frame.Frame) {
	// 验证帧
	if c.validator != nil {
		if err := c.validator.Validate(f); err != nil {
			c.log.Warningf("client=%s user=%s: validation failed for %s frame: %v", c.id, c.login, f.Command, err)
			c.SendError(err.Error(), f.Header.Get(frame.Receipt))
			return
		}
	}

	// 检查连接状态
	if !c.connected && f.Command != frame.CONNECT && f.Command != frame.STOMP {
		c.SendError(errNotConnected.Error(), "")
		return
	}

	if c.connected && (f.Command == frame.CONNECT || f.Command == frame.STOMP) {
		c.SendError(errAlreadyConnected.Error(), "")
		return
	}

	switch f.Command {
	case frame.CONNECT, frame.STOMP:
		s.handleConnect(c, f)
	case frame.DISCONNECT:
		s.handleDisconnect(c, f)
	case frame.SUBSCRIBE:
		s.handleSubscribe(c, f)
	case frame.SEND:
		s.handleSend(c, f)
	case frame.UNSUBSCRIBE:
		s.handleUnsubscribe(c, f)
	case frame.ACK:
		s.handleAck(c, f)
	case frame.NACK:
		s.handleNack(c, f)
	case frame.BEGIN:
		s.handleBegin(c, f)
	case frame.COMMIT:
		s.handleCommit(c, f)
	case frame.ABORT:
		s.handleAbort(c, f)
	default:
		c.SendError(errUnknownCommand.Error(), f.Header.Get(frame.Receipt))
	}
}

func (c *Client) close() {
	c.closeOnce.Do(func() {
		close(c.closeChan)
		c.conn.Close()
	})
}

func (c *Client) SendFrame(f *frame.Frame) {
	select {
	case c.writeChan <- f:
	case <-time.After(time.Second):
		c.close()
	case <-c.closeChan:
	}
}

func (c *Client) SendError(message string, receiptId string) {
	errFrame := frame.New(frame.ERROR, frame.Message, message)
	if receiptId != "" {
		errFrame.Header.Add(frame.ReceiptId, receiptId)
	}
	c.SendFrame(errFrame)
}

func (c *Client) allocateMessageId(f *frame.Frame, sub *Subscription) {
	if f.Command == frame.MESSAGE {
		c.lastMsgId++
		msgId := uint64ToString(c.lastMsgId)
		f.Header.Set(frame.MessageId, msgId)

		if sub != nil {
			sub.msgId = c.lastMsgId
			f.Header.Set(frame.Subscription, sub.id)
			if sub.ack != frame.AckAuto {
				f.Header.Set(frame.Ack, msgId)
			}
		}
	}
}

// GetClientInfo returns client information for use in callbacks
func (c *Client) GetClientInfo() (id, login, ip, version string) {
	return c.id, c.login, c.ip, c.version
}
