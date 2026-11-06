package job

import (
	"context"
	"reflect"
	"testing"
)

func TestRunnerStartsInOrderAndStopsInReverse(t *testing.T) {
	var events []string
	newTask := func(name string) Task {
		return Task{
			Name:  name,
			Start: func() { events = append(events, "start:"+name) },
			Stop:  func(context.Context) { events = append(events, "stop:"+name) },
		}
	}

	runner := New()
	runner.Add(newTask("first"), newTask("second"))
	runner.Start()
	runner.Start() // 重复 Start 不得重复执行
	runner.Stop(context.Background())
	runner.Stop(context.Background()) // 重复 Stop 不得重复停止

	want := []string{"start:first", "start:second", "stop:second", "stop:first"}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("events = %v, want %v", events, want)
	}
}
