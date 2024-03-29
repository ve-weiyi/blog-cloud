package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr"
)

// JwtToken jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		uid := c.Request.Header.Get("uid")
		claims, err := global.JWT.VerifyToken(token, uid)
		//必须要token才能过
		if err != nil {
			global.LOG.Error("token-->", token)
			global.LOG.Error("-->", err)
			// 有错误，直接返回给前端错误，前端直接报错500
			//c.AbortWithStatus(http.StatusInternalServerError)
			// 该方式前端不报错
			c.JSON(http.StatusOK, response.Response{
				Code:    apierr.ErrorUserUnLogin.Code(),
				Message: err.Error(),
				Data:    nil,
			})
			c.Abort()
			return
		}

		global.LOG.Infof("user login-->%v", claims.Ext.Username)
		//global.LOG.JsonIndent(claims)

		c.Set("token", token)
		c.Set("uid", uid)
		c.Set("domain", claims.Issuer)
		c.Set("username", claims.Ext.Username)
		c.Set("login_type", claims.Ext.LoginType)
		c.Next()

	}
}
