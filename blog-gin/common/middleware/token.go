package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/jwtx"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
)

// 未登录token
// 未登录时，token = md5(tm,ts)
func TerminalToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk := c.Request.Header.Get(restx.HeaderXTerminalToken)
		tm := c.Request.Header.Get(restx.HeaderXTerminalId)
		ts := c.Request.Header.Get(restx.HeaderTimestamp)

		//logz.Infof("api is no login required. tk:%v, tm:%v,ts:%v", tk, tm, ts)
		// 请求头缺少参数
		if tk == "" {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderXAdminToken)))
			c.Abort()
			return
		}

		if tm == "" {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderXTerminalId)))
			c.Abort()
			return
		}

		if ts == "" {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderTimestamp)))
			c.Abort()
			return
		}

		// 判断 token = md5(tm,ts)
		if tk != crypto.Md5v(tm, ts) {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeUserLoginExpired, "无效请求"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// 登录token
// 登录时，token = md5(uid,ts)，从redis中获取token对应的用户信息
func UserToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk := c.Request.Header.Get(restx.HeaderToken)
		uid := c.Request.Header.Get(restx.HeaderUid)

		// 请求头缺少参数
		// token为空或者uid为空
		if tk == "" {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderToken)))
			c.Abort()
			return
		}

		if uid == "" {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderUid)))
			c.Abort()
			return
		}

		// 判断 uid = cache.get(token)

		c.Next()
	}
}

// 管理员token
// jwt token
func AdminToken(tk *jwtx.JwtInstance) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(restx.HeaderAuthorization)
		uid := c.Request.Header.Get(restx.HeaderUid)

		// token为空或者uid为空
		if token == "" {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderAuthorization)))
			c.Abort()
			return
		}

		if uid == "" {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderUid)))
			c.Abort()
			return
		}

		// 解析token
		tok, err := tk.ParseToken(token)
		if err != nil {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeUserLoginExpired, "token parse error"))
			c.Abort()
			return
		}

		// token不合法
		if !tok.Valid {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeUserLoginExpired, "token is invalid"))
			c.Abort()
			return
		}

		// 获取claims
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeUserLoginExpired, "token claims is not jwt.MapClaims"))
			c.Abort()
			return
		}

		// uid不一致
		if uid != cast.ToString(claims["uid"]) {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeUserLoginExpired, "token uid is not equal"))
			c.Abort()
			return
		}

		// token验证成功,但用户在别处登录或退出登录
		// if jwtService.IsBlacklist(token) {
		//
		// }

		// 写入上下文
		ctx := c.Request.Context()
		for k, v := range claims {
			switch k {
			case jwtx.JwtAudience, jwtx.JwtExpire, jwtx.JwtId, jwtx.JwtIssueAt, jwtx.JwtIssuer, jwtx.JwtNotBefore, jwtx.JwtSubject:
				// ignore the standard claims
			default:
				ctx = context.WithValue(ctx, k, v)
			}
		}

		c.Next()

	}
}
