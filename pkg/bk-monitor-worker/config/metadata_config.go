package config

var (
	MetadataMetricDimensionMetricKeyPrefix             string
	MetadataMetricDimensionKeyPrefix                   string
	MetadataMetricDimensionMaxMetricFetchStep          int
	MetadataMetricDimensionTimeSeriesMetricExpiredDays int
)

func initMetadataVariables() {
	MetadataMetricDimensionMetricKeyPrefix = GetValue("taskConfig.metadata.metricDimension.metricKeyPrefix", "bkmonitor:metrics_")
	MetadataMetricDimensionKeyPrefix = GetValue("taskConfig.metadata.metricDimension.metricDimensionKeyPrefix", "bkmonitor:metric_dimensions_")
	MetadataMetricDimensionMaxMetricFetchStep = GetValue("taskConfig.metadata.metricDimension.maxMetricsFetchStep", 500)
	MetadataMetricDimensionTimeSeriesMetricExpiredDays = GetValue("taskConfig.metadata.metricDimension.timeSeriesMetricExpiredDays", 30)
}
