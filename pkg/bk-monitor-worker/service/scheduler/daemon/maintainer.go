package daemon

import (
	"context"
	"encoding/json"
	rdb "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/broker/redis"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/common"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/config"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type runningBinding struct {
	errorReceiveChan chan error
	TaskBinding
	retryCount    int
	startTime     time.Time
	lastRetryTime time.Time
	nextRetryTime time.Time
	baseCtx       context.Context
	baseCtxCancel context.CancelFunc
}

type RunMaintainerOptions struct {
	checkInterval         time.Duration
	RetryTolerateCount    int
	RetryTolerateInterval time.Duration
	RetryIntolerantFactor int
}

type RunMaintainer struct {
	ctx    context.Context
	config RunMaintainerOptions

	listenWorkerId string
	listenTaskKey  string

	redisClient           redis.UniversalClient
	methodOperatorMapping map[string]Operator
	runningInstance       *sync.Map
}

func (r *RunMaintainer) Run() {
	logger.Infof("\nDaemonTask maintainer started. \n - \t workerId: %s \n - \t listen: %s \n - \t queue: %s \n - \t interval: %s", r.listenWorkerId, r.listenWorkerId, r.listenTaskKey, r.config.checkInterval)
	taskMarkMapping := make(map[string]bool)
	for {
		select {
		case <-r.ctx.Done():
			logger.Infof("DaemonTask maintainer stopped.")
		default:
			currentTask := make(map[string]bool)
			taskHash, err := r.redisClient.HGetAll(r.ctx, r.listenTaskKey).Result()
			if err != nil {
				logger.Errorf("DaemonTask maintainer(%s) check %s change failed. error: %s", r.listenWorkerId, r.listenTaskKey, err)
				goto cont
			}

			for taskUniId, taskStr := range taskHash {
				var taskBinding TaskBinding
				if err = json.Unmarshal([]byte(taskStr), &taskBinding); err != nil {
					logger.Errorf("failed to parse value to TaskBinding on key: %s field: %s. error: %s", r.listenTaskKey, taskUniId, err)
					continue
				}

				currentTask[taskBinding.UniId] = true
				if !taskMarkMapping[taskBinding.UniId] {
					r.handleAddTaskBinding(taskBinding)
				}
			}

			for t := range taskMarkMapping {
				if !currentTask[t] {
					r.handleDeleteTaskBinding(t)
				}
			}

			taskMarkMapping = currentTask
		cont:
			time.Sleep(r.config.checkInterval)
		}
	}
}

func (r *RunMaintainer) handleAddTaskBinding(taskBinding TaskBinding) {
	define, exist := r.methodOperatorMapping[taskBinding.Kind]
	if !exist {
		logger.Errorf("Failed to run method: %s which not exist in task defines", taskBinding.Kind)
		return
	}

	errorReceiveChan := make(chan error, 1)
	baseCtx, baseCtxCancel := context.WithCancel(r.ctx)
	go define.Start(baseCtx, errorReceiveChan, taskBinding.Payload)
	now := time.Now()
	r.runningInstance.Store(taskBinding.UniId, &runningBinding{
		baseCtx:          baseCtx,
		baseCtxCancel:    baseCtxCancel,
		errorReceiveChan: errorReceiveChan,
		TaskBinding:      taskBinding,
		retryCount:       0,
		startTime:        now,
	})
	go r.listenRunningState(taskBinding.UniId, errorReceiveChan, baseCtx)
	logger.Infof("Binding(%s <------> %s) is discovered, task is started", taskBinding.UniId, r.listenWorkerId)
}

func (r *RunMaintainer) listenRunningState(taskUniId string, errorReceiveChan chan error, baseCtx context.Context) {
	retryTicker := &time.Ticker{}

	for {
		select {
		case receiveErr := <-errorReceiveChan:
			v, _ := r.runningInstance.Load(taskUniId)
			rB := v.(*runningBinding)
			rB.retryCount++
			rB.lastRetryTime = time.Now()

			var nextRetryTime time.Duration
			if rB.retryCount < r.config.RetryTolerateCount {
				nextRetryTime = r.config.RetryTolerateInterval
			} else {
				nextRetryTime = r.config.RetryTolerateInterval * time.Duration(1<<(rB.retryCount-r.config.RetryIntolerantFactor))
			}
			rB.nextRetryTime = time.Now().Add(nextRetryTime)
			r.runningInstance.Store(taskUniId, rB)
			logger.Warnf("[FAILED RETRY] ERROR: %s. Task: %s, %d retry failed. The retry time of the next attempt is: %s, (%.2f seconds later)", receiveErr, taskUniId, rB.retryCount, rB.nextRetryTime, nextRetryTime.Seconds())
			retryTicker = time.NewTicker(nextRetryTime)
		case <-retryTicker.C:
			v, _ := r.runningInstance.Load(taskUniId)
			rB := v.(*runningBinding)
			define, _ := r.methodOperatorMapping[rB.TaskBinding.Kind]
			go define.Start(rB.baseCtx, errorReceiveChan, rB.Task.Payload)
			logger.Infof("[FAILED RETRY] Task: %s retry performed", taskUniId)
			retryTicker = &time.Ticker{}
		case <-baseCtx.Done():
			logger.Infof("[RetryListen] stopped.")
			return
		}
	}
}

func (r *RunMaintainer) handleDeleteTaskBinding(taskUniId string) {
	ins, exist := r.runningInstance.Load(taskUniId)
	if !exist {
		logger.Errorf("Attempt to delete a task: %s that is not executing", taskUniId)
		return
	}
	ins.(*runningBinding).baseCtxCancel()
	r.runningInstance.Delete(taskUniId)
	logger.Infof("Binding runInstance removed. taskUniId: %s", taskUniId)
}

func NewDaemonTaskRunMaintainer(ctx context.Context, workerId string) *RunMaintainer {

	operatorMapping := make(map[string]Operator, len(taskDefine))

	for taskKind, define := range taskDefine {
		op, err := define.initialFunc(ctx)
		if err != nil {
			logger.Errorf("[!WARNING!] Task: %s implementation initialization failed, this task type will not be executed! error: %s", taskKind, err)
			continue
		}
		operatorMapping[taskKind] = op
	}

	options := RunMaintainerOptions{
		checkInterval:         time.Duration(config.WorkerDaemonTaskMaintainerInterval) * time.Second,
		RetryTolerateCount:    config.WorkerDaemonTaskRetryTolerateCount,
		RetryTolerateInterval: time.Duration(config.WorkerDaemonTaskRetryTolerateInterval) * time.Second,
		RetryIntolerantFactor: config.WorkerDaemonTaskRetryIntolerantFactor,
	}

	return &RunMaintainer{
		ctx:                   ctx,
		config:                options,
		listenWorkerId:        workerId,
		listenTaskKey:         common.DaemonBindingWorker(workerId),
		redisClient:           rdb.GetRDB().Client(),
		methodOperatorMapping: operatorMapping,
		runningInstance:       &sync.Map{},
	}
}
