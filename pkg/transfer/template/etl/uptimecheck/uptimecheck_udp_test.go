// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package uptimecheck_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/template/etl/uptimecheck"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/testsuite"
)

//go:embed fixture/uptimecheck_udp_test_data.json
var uptimecheckUDPTestData string

//go:embed fixture/uptimecheck_udp_test_consul_data.json
var uptimecheckUDPTestConsulData string

// UptimecheckUDPTest :
type UptimecheckUDPTest struct {
	testsuite.ETLSuite
}

// TestUsage :
func (s *UptimecheckUDPTest) TestUsage() {
	s.CTX = testsuite.PipelineConfigStringInfoContext(s.CTX, s.PipelineConfig, uptimecheckUDPTestConsulData)
	processor, err := uptimecheck.NewUDPProcessor(s.CTX, "test")
	s.NoError(err)
	s.Run(
		uptimecheckUDPTestData,
		processor,
		func(result map[string]interface{}) {
			s.EqualRecord(result, map[string]interface{}{
				"dimensions": map[string]interface{}{
					"bk_cloud_id":    "2",
					"error_code":     "0",
					"status":         "0",
					"node_id":        "4",
					"task_id":        "109",
					"task_type":      "udp",
					"bk_agent_id":    "010000525400c48bdc1670385834306k",
					"bk_host_id":     "30145",
					"bk_biz_id":      2.0,
					"ip":             "127.0.0.1",
					"bk_supplier_id": "0",
					"target_host":    "127.0.0.1",
					"target_port":    "9211",
				},
				"metrics": map[string]interface{}{
					"task_duration": 5000.0,
					"available":     1.000000,
					"times":         0.0,
				},
				"time": 1552967059,
			})
		},
	)
}

// TestServletTest :
func TestUptimecheckUDPTest(t *testing.T) {
	suite.Run(t, new(UptimecheckUDPTest))
}
