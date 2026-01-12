package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jwtx"
)

// 未登录token
// 未登录时，token = md5(tm,ts)
func TerminalToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk := c.Request.Header.Get(bizheader.HeaderXTerminalToken)
		tm := c.Request.Header.Get(bizheader.HeaderXTerminalId)
		ts := c.Request.Header.Get(bizheader.HeaderTimestamp)

		//logz.Infof("api is no login required. tk:%v, tm:%v,ts:%v", tk, tm, ts)
		// 请求头缺少参数
		if tk == "" {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderXAdminToken)))
			c.Abort()
			return
		}

		if tm == "" {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderXTerminalId)))
			c.Abort()
			return
		}

		if ts == "" {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderTimestamp)))
			c.Abort()
			return
		}

		// 判断 token = md5(tm,ts)
		if tk != cryptox.Md5v(tm, ts) {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeUserLoginExpired, "无效请求"))
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
		tk := c.Request.Header.Get(bizheader.HeaderToken)
		uid := c.Request.Header.Get(bizheader.HeaderUid)

		// 请求头缺少参数
		// token为空或者uid为空
		if tk == "" {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderToken)))
			c.Abort()
			return
		}

		if uid == "" {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderUid)))
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
		token := c.Request.Header.Get(bizheader.HeaderAuthorization)
		uid := c.Request.Header.Get(bizheader.HeaderUid)

		// token为空或者uid为空
		if token == "" {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderAuthorization)))
			c.Abort()
			return
		}

		if uid == "" {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderUid)))
			c.Abort()
			return
		}

		// 解析token
		tok, err := tk.ParseToken(token)
		if err != nil {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeUserLoginExpired, "token parse error"))
			c.Abort()
			return
		}

		// token不合法
		if !tok.Valid {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeUserLoginExpired, "token is invalid"))
			c.Abort()
			return
		}

		// 获取claims
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeUserLoginExpired, "token claims is not jwt.MapClaims"))
			c.Abort()
			return
		}

		// uid不一致
		if uid != cast.ToString(claims["uid"]) {
			response.ResponseError(c, bizerr.NewBizError(bizcode.CodeUserLoginExpired, "token uid is not equal"))
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
