// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/config"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/logging"
	service "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/service"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
)

const (
	serviceListenPath = "service.listen"
	servicePortPath   = "service.port"
)

func init() {
	viper.SetDefault(serviceListenPath, "127.0.0.1")
	viper.SetDefault(servicePortPath, 10211)
}

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "bk-monitor workers",
	Long:  "worker module for blueking monitor",
	Run:   start,
}

// start 启动服务
func start(cmd *cobra.Command, args []string) {
	fmt.Println("start service...")
	// 初始化配置
	config.InitConfig()

	// 初始化日志
	logging.InitLogger()

	// 启动 worker

	// 启动 http 服务
	r := service.NewHTTPService()
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", viper.GetString(serviceListenPath), viper.GetInt(servicePortPath)),
		Handler: r,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen addr error, %v", err)
			panic(err)
		}
	}()
	// 信号处理
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch <-s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("shutdown service error : %s", err)
			}
			logger.Warn("service exit by syscall SIGQUIT, SIGTERM or SIGINT")
			return
		}
	}
}

// Execute 执行命令
func Execute() {
	rootCmd.Flags().StringVarP(
		&config.ConfigPath, "config", "c", "", "path of project service config files",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("start bmw service error, %s", err)
		os.Exit(1)
	}
}