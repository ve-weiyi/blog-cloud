package client

import (
	"strings"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

func (s *StompHubServer) handleConnect(c *Client, f *frame.Frame) {
	// 身份验证
	clientId, login, err := s.authenticator.Authenticate(c, f)
	if err != nil {
		s.log.Errorf("authentication failed: %v", err)
		time.Sleep(authFailedDelay) // Prevent brute force
		c.SendError(err.Error(), "")
		return
	}

	c.id = clientId
	c.login = login

	// 版本协商
	version, err := negotiateVersion(f.Header.Get(frame.AcceptVersion))
	if err != nil {
		c.SendError(err.Error(), "")
		return
	}
	c.version = version

	// 心跳协商
	cx, cy := parseHeartBeat(f.Header.Get(frame.HeartBeat))
	minHeartBeat := 30 * time.Second

	if cx > 0 {
		if cx < minHeartBeat {
			cx = minHeartBeat
		}
		c.readTimeout = cx
	}

	if cy > 0 {
		if cy < minHeartBeat {
			cy = minHeartBeat
		}
		c.writeTimeout = cy
	}

	// 初始化验证器
	if err := stomp.Version(version).CheckSupported(); err == nil {
		c.validator = newValidator(stomp.Version(version))
	}

	c.connected = true
	s.clients.Store(c.id, c)

	s.log.Infof("client=%s login=%s: connected, version=%s, heartbeat=%s", c.id, c.login, version, formatHeartBeat(cy, cx))

	response := frame.New(frame.CONNECTED,
		frame.Version, version,
		frame.Session, c.id,
		frame.HeartBeat, formatHeartBeat(cy, cx))
	c.SendFrame(response)

	s.sendReceipt(c, f)

	// Trigger OnConnect callback
	for _, hook := range s.eventHooks {
		hook.OnConnect(s, c)
	}
}

func (s *StompHubServer) handleDisconnect(c *Client, f *frame.Frame) {
	// 先发送 receipt
	s.sendReceipt(c, f)

	// 执行断开连接的清理逻辑
	s.cleanupClient(c)

}

// cleanupClient 执行客户端断开连接时的清理工作
func (s *StompHubServer) cleanupClient(c *Client) {
	if !c.connected {
		return
	}

	// 清理事务
	c.txStore.Clear()

	s.clients.Delete(c.id)

	// 重新入队未确认的消息
	for sub := c.subList.Get(); sub != nil; sub = c.subList.Get() {
		if sub.frame != nil && strings.HasPrefix(sub.destination, "/queue/") {
			q := s.queueManager.Find(sub.destination)
			q.Requeue(sub.frame)
		}
	}

	// 取消订阅
	for _, sub := range c.subscriptions {
		s.UnsubscribeFromDestination(sub)
		// 订阅成功后触发回调
		for _, hook := range s.eventHooks {
			hook.OnUnsubscribe(s, c, sub.destination, sub.id)
		}
	}

	s.log.Infof("client=%s login=%s: disconnected", c.id, c.login)

	// Trigger OnDisconnect callback
	for _, hook := range s.eventHooks {
		hook.OnDisconnect(s, c)
	}

	c.connected = false
}

func (s *StompHubServer) handleSubscribe(c *Client, f *frame.Frame) {
	dest := f.Header.Get(frame.Destination)
	id := f.Header.Get(frame.Id)
	ack := f.Header.Get(frame.Ack)

	// 检查订阅冲突
	if _, exists := c.subscriptions[id]; exists {
		c.SendError(errSubscriptionExists.Error(), f.Header.Get(frame.Receipt))
		return
	}

	if ack == "" {
		ack = frame.AckAuto
	}

	sub := newSubscription(c, dest, id, ack)
	c.subscriptions[id] = sub

	if err := s.SubscribeToDestination(sub); err != nil {
		c.SendError(err.Error(), f.Header.Get(frame.Receipt))
		return
	}
	// 订阅成功后触发回调
	for _, hook := range s.eventHooks {
		hook.OnSubscribe(s, c, dest, id)
	}

	s.log.Infof("client=%s login=%s: subscribed to %s (ack=%s)", c.id, c.login, dest, ack)
	s.sendReceipt(c, f)
}

func (s *StompHubServer) handleUnsubscribe(c *Client, f *frame.Frame) {
	id := f.Header.Get(frame.Id)
	sub, ok := c.subscriptions[id]
	if !ok {
		c.SendError(errSubscriptionNotFound.Error(), "")
		return
	}

	delete(c.subscriptions, id)

	s.UnsubscribeFromDestination(sub)
	s.log.Infof("client=%s login=%s: unsubscribed from %s", c.id, c.login, sub.destination)

	// 先发送 receipt
	s.sendReceipt(c, f)

	// 取消订阅后触发回调
	for _, hook := range s.eventHooks {
		hook.OnUnsubscribe(s, c, sub.destination, id)
	}
}

func (s *StompHubServer) handleSend(c *Client, f *frame.Frame) {
	// 先发送 receipt
	s.sendReceipt(c, f)

	// 检查是否在事务中
	if tx, ok := f.Header.Contains(frame.Transaction); ok {
		if err := c.txStore.Add(tx, f); err != nil {
			c.SendError(err.Error(), "")
		}
		return
	}

	// Trigger OnSend callback
	for _, hook := range s.eventHooks {
		if !hook.OnSend(s, c, f) {
			return
		}
	}
}

func (s *StompHubServer) handleAck(c *Client, f *frame.Frame) {
	ackId := f.Header.Get(frame.Ack)
	if ackId == "" {
		ackId = f.Header.Get(frame.MessageId)
	}

	if ackId == "" {
		return
	}

	// 先发送 receipt
	s.sendReceipt(c, f)

	// 检查是否在事务中
	if tx, ok := f.Header.Contains(frame.Transaction); ok {
		if err := c.txStore.Add(tx, f); err != nil {
			c.SendError(err.Error(), "")
		}
		return
	}

	msgId := stringToUint64(ackId)

	// Trigger OnAck callback
	for _, hook := range s.eventHooks {
		hook.OnAck(c, ackId)
	}

	// 使用 SubscriptionList 处理 ACK
	c.subList.Ack(msgId, func(sub *Subscription) {
		sub.frame = nil
		// 重新订阅队列，传入Subscription对象
		if strings.HasPrefix(sub.destination, "/queue/") {
			s.queueManager.Find(sub.destination).Subscribe(sub)
		}
	})
}

func (s *StompHubServer) handleNack(c *Client, f *frame.Frame) {
	ackId := f.Header.Get(frame.Ack)
	if ackId == "" {
		ackId = f.Header.Get(frame.MessageId)
	}

	// 先发送 receipt
	s.sendReceipt(c, f)

	// 检查是否在事务中
	if tx, ok := f.Header.Contains(frame.Transaction); ok {
		if err := c.txStore.Add(tx, f); err != nil {
			c.SendError(err.Error(), "")
		}
		return
	}

	msgId := stringToUint64(ackId)

	// Trigger OnNack callback
	for _, hook := range s.eventHooks {
		hook.OnNack(c, ackId)
	}

	// 使用 SubscriptionList 处理 NACK
	c.subList.Nack(msgId, func(sub *Subscription) {
		// 重新入队到队首
		if strings.HasPrefix(sub.destination, "/queue/") {
			q := s.queueManager.Find(sub.destination)
			q.Requeue(sub.frame)
		}

		sub.frame = nil

		// 重新订阅队列，传入Subscription对象
		if strings.HasPrefix(sub.destination, "/queue/") {
			s.queueManager.Find(sub.destination).Subscribe(sub)
		}
	})
}

func (s *StompHubServer) handleBegin(c *Client, f *frame.Frame) {
	tx := f.Header.Get(frame.Transaction)
	if err := c.txStore.Begin(tx); err != nil {
		c.SendError(err.Error(), f.Header.Get(frame.Receipt))
		return
	}

	s.log.Infof("client=%s login=%s: transaction %s begun", c.id, c.login, tx)
	s.sendReceipt(c, f)
}

func (s *StompHubServer) handleCommit(c *Client, f *frame.Frame) {
	tx := f.Header.Get(frame.Transaction)
	// 先发送 receipt
	s.sendReceipt(c, f)

	// 提交事务，执行所有缓存的帧
	err := c.txStore.Commit(tx, func(frame *frame.Frame) error {
		// 重新处理每个帧
		c.handleFrame(s, frame)
		return nil
	})

	if err != nil {
		s.log.Errorf("client=%s login=%s: transaction %s commit failed: %v", c.id, c.login, tx, err)
		c.SendError(err.Error(), "")
	} else {
		s.log.Infof("client=%s login=%s: transaction %s committed", c.id, c.login, tx)
	}
}

func (s *StompHubServer) handleAbort(c *Client, f *frame.Frame) {
	tx := f.Header.Get(frame.Transaction)
	if err := c.txStore.Abort(tx); err != nil {
		c.SendError(err.Error(), f.Header.Get(frame.Receipt))
		return
	}

	s.log.Infof("client=%s login=%s: transaction %s aborted", c.id, c.login, tx)
	s.sendReceipt(c, f)
}

func (s *StompHubServer) sendReceipt(c *Client, f *frame.Frame) {
	if receipt, ok := f.Header.Contains(frame.Receipt); ok {
		receiptFrame := frame.New(frame.RECEIPT, frame.ReceiptId, receipt)
		c.SendFrame(receiptFrame)
	}
}
