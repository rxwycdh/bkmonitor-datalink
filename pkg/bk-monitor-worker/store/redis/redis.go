// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
	redisUtils "github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/register/redis"
)

const (
	redisModePath             = "store.redis.mode"
	redisMasterNamePath       = "store.redis.master_name"
	redisAddressPath          = "store.redis.address"
	redisHostPath             = "store.redis.host"
	redisPortPath             = "store.redis.port"
	redisUsernamePath         = "store.redis.username"
	redisSentinelPasswordPath = "store.redis.sentinel_password"
	redisPasswordPath         = "store.redis.password"
	redisDatabasePath         = "store.redis.database"
	redisDialTimeoutPath      = "store.redis.dial_timeout"
	redisReadTimeoutPath      = "store.redis.read_timeout"
)

func init() {
	viper.SetDefault(redisMasterNamePath, "")
	viper.SetDefault(redisAddressPath, []string{"127.0.0.1:6379"})
	viper.SetDefault(redisHostPath, "127.0.0.1")
	viper.SetDefault(redisPortPath, 6379)
	viper.SetDefault(redisUsernamePath, "root")
	viper.SetDefault(redisPasswordPath, "")
	viper.SetDefault(redisSentinelPasswordPath, "")
	viper.SetDefault(redisDatabasePath, 0)
	viper.SetDefault(redisDialTimeoutPath, time.Second*10)
	viper.SetDefault(redisReadTimeoutPath, time.Second*10)
}

type Instance struct {
	ctx    context.Context
	Client redis.UniversalClient
}

var instance *Instance

func NewInstance(ctx context.Context) (*Instance, error) {
	client, err := redisUtils.NewRedisClient(
		ctx,
		&redisUtils.Option{
			Mode:             viper.GetString(redisModePath),
			Host:             viper.GetString(redisHostPath),
			Port:             viper.GetInt(redisPortPath),
			SentinelAddress:  viper.GetStringSlice(redisAddressPath),
			MasterName:       viper.GetString(redisMasterNamePath),
			Password:         viper.GetString(redisPasswordPath),
			SentinelPassword: viper.GetString(redisSentinelPasswordPath),
			Db:               viper.GetInt(redisDatabasePath),
			DialTimeout:      viper.GetDuration(redisDialTimeoutPath),
			ReadTimeout:      viper.GetDuration(redisReadTimeoutPath),
		},
	)
	if err != nil {
		return nil, err
	}
	return &Instance{ctx: ctx, Client: client}, nil
}

// GetInstance get a redis instance
func GetInstance(ctx context.Context) (*Instance, error) {
	if instance != nil {
		return instance, nil
	}
	return NewInstance(ctx)
}

// Open new a instance
func (r *Instance) Open() error {
	return nil
}

// Put put a key-val
func (r *Instance) Put(key, val string, expiration time.Duration) error {
	if err := r.Client.Set(r.ctx, key, val, expiration).Err(); err != nil {
		logger.Errorf("put redis error, key: %s, val: %s, err: %v", key, val, err)
		return err
	}
	return nil
}

// Get get a val from key
func (r *Instance) Get(key string) ([]byte, error) {
	data, err := r.Client.Get(r.ctx, key).Bytes()
	if err != nil {
		logger.Errorf("get redis key: %s error, %v", key, err)
		return nil, err
	}
	return data, nil
}

// Delete delete a key
func (r *Instance) Delete(key string) error {
	exist, err := r.Client.Exists(r.ctx, key).Result()
	if err != nil {
		logger.Errorf("check redis key: %s exist error, %v", key, err)
		return err
	}
	if exist == 0 {
		logger.Warnf("key: %s not exist from redis", key)
		return nil
	}
	if err := r.Client.Del(r.ctx, key).Err(); err != nil {
		logger.Errorf("delete key: %s error, %v", key, err)
		return err
	}
	return nil
}

// Close close connection
func (r *Instance) Close() error {
	return r.Client.Close()
}
