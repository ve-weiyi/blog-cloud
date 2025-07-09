package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/permissionx"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/responsex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
)

type PermissionMiddleware struct {
	holder permissionx.PermissionHolder
}

func NewPermissionMiddleware(holder permissionx.PermissionHolder) *PermissionMiddleware {
	return &PermissionMiddleware{
		holder: holder,
	}
}

// 权限拦截
func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("PermissionMiddleware Handle path: %v", r.URL.Path)
		var uid string
		uid = r.Header.Get(restx.HeaderUid)
		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, "request header field 'uid' is missing"))
			return
		}

		// 验证用户是否有权限访问资源
		err := m.holder.Enforce(uid, r.URL.Path, r.Method)
		if err != nil {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeUserNotPermission, err.Error()))
			return
		}

		// 调用下一层的处理
		next.ServeHTTP(w, r)
	}
}
