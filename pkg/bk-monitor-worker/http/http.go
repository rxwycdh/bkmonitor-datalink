package http

import (
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/config"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/metrics"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prometheusHandler() gin.HandlerFunc {
	ph := promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{Registry: metrics.Registry})

	return func(c *gin.Context) {
		ph.ServeHTTP(c.Writer, c.Request)
	}
}

// NewHTTPService new a http service
func NewHTTPService(enableApi bool) *gin.Engine {
	svr := gin.Default()
	gin.SetMode(config.HttpGinMode)

	if config.HttpEnabledPprof {
		pprof.Register(svr)
		logger.Info("Pprof started")
	}

	if enableApi {
		// 注册任务
		svr.POST("/task/", CreateTask)
		// 获取运行中的任务列表
		svr.GET("/task/", ListTask)
		// 删除任务
		svr.DELETE("/task/", RemoveTask)
	}

	// metrics
	svr.GET("/metrics", prometheusHandler())

	return svr
}
