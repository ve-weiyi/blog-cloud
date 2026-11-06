package job

import (
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
)

// Guard 执行一轮后台工作，捕获 panic 并记录调用栈。
//
// 后台循环通常是长驻 goroutine：单轮里的空指针或越界会直接终止进程，
// 把一次可重试的批处理失败放大成整服务不可用。Guarded 后失败的只有这一轮，
// 下一轮仍会继续；真正的错误处理仍由循环内的 err 分支负责。
func Guard(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("[job] %s panic: %v\n%s", name, r, debug.Stack())
		}
	}()
	fn()
}
