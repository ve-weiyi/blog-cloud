package plugins

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/docs"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/http/staticfile"
	"github.com/ve-weiyi/vkit/plugins/music"
	"github.com/ve-weiyi/vkit/plugins/swagger"
)

func RegisterPluginHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 注册swagger服务
	var swaggerPrefix = "/api/v1/swagger"
	server.AddRoutes(staticfile.PrefixRoutes(swaggerPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		swagger.NewSwaggerPlugin(docs.Docs).Handler(swaggerPrefix).ServeHTTP(w, r)
	}))

	// 注册music服务。外部站点代理必须限流，避免被放大利用。
	// 插件在启动时构造一次：缓存、singleflight 与 HTTP 连接池都挂在实例上，
	// 放在请求里构造会让它们全部失效。
	var musicPrefix = "/api/v1/music"
	musicHandler := music.NewMusicPlugin().Handler(musicPrefix)
	musicRoutes := staticfile.PrefixRoutes(musicPrefix, http.MethodGet, musicHandler)
	server.AddRoutes(rest.WithMiddlewares([]rest.Middleware{serverCtx.MusicRateLimit}, musicRoutes...))
}
