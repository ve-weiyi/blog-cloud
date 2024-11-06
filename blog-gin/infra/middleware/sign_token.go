package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/glog"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
)

// 未登录token
// 未登录时，token = md5(tm,ts)
func SignToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk := c.Request.Header.Get(constant.HeaderToken)
		tm := c.Request.Header.Get(constant.HeaderTerminal)
		ts := c.Request.Header.Get(constant.HeaderTimestamp)

		//glog.Infof("api is no login required. tk:%v, tm:%v,ts:%v", tk, tm, ts)
		// 请求头缺少参数
		if tk == "" || tm == "" || ts == "" {
			c.JSON(http.StatusOK, apierr.NewApiError(codex.CodeUserNotPermission, "无效请求"))
			c.Abort()
			return
		}
		// 判断 token = md5(tm,ts)
		if tk != crypto.Md5v(tm, ts) {
			c.JSON(http.StatusOK, apierr.NewApiError(codex.CodeUserNotPermission, "无效请求"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// 登录token
// 登录时，token = md5(uid,ts)，从redis中获取token对应的用户信息
func LoginToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk := c.Request.Header.Get(constant.HeaderToken)
		uid := c.Request.Header.Get(constant.HeaderUid)

		glog.Infof("api is login required. tk:%v, uid:%v", tk, uid)
		// 请求头缺少参数
		if tk == "" || uid == "" {
			c.JSON(http.StatusOK, apierr.NewApiError(codex.CodeUserNotPermission, "无效请求"))
			c.Abort()
			return
		}
		// 判断 uid = cache.get(token)

		c.Next()
	}
}
