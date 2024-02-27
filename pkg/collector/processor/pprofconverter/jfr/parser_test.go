// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package jfr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/define"
)

func TestJfrConvertBody(t *testing.T) {
	c := &Converter{}

	t.Run("invalid data type", func(t *testing.T) {
		_, _, err := c.convertBody("invalid")
		assert.Error(t, err)
	})

	t.Run("valid data", func(t *testing.T) {
		jfrData := define.ProfileJfrFormatOrigin{
			Jfr:    []byte("jfr data"),
			Labels: []byte("jfr labels"),
		}
		_, _, err := c.convertBody(jfrData)
		assert.NoError(t, err)
	})
}

func TestConverter_ParseToPprof(t *testing.T) {
	c := &Converter{}
	data, err := ReadGzipFile("../testdata/jfr_cortex-dev-01__kafka-0__cpu_lock_alloc__0.jfr.gz")
	assert.NoError(t, err)

	pd := define.ProfilesRawData{
		Metadata: define.ProfileMetadata{
			StartTime:       time.Now(),
			EndTime:         time.Now(),
			AppName:         "testApp",
			BkBizID:         1,
			SpyName:         "testSpy",
			Format:          define.FormatJFR,
			SampleRate:      100,
			Units:           UnitNanoseconds,
			AggregationType: "testAggregation",
			Tags:            map[string]string{"tag1": "value1"},
		},
		Data: define.ProfileJfrFormatOrigin{Jfr: data},
	}

	t.Run("Test Parse", func(t *testing.T) {
		result, err := c.Parse(pd)

		assert.Nil(t, err)

		assert.Equal(t, pd.Metadata, result.Metadata)

		assert.NotNil(t, result.Profiles)
	})
}
