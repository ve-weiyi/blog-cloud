package plugins

import "net/http"

// 服务代理插件
// 功能: 快速注册整套服务
type PluginProvider interface {
	Handler(prefix string) http.HandlerFunc
}
