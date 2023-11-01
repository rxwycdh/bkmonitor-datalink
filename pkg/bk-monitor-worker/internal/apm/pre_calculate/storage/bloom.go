// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package storage

import (
	"time"

	redisBloom "github.com/RedisBloom/redisbloom-go"
	"github.com/gomodule/redigo/redis"
	boom "github.com/tylertreat/BoomFilters"
)

type BloomStorageData struct {
	Key []byte
}

type BloomOperator interface {
	Add(BloomStorageData) error
	Exist([]byte) (bool, error)
}

type BloomOptions struct {
	fpRate    float64
	autoClean time.Duration
	initCap   int
}

type BloomOption func(*BloomOptions)

func BloomFpRate(s float64) BloomOption {
	return func(options *BloomOptions) {
		options.fpRate = s
	}
}

func BloomAutoClean(s int) BloomOption {
	return func(options *BloomOptions) {
		options.autoClean = time.Duration(s) * time.Minute
	}
}

func InitCap(s int) BloomOption {
	return func(options *BloomOptions) {
		options.initCap = s
	}
}

type Bloom struct {
	filterName string
	config     BloomOptions
	c          *redisBloom.Client
}

func (b *Bloom) Add(data BloomStorageData) error {
	_, err := b.c.Add(b.filterName, string(data.Key))
	return err
}

func (b *Bloom) Exist(k []byte) (bool, error) {
	return b.c.Exists(b.filterName, string(k))
}

func newRedisBloomClient(rConfig RedisCacheOptions, opts BloomOptions) (BloomOperator, error) {
	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", rConfig.host, redis.DialPassword(rConfig.password), redis.DialDatabase(rConfig.db))
	}}
	c := redisBloom.NewClientFromPool(pool, "bloom-client")

	return &Bloom{filterName: "traceMeta", config: opts, c: c}, nil
}

type MemoryBloom struct {
	config        BloomOptions
	c             *boom.ScalableBloomFilter
	nextCleanDate time.Time
}

func (m *MemoryBloom) Add(data BloomStorageData) error {
	m.c.Add(data.Key)
	return nil
}

func (m *MemoryBloom) Exist(key []byte) (bool, error) {
	return m.c.Test(key), nil
}

func (m *MemoryBloom) AutoReset() {
	// Prevent the memory from being too large.
	// Data will be cleared after a specified time.
	logger.Infof("Bloom-filter will reset data every %s", m.config.autoClean)
	for {
		logger.Infof("Next time the filter reset data is %s", m.nextCleanDate)
		if time.Now().After(m.nextCleanDate) {
			logger.Infof("Bloom-filter reset data trigger")
			m.c.Reset()
		}
		time.Sleep(24 * time.Hour)
	}
}

func newMemoryCacheBloomClient(options BloomOptions) (BloomOperator, error) {
	sbf := boom.NewScalableBloomFilter(10000, options.fpRate, 0.8)
	bloom := &MemoryBloom{c: sbf, config: options, nextCleanDate: time.Now().Add(options.autoClean)}
	go bloom.AutoReset()
	return bloom, nil
}
