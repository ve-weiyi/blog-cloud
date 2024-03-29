package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/rbac"
)

func CasbinHandler() gin.HandlerFunc {
	enforcer := rbac.NewCachedEnforcer(global.DB) // 判断策略中是否存在
	return func(c *gin.Context) {
		if global.CONFIG.System.Env != "develop" {
			//claims, err := global.JWT.ParseTokenByGin(c)
			//获取请求的PATH
			obj := c.Request.URL.Path
			// 获取请求方法
			act := c.Request.Method
			//是白名单
			if enforcer.IsWhileList(obj, act) {
				c.Next()
				return
			}
			// 请求域
			domain := c.GetString("domain")
			// 获取用户的角色
			value, ok := c.Get("roles")
			if !ok {
				global.LOG.Println("err ", value)
			}
			roles := value.([]string)

			for _, role := range roles {
				//超级权限
				if role == "super-admin" {
					c.Next()
					return
				}
				sub := role
				success, _ := enforcer.Enforce(sub, domain, obj, act)
				global.LOG.Infof("sub:%v domain:%v obj:%v act:%v", sub, domain, obj, act)
				global.LOG.Println("success ", success)
				if success {
					c.Next()
					return
				}
				//重新加载，一般在添加时才需要做
				enforcer.LoadPolicy()
			}

			c.JSON(http.StatusOK,
				response.Response{
					Code:    apierr.ErrorUserNotPermission.Code(),
					Message: "角色权限不足",
					Data:    nil,
				})
			c.Abort()
			return
		}
		c.Next()
	}
}
