// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Version = "v1beta1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DataID struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec DataIDSpec `json:"spec,omitempty"`
}

type DataIDSpec struct {
	DataID           int               `json:"dataID,omitempty"`
	MonitorResource  MonitorResource   `json:"monitorResource,omitempty"`
	Report           Report            `json:"report,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	MetricReplace    map[string]string `json:"metricReplace,omitempty"`
	DimensionReplace map[string]string `json:"dimensionReplace,omitempty"`
}

type MonitorResource struct {
	// 资源所属命名空间
	NameSpace string `json:"namespace,omitempty"`
	// 资源所属类型 serviceMonitor podMonitor probe
	Kind string `json:"kind,omitempty"`
	// 资源定义的name
	Name string `json:"name,omitempty"`
}

type Report struct {
	MaxMetric   int `json:"maxMetric,omitempty"`
	MaxSizeByte int `json:"maxSizeByte,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DataIDList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []*DataID `json:"items"`
}
