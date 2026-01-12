package online

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/rediskey"
)

type OnlineUserService struct {
	rdb          *redis.Client
	key          string
	expirePeriod time.Duration
}

func NewOnlineUserService(rdb *redis.Client, expire int) *OnlineUserService {
	return &OnlineUserService{
		rdb:          rdb,
		key:          rediskey.GetOnlineUserKey(),
		expirePeriod: time.Duration(expire) * time.Second,
	}
}

// 用户登录（加入 ZSet，score 为当前时间戳）
func (s *OnlineUserService) Login(ctx context.Context, userID string) error {
	score := float64(time.Now().Unix())
	return s.rdb.ZAdd(ctx, s.key, redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
}

// 用户退出（从 ZSet 移除）
func (s *OnlineUserService) Logout(ctx context.Context, userID string) error {
	return s.rdb.ZRem(ctx, s.key, userID).Err()
}

// 清除过期用户（30分钟未活动）
func (s *OnlineUserService) CleanExpired(ctx context.Context) error {
	expireBefore := float64(time.Now().Add(-s.expirePeriod).Unix())
	return s.rdb.ZRemRangeByScore(ctx, s.key, "-inf", fmt.Sprintf("%f", expireBefore)).Err()
}

// 获取在线用户 ID 列表（先清理，再返回）
func (s *OnlineUserService) GetOnlineUsers(ctx context.Context, page, pageSize int64) ([]string, error) {
	offset := (page - 1) * pageSize
	limit := pageSize

	if err := s.CleanExpired(ctx); err != nil {
		return nil, err
	}

	result, err := s.rdb.ZRevRangeWithScores(ctx, s.key, offset, offset+limit-1).Result()
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, item := range result {
		uids = append(uids, item.Member.(string))
	}

	return uids, nil
}

// 获取在线用户数量（先清理）
func (s *OnlineUserService) GetOnlineUserCount(ctx context.Context) (int64, error) {
	if err := s.CleanExpired(ctx); err != nil {
		return 0, err
	}
	return s.rdb.ZCard(ctx, s.key).Result()
}
