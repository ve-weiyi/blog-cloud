package plugins

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/blog-cloud/api/admin/docs"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/vkit/plugins/knife4j"

	"github.com/ve-weiyi/blog-cloud/infra/staticfile"
)

func RegisterPluginHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 注册knife4j服务
	knife4jPrefix := "/admin-api/v1/swagger"
	server.AddRoutes(staticfile.PrefixRoutes(knife4jPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		knife4j.NewKnife4jPlugin(docs.Docs).Handler(knife4jPrefix).ServeHTTP(w, r)
	}))
}
