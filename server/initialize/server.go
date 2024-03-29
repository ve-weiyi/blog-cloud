package initialize

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
)

type server interface {
	ListenAndServe() (err error)
	RegisterOnShutdown(f func())
	// 优雅的停机重启
	Shutdown(ctx context.Context) error
}

func RunWindowsServer() {

	routers := Routers()

	address := fmt.Sprintf(":%d", global.CONFIG.System.Port)
	var s server
	s = &http.Server{
		Addr:           address,
		Handler:        routers,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.LOG.Infof("server run success on %v", address)

	fmt.Printf(`
	欢迎使用 ve-blog-golang
	当前版本: v1.0.0
	微信号：wy791422171 QQ：791422171
	默认接口文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端运行地址:http://127.0.0.1:9090
`, address)
	global.LOG.Error(s.ListenAndServe().Error())
}
