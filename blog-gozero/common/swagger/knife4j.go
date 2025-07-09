package swagger

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/ve-blog-golang/kit/knife4j"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/static"
)

// prefix: "/api/v1/swagger/"
func RegisterKnife4jSwagHandler(server *rest.Server, prefix string, docs []byte) {

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s%s", prefix, "swagger-resources"),
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[
			  {
				"name": "2.X版本",
				"url": "/v2/api-docs",
				"swaggerVersion": "2.0",
				"location": "v2"
			  }
			]`))
		},
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s%s", prefix, "v2/api-docs"),
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(docs)
		},
	})

	server.AddRoutes(static.PrefixRoutes(prefix, func(w http.ResponseWriter, r *http.Request) {
		knife4j.Knife4jSwagHandler(prefix).ServeHTTP(w, r)
	}))
}
