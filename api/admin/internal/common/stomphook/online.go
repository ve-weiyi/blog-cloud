package stomphook

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-stomp/stomp/v3/frame"
	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/stompws/server/client"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
)

const defaultKeyPrefix = cachekey.OnlineAdminKey

// RedisOnlineTracker implements client.OnlineTracker using Redis ZSET for
// distributed online user tracking.
type RedisOnlineTracker struct {
	rdb *redis.Client
	key string
	hub *client.StompHubServer

	lastCount int64

	cleanupCtx    context.Context
	cleanupCancel context.CancelFunc
}

func NewRedisOnlineTracker(rdb *redis.Client, keyPrefix string) *RedisOnlineTracker {
	if keyPrefix == "" {
		keyPrefix = defaultKeyPrefix
	}
	ctx, cancel := context.WithCancel(context.Background())
	t := &RedisOnlineTracker{
		rdb:           rdb,
		key:           keyPrefix,
		cleanupCtx:    ctx,
		cleanupCancel: cancel,
	}
	go t.startCleanup()
	return t
}

func (t *RedisOnlineTracker) SetHub(hub *client.StompHubServer) {
	t.hub = hub
}

// ──────────── client.OnlineTracker interface ────────────

func (t *RedisOnlineTracker) OnConnect(login string) {
	_ = t.rdb.ZAdd(context.Background(), t.key, redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: login,
	}).Err()
	t.broadcastIfChanged(context.Background())
}

func (t *RedisOnlineTracker) OnDisconnect(login string) {
	_ = t.rdb.ZRem(context.Background(), t.key, login).Err()
	t.broadcastIfChanged(context.Background())
}

func (t *RedisOnlineTracker) OnActive(login string) {
	_ = t.rdb.ZAdd(context.Background(), t.key, redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: login,
	}).Err()
}

func (t *RedisOnlineTracker) GetOnlineCount(ctx context.Context) (int64, error) {
	return t.rdb.ZCard(ctx, t.key).Result()
}

func (t *RedisOnlineTracker) GetOnlineUsers(ctx context.Context, maxAgeSec int64) ([]*client.OnlineUser, error) {
	if maxAgeSec <= 0 {
		maxAgeSec = 60
	}
	threshold := time.Now().Add(-time.Duration(maxAgeSec) * time.Second).UnixMilli()
	results, err := t.rdb.ZRangeByScoreWithScores(ctx, t.key, &redis.ZRangeBy{
		Min: strconv.FormatInt(threshold, 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, err
	}
	users := make([]*client.OnlineUser, 0, len(results))
	for _, z := range results {
		userId, ok := z.Member.(string)
		if !ok {
			continue
		}
		users = append(users, &client.OnlineUser{
			UserId:       userId,
			LastActiveAt: int64(z.Score),
		})
	}
	return users, nil
}

func (t *RedisOnlineTracker) Close() error {
	t.cleanupCancel()
	return nil
}

// ──────────────────── internal ────────────────────

func (t *RedisOnlineTracker) broadcastIfChanged(ctx context.Context) {
	if t.hub == nil {
		return
	}
	count, err := t.GetOnlineCount(ctx)
	if err != nil {
		return
	}
	if count == t.lastCount {
		return
	}
	t.lastCount = count
	msg := frame.New(frame.MESSAGE, frame.Destination, "/topic/system/online", frame.MessageId, "0")
	msg.Body = []byte(fmt.Sprintf(`{"type":"online_count","count":%d}`, count))
	t.hub.RouteMessage(nil, msg)
}

func (t *RedisOnlineTracker) startCleanup() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			threshold := time.Now().Add(-60 * time.Second).UnixMilli()
			removed, _ := t.rdb.ZRemRangeByScore(context.Background(), t.key, "0", strconv.FormatInt(threshold, 10)).Result()
			if removed > 0 {
				t.broadcastIfChanged(context.Background())
			}
		case <-t.cleanupCtx.Done():
			return
		}
	}
}

// ──────────────────── OnlineCatchupHook ────────────────────

type OnlineCatchupHook struct {
	client.DefaultEventHook
	Tracker *RedisOnlineTracker
}

func NewOnlineCatchupHook(tracker *RedisOnlineTracker) *OnlineCatchupHook {
	return &OnlineCatchupHook{Tracker: tracker}
}

func (h *OnlineCatchupHook) OnSubscribe(server *client.StompHubServer, c *client.Client, destination string, subscriptionId string) {
	_, login, _ := c.GetClientInfo()
	if strings.HasPrefix(destination, "/topic/system/online") {
		count, _ := h.Tracker.GetOnlineCount(context.Background())
		msg := frame.New(frame.MESSAGE, frame.Destination, "/topic/system/online", frame.MessageId, "0", frame.Subscription, subscriptionId)
		msg.Body = []byte(fmt.Sprintf(`{"type":"online_count","count":%d}`, count))
		c.SendFrame(msg)
	}
	_ = login
}
