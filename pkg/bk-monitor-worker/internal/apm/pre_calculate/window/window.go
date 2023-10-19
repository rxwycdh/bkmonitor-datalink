// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package window

import (
	"fmt"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/core"
	monitorLogger "github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	"go.uber.org/zap"
	"strconv"
)

type OriginMessage struct {
	DataId   int    `json:"dataid"`
	Items    []Span `json:"items"`
	Datetime string `json:"datetime"`
}

type SpanStatus struct {
	Code    core.SpanStatusCode `json:"code"`
	Message string              `json:"message"`
}

type Span struct {
	TraceId      string         `json:"trace_id"`
	ParentSpanId string         `json:"parent_span_id"`
	EndTime      int            `json:"end_time"`
	ElapsedTime  int            `json:"elapsed_time"`
	Attributes   map[string]any `json:"attributes"`
	Status       SpanStatus     `json:"status"`
	SpanName     string         `json:"span_name"`
	Resource     map[string]any `json:"resource"`
	SpanId       string         `json:"span_id"`
	Kind         int            `json:"kind"`
	StartTime    int            `json:"start_time"`
}

func (s *Span) ToStandardSpan() *StandardSpan {

	return &StandardSpan{
		SpanId:       s.SpanId,
		SpanName:     s.SpanName,
		ParentSpanId: s.ParentSpanId,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		ElapsedTime:  s.ElapsedTime,
		StatusCode:   s.Status.Code,
		Kind:         s.Kind,
		Collections:  s.exactStandardFields(),
	}
}

func (s *Span) exactStandardFields() map[string]string {
	res := make(map[string]string, len(core.StandardFields))

	for _, f := range core.StandardFields {
		switch f.Source {
		case core.SourceAttributes:
			v, exist := s.Attributes[f.Key]
			if exist {
				switch v.(type) {
				case float64:
					res[fmt.Sprintf("attributes.%s", f.Key)] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
				default:
					res[fmt.Sprintf("attributes.%s", f.Key)] = v.(string)
				}
			}
		case core.SourceResource:
			v, exist := s.Resource[f.Key]
			if exist {
				switch v.(type) {
				case float64:
					res[fmt.Sprintf("attributes.%s", f.Key)] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
				default:
					res[fmt.Sprintf("resource.%s", f.Key)] = v.(string)
				}
			}
		case core.SourceOuter:
			switch f.Key {
			case "kind":
				res[f.Key] = strconv.Itoa(s.Kind)
			case "span_name":
				res[f.Key] = s.SpanName
			default:
				logger.Warnf("Try to get a standard field: %s that does not exist. Is the standard field been updated?", f.Key)
			}
		}
	}

	return res
}

type CollectTrace struct {
	TraceId string
	Spans   []*StandardSpan
	Graph   *DiGraph

	Runtime Runtime
}
type StandardSpan struct {
	SpanId       string
	SpanName     string
	ParentSpanId string
	StartTime    int
	EndTime      int
	ElapsedTime  int

	StatusCode  core.SpanStatusCode
	Kind        int
	Collections map[string]string
}

func (s *StandardSpan) GetFieldValue(f core.CommonField) string {
	return s.Collections[f.DisplayKey()]
}

// Handler window handle logic
type Handler interface {
	add(Span)
}

type OperatorMetricKey string

var (
	TraceCount OperatorMetricKey = "traceCount"
	SpanCount  OperatorMetricKey = "spanCount"
)

type OperatorMetric struct {
	Dimension map[string]string
	Value     int
}

// Operator Window processing strategy
type Operator interface {
	Start(spanChan <-chan []Span, runtimeOpt ...RuntimeConfigOption)
	ReportMetric() map[OperatorMetricKey][]OperatorMetric
}

type Operation struct {
	Operator Operator
}

func (o *Operation) Run(spanChan <-chan []Span, runtimeOpt ...RuntimeConfigOption) {
	o.Operator.Start(spanChan, runtimeOpt...)
}

// SpanExistHandler This interface determines how to process existing spans when a span received
type SpanExistHandler interface {
	handleExist(CollectTrace, StandardSpan)
}

// OperatorForm different window implements
type OperatorForm int

const (
	Distributive OperatorForm = 1 << iota
)

var logger = monitorLogger.With(
	zap.String("location", "window"),
)
