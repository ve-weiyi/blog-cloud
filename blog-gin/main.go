package main

import (
	"log"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/cmd"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// @title			ve-blog-golang API
// @version			1.0.0
// @description		ve-blog-golang 博客系统 API 文档
// @termsOfService	https://github.com/ve-weiyi/ve-blog-golang/blog-gin
// @contact.name	API Support
// @contact.url		https://github.com/ve-weiyi
// @contact.email	support@swagger.io
// @license.name	MIT
// @license.url		https://opensource.org/licenses/MIT
// @host			localhost:9090
// @BasePath
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						x-token
func main() {
	cmd.Execute()
}
