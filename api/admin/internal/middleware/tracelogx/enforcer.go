package tracelogx

type Enforcer interface {
	// 加载规则
	LoadPolicy() error
	// 重新加载规则
	ReloadPolicy() error
	// 是否记录日志
	IsTraceLog(url string, method string) (bool, error)
}
