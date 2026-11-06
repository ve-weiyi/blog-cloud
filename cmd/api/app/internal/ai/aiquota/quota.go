// Package aiquota 统一 AI 配额的计数维度与读取口径。
package aiquota

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aictx"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/config"
)

const (
	// DailyPrefix 按身份计数的日配额前缀。
	DailyPrefix = "blog:ai:quota:"
	// IPPrefix 按来源计数的日配额前缀。
	IPPrefix = "blog:ai:ip:"
)

// Subject 返回身份维度：登录用户按用户 id，游客按设备标识。
func Subject(identity aictx.Identity) string {
	if identity.Guest {
		return "guest:" + identity.DeviceId
	}
	return "user:" + identity.UserId
}

// DailyKey 返回身份维度的日计数键。
func DailyKey(subject, day string) string {
	return DailyPrefix + subject + ":" + day
}

// IPDailyKey 返回来源维度的日计数键。
func IPDailyKey(ip, day string) string {
	return IPPrefix + ip + ":" + day
}

// Limit 返回身份维度的日额度；0 表示不限制。
func Limit(conf config.AiQuotaConf, guest bool) int64 {
	if guest {
		return conf.GuestDaily
	}
	return conf.UserDaily
}

// Used 读取当前已用次数；键不存在或缓存不可用时返回 0。
func Used(ctx context.Context, rdb *redis.Client, key string) (int64, error) {
	if rdb == nil {
		return 0, nil
	}
	used, err := rdb.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return used, nil
}
