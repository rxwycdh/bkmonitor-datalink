package storage

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/prometheus/prompb"
)

// MetricsName
const (
	ApmServiceInstanceRelation = "apm_service_with_apm_service_instance_relation"
	ApmServiceK8sRelation      = "apm_service_instance_with_k8s_address_relation"
	ApmServiceSystemRelation   = "apm_service_instance_with_system_relation"

	ApmServiceFlow       = "apm_service_to_apm_service_flow"
	SystemApmServiceFlow = "system_to_apm_service_flow"
	ApmServiceSystemFlow = "apm_service_to_system_flow"
	SystemFlow           = "system_to_system_flow"
)

var relationMetricsMetadataMapping = map[string]prompb.MetricMetadata{
	ApmServiceInstanceRelation: {
		MetricFamilyName: ApmServiceInstanceRelation,
		Type:             prompb.MetricMetadata_GAUGE,
		Help:             "APM 服务实例关联指标",
	},
	ApmServiceK8sRelation: {
		MetricFamilyName: ApmServiceK8sRelation,
		Type:             prompb.MetricMetadata_GAUGE,
		Help:             "APM 服务与 K8s 关联指标",
	},
	ApmServiceSystemRelation: {
		MetricFamilyName: ApmServiceSystemRelation,
		Type:             prompb.MetricMetadata_GAUGE,
		Help:             "APM 服务与主机关联指标",
	},
}

var flowMetricsMetadataMapping = map[string]prompb.MetricMetadata{
	"min": {
		Type: prompb.MetricMetadata_GAUGE,
		Help: "聚合周期内耗时最小值",
		Unit: "microseconds",
	},
	"max": {
		Type: prompb.MetricMetadata_GAUGE,
		Help: "聚合周期内耗时最大值",
		Unit: "microseconds",
	},
	"sum": {
		Type: prompb.MetricMetadata_COUNTER,
		Help: "聚合周期内耗时总和",
		Unit: "microseconds",
	},
	"count": {
		Type: prompb.MetricMetadata_COUNTER,
		Help: "聚合周期内耗时总数",
	},
	"bucket": {
		Type: prompb.MetricMetadata_HISTOGRAM,
		Help: "聚合周期内耗时bucket",
		Unit: "microseconds",
	},
}

type MetricCollector interface {
	Observe(value any)
	Collect() prompb.WriteRequest
	Ttl() time.Duration
}

func dimensionKeyToNameAndLabel(dimensionKey string, ignoreName bool) (string, []prompb.Label) {
	pairs := strings.Split(dimensionKey, ",")
	var labels []prompb.Label
	var name string
	for _, pair := range pairs {
		composition := strings.Split(pair, "=")
		if len(composition) == 2 {
			if composition[0] == "__name__" {
				name = composition[1]
				if !ignoreName {
					labels = append(labels, prompb.Label{Name: composition[0], Value: composition[1]})
				}
			} else {
				labels = append(labels, prompb.Label{Name: composition[0], Value: composition[1]})
			}
		}
	}
	return name, labels
}

type FlowMetricRecordStats struct {
	DurationValues []float64
}

type flowMetricStats struct {
	FlowDurationMax, FlowDurationMin, FlowDurationSum, FlowDurationCount float64
	FlowDurationBucket                                                   map[float64]int
}

type flowMetricsCollector struct {
	mu  sync.Mutex
	ttl time.Duration

	data    map[string]*flowMetricStats
	buckets []float64
}

func (c *flowMetricsCollector) Ttl() time.Duration { return c.ttl }

func (c *flowMetricsCollector) Observe(value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	mapping := value.(map[string]*FlowMetricRecordStats)
	for dimensionKey, v := range mapping {
		for _, duration := range v.DurationValues {
			if _, exist := c.data[dimensionKey]; !exist {
				c.data[dimensionKey] = &flowMetricStats{
					FlowDurationMax:    math.Inf(-1),
					FlowDurationMin:    math.Inf(1),
					FlowDurationBucket: make(map[float64]int),
				}
			}

			c.data[dimensionKey].FlowDurationCount++
			c.data[dimensionKey].FlowDurationSum += duration

			if duration > c.data[dimensionKey].FlowDurationMax {
				c.data[dimensionKey].FlowDurationMax = duration
			}

			if duration < c.data[dimensionKey].FlowDurationMin {
				c.data[dimensionKey].FlowDurationMin = duration
			}
			for _, bucket := range c.buckets {
				if duration <= bucket {
					c.data[dimensionKey].FlowDurationBucket[bucket]++
				}
			}
		}
	}
}

func (c *flowMetricsCollector) Collect() prompb.WriteRequest {
	c.mu.Lock()
	defer c.mu.Unlock()

	res := c.convert()
	c.data = make(map[string]*flowMetricStats)

	return res
}

