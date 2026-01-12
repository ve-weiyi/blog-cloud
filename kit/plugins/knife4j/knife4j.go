package knife4j

import (
	"embed"
	_ "embed"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed static/doc.html
var index string

//go:embed static/*
var content embed.FS

var SwaggerResourcesText = `
[
  {
	"name": "2.X版本",
	"url": "/v2/api-docs",
	"swaggerVersion": "2.0",
	"location": "v2"
  }
]`

type Knife4jPlugin struct {
	docs string
}

func NewKnife4jPlugin(docs string) *Knife4jPlugin {
	return &Knife4jPlugin{
		docs: docs,
	}
}

func (p *Knife4jPlugin) Handler(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, prefix)
		path = strings.TrimPrefix(path, "/")

		switch filepath.Ext(path) {
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}

		switch path {
		case "swagger-resources":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(SwaggerResourcesText))
		case "v2/api-docs":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(p.docs))
		case "index.html":
			w.Write([]byte(index))
		case "doc.html":
			w.Write([]byte(index))
		case "":
			http.Redirect(w, r, prefix+"index.html", http.StatusMovedPermanently)
		default:
			data, err := content.ReadFile("static/" + path)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			w.Write(data)
		}
	}
}
