// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package apdexcalculator

import (
	"errors"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/define"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/internal/foreach"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/internal/generator"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/internal/testkits"
)

const dst = "bk_apm_duration_destination"

func TestFactory(t *testing.T) {
	content := `
processor:
  - name: "apdex_calculator/standard"
    config:
      calculator:
        type: "standard"
      rules:
        - kind: ""
          metric_name: "bk_apm_duration"
          destination: "apdex_type"
          apdex_t: 20 # ms
`

	psc := testkits.MustLoadProcessorConfigs(content)
	factory, err := newFactory(psc[0].Config, nil)
	assert.NoError(t, err)
	assert.Equal(t, psc[0].Config, factory.MainConfig())

	var c Config
	err = mapstructure.Decode(psc[0].Config, &c)
	assert.NoError(t, err)

	actualConfig := factory.configs.Get("", "", "").(*Config)
	assert.Equal(t, c.Rules, actualConfig.Rules)

	assert.Equal(t, define.ProcessorApdexCalculator, factory.Name())
	assert.False(t, factory.IsDerived())
	assert.False(t, factory.IsPreCheck())
}

func TestProcessMetricsFixedCalculatorNotExistRules(t *testing.T) {
	g := generator.NewMetricsGenerator(define.MetricsOptions{
		MetricName: "bk_apm_duration",
		GaugeCount: 1,
		GeneratorOptions: define.GeneratorOptions{
			Attributes: map[string]string{"kind": "SPAN_KIND_CLIENT"},
		},
	})

	data := g.Generate()
	config := &Config{
		Calculator: CalculatorConfig{
			Type:        "fixed",
			ApdexStatus: apdexSatisfied,
		},
		Rules: []RuleConfig{{
			Kind:        "SPAN_KIND_SERVER",
			MetricName:  "bk_apm_duration",
			Destination: dst,
		}},
	}

	confMap := make(map[string]interface{})
	assert.NoError(t, mapstructure.Decode(config, &confMap))

	factory, err := NewFactory(confMap, nil)
	assert.NoError(t, err)

	record := &define.Record{
		RecordType: define.RecordMetrics,
		Data:       data,
	}
	_, err = factory.Process(record)
	assert.NoError(t, err)

	foreach.Metrics(record.Data.(pmetric.Metrics).ResourceMetrics(), func(metric pmetric.Metric) {
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			dps := metric.Gauge().DataPoints()
			for n := 0; n < dps.Len(); n++ {
				dp := dps.At(n)
				_, ok := dp.Attributes().Get(dst)
				assert.False(t, ok)
			}
		}
	})
}

func TestProcessMetricsFixedCalculatorGlobalDefaultRules(t *testing.T) {
	g := generator.NewMetricsGenerator(define.MetricsOptions{
		MetricName: "bk_apm_duration",
		GaugeCount: 1,
		GeneratorOptions: define.GeneratorOptions{
			Attributes: map[string]string{"kind": "SPAN_KIND_CLIENT"},
		},
	})

	data := g.Generate()
	config := &Config{
		Calculator: CalculatorConfig{
			Type:        "fixed",
			ApdexStatus: apdexTolerating,
		},
		Rules: []RuleConfig{{
			MetricName:  "bk_apm_duration",
			Destination: dst,
		}},
	}

	confMap := make(map[string]interface{})
	assert.NoError(t, mapstructure.Decode(config, &confMap))

	factory, err := NewFactory(confMap, nil)
	assert.NoError(t, err)

	record := &define.Record{
		RecordType: define.RecordMetrics,
		Data:       data,
	}
	_, err = factory.Process(record)
	assert.NoError(t, err)

	foreach.Metrics(record.Data.(pmetric.Metrics).ResourceMetrics(), func(metric pmetric.Metric) {
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			dps := metric.Gauge().DataPoints()
			for n := 0; n < dps.Len(); n++ {
				dp := dps.At(n)
				v, ok := dp.Attributes().Get(dst)
				assert.True(t, ok)
				assert.Equal(t, apdexTolerating, v.AsString())
			}
		}
	})
}

