package swagger

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	httpSwagger "github.com/swaggo/http-swagger"
)

type SwaggerPlugin struct {
	docs string
}

func NewSwaggerPlugin(docs string) *SwaggerPlugin {
	return &SwaggerPlugin{
		docs: docs,
	}
}

func (p *SwaggerPlugin) Handler(prefix string) http.HandlerFunc {
	var once sync.Once
	var handler http.Handler

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, prefix)
		path = strings.TrimPrefix(path, "/")

		switch path {
		case "docs.json":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(p.docs))
		default:
			once.Do(func() {
				handler = httpSwagger.Handler(
					httpSwagger.URL(fmt.Sprintf("%s/%s", prefix, "docs.json")), //The url pointing to API definition
				)
			})
			handler.ServeHTTP(w, r)
		}
	}
}
