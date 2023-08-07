// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package service

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/http"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/metrics"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/worker"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
)

const (
	ginModePath           = "service.gin_mode"
	workerConcurrencyPath = "worker.concurrency"
)

func init() {
	viper.SetDefault(ginModePath, "release")
	// 默认为0，通过动态获取逻辑 cpu 核数
	viper.SetDefault(workerConcurrencyPath, 0)
}

func prometheusHandler() gin.HandlerFunc {
	ph := promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{Registry: metrics.Registry})

	return func(c *gin.Context) {
		ph.ServeHTTP(c.Writer, c.Request)
	}
}

// NewHTTPService new a http service
func NewHTTPService() *gin.Engine {
	svr := gin.Default()
	gin.SetMode(viper.GetString(ginModePath))

	pprof.Register(svr)

	// 注册任务
	svr.POST("/task/", http.CreateTask)

	// metrics
	svr.GET("/metrics", prometheusHandler())

	return svr
}

// NewWorkerService new a worker service
func NewWorkerService() error {
	// TODO: 暂时不指定队列
	w, err := worker.NewWorker(
		worker.WorkerConfig{
			Concurrency: viper.GetInt(workerConcurrencyPath),
		},
	)
	if err != nil {
		logger.Errorf("start a worker service error, %v", err)
		return err
	}
	// 加载 handle
	mux := worker.NewServeMux()
	for p, h := range internal.RegisterTaskHandleFunc {
		mux.HandleFunc(p, h)
	}
	if err := w.Run(mux); err != nil {
		logger.Errorf("run worker run")
		return err
	}
	return nil
}