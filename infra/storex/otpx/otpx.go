package otpx

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
	"github.com/ve-weiyi/vkit/x/randomx"
)

// Store OTP 验证码存储
type Store struct {
	store storex.KVStore
}

// New 创建 OTP 存储实例
func New(store storex.KVStore) *Store {
	return &Store{store: store}
}

// Key 构建存储 key：{channel}:otp:{scene}:{target}
// channel / scene / target 均不应包含冒号。
//
// 示例: Key("email", "register", "user@example.com") → "email:otp:register:user@example.com"
func Key(channel, scene, target string) string {
	return fmt.Sprintf("%s:otp:%s:%s", channel, scene, target)
}

// ---- Generate ----------------------------------------------------------------

// GenerateOption 生成参数选项
type GenerateOption func(*generateOptions)

type generateOptions struct {
	length int
	expire time.Duration
}

// WithLength 设置验证码长度（默认 6，范围 4-10）
func WithLength(n int) GenerateOption {
	return func(o *generateOptions) {
		o.length = n
	}
}

// WithExpire 设置过期时间（默认 10 分钟）
func WithExpire(d time.Duration) GenerateOption {
	return func(o *generateOptions) {
		o.expire = d
	}
}

const (
	defaultLength = 6
	minLength     = 4
	maxLength     = 10
	defaultExpire = 10 * time.Minute
)

// Generate 生成随机验证码并存储。
// 默认 6 位数字、10 分钟过期，可通过 GenerateOption 覆盖。
// 并发安全取决于底层 KVStore 实现（Redis 版本安全）。
func (s *Store) Generate(ctx context.Context, key string, opts ...GenerateOption) (string, error) {
	if key == "" {
		return "", errors.New("key is empty")
	}

	o := generateOptions{
		length: defaultLength,
		expire: defaultExpire,
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.length < minLength || o.length > maxLength {
		return "", fmt.Errorf("length must be between %d and %d", minLength, maxLength)
	}
	if o.expire <= 0 {
		o.expire = defaultExpire
	}

	code, err := randomx.GenerateCode(o.length)
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}
	if err := s.store.Set(ctx, key, code, o.expire); err != nil {
		return "", fmt.Errorf("store code: %w", err)
	}

	return code, nil
}

// Verify 校验验证码，匹配后删除（一次性使用）。
// 使用 constant-time 比较防止时序攻击。
// 并发安全取决于底层 KVStore 实现（Redis 版本安全）。
func (s *Store) Verify(ctx context.Context, key, code string) (bool, error) {
	if key == "" || code == "" {
		return false, errors.New("key or code is empty")
	}

	stored, found, err := s.store.Get(ctx, key)
	if err != nil {
		return false, fmt.Errorf("get code: %w", err)
	}
	if !found {
		return false, nil
	}

	if subtle.ConstantTimeCompare([]byte(stored), []byte(code)) != 1 {
		return false, nil
	}

	if err := s.store.Delete(ctx, key); err != nil {
		return false, err
	}
	return true, nil
}
