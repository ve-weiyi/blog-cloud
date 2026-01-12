package client

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-stomp/stomp/v3/frame"
	"github.com/gorilla/websocket"

	"github.com/ve-weiyi/ve-blog-golang/stompws/logws"
	"github.com/ve-weiyi/ve-blog-golang/stompws/server/queue"
	"github.com/ve-weiyi/ve-blog-golang/stompws/server/topic"
)

type StompHubServer struct {
	clients sync.Map // clientId -> *Client

	topicManager  *topic.Manager
	queueManager  *queue.Manager
	upgrader      websocket.Upgrader
	log           logws.Logger
	authenticator Authenticator
	eventHooks    []EventHook
}

type ServerOption func(*StompHubServer)

func WithAuthenticator(auth Authenticator) ServerOption {
	return func(s *StompHubServer) {
		s.authenticator = auth
	}
}

func WithEventHooks(handlers ...EventHook) ServerOption {
	return func(s *StompHubServer) {
		s.eventHooks = append(s.eventHooks, handlers...)
	}
}

func WithLogger(logger logws.Logger) ServerOption {
	return func(s *StompHubServer) {
		s.log = logger
	}
}

func WithQueueStorage(storage queue.Storage) ServerOption {
	return func(s *StompHubServer) {
		s.queueManager = queue.NewManager(storage)
	}
}

func WithCheckOrigin(checkOrigin func(*http.Request) bool) ServerOption {
	return func(s *StompHubServer) {
		s.upgrader.CheckOrigin = checkOrigin
	}
}

func NewStompHubServer(opts ...ServerOption) *StompHubServer {
	s := &StompHubServer{
		topicManager: topic.NewManager(),
		queueManager: queue.NewManager(queue.NewMemoryQueueStorage()),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols: []string{"stomp", "mqtt", "v12.stomp", "v11.stomp", "v10.stomp"},
		},
		log:           logws.NewDefaultLogger(),
		authenticator: NewNoAuthenticator(),
		eventHooks:    []EventHook{NewDefaultEventHook()},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *StompHubServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("websocket upgrade failed: %v", err)
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	c := newClient(conn)
	c.ip = ip
	c.log = s.log
	s.log.Infof("new connection from %s", ip)

	go c.readLoop()
	go c.writeLoop(s)
	go c.processLoop(s)
}

func (s *StompHubServer) disconnect(c *Client) {
	// 执行清理逻辑
	s.cleanupClient(c)
}

// RouteMessage routes a message based on its destination
func (s *StompHubServer) RouteMessage(from *Client, msg *frame.Frame) {
	dest := msg.Header.Get(frame.Destination)
	switch {
	case strings.HasPrefix(dest, "/topic/"):
		s.topicManager.Find(dest).Enqueue(msg)
	case strings.HasPrefix(dest, "/queue/"):
		s.queueManager.Find(dest).Enqueue(msg)
	default:
		if from != nil {
			errorMsg := frame.New(frame.MESSAGE, frame.Destination, "/topic/system", frame.MessageId, "0")
			errorMsg.Body = []byte(fmt.Sprintf(`{"username":"System","content":"Invalid destination: %s"}`, dest))
			from.SendFrame(errorMsg)
		}
	}
}

// SubscribeToDestination handles subscription based on destination type
func (s *StompHubServer) SubscribeToDestination(sub *Subscription) error {
	dest := sub.Destination()
	switch {
	case strings.HasPrefix(dest, "/topic/"):
		s.topicManager.Find(dest).Subscribe(sub)
	case strings.HasPrefix(dest, "/queue/"):
		s.queueManager.Find(dest).Subscribe(sub)
	default:
		return errInvalidDestination
	}
	return nil
}

// UnsubscribeFromDestination handles unsubscription based on destination type
func (s *StompHubServer) UnsubscribeFromDestination(sub *Subscription) error {
	dest := sub.Destination()
	switch {
	case strings.HasPrefix(dest, "/topic/"):
		s.topicManager.Find(dest).Unsubscribe(sub)
	case strings.HasPrefix(dest, "/queue/"):
		s.queueManager.Find(dest).Unsubscribe(sub)
	default:
		return errInvalidDestination
	}
	return nil
}
