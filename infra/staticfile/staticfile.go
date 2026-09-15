// Package staticfile go-zero 静态文件服务辅助。
// go-zero 路由不支持 * 通配，用命名参数 :a/:b/:c 生成多级路由模拟任意深度目录。
package staticfile

import (
	"net/http"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

// dirLevel 命名参数占位（匹配顺序：:a / :a/:b / :a/:b/:c / ...）
var dirLevel = []string{":a", ":b", ":c", ":d", ":e"}

// maxLevel 支持的最大目录层级（图片 key 形如 batch_code/station/file，2 级，预留余量到 5）
const maxLevel = 5

// PrefixRoutes 生成 prefix 前缀的多级静态文件路由（method + handler）。
// 每个路由匹配一层更深目录，最终由 handler 统一处理任意深度请求。
func PrefixRoutes(prefix, method string, handler http.HandlerFunc) []rest.Route {
	routes := make([]rest.Route, 0, maxLevel)
	for i := 1; i <= maxLevel; i++ {
		routes = append(routes, rest.Route{
			Method:  method,
			Path:    path.Join(prefix, strings.Join(dirLevel[:i], "/")),
			Handler: handler,
		})
	}
	return routes
}

// FileServerRoutes 注册 prefix 前缀的多级静态文件服务路由，文件根目录为 dir。
// 请求经 http.StripPrefix 剥离 prefix 后交给 http.FileServer 从 dir 读取。
func FileServerRoutes(prefix, dir string) []rest.Route {
	handler := http.StripPrefix(prefix, http.FileServer(http.Dir(dir))).ServeHTTP
	return PrefixRoutes(prefix, http.MethodGet, handler)
}
