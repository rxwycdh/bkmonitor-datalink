package daemon

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	rdb "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/broker/redis"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/common"
	apmTasks "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/service"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/task"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"math/rand"
)

type Runnable interface {
	start()
}

type Watcher interface {
	Runnable

	handleAddTask(watchTaskMark)
	handleDeleteTask(watchTaskMark)

	handleAddWorker(watchWorkerMark)
	handleDeleteWorker(string, watchWorkerMark)
}

type Numerator interface {
	Runnable
}

type Operator interface {
	Start(stopParentContext context.Context, errorReceiveChan chan<- error, payload []byte)
}

type OperatorDefine struct {
	handler     func() Operator
	initialFunc func(ctx context.Context) (Operator, error)
}

var taskDefine = map[string]OperatorDefine{
	"daemon:apm:pre_calculate": {initialFunc: func(ctx context.Context) (Operator, error) {
		op, err := apmTasks.Initial(ctx)
		if err != nil {
			return nil, err
		}
		go op.Run()
		return op, err
	}},
}

type TaskScheduler struct {
	ctx context.Context

	watcher   Watcher
	numerator Numerator
}

func (d *TaskScheduler) Run() {
	d.watcher.start()
	d.numerator.start()

	for {
		select {
		case <-d.ctx.Done():
			logger.Info("Scheduler received the termination signal, stopped.")
			return
		}
	}
}

func NewDaemonTaskScheduler(ctx context.Context) *TaskScheduler {
	watcher := NewDefaultWatcher(ctx)
	numerator := NewDefaultNumerator(ctx)

	return &TaskScheduler{
		ctx:       ctx,
		watcher:   watcher,
		numerator: numerator,
	}
}

func ComputeTaskUniId(task task.Task) string {
	// 统一名称+参数视为同一个常驻任务
	return fmt.Sprintf("%s-%s", task.Kind, hex.EncodeToString(task.Payload))
}

func computeWorker(t task.Task) (service.WorkerInfo, error) {
	ctx := context.Background()

	redisClient := rdb.GetRDB().Client()
	// 获取task是否有指定队列
	var queueOpt task.Option
	for _, o := range t.Options {
		if o.Type() == task.QueueOpt {
			queueOpt = o
			break
		}
	}
	var queue string
	if queueOpt != nil {
		queue = queueOpt.Value().(string)
	} else {
		queue = common.DefaultQueueName
	}

	var res service.WorkerInfo
	// 获取此队列下的所有worker
	queueWorkerPrefix := fmt.Sprintf("%s*", common.WorkerKeyQueuePrefix(queue))
	keys, err := redisClient.Keys(ctx, queueWorkerPrefix).Result()
	if err != nil {
		return res, fmt.Errorf("failed to obtain the workers with the prefix: %s from redis. Task: %s will not be attempted to schedule until the next numerator check", queueWorkerPrefix, t.Kind)
	}
	if len(keys) == 0 {
		return res, fmt.Errorf("the list of workers with prefix: %s from redis is empty, is no worker listening to this queue: %s?. Task: %s will not be attempted to schedule until the next numerator check", queueWorkerPrefix, queue, t.Kind)
	}

	// TODO 从worker列表中选择worker进行调度 待补充更多的调度规则 目前暂时使用随机选择
	data, _ := redisClient.Get(ctx, keys[rand.Intn(len(keys))]).Bytes()
	if err = json.Unmarshal(data, &res); err != nil {
		return res, fmt.Errorf("parse workerInfo failed. error: %s", err)
	}
	return res, nil
}
