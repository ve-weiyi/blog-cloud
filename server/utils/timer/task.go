package timer

import (
	"fmt"
	"testing"
	"time"
)

const (
	Start = 1
	Stop  = 2
	Pause = 3
)

type Task struct {
	id          int           //key
	interval    time.Duration //中断时间
	MaxRunCount int           //最大运行数
	callback    func()        //运行回调
}

type TaskQueue struct {
	tasks  []*Task
	stopCh chan struct{}
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks:  make([]*Task, 0),
		stopCh: make(chan struct{}),
	}
}

func (q *TaskQueue) AddTask(task *Task) {
	q.tasks = append(q.tasks, task)
}

func (q *TaskQueue) RemoveTask(id int) {
	for i, task := range q.tasks {
		if task.id == id {
			q.tasks = append(q.tasks[:i], q.tasks[i+1:]...)
			break
		}
	}
}

func (q *TaskQueue) Start() {
	go func() {
		for {
			select {
			case <-q.stopCh:
				return
			default:
				for _, task := range q.tasks {
					select {
					case <-q.stopCh:
						return
					default:
						task.callback()
						time.Sleep(task.interval)
					}
				}
			}
		}
	}()
}

func (q *TaskQueue) Pause() {
	q.stopCh <- struct{}{}
}

func (q *TaskQueue) Stop() {
	close(q.stopCh)
}

func TestTask(t *testing.T) {
	queue := NewTaskQueue()

	task1 := &Task{
		id:       1,
		interval: 1000 * time.Millisecond,
		callback: func() {
			fmt.Println("Task 1 is running...")
		},
	}

	task2 := &Task{
		id:       2,
		interval: 2000 * time.Millisecond,
		callback: func() {
			fmt.Println("Task 2 is running...")
		},
	}

	queue.AddTask(task1)
	queue.AddTask(task2)

	queue.Start()
	time.Sleep(5 * time.Second)

	queue.Pause()
	fmt.Println("Queue is paused")
	time.Sleep(5 * time.Second)

	queue.Start()
	fmt.Println("Queue is resumed")
	time.Sleep(5 * time.Second)

	queue.RemoveTask(1)
	fmt.Println("Task 1 is removed from queue")
	time.Sleep(5 * time.Second)

	queue.Stop()
	fmt.Println("Queue is stopped")
}
