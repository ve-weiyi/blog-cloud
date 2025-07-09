package swagger

import (
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/static"
)

func RegisterHttpSwagHandler(server *rest.Server, prefix string, docs []byte) {
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s%s", prefix, "docs.json"),
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(docs)
		},
	})

	server.AddRoutes(static.PrefixRoutes(prefix, func(w http.ResponseWriter, r *http.Request) {
		httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%s%s", prefix, "docs.json")), //The url pointing to API definition
		).ServeHTTP(w, r)
	}))
}
