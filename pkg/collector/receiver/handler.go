// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package receiver

import (
	"net/http"
	"time"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/define"
)

// ResponseHandler 请求响应处理器
type ResponseHandler interface {
	// ContentType 返回相应内容类型
	ContentType() string

	// Response 根据 recordType 返回响应内容
	Response(rtype define.RecordType) ([]byte, error)

	// ErrorStatus 返回 status 序列化后内容
	ErrorStatus(status interface{}) ([]byte, error)

	// Unmarshal 根据 recordType 反序列化数据
	Unmarshal(rtype define.RecordType, b []byte) (interface{}, error)
}

func WriteResponse(w http.ResponseWriter, contentType string, statusCode int, msg []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, _ = w.Write(msg)
}

func RecordHandleMetrics(mm *metricMonitor, token define.Token, protocol define.RequestType, rtype define.RecordType, bs int, t time.Time) {
	mm.AddReceivedBytesCounter(float64(bs), protocol, rtype, token.Original)
	mm.ObserveBytesDistribution(float64(bs), protocol, rtype, token.Original)
	mm.ObserveHandledDuration(t, protocol, rtype, token.Original)
	mm.IncHandledCounter(protocol, rtype, token.Original)
	DefaultMetricMonitor.SetTokenInfo(token)
}

// Encoder 负责解析 Traces/Metrics/Logs 数据至 OT 标准数据模型
type Encoder interface {
	Type() string
	UnmarshalTraces(b []byte) (ptrace.Traces, error)
	UnmarshalMetrics(b []byte) (pmetric.Metrics, error)
	UnmarshalLogs(b []byte) (plog.Logs, error)
}

type NoopEncoder struct{}

func (NoopEncoder) UnmarshalTraces(_ []byte) (ptrace.Traces, error) {
	return ptrace.NewTraces(), nil
}

func (NoopEncoder) UnmarshalMetrics(_ []byte) (pmetric.Metrics, error) {
	return pmetric.NewMetrics(), nil
}

func (NoopEncoder) UnmarshalLogs(_ []byte) (plog.Logs, error) {
	return plog.NewLogs(), nil
}

func UnmarshalRecordData(encoder Encoder, rtype define.RecordType, b []byte) (interface{}, error) {
	switch rtype {
	case define.RecordTraces:
		return encoder.UnmarshalTraces(b)
	case define.RecordMetrics:
		return encoder.UnmarshalMetrics(b)
	case define.RecordLogs:
		return encoder.UnmarshalLogs(b)
	}
	return nil, define.ErrUnknownRecordType
}
