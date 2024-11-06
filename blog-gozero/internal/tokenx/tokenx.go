package tokenx

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/jtoken"
)

type JwtTokenHolder struct {
	issuer string

	token *jtoken.JwtInstance
	cache *redis.Redis
}

func NewJwtTokenHolder(issuer string, secret string, cache *redis.Redis) *JwtTokenHolder {
	return &JwtTokenHolder{
		issuer: issuer,
		token:  jtoken.NewJwtInstance([]byte(secret)),
		cache:  cache,
	}
}

func (j *JwtTokenHolder) CreateToken(ctx context.Context, uid string, perms string, expiresIn int64) (string, error) {
	issuer := "blog"
	now := time.Now().Unix()

	opts := []jtoken.Option{
		jtoken.WithIssuer(issuer),
		jtoken.WithIssuedAt(now),
		jtoken.WithExpiresAt(expiresIn),
		jtoken.WithClaimExt("uid", uid),
		jtoken.WithClaimExt("roles", perms),
	}

	tk, err := j.token.CreateToken(opts...)
	if err != nil {
		return "", err
	}

	return tk, nil
}

func (j *JwtTokenHolder) VerifyToken(ctx context.Context, token string, uid string) (jwt.MapClaims, error) {
	//token为空或者uid为空
	if token == "" || uid == "" {
		return nil, fmt.Errorf("token or uid is empty")
	}

	// 解析token
	tok, err := j.token.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	// token不合法
	if !tok.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// 获取claims
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token claims is not jwt.MapClaims")
	}

	// uid不一致
	if uid != cast.ToString(claims[constant.HeaderUid]) {
		return nil, fmt.Errorf("token cannot use by uid")
	}

	//token验证成功,但用户在别处登录或退出登录
	if j.IsLogout(ctx, cast.ToString(claims[constant.HeaderUid]), cast.ToInt64(claims[jtoken.JwtIssueAt])) {
		return nil, fmt.Errorf("user already logout or login in other place")
	}

	return claims, nil
}

// 设置退出登录
func (j *JwtTokenHolder) SetLogout(ctx context.Context, uid string, loginOut int64) error {
	redisKey := GetUserLogoutKey(uid)
	return j.cache.SetexCtx(ctx, redisKey, fmt.Sprintf("%d", loginOut), 7*24*60*60)
}

// 已退出登录
func (j *JwtTokenHolder) IsLogout(ctx context.Context, uid string, loginAt int64) bool {
	redisKey := GetUserLogoutKey(uid)
	at, err := j.cache.GetCtx(ctx, redisKey)
	if err != nil {
		return false
	}

	logoutAt := cast.ToInt64(at)
	if loginAt < logoutAt {
		logx.Infof("loginAt=%d, at.LogoutAt=%d", loginAt, logoutAt)
		return true
	}

	return false
}

func GetUserLogoutKey(uid string) string {
	return fmt.Sprintf("user:logout:%d", uid)
}
