package pre_calculate

import (
	"context"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/config"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/notifier"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/storage"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/window"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"go.uber.org/zap"
)

func Initial(parentCtx context.Context) (PreCalculateProcessor, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	return NewPrecalculate().
		WithContext(ctx, cancel).
		WithNotifierConfig(
			notifier.BufferSize(config.NotifierChanBufferSize),
		).
		WithWindowRuntimeConfig(
			window.RuntimeConfigMaxSize(config.WindowMaxSize),
			window.RuntimeConfigExpireInterval(config.WindowExpireInterval),
			window.RuntimeConfigMaxDuration(config.WindowMaxDuration),
			window.ExpireIntervalIncrement(config.WindowExpireIntervalIncrement),
			window.NoDataMaxDuration(config.WindowNoDataMaxDuration),
		).
		WithDistributiveWindowConfig(
			window.DistributiveWindowSubSize(config.DistributiveWindowSubSize),
			window.DistributiveWindowWatchExpiredInterval(config.DistributiveWindowWatchExpireInterval),
			window.ConcurrentProcessCount(config.DistributiveWindowConcurrentCount),
			window.ConcurrentExpirationMaximum(config.DistributiveWindowConcurrentExpirationMaximum),
		).
		WithProcessorConfig(
			window.EnabledTraceInfoCache(config.EnabledTraceInfoCache != 0),
			window.TraceMetaCutLength(config.TraceMetaBloomCutLength),
		).
		WithStorageConfig(
			storage.WorkerCount(config.StorageWorkerCount),
			storage.SaveHoldMaxCount(config.StorageSaveHoldMaxCount),
			storage.SaveHoldDuration(config.StorageSaveHoldMaxDuration),
			storage.CacheBackend(storage.CacheTypeRedis),
			storage.CacheRedisConfig(
				storage.RedisCacheMode(config.StorageRedisMode),
				storage.RedisCacheHost(config.StorageRedisStandaloneHost),
				storage.RedisCachePort(config.StorageRedisStandalonePort),
				storage.RedisCacheSentinelAddress(config.StorageRedisSentinelAddress...),
				storage.RedisCacheMasterName(config.StorageRedisSentinelMasterName),
				storage.RedisCacheSentinelPassword(config.StorageRedisSentinelPassword),
				storage.RedisCachePassword(config.StorageRedisStandalonePassword),
				storage.RedisCacheDb(config.StorageRedisDatabase),
				storage.RedisCacheDialTimeout(config.StorageRedisDialTimeout),
				storage.RedisCacheReadTimeout(config.StorageRedisReadTimeout),
			),
			storage.RedisBloomConfig(
				storage.BloomFpRate(config.StorageBloomFpRate),
				storage.BloomAutoClean(config.StorageBloomAutoClean),
			),
			storage.SaveReqBufferSize(config.StorageSaveRequestBufferSize),
		).
		WithMetricReport(
			EnabledMetric(config.MetricEnabled),
			EnabledProfile(config.ProfileEnabled),
			ProfileAddress(config.ProfileHost),
			ReportHost(config.MetricReportHost),
			SaveRequestCountMetric(config.SaveRequestCountMetricDataId, config.SaveRequestCountMetricAccessToken),
			MessageChanCountMetric(config.MessageCountMetricDataId, config.MessageCountMetricAccessToken),
			WindowTraceAndSpanCountMetric(
				config.WindowSpanCountMetricDataId, config.WindowSpanCountMetricAccessToken,
				config.WindowTraceCountMetricDataId, config.WindowTraceCountMetricAccessToken,
			),
			EsTraceCountMetric(
				config.EsOriginTraceCountMetricDataId, config.EsOriginTraceCountMetricAccessToken,
				config.EsPreCalTraceCountMetricDataId, config.EsPreCalTraceCountMetricAccessToken,
			),
		).
		Build(), nil
}

var apmLogger = logger.With(zap.String("package", "apm_precalculate"))
