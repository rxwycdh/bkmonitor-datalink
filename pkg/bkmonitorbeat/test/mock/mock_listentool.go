// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/TencentBlueKing/bkmonitorbeat/task/custom (interfaces: ListenTool)

// Package mock_custom is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	define "github.com/TencentBlueKing/bkmonitor-datalink/pkg/bkmonitorbeat/define"
)

// MockListenTool is a mock of ListenTool interface
type MockListenTool struct {
	ctrl     *gomock.Controller
	recorder *MockListenToolMockRecorder
}

// MockListenToolMockRecorder is the mock recorder for MockListenTool
type MockListenToolMockRecorder struct {
	mock *MockListenTool
}

// NewMockListenTool creates a new mock instance
func NewMockListenTool(ctrl *gomock.Controller) *MockListenTool {
	mock := &MockListenTool{ctrl: ctrl}
	mock.recorder = &MockListenToolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockListenTool) EXPECT() *MockListenToolMockRecorder {
	return m.recorder
}

// Init mocks base method
func (m *MockListenTool) Init(arg0 define.TaskConfig, arg1 chan<- define.Event) error {
	ret := m.ctrl.Call(m, "Init", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockListenToolMockRecorder) Init(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockListenTool)(nil).Init), arg0, arg1)
}

// Start mocks base method
func (m *MockListenTool) Start() error {
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockListenToolMockRecorder) Start() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockListenTool)(nil).Start))
}

// Stop mocks base method
func (m *MockListenTool) Stop(arg0 context.Context) error {
	ret := m.ctrl.Call(m, "Stop", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockListenToolMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockListenTool)(nil).Stop), arg0)
}
