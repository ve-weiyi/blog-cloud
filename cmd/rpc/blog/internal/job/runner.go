// Package job 管理本进程的后台任务生命周期。
package job

import (
	"context"
	"sync"
)

// Task 是一个受 Runner 管理的后台任务。
// Start 必须立即返回；Stop 停止任务并等待在途工作退出。
type Task struct {
	Name  string
	Start func()
	Stop  func(ctx context.Context)
}

// Runner 按注册顺序启动任务，停止时逆序。
type Runner struct {
	mu      sync.Mutex
	tasks   []Task
	started bool
}

// New 创建任务容器。
func New() *Runner {
	return &Runner{}
}

// Add 注册任务；应在 Start 之前调用。
func (r *Runner) Add(tasks ...Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, task := range tasks {
		if task.Start == nil || task.Stop == nil {
			panic("job: task " + task.Name + " must define start and stop")
		}
		r.tasks = append(r.tasks, task)
	}
}

// Start 启动全部任务，重复调用无副作用。
func (r *Runner) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.started {
		return
	}
	r.started = true
	for _, task := range r.tasks {
		task.Start()
	}
}

// Stop 逆序停止全部任务，重复调用无副作用。nil 安全。
func (r *Runner) Stop(ctx context.Context) {
	if r == nil {
		return
	}
	r.mu.Lock()
	tasks := r.tasks
	r.tasks = nil
	r.started = false
	r.mu.Unlock()

	for i := len(tasks) - 1; i >= 0; i-- {
		tasks[i].Stop(ctx)
	}
}
