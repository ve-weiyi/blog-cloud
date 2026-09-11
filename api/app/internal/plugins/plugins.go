package plugins

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/blog-cloud/api/app/docs"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/vkit/plugins/music"
	"github.com/ve-weiyi/vkit/plugins/swagger"

	"github.com/ve-weiyi/blog-cloud/infra/staticx"
)

func RegisterPluginHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 注册swagger服务
	var swaggerPrefix = "/api/v1/swagger"
	server.AddRoutes(staticx.PrefixRoutes(swaggerPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		swagger.NewSwaggerPlugin(docs.Docs).Handler(swaggerPrefix).ServeHTTP(w, r)
	}))

	// 注册music服务
	var musicPrefix = "/api/v1/music"
	server.AddRoutes(staticx.PrefixRoutes(musicPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		music.NewMusicPlugin().Handler(musicPrefix).ServeHTTP(w, r)
	}))
}
