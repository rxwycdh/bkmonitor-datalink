// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/config"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/bk-monitor-worker/store/consul"
)

// MetadataCenter The configuration center uses dataId as the key and stores basic information,
// including app_name, bk_biz_name, app_id, etc...
type MetadataCenter struct {
	mapping map[string]DataIdInfo
	consul  consul.Instance
}

type ConsulInfo struct {
	BkBizId     int             `json:"bk_biz_id"`
	BkBizName   string          `json:"bk_biz_name"`
	AppId       int             `json:"app_id"`
	AppName     string          `json:"app_name"`
	KafkaInfo   consulKafkaInfo `json:"kafka_info"`
	TraceEsInfo consulEsInfo    `json:"trace_es_info"`
	SaveEsInfo  consulEsInfo    `json:"save_es_info"`
}

type consulKafkaInfo struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
}

type consulEsInfo struct {
	IndexName string `json:"index_name"`
	Host      string `json:"host"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type DataIdInfo struct {
	dataId string

	BaseInfo BaseInfo

	TraceEs    TraceEsConfig
	SaveEs     TraceEsConfig
	TraceKafka TraceKafkaConfig
}

type BaseInfo struct {
	BkBizId   string
	BkBizName string
	AppId     string
	AppName   string
}

type TraceEsConfig struct {
	IndexName string
	Host      string
	Username  string
	Password  string
}

type TraceKafkaConfig struct {
	Topic    string
	Host     string
	Username string
	Password string
}

var (
	metadataOnce   sync.Once
	centerInstance *MetadataCenter
)

func CreateMetadataCenter() {
	metadataOnce.Do(func() {
		consulClient, err := consul.GetInstance(context.Background())
		if err != nil {
			logger.Errorf("Failed to create consul client. error: %s", err)
		}
		centerInstance = &MetadataCenter{
			mapping: make(map[string]DataIdInfo),
			consul:  *consulClient,
		}
		logger.Infof("Create metadata-center successfully")
	})
}

func (c *MetadataCenter) AddDataIdAndInfo(dataId string, info DataIdInfo) {
	info.dataId = dataId
	c.mapping[dataId] = info
}

func (c *MetadataCenter) AddDataId(dataId string) error {
	info := DataIdInfo{dataId: dataId}
	if err := c.fillInfo(dataId, &info); err != nil {
		return err
	}

	c.mapping[dataId] = info
	logger.Infof("get dataId info successfully, dataId: %s, info: %+v", dataId, info)
	return nil
}

func (c *MetadataCenter) fillInfo(dataId string, info *DataIdInfo) error {
	key := fmt.Sprintf("%s/apm/data_id/%s", config.StorageConsulPathPrefix, dataId)
	bytesData, err := c.consul.Get(key)
	if err != nil {
		return fmt.Errorf("failed to get key: %s from consul. error: %s", key, err)
	}
	if bytesData == nil {
		return fmt.Errorf("failed to get value as key: %s maybe not exist", key)
	}

	var apmInfo ConsulInfo
	if err = json.Unmarshal(bytesData, &apmInfo); err != nil {
		return fmt.Errorf("failed to parse value to ApmInfo, value: %s. error: %s", bytesData, err)
	}
	info.BaseInfo = BaseInfo{
		BkBizId:   strconv.Itoa(apmInfo.BkBizId),
		BkBizName: apmInfo.BkBizName,
		AppId:     strconv.Itoa(apmInfo.AppId),
		AppName:   apmInfo.AppName,
	}
	info.TraceKafka = TraceKafkaConfig{
		Topic:    apmInfo.KafkaInfo.Topic,
		Host:     apmInfo.KafkaInfo.Host,
		Username: apmInfo.KafkaInfo.Username,
		Password: apmInfo.KafkaInfo.Password,
	}
	info.TraceEs = TraceEsConfig{
		IndexName: apmInfo.TraceEsInfo.IndexName,
		Host:      apmInfo.TraceEsInfo.Host,
		Username:  apmInfo.TraceEsInfo.Username,
		Password:  apmInfo.TraceEsInfo.Password,
	}
	info.SaveEs = TraceEsConfig{
		IndexName: apmInfo.SaveEsInfo.IndexName,
		Host:      apmInfo.SaveEsInfo.Host,
		Username:  apmInfo.SaveEsInfo.Username,
		Password:  apmInfo.SaveEsInfo.Password,
	}
	return nil
}

func (c *MetadataCenter) GetKafkaConfig(dataId string) TraceKafkaConfig {
	return c.mapping[dataId].TraceKafka
}

func (c *MetadataCenter) GetTraceEsConfig(dataId string) TraceEsConfig {
	return c.mapping[dataId].TraceEs
}

func (c *MetadataCenter) GetSaveEsConfig(dataId string) TraceEsConfig {
	return c.mapping[dataId].SaveEs
}

func (c *MetadataCenter) GetBaseInfo(dataId string) BaseInfo {
	return c.mapping[dataId].BaseInfo
}

func GetMetadataCenter() *MetadataCenter {
	return centerInstance
}
