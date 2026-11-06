// Package aictx 在 AI 请求链路里传递身份与配额信息。
package aictx

import "context"

// Identity 是发起问答的主体：登录用户或带设备标识的游客。
type Identity struct {
	UserId   string
	DeviceId string
	Guest    bool
}

type identityKey struct{}

type quotaKey struct{}

// WithIdentity 写入当前请求的身份。
func WithIdentity(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, identityKey{}, id)
}

// IdentityFrom 读取当前请求的身份。
func IdentityFrom(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(identityKey{}).(Identity)
	return id, ok
}

// WithQuotaLeft 写入本轮请求后的剩余额度。
func WithQuotaLeft(ctx context.Context, left int64) context.Context {
	return context.WithValue(ctx, quotaKey{}, left)
}

// QuotaLeftFrom 读取剩余额度；未知时返回 0。
func QuotaLeftFrom(ctx context.Context) int64 {
	left, _ := ctx.Value(quotaKey{}).(int64)
	return left
}
