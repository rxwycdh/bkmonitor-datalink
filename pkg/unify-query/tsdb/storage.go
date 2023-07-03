// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package tsdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/consul"
)

var (
	storageMap  = make(map[string]*Storage)
	storageLock = new(sync.RWMutex)
)

// ReloadTsDBStorage 重新加载存储实例到内存里面
func ReloadTsDBStorage(ctx context.Context, tsDBs map[string]*consul.Storage, opt *Options) error {
	newStorageMap := make(map[string]*Storage, len(tsDBs))

	for storageID, tsDB := range tsDBs {
		var (
			storage *Storage
		)
		if tsDB.Type == consul.VictoriaMetricsStorageType {
			storage = &Storage{
				Type:     tsDB.Type,
				Address:  tsDB.Address,
				Username: tsDB.Username,
				Password: tsDB.Password,

				UriPath: opt.VM.UriPath,
				Timeout: opt.VM.Timeout,
			}
		} else {
			storage = &Storage{
				Type:     tsDB.Type,
				Address:  tsDB.Address,
				Username: tsDB.Username,
				Password: tsDB.Password,
				Timeout:  opt.InfluxDB.Timeout,

				MaxLimit:      opt.InfluxDB.MaxLimit,
				MaxSLimit:     opt.InfluxDB.MaxSLimit,
				Toleration:    opt.InfluxDB.Tolerance,
				ReadRateLimit: opt.InfluxDB.ReadRateLimit,

				ContentType: opt.InfluxDB.ContentType,
				ChunkSize:   opt.InfluxDB.ChunkSize,

				UriPath:        opt.InfluxDB.RawUriPath,
				Accept:         opt.InfluxDB.Accept,
				AcceptEncoding: opt.InfluxDB.AcceptEncoding,
			}
		}

		newStorageMap[storageID] = storage
	}
	storageLock.Lock()
	defer storageLock.Unlock()

	storageMap = newStorageMap
	return nil
}

func Print() string {
	storageLock.RLock()
	defer storageLock.RUnlock()
	str := "--------------------------- storage list --------------------------------------\n"
	for k, s := range storageMap {
		str += fmt.Sprintf("%s: %+v \n", k, s)
	}
	return str
}

// GetStorage 初始化全局 tsdb 标准实例
func GetStorage(storageID string) (*Storage, error) {
	storageLock.Lock()
	defer storageLock.Unlock()
	storage, ok := storageMap[storageID]
	if !ok {
		return nil, fmt.Errorf("%s: storageID: %s", ErrStorageNotFound, storageID)
	}
	return storage, nil
}

// SetStorage 写入实例到内存去
func SetStorage(storageID string, storage *Storage) {
	storageLock.Lock()
	defer storageLock.Unlock()

	storageMap[storageID] = storage
}
