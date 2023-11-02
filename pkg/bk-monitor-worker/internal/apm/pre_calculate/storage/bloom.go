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
	"crypto/md5"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/internal/apm/pre_calculate/core"
	"github.com/minio/highwayhash"
	"time"

	redisBloom "github.com/RedisBloom/redisbloom-go"
	"github.com/gomodule/redigo/redis"
	boom "github.com/tylertreat/BoomFilters"
)

type BloomStorageData struct {
	Key string
}

type BloomOperator interface {
	Add(BloomStorageData) error
	Exist(string) (bool, error)
}

type BloomOptions struct {
	fpRate    float64
	autoClean time.Duration
	initCap   int

	layersBloomOptions LayersBloomOptions
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

func LayerBloomConfig(opts ...LayersBloomOption) BloomOption {
	return func(options *BloomOptions) {
		option := LayersBloomOptions{}
		for _, setter := range opts {
			setter(&option)
		}
		options.layersBloomOptions = option
	}
}

type Bloom struct {
	filterName string
	config     BloomOptions
	c          *redisBloom.Client
}

func (b *Bloom) Add(data BloomStorageData) error {
	_, err := b.c.Add(b.filterName, data.Key)
	return err
}

func (b *Bloom) Exist(k string) (bool, error) {
	return b.c.Exists(b.filterName, k)
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
	cleanDuration time.Duration
}

func (m *MemoryBloom) Add(data BloomStorageData) error {
	m.c.Add([]byte(data.Key))
	return nil
}

func (m *MemoryBloom) Exist(key string) (bool, error) {
	return m.c.Test([]byte(key)), nil
}

func (m *MemoryBloom) AutoReset() {
	// Prevent the memory from being too large.
	// Data will be cleared after a specified time.
	logger.Infof("Bloom-filter will reset every %s", m.config.autoClean)
	for {
		if time.Now().After(m.nextCleanDate) {
			m.c = m.c.Reset()
			m.nextCleanDate = time.Now().Add(m.cleanDuration)
			logger.Infof("Bloom-filter reset data trigger, next time the filter reset data is %s", m.nextCleanDate)
		}
		time.Sleep(1 * time.Minute)
	}
}

func newMemoryCacheBloomClient(options BloomOptions) (BloomOperator, error) {
	sbf := boom.NewScalableBloomFilter(uint(options.initCap), options.fpRate, 0.8)
	bloom := &MemoryBloom{c: sbf, config: options, nextCleanDate: time.Now().Add(options.autoClean), cleanDuration: options.autoClean}
	go bloom.AutoReset()
	return bloom, nil
}

type LayersBloomOptions struct {
	layers int
}

type LayersBloomOption func(*LayersBloomOptions)

func Layers(s int) LayersBloomOption {
	return func(options *LayersBloomOptions) {
		if s > len(strategies) {
			logger.Warnf("layer: %d > strategies count, set to %d", s, len(strategies))
			s = len(strategies)
		}
		options.layers = s
	}
}

type layerStrategy func(string) []byte

var (
	strategies = []layerStrategy{
		// truncated 16
		func(s string) []byte {
			return []byte(s[16:])
		},
		// truncated 8
		func(s string) []byte {
			return []byte(s[24:])
		},
		// full
		func(s string) []byte {
			return []byte(s)
		},
		// md5
		func(s string) []byte {
			hash := md5.New()
			hash.Write([]byte(s))
			return hash.Sum(nil)
		},
		// hash
		func(s string) []byte {
			h, _ := highwayhash.New([]byte(core.HashSecret))
			h.Write([]byte(s))
			return h.Sum(nil)
		},
	}
)

type LayersMemoryBloom struct {
	blooms     []*MemoryBloom
	strategies []layerStrategy
}

func newLayersBloomClient(options BloomOptions) (BloomOperator, error) {
	var blooms []*MemoryBloom

	for i := 0; i < options.layersBloomOptions.layers; i++ {
		sbf := boom.NewScalableBloomFilter(uint(options.initCap), options.fpRate, 0.8)
		bloom := &MemoryBloom{c: sbf, config: options, nextCleanDate: time.Now().Add(options.autoClean), cleanDuration: options.autoClean}
		go bloom.AutoReset()
		blooms = append(blooms, bloom)
	}

	return &LayersMemoryBloom{blooms: blooms, strategies: strategies}, nil
}

func (l *LayersMemoryBloom) Add(data BloomStorageData) error {
	for index, b := range l.blooms {
		key := l.strategies[index](data.Key)
		b.c.Add(key)
	}
	return nil
}

func (l *LayersMemoryBloom) Exist(originKey string) (bool, error) {

	for index, b := range l.blooms {
		key := l.strategies[index](originKey)
		e := b.c.Test(key)
		if !e {
			return false, nil
		}
	}

	return true, nil
}
