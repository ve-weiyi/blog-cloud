package tokenx

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
)

type TokenHolder interface {
	TokenType() string
	VerifyToken(ctx context.Context, token string, uid string) error
	CreateToken(ctx context.Context, uid string, expires time.Duration) (string, error)
	RemoveToken(ctx context.Context, uid string) error
}

var (
	ErrTokenEmpty   = fmt.Errorf("token is empty")
	ErrTokenInvalid = fmt.Errorf("token is invalid")
	ErrTokenExpired = fmt.Errorf("token is expired")
)

type SignTokenHolder struct {
	issuer string
	secret string

	cache *redis.Redis
}

func NewSignTokenHolder(issuer string, secret string, cache *redis.Redis) *SignTokenHolder {
	return &SignTokenHolder{
		issuer: issuer,
		secret: secret,
		cache:  cache,
	}
}

func (j *SignTokenHolder) TokenType() string {
	return "Sign"
}

func (j *SignTokenHolder) VerifyToken(ctx context.Context, token string, uid string) error {
	key := GetSignTokenKey(uid)
	tk, err := j.cache.GetCtx(ctx, key)
	if err != nil {
		return err
	}

	if tk != token {
		return ErrTokenExpired
	}

	return nil
}

func (j *SignTokenHolder) CreateToken(ctx context.Context, uid string, expires time.Duration) (string, error) {
	ts := cast.ToString(time.Now().Unix())
	tk := crypto.Md5v(uid+ts, j.secret)

	key := GetSignTokenKey(uid)
	err := j.cache.SetexCtx(ctx, key, tk, int(expires.Seconds()))
	if err != nil {
		return "", err
	}

	return tk, nil
}

func (j *SignTokenHolder) RemoveToken(ctx context.Context, uid string) error {
	key := GetSignTokenKey(uid)
	_, err := j.cache.DelCtx(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func GetSignTokenKey(uid string) string {
	return fmt.Sprintf("blog:user:token:%v", uid)
}
