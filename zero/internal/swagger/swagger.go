package swagger

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/ve-blog-golang/kit/knife4j"

	"github.com/ve-weiyi/ve-blog-golang/zero/internal/static"
)

func RegisterHttpSwagHandler(server *rest.Server, prefix, swaggerFile string) {
	f, err := os.ReadFile(swaggerFile)
	if err != nil {
		log.Println(err)
		return
	}

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s%s", prefix, "docs/blog.json"),
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(f)
		},
	})

	server.AddRoutes(static.PrefixRoutes(prefix, func(w http.ResponseWriter, r *http.Request) {
		httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%s%s", prefix, "docs/blog.json")), //The url pointing to API definition
		).ServeHTTP(w, r)
	}))
}

// prefix: "/api/v1/swagger/"
func RegisterKnife4jSwagHandler(server *rest.Server, prefix, swaggerFile string) {
	f, err := os.ReadFile(swaggerFile)
	if err != nil {
		log.Fatal(err)
		return
	}

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
			w.Write(f)
		},
	})

	server.AddRoutes(static.PrefixRoutes(prefix, func(w http.ResponseWriter, r *http.Request) {
		knife4j.Knife4jSwagHandler(prefix).ServeHTTP(w, r)
	}))
}
