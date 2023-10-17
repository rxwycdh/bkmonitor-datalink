package scheduler

import (
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/example"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/processor"
)

type Task struct {
	Handler processor.HandlerFunc
}

var (
	exampleTask = "async:test_example"

	asyncTaskDefine = map[string]Task{
		exampleTask: {
			Handler: example.HandleExampleTask,
		},
	}
)

func GetAsyncTaskMapping() map[string]Task {
	return asyncTaskDefine
}
