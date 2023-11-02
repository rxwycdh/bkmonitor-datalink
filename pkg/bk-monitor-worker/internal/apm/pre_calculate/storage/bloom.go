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

	layersBloomOptions            LayersBloomOptions
	layersCapDecreaseBloomOptions LayersCapDecreaseBloomOptions
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

func LayerCapDecreaseBloomConfig(opts ...LayersCapDecreaseBloomOption) BloomOption {
	return func(options *BloomOptions) {
		option := LayersCapDecreaseBloomOptions{}
		for _, setter := range opts {
			setter(&option)
		}
		options.layersCapDecreaseBloomOptions = option
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
	c             boom.Filter
	nextCleanDate time.Time
	cleanDuration time.Duration
	resetFunc     func()
}

func (m *MemoryBloom) Add(data []byte) boom.Filter {
	return m.c.Add(data)
}

func (m *MemoryBloom) Test(key []byte) bool {
	return m.c.Test(key)
}

func (m *MemoryBloom) TestAndAdd(key []byte) bool {
	return m.c.TestAndAdd(key)
}

func (m *MemoryBloom) AutoReset() {
	// Prevent the memory from being too large.
	// Data will be cleared after a specified time.
	logger.Infof("Bloom-filter will reset every %s", m.config.autoClean)
	for {
		if time.Now().After(m.nextCleanDate) {
			m.resetFunc()
			m.nextCleanDate = time.Now().Add(m.cleanDuration)
			logger.Infof("Bloom-filter reset data trigger, next time the filter reset data is %s", m.nextCleanDate)
		}
		time.Sleep(1 * time.Minute)
	}
}

func newBloomClient(f boom.Filter, resetFunc func(), options BloomOptions) boom.Filter {
	bloom := &MemoryBloom{
		c:             f,
		config:        options,
		nextCleanDate: time.Now().Add(options.autoClean),
		cleanDuration: options.autoClean,
		resetFunc:     resetFunc,
	}
	go bloom.AutoReset()
	return bloom
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
	blooms     []boom.Filter
	strategies []layerStrategy
}

func newLayersBloomClient(options BloomOptions) (BloomOperator, error) {
	var blooms []boom.Filter

	for i := 0; i < options.layersBloomOptions.layers; i++ {
		sbf := boom.NewScalableBloomFilter(uint(options.initCap), options.fpRate, 0.8)
		bloom := newBloomClient(sbf, func() { sbf.Reset() }, options)
		blooms = append(blooms, bloom)
	}
	logger.Infof("bloom-filter layers: %d", options.layersBloomOptions.layers)
	return &LayersMemoryBloom{blooms: blooms, strategies: strategies}, nil
}

func (l *LayersMemoryBloom) Add(data BloomStorageData) error {
	for index, b := range l.blooms {
		key := l.strategies[index](data.Key)
		if err := b.Add(key); err != nil {
			logger.Errorf("failed to add data in blooms[%d]. error: %s", index, err)
		}

	}
	return nil
}

func (l *LayersMemoryBloom) Exist(originKey string) (bool, error) {

	for index, b := range l.blooms {
		key := l.strategies[index](originKey)
		e := b.Test(key)
		if !e {
			return false, nil
		}
	}

	return true, nil
}

type LayersCapDecreaseBloomOption func(*LayersCapDecreaseBloomOptions)

type LayersCapDecreaseBloomOptions struct {
	cap     int
	layers  int
	divisor int
}

func CapDecreaseBloomCap(c int) LayersCapDecreaseBloomOption {
	return func(options *LayersCapDecreaseBloomOptions) {
		options.cap = c
	}
}

func CapDecreaseBloomLayers(c int) LayersCapDecreaseBloomOption {
	return func(options *LayersCapDecreaseBloomOptions) {
		options.layers = c
	}
}

func CapDecreaseBloomDivisor(c int) LayersCapDecreaseBloomOption {
	return func(options *LayersCapDecreaseBloomOptions) {
		options.divisor = c
	}
}

type LayersCapDecreaseBloom struct {
	blooms []boom.Filter
}

func newLayersCapDecreaseBloomClient(options BloomOptions) (BloomOperator, error) {
	var blooms []boom.Filter

	curCap := options.layersCapDecreaseBloomOptions.cap
	for i := 0; i < options.layersCapDecreaseBloomOptions.layers; i++ {
		sbf := boom.NewBloomFilter(uint(curCap), options.fpRate)
		bloom := newBloomClient(sbf, func() { sbf.Reset() }, options)
		blooms = append(blooms, bloom)
		curCap = curCap / options.layersCapDecreaseBloomOptions.divisor
	}

	return &LayersCapDecreaseBloom{blooms: blooms}, nil
}

func (l *LayersCapDecreaseBloom) Add(data BloomStorageData) error {
	key := []byte(data.Key)
	for _, b := range l.blooms {
		b.Add(key)
	}
	return nil
}

func (l *LayersCapDecreaseBloom) Exist(originKey string) (bool, error) {
	key := []byte(originKey)

	for _, b := range l.blooms {
		exist := b.Test(key)
		if !exist {
			return false, nil
		}
	}

	return true, nil
}
