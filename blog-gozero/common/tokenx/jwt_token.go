package tokenx

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/jwtx"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
)

type JwtTokenHolder struct {
	issuer string

	token *jwtx.JwtInstance
	cache *redis.Redis
}

func NewJwtTokenHolder(issuer string, secret string, cache *redis.Redis) *JwtTokenHolder {
	return &JwtTokenHolder{
		issuer: issuer,
		token:  jwtx.NewJwtInstance([]byte(secret)),
		cache:  cache,
	}
}

func (j *JwtTokenHolder) TokenType() string {
	return "Bearer"
}

func (j *JwtTokenHolder) VerifyToken(ctx context.Context, token string, uid string) error {
	//token为空或者uid为空
	if token == "" || uid == "" {
		return ErrTokenEmpty
	}

	// 解析token
	tok, err := j.token.ParseToken(token)
	if err != nil {
		return err
	}

	// token不合法
	if !tok.Valid {
		return ErrTokenInvalid
	}

	// 获取claims
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("token claims is not jwt.MapClaims")
	}

	// uid不一致
	if uid != cast.ToString(claims[restx.HeaderUid]) {
		return fmt.Errorf("token cannot use by uid")
	}

	//token验证成功,但用户在别处登录或退出登录
	redisKey := GetUserLogoutKey(uid)
	at, err := j.cache.GetCtx(ctx, redisKey)
	if err != nil {
		return err
	}

	loginAt := cast.ToInt64(claims[jwtx.JwtIssueAt])
	logoutAt := cast.ToInt64(at)
	if loginAt < logoutAt {
		logx.Infof("loginAt=%d, at.LogoutAt=%d", loginAt, logoutAt)
		return ErrTokenExpired
	}

	return nil
}

func (j *JwtTokenHolder) CreateToken(ctx context.Context, uid string, expires time.Duration) (string, error) {
	issuer := "blog"
	now := time.Now().Unix()
	expiresAt := time.Now().Add(expires).Unix()

	opts := []jwtx.Option{
		jwtx.WithIssuer(issuer),
		jwtx.WithIssuedAt(now),
		jwtx.WithExpiresAt(expiresAt),
		jwtx.WithClaimExt(restx.HeaderUid, uid),
	}

	tk, err := j.token.CreateToken(opts...)
	if err != nil {
		return "", err
	}

	return tk, nil
}

func (j *JwtTokenHolder) RemoveToken(ctx context.Context, uid string) error {
	redisKey := GetSignTokenKey(uid)

	ts := cast.ToString(time.Now().Unix())
	return j.cache.SetexCtx(ctx, redisKey, ts, 7*24*60*60)
}

func GetUserLogoutKey(uid string) string {
	return fmt.Sprintf("blog:user:logout:%v", uid)
}
