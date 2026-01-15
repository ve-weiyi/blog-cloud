package plugins

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/staticx"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/docs"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/knife4j"
)

func RegisterPluginHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 注册knife4j服务
	var knife4jPrefix = "/admin-api/v1/swagger"
	server.AddRoutes(staticx.PrefixRoutes(knife4jPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		knife4j.NewKnife4jPlugin(docs.Docs).Handler(knife4jPrefix).ServeHTTP(w, r)
	}))
}
