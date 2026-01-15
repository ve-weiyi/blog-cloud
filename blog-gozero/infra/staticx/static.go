package staticx

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

// go-zero实现路径通配方式
func PrefixRoutes(prefix string, method string, handler http.HandlerFunc) []rest.Route {
	var rt []rest.Route
	dirLevel := []string{":a", ":b", ":c", ":d", ":e"}
	for i := 1; i <= len(dirLevel); i++ {
		p := path.Join(prefix, strings.Join(dirLevel[:i], "/"))
		rt = append(rt, rest.Route{
			Method:  method,
			Path:    p,
			Handler: handler,
		})
	}
	return rt
}

// 静态文件处理
// staticFileHandler(server)
func staticFileHandler(server *rest.Server) {
	//这里注册
	pattern := "web"
	dir := "web/assets/"

	rd, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// 不能匹配多级目录
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/static/:file",
		// prefix 前缀会被删除，最后匹配到的路径会被传递给 http.FileServer 目录下
		Handler: http.StripPrefix("/static/", http.FileServer(http.Dir("runtime/resource"))).ServeHTTP,
	})

	//添加进路由最后生成 /asset
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/index.html",
				Handler: DirHandler("index.html", pattern),
			},
			{
				Method:  http.MethodGet,
				Path:    "/",
				Handler: DirHandler("/", pattern),
			},
			{
				Method:  http.MethodGet,
				Path:    "/favicon.ico",
				Handler: DirHandler("/favicon.ico", pattern),
			},
		})
	for _, f := range rd {
		filename := f.Name()
		//最后生成 /asset
		server.AddRoute(
			rest.Route{
				Method:  http.MethodGet,
				Path:    "/assets/" + filename,
				Handler: DirHandler("/assets/", dir),
			})
	}

}

// 不能匹配多级目录
// prefix 前缀会被删除，最后匹配到的路径会被传递给 http.FileServer 目录下
func DirHandler(prefix, dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
		handler.ServeHTTP(w, req)
	}
}
