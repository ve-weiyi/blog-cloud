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

func Knife4jSwagHandler(prefix string) http.HandlerFunc {
	//handler := &webdav.Handler{
	//	Prefix:     prefix,
	//	FileSystem: webdav.Dir("static"),
	//	LockSystem: webdav.NewMemLS(),
	//}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		path := r.RequestURI

		path = strings.Replace(path, prefix, "", 1)
		if strings.Contains(path, "?") {
			path = path[:strings.Index(path, "?")]
		}
		//log.Println("path", path)
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
			//handler.ServeHTTP(w, r)
		}
	}
}