func TestProcessMetricsFixedCalculatorKindDefaultRules(t *testing.T) {
	g := generator.NewMetricsGenerator(define.MetricsOptions{
		MetricName: "bk_apm_duration",
		GaugeCount: 1,
		GeneratorOptions: define.GeneratorOptions{
			Attributes: map[string]string{"kind": "SPAN_KIND_UNSPECIFIED"},
		},
	})

	data := g.Generate()
	config := &Config{
		Calculator: CalculatorConfig{
			Type:        "fixed",
			ApdexStatus: apdexTolerating,
		},
		Rules: []RuleConfig{{
			Kind:        "SPAN_KIND_UNSPECIFIED",
			MetricName:  "bk_apm_duration",
			Destination: dst,
		}},
	}

	confMap := make(map[string]interface{})
	assert.NoError(t, mapstructure.Decode(config, &confMap))

	factory, err := NewFactory(confMap, nil)
	assert.NoError(t, err)

	record := &define.Record{
		RecordType: define.RecordMetrics,
		Data:       data,
	}
	_, err = factory.Process(record)
	assert.NoError(t, err)

	foreach.Metrics(record.Data.(pmetric.Metrics).ResourceMetrics(), func(metric pmetric.Metric) {
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			dps := metric.Gauge().DataPoints()
			for n := 0; n < dps.Len(); n++ {
				dp := dps.At(n)
				v, ok := dp.Attributes().Get(dst)
				assert.True(t, ok)
				assert.Equal(t, apdexTolerating, v.AsString())
			}
		}
	})
}

func TestProcessMetricsFixedCalculator(t *testing.T) {
	g := generator.NewMetricsGenerator(define.MetricsOptions{
		MetricName: "bk_apm_duration",
		GaugeCount: 1,
		GeneratorOptions: define.GeneratorOptions{
			Attributes: map[string]string{
				"kind":        "SPAN_KIND_UNSPECIFIED",
				"http.method": "GET",
			},
		},
	})

	data := g.Generate()
	config := &Config{
		Calculator: CalculatorConfig{
			Type:        "fixed",
			ApdexStatus: apdexSatisfied,
		},
		Rules: []RuleConfig{{
			Kind:         "SPAN_KIND_UNSPECIFIED",
			PredicateKey: "attributes.http.method",
			MetricName:   "bk_apm_duration",
			Destination:  dst,
		}},
	}

	confMap := make(map[string]interface{})
	assert.NoError(t, mapstructure.Decode(config, &confMap))

	factory, err := NewFactory(confMap, nil)
	assert.NoError(t, err)

	record := &define.Record{
		RecordType: define.RecordMetrics,
		Data:       data,
	}
	_, err = factory.Process(record)
	assert.NoError(t, err)

	foreach.Metrics(record.Data.(pmetric.Metrics).ResourceMetrics(), func(metric pmetric.Metric) {
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			dps := metric.Gauge().DataPoints()
			for n := 0; n < dps.Len(); n++ {
				dp := dps.At(n)
				v, ok := dp.Attributes().Get(dst)
				assert.True(t, ok)
				assert.Equal(t, apdexSatisfied, v.AsString())
			}
		}
	})
}

func TestProcessMetricsStandardCalculator(t *testing.T) {
	// apdexSatisfied: val <= threshold
	ok, err := testProcessMetricsStandardCalculator(1e6, 2, apdexSatisfied)
	assert.NoError(t, err)
	assert.True(t, ok)

	// apdexTolerating: val <= 4*threshold
	ok, err = testProcessMetricsStandardCalculator(7e6, 2, apdexTolerating)
	assert.NoError(t, err)
	assert.True(t, ok)

	// apdexFrustrated: val > 4*threshold
	ok, err = testProcessMetricsStandardCalculator(10e6, 2, apdexFrustrated)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func testProcessMetricsStandardCalculator(val float64, threshold float64, status string) (bool, error) {
	g := generator.NewMetricsGenerator(define.MetricsOptions{
		MetricName: "bk_apm_duration",
		GaugeCount: 1,
		Value:      &val,
	})

	data := g.Generate()
	config := &Config{
		Calculator: CalculatorConfig{
			Type: "standard",
		},
		Rules: []RuleConfig{{
			MetricName:  "bk_apm_duration",
			Destination: dst,
			ApdexT:      threshold,
		}},
	}

	confMap := make(map[string]interface{})
	if err := mapstructure.Decode(config, &confMap); err != nil {
		return false, err
	}

	factory, err := NewFactory(confMap, nil)
	if err != nil {
		return false, err
	}

	record := &define.Record{
		RecordType: define.RecordMetrics,
		Data:       data,
	}
	_, err = factory.Process(record)
	if err != nil {
		return false, err
	}

	var errs []error
	foreach.Metrics(record.Data.(pmetric.Metrics).ResourceMetrics(), func(metric pmetric.Metric) {
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			dps := metric.Gauge().DataPoints()
			for n := 0; n < dps.Len(); n++ {
				dp := dps.At(n)
				v, ok := dp.Attributes().Get(dst)
				if !ok {
					errs = append(errs, errors.New("attribute does not exist"))
				}
				if status != v.AsString() {
					errs = append(errs, errors.New("attribute does not exist"))
				}
			}
		}
	})
	if len(errs) > 0 {
		return false, errs[0]
	}
	return true, nil
}