func (c *flowMetricsCollector) convert() prompb.WriteRequest {
	logger.Infof("🌟🌟🌟convert start 🌟🌟🌟")

	copyLabels := func(labels []prompb.Label) []prompb.Label {
		newLabels := make([]prompb.Label, len(labels))
		copy(newLabels, labels)
		return newLabels
	}
	var series []prompb.TimeSeries
	var metricsName []string

	ts := time.Now().UnixNano() / int64(time.Millisecond)
	for key, stats := range c.data {
		name, labels := dimensionKeyToNameAndLabel(key, true)
		// todo remove test log
		logger.Infof("[FlowMetricsCollector] fromService: %s toService: %s minValue: %+v", labels[0].Value, labels[2].Value, stats.FlowDurationMin)

		metricsName = append(metricsName, name)

		series = append(series, prompb.TimeSeries{
			Labels:  append(copyLabels(labels), prompb.Label{Name: "__name__", Value: fmt.Sprintf("%s%s", name, "_min")}),
			Samples: []prompb.Sample{{Value: stats.FlowDurationMin, Timestamp: ts}},
		})

		series = append(series, prompb.TimeSeries{
			Labels:  append(copyLabels(labels), prompb.Label{Name: "__name__", Value: fmt.Sprintf("%s%s", name, "_max")}),
			Samples: []prompb.Sample{{Value: stats.FlowDurationMax, Timestamp: ts}},
		})

		series = append(series, prompb.TimeSeries{
			Labels:  append(copyLabels(labels), prompb.Label{Name: "__name__", Value: fmt.Sprintf("%s%s", name, "_sum")}),
			Samples: []prompb.Sample{{Value: stats.FlowDurationSum, Timestamp: ts}},
		})

		series = append(series, prompb.TimeSeries{
			Labels:  append(copyLabels(labels), prompb.Label{Name: "__name__", Value: fmt.Sprintf("%s%s", name, "_count")}),
			Samples: []prompb.Sample{{Value: stats.FlowDurationCount, Timestamp: ts}},
		})

		for bucket, count := range stats.FlowDurationBucket {
			series = append(series, prompb.TimeSeries{
				Labels: append(
					copyLabels(labels), []prompb.Label{
						{Name: "__name__", Value: name + "_bucket"},
						{Name: "le", Value: fmt.Sprintf("%f", bucket)},
					}...),
				Samples: []prompb.Sample{{Value: float64(count), Timestamp: ts}},
			})
		}
	}

	var infos []prompb.MetricMetadata
	for _, n := range metricsName {
		for k, info := range flowMetricsMetadataMapping {
			infos = append(infos, prompb.MetricMetadata{
				MetricFamilyName: fmt.Sprintf("%s_%s", n, k),
				Type:             info.Type,
				Help:             info.Help,
				Unit:             info.Unit,
			})
		}
	}

	logger.Infof("🌟🌟🌟convert end 🌟🌟🌟")
	return prompb.WriteRequest{Timeseries: series, Metadata: infos}
}

func newFlowMetricCollector(buckets []float64, ttl time.Duration) *flowMetricsCollector {
	return &flowMetricsCollector{
		ttl:     ttl,
		data:    make(map[string]*flowMetricStats),
		buckets: buckets,
	}
}

type relationMetricsCollector struct {
	mu   sync.Mutex
	data map[string]time.Time
	ttl  time.Duration
}

func (r *relationMetricsCollector) Ttl() time.Duration { return r.ttl }

func (r *relationMetricsCollector) Observe(value any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	labels := value.([]string)
	for _, s := range labels {
		if _, exist := r.data[s]; !exist {
			r.data[s] = time.Now()
		}
	}
}

func (r *relationMetricsCollector) Collect() prompb.WriteRequest {
	r.mu.Lock()
	defer r.mu.Unlock()

	edge := time.Now().Add(-r.ttl)
	var keys []string
	for dimensionKey, ts := range r.data {
		if ts.Before(edge) {
			logger.Debugf("[RelationMetricsCollector] key: %s expired, timestamp: %s", dimensionKey, ts)
			keys = append(keys, dimensionKey)
		}
	}
	res := r.convert(keys)
	for _, k := range keys {
		delete(r.data, k)
	}
	return res
}

func (r *relationMetricsCollector) convert(dimensionKeys []string) prompb.WriteRequest {
	logger.Infof("🌈🌈🌈 convert start 🌈🌈🌈")

	var series []prompb.TimeSeries
	metricName := make(map[string]int, len(dimensionKeys))
	var infos []prompb.MetricMetadata

	ts := time.Now().UnixNano() / int64(time.Millisecond)
	for _, key := range dimensionKeys {
		name, labels := dimensionKeyToNameAndLabel(key, false)
		series = append(series, prompb.TimeSeries{
			Labels:  labels,
			Samples: []prompb.Sample{{Value: 1, Timestamp: ts}},
		})
		metricName[name]++
	}

	for k := range metricName {
		metadata, exist := relationMetricsMetadataMapping[k]
		if exist {
			infos = append(infos, metadata)
		}
	}

	logger.Infof("🌈🌈🌈 convert end 🌈🌈🌈")
	return prompb.WriteRequest{Timeseries: series, Metadata: infos}
}

func newRelationMetricCollector(ttl time.Duration) *relationMetricsCollector {
	return &relationMetricsCollector{ttl: ttl, data: make(map[string]time.Time)}
}
