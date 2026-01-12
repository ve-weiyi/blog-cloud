package plugins

import (
	"net/http"
	"time"

	"github.com/openai/openai-go/v3/option"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/static"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/docs"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/ai"
	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/music"
	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/swagger"
)

func RegisterPluginHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 注册swagger服务
	var swaggerPrefix = "/blog-api/v1/swagger"
	server.AddRoutes(static.PrefixRoutes(swaggerPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		swagger.NewSwaggerPlugin(docs.Docs).Handler(swaggerPrefix).ServeHTTP(w, r)
	}))

	// 注册music服务
	var musicPrefix = "/blog-api/v1/music"
	server.AddRoutes(static.PrefixRoutes(musicPrefix, http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		music.NewMusicPlugin().Handler(musicPrefix).ServeHTTP(w, r)
	}))

	// 注册ai服务
	var aiPrefix = "/blog-api/v1/ai"
	server.AddRoutes(static.PrefixRoutes(aiPrefix, http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		ai.NewAiPlugin(
			option.WithAPIKey(serverCtx.Config.AiProxyConf.ApiKey),
			option.WithBaseURL(serverCtx.Config.AiProxyConf.ApiHost),
		).Handler(aiPrefix).ServeHTTP(w, r)
	}),
		// 超时时间设置长一点，否则会截断ai的回答
		rest.WithTimeout(30*time.Second),
	)
}
