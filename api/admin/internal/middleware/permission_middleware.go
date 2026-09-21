package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/middleware/permissionx"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

type PermissionMiddleware struct {
	holder permissionx.Enforcer
}

func NewPermissionMiddleware(holder permissionx.Enforcer) *PermissionMiddleware {
	return &PermissionMiddleware{
		holder: holder,
	}
}

// 权限拦截
func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var uid string
		uid = r.Header.Get(bizheader.HeaderUid)
		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "request header field 'uid' is missing"))
			return
		}

		// 验证用户是否有权限访问资源
		ok, err := m.holder.Enforce(uid, r.URL.Path, r.Method)
		if err != nil {
			logx.Infof("[Perm] user=%s, path=%s, method=%s, allowed=%v, err=%v", uid, r.URL.Path, r.Method, ok, err)
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeNoPermission, err.Error()))
			return
		}

		if ok != true {
			logx.Infof("[Perm] user=%s, path=%s, method=%s, allowed=%v", uid, r.URL.Path, r.Method, ok)
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeNoPermission, "user does not have permission to access this resource"))
			return
		}

		// 调用下一层的处理
		next.ServeHTTP(w, r)
	}
}
