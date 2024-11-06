package swagger

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zeromicro/go-zero/rest"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/internal/static"
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
