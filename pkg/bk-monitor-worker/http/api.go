// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	rdb "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/broker/redis"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/common"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/metrics"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/service/scheduler/daemon"
	storeRedis "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/store/redis"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/task"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/utils/errors"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/utils/timex"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/worker"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
)

type taskOptions struct {
	Retry     int    `json:"retry,omitempty"`
	Queue     string `json:"queue,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
	Deadline  string `json:"deadline,omitempty"`
	UniqueTTL int    `json:"unique_ttl,omitempty"`
}

type taskParams struct {
	Kind    string                 `binding:"required" json:"kind"`
	Payload map[string]interface{} `json:"payload"`
	Options taskOptions            `json:"options"`
}

type daemonTaskItem struct {
	task.SerializerTask
	UniId string `json:"uni_id"`
}

type removeTaskParams struct {
	TaskType  string `json:"task_type"`
	TaskUniId string `json:"task_uni_id"`
}

// CreateTask create a delay task
func CreateTask(c *gin.Context) {
	// get data
	params := new(taskParams)
	if err := BindJSON(c, params); err != nil {
		BadReqResponse(c, "parse params error: %v", err)
		return
	}
	// compose task
	payload, err := json.Marshal(params.Payload)
	if err != nil {
		ServerErrResponse(c, "json marshal error: %v", err)
		return
	}
	// 如果是周期任务，则写入到 redis 中，周期性取任务写入到队列执行
	// 如果是异步任务，则直接写入到队列，然后执行任务
	// 如果是常驻任务，则直接写入到常驻任务队列中即可
	kind := params.Kind
	if err = metrics.RegisterTaskCount(kind); err != nil {
		logger.Errorf("Report task count metric failed: %s", err)
	}
	// 组装 task
	newedTask := &task.Task{
		Kind:    kind,
		Payload: payload,
		Options: composeOption(params.Options),
	}
	// 根据类型做判断
	if strings.HasPrefix(kind, AsyncTask) {
		if err = enqueueAsyncTask(newedTask); err != nil {
			ServerErrResponse(c, "enqueue async task error, %v", err)
			return
		}
	} else if strings.HasPrefix(kind, PeriodicTask) {
		if err = pushPeriodicTaskToRedis(c, newedTask); err != nil {
			ServerErrResponse(c, "push task to redis error, %v", err)
			return
		}
	} else if strings.HasPrefix(kind, DaemonTask) {
		if err = enqueueDaemonTask(newedTask); err != nil {
			ServerErrResponse(c, "enqueue daemon task error error, %v", err)
			return
		}
	} else {
		BadReqResponse(c, "task kind: %s not support", kind)
		return
	}

	// success response
	Response(c, nil)
}

// 组装传递的 option， 如 retry、deadline 等
func composeOption(opt taskOptions) []task.Option {
	var opts []task.Option
	// 添加 option
	if opt.Retry != 0 {
		opts = append(opts, task.MaxRetry(opt.Retry))
	}
	if opt.Queue != "" {
		opts = append(opts, task.Queue(opt.Queue))
	}
	if opt.Timeout != 0 {
		timeoutOpt := IntToSecond(opt.Timeout)
		opts = append(opts, task.Timeout(timeoutOpt))
	}
	if opt.Deadline != "" {
		deadlineOpt, _ := timex.StringToTime(opt.Deadline)
		opts = append(opts, task.Deadline(deadlineOpt))
	}
	if opt.UniqueTTL != 0 {
		uniqueTTLOpt := IntToSecond(opt.UniqueTTL)
		opts = append(opts, task.Timeout(uniqueTTLOpt))
	}
	return opts
}

// 写入任务队列
func enqueueAsyncTask(t *task.Task) error {
	// new client
	client := worker.GetClient()
	defer client.Close()

	// 入队列
	if _, err := client.Enqueue(t); err != nil {
		return errors.New(fmt.Sprintf("enqueue task error, %v", err))
	}

	return nil
}

// 推送任务到 redis 中
func pushPeriodicTaskToRedis(c *gin.Context, t *task.Task) error {
	r := storeRedis.GetInstance(c)

	// expiration set zero，means the key has no expiration time
	if err := r.HSet(storeRedis.StoragePeriodicTaskKey, t.Kind, string(t.Payload)); err != nil {
		return err
	}

	// public msg
	if err := r.Publish(storeRedis.StoragePeriodicTaskChannelKey, t.Kind); err != nil {
		return err
	}

	return nil
}

// 推送任务到 task队列中
func enqueueDaemonTask(t *task.Task) error {
	client := rdb.GetRDB().Client()

	serializerTask, err := task.NewSerializerTask(*t)
	if err != nil {
		return err
	}
	data, err := json.Marshal(serializerTask)
	if err != nil {
		return err
	}

	return client.SAdd(context.Background(), common.DaemonTaskKey(), data).Err()
}

func RemoveTask(c *gin.Context) {
	params := new(removeTaskParams)
	if err := BindJSON(c, params); err != nil {
		BadReqResponse(c, "parse params error: %v", err)
		return
	}

	switch params.TaskType {
	case DaemonTask:
		client := rdb.GetRDB().Client()
		tasks, err := client.SMembers(context.Background(), common.DaemonTaskKey()).Result()
		if err != nil {
			ServerErrResponse(c, fmt.Sprintf("failed to list task by key: %s.", common.DaemonTaskKey()), err)
			return
		}
		for _, i := range tasks {
			var item task.SerializerTask
			if err = json.Unmarshal([]byte(i), &item); err != nil {
				ServerErrResponse(c, fmt.Sprintf("failed to parse key: %v to Task on value: %s", common.DaemonTaskKey(), i), err)
				return
			}
			taskUniId := daemon.ComputeTaskUniId(item)
			if taskUniId == params.TaskUniId {
				client.SRem(context.Background(), common.DaemonTaskKey(), i)
				Response(c, &gin.H{"data": taskUniId})
				return
			}
		}
		ServerErrResponse(c, fmt.Sprintf("failed to remove TaskUniId: %s, not found in key: %s", params.TaskUniId, common.DaemonTaskKey()))
		return
	default:
		ServerErrResponse(c, fmt.Sprintf("Task remove not support type: %s", params.TaskType))
		return
	}
}

func ListTask(c *gin.Context) {
	taskType := c.DefaultQuery("task_type", "empty")

	switch taskType {
	case DaemonTask:
		client := rdb.GetRDB().Client()
		tasks, err := client.SMembers(context.Background(), common.DaemonTaskKey()).Result()
		if err != nil {
			ServerErrResponse(c, fmt.Sprintf("failed to list task by key: %s.", common.DaemonTaskKey()), err)
			return
		}
		var res []daemonTaskItem
		for _, i := range tasks {
			var item task.SerializerTask
			if err = json.Unmarshal([]byte(i), &item); err != nil {
				ServerErrResponse(c, fmt.Sprintf("failed to parse key: %v to Task on value: %s", common.DaemonTaskKey(), i), err)
				return
			}
			res = append(res, daemonTaskItem{SerializerTask: item, UniId: daemon.ComputeTaskUniId(item)})
		}
		Response(c, &gin.H{"data": res})
	default:
		ServerErrResponse(c, fmt.Sprintf("Task list not support type: %s", taskType))
		return
	}
}
