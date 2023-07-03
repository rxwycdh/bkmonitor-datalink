// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package logging

import (
	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
)

const (
	logFileNamePath   = "log.filename"
	logMaxSizePath    = "log.max_size"
	logMaxAgePath     = "log.max_age"
	logMaxBackupsPath = "log.max_backups"
	logLevelPath      = "log.level"
)

func init() {
	viper.SetDefault(logFileNamePath, "./task.log")
	viper.SetDefault(logMaxSizePath, 200)  // default 200M
	viper.SetDefault(logMaxAgePath, 1)     // default 1
	viper.SetDefault(logMaxBackupsPath, 5) // default 5
	viper.SetDefault(logLevelPath, "info") // default info level
}

// InitLogger init a bkmonitor logger
func InitLogger() {
	logger.SetOptions(logger.Options{
		Filename:   viper.GetString(logFileNamePath),
		MaxSize:    viper.GetInt(logMaxSizePath),
		MaxAge:     viper.GetInt(logMaxAgePath),
		MaxBackups: viper.GetInt(logMaxBackupsPath),
		Level:      viper.GetString(logLevelPath),
	})
}
