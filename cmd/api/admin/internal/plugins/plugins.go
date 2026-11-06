package plugins

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/docs"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/http/staticfile"
	"github.com/ve-weiyi/vkit/plugins/knife4j"
)

func RegisterPluginHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 生产环境不暴露接口文档
	if serverCtx.Config.Mode == service.ProMode {
		return
	}

	// 注册knife4j服务
	knife4jPrefix := "/admin-api/v1/swagger"
	server.AddRoutes(staticfile.PrefixRoutes(knife4jPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		knife4j.NewKnife4jPlugin(docs.Docs).Handler(knife4jPrefix).ServeHTTP(w, r)
	}))
}
