// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Shopify/sarama (interfaces: Client,ConsumerGroup,ConsumerGroupSession,ConsumerGroupClaim)

// Package testsuite is a generated GoMock package.
package testsuite

import (
	context "context"
	reflect "reflect"

	sarama "github.com/Shopify/sarama"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Brokers mocks base method.
func (m *MockClient) Brokers() []*sarama.Broker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Brokers")
	ret0, _ := ret[0].([]*sarama.Broker)
	return ret0
}

// Brokers indicates an expected call of Brokers.
func (mr *MockClientMockRecorder) Brokers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Brokers", reflect.TypeOf((*MockClient)(nil).Brokers))
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// Closed mocks base method.
func (m *MockClient) Closed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Closed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Closed indicates an expected call of Closed.
func (mr *MockClientMockRecorder) Closed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Closed", reflect.TypeOf((*MockClient)(nil).Closed))
}

// Config mocks base method.
func (m *MockClient) Config() *sarama.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*sarama.Config)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockClientMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockClient)(nil).Config))
}

// Controller mocks base method.
func (m *MockClient) Controller() (*sarama.Broker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Controller")
	ret0, _ := ret[0].(*sarama.Broker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Controller indicates an expected call of Controller.
func (mr *MockClientMockRecorder) Controller() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Controller", reflect.TypeOf((*MockClient)(nil).Controller))
}

// Coordinator mocks base method.
func (m *MockClient) Coordinator(arg0 string) (*sarama.Broker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Coordinator", arg0)
	ret0, _ := ret[0].(*sarama.Broker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Coordinator indicates an expected call of Coordinator.
func (mr *MockClientMockRecorder) Coordinator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Coordinator", reflect.TypeOf((*MockClient)(nil).Coordinator), arg0)
}

// GetOffset mocks base method.
func (m *MockClient) GetOffset(arg0 string, arg1 int32, arg2 int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOffset", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOffset indicates an expected call of GetOffset.
func (mr *MockClientMockRecorder) GetOffset(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOffset", reflect.TypeOf((*MockClient)(nil).GetOffset), arg0, arg1, arg2)
}

// InSyncReplicas mocks base method.
func (m *MockClient) InSyncReplicas(arg0 string, arg1 int32) ([]int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InSyncReplicas", arg0, arg1)
	ret0, _ := ret[0].([]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InSyncReplicas indicates an expected call of InSyncReplicas.
func (mr *MockClientMockRecorder) InSyncReplicas(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InSyncReplicas", reflect.TypeOf((*MockClient)(nil).InSyncReplicas), arg0, arg1)
}

// InitProducerID mocks base method.
func (m *MockClient) InitProducerID() (*sarama.InitProducerIDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitProducerID")
	ret0, _ := ret[0].(*sarama.InitProducerIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InitProducerID indicates an expected call of InitProducerID.
func (mr *MockClientMockRecorder) InitProducerID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitProducerID", reflect.TypeOf((*MockClient)(nil).InitProducerID))
}

// Leader mocks base method.
func (m *MockClient) Leader(arg0 string, arg1 int32) (*sarama.Broker, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leader", arg0, arg1)
	ret0, _ := ret[0].(*sarama.Broker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Leader indicates an expected call of Leader.
func (mr *MockClientMockRecorder) Leader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leader", reflect.TypeOf((*MockClient)(nil).Leader), arg0, arg1)
}

// OfflineReplicas mocks base method.
func (m *MockClient) OfflineReplicas(arg0 string, arg1 int32) ([]int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OfflineReplicas", arg0, arg1)
	ret0, _ := ret[0].([]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OfflineReplicas indicates an expected call of OfflineReplicas.
func (mr *MockClientMockRecorder) OfflineReplicas(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OfflineReplicas", reflect.TypeOf((*MockClient)(nil).OfflineReplicas), arg0, arg1)
}

// Partitions mocks base method.
func (m *MockClient) Partitions(arg0 string) ([]int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Partitions", arg0)
	ret0, _ := ret[0].([]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Partitions indicates an expected call of Partitions.
func (mr *MockClientMockRecorder) Partitions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Partitions", reflect.TypeOf((*MockClient)(nil).Partitions), arg0)
}

// RefreshCoordinator mocks base method.
func (m *MockClient) RefreshCoordinator(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshCoordinator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RefreshCoordinator indicates an expected call of RefreshCoordinator.
func (mr *MockClientMockRecorder) RefreshCoordinator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshCoordinator", reflect.TypeOf((*MockClient)(nil).RefreshCoordinator), arg0)
}

// RefreshMetadata mocks base method.
func (m *MockClient) RefreshMetadata(arg0 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RefreshMetadata", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// RefreshMetadata indicates an expected call of RefreshMetadata.
func (mr *MockClientMockRecorder) RefreshMetadata(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshMetadata", reflect.TypeOf((*MockClient)(nil).RefreshMetadata), arg0...)
}

// Replicas mocks base method.
func (m *MockClient) Replicas(arg0 string, arg1 int32) ([]int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replicas", arg0, arg1)
	ret0, _ := ret[0].([]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Replicas indicates an expected call of Replicas.
func (mr *MockClientMockRecorder) Replicas(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replicas", reflect.TypeOf((*MockClient)(nil).Replicas), arg0, arg1)
}

// Topics mocks base method.
func (m *MockClient) Topics() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Topics")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Topics indicates an expected call of Topics.
func (mr *MockClientMockRecorder) Topics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Topics", reflect.TypeOf((*MockClient)(nil).Topics))
}

// WritablePartitions mocks base method.
func (m *MockClient) WritablePartitions(arg0 string) ([]int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WritablePartitions", arg0)
	ret0, _ := ret[0].([]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WritablePartitions indicates an expected call of WritablePartitions.
func (mr *MockClientMockRecorder) WritablePartitions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WritablePartitions", reflect.TypeOf((*MockClient)(nil).WritablePartitions), arg0)
}

// MockConsumerGroup is a mock of ConsumerGroup interface.
type MockConsumerGroup struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerGroupMockRecorder
}

// MockConsumerGroupMockRecorder is the mock recorder for MockConsumerGroup.
type MockConsumerGroupMockRecorder struct {
	mock *MockConsumerGroup
}

// NewMockConsumerGroup creates a new mock instance.
func NewMockConsumerGroup(ctrl *gomock.Controller) *MockConsumerGroup {
	mock := &MockConsumerGroup{ctrl: ctrl}
	mock.recorder = &MockConsumerGroupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerGroup) EXPECT() *MockConsumerGroupMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConsumerGroup) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConsumerGroupMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConsumerGroup)(nil).Close))
}

// Consume mocks base method.
func (m *MockConsumerGroup) Consume(arg0 context.Context, arg1 []string, arg2 sarama.ConsumerGroupHandler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Consume indicates an expected call of Consume.
func (mr *MockConsumerGroupMockRecorder) Consume(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockConsumerGroup)(nil).Consume), arg0, arg1, arg2)
}

// Errors mocks base method.
func (m *MockConsumerGroup) Errors() <-chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Errors")
	ret0, _ := ret[0].(<-chan error)
	return ret0
}

// Errors indicates an expected call of Errors.
func (mr *MockConsumerGroupMockRecorder) Errors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errors", reflect.TypeOf((*MockConsumerGroup)(nil).Errors))
}

// MockConsumerGroupSession is a mock of ConsumerGroupSession interface.
type MockConsumerGroupSession struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerGroupSessionMockRecorder
}

func (m *MockConsumerGroupSession) Commit() {
	//TODO implement me
	panic("implement me")
}

// MockConsumerGroupSessionMockRecorder is the mock recorder for MockConsumerGroupSession.
type MockConsumerGroupSessionMockRecorder struct {
	mock *MockConsumerGroupSession
}

// NewMockConsumerGroupSession creates a new mock instance.
func NewMockConsumerGroupSession(ctrl *gomock.Controller) *MockConsumerGroupSession {
	mock := &MockConsumerGroupSession{ctrl: ctrl}
	mock.recorder = &MockConsumerGroupSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerGroupSession) EXPECT() *MockConsumerGroupSessionMockRecorder {
	return m.recorder
}

// Claims mocks base method.
func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Claims")
	ret0, _ := ret[0].(map[string][]int32)
	return ret0
}

// Claims indicates an expected call of Claims.
func (mr *MockConsumerGroupSessionMockRecorder) Claims() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Claims", reflect.TypeOf((*MockConsumerGroupSession)(nil).Claims))
}

// Context mocks base method.
func (m *MockConsumerGroupSession) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockConsumerGroupSessionMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockConsumerGroupSession)(nil).Context))
}

// GenerationID mocks base method.
func (m *MockConsumerGroupSession) GenerationID() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerationID")
	ret0, _ := ret[0].(int32)
	return ret0
}

// GenerationID indicates an expected call of GenerationID.
func (mr *MockConsumerGroupSessionMockRecorder) GenerationID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerationID", reflect.TypeOf((*MockConsumerGroupSession)(nil).GenerationID))
}

// MarkMessage mocks base method.
func (m *MockConsumerGroupSession) MarkMessage(arg0 *sarama.ConsumerMessage, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkMessage", arg0, arg1)
}

// MarkMessage indicates an expected call of MarkMessage.
func (mr *MockConsumerGroupSessionMockRecorder) MarkMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessage", reflect.TypeOf((*MockConsumerGroupSession)(nil).MarkMessage), arg0, arg1)
}

// MarkOffset mocks base method.
func (m *MockConsumerGroupSession) MarkOffset(arg0 string, arg1 int32, arg2 int64, arg3 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkOffset", arg0, arg1, arg2, arg3)
}

// MarkOffset indicates an expected call of MarkOffset.
func (mr *MockConsumerGroupSessionMockRecorder) MarkOffset(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkOffset", reflect.TypeOf((*MockConsumerGroupSession)(nil).MarkOffset), arg0, arg1, arg2, arg3)
}

// MemberID mocks base method.
func (m *MockConsumerGroupSession) MemberID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MemberID")
	ret0, _ := ret[0].(string)
	return ret0
}

// MemberID indicates an expected call of MemberID.
func (mr *MockConsumerGroupSessionMockRecorder) MemberID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MemberID", reflect.TypeOf((*MockConsumerGroupSession)(nil).MemberID))
}

// ResetOffset mocks base method.
func (m *MockConsumerGroupSession) ResetOffset(arg0 string, arg1 int32, arg2 int64, arg3 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetOffset", arg0, arg1, arg2, arg3)
}

// ResetOffset indicates an expected call of ResetOffset.
func (mr *MockConsumerGroupSessionMockRecorder) ResetOffset(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetOffset", reflect.TypeOf((*MockConsumerGroupSession)(nil).ResetOffset), arg0, arg1, arg2, arg3)
}

// MockConsumerGroupClaim is a mock of ConsumerGroupClaim interface.
type MockConsumerGroupClaim struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerGroupClaimMockRecorder
}

// MockConsumerGroupClaimMockRecorder is the mock recorder for MockConsumerGroupClaim.
type MockConsumerGroupClaimMockRecorder struct {
	mock *MockConsumerGroupClaim
}

// NewMockConsumerGroupClaim creates a new mock instance.
func NewMockConsumerGroupClaim(ctrl *gomock.Controller) *MockConsumerGroupClaim {
	mock := &MockConsumerGroupClaim{ctrl: ctrl}
	mock.recorder = &MockConsumerGroupClaimMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerGroupClaim) EXPECT() *MockConsumerGroupClaimMockRecorder {
	return m.recorder
}

// HighWaterMarkOffset mocks base method.
func (m *MockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HighWaterMarkOffset")
	ret0, _ := ret[0].(int64)
	return ret0
}

// HighWaterMarkOffset indicates an expected call of HighWaterMarkOffset.
func (mr *MockConsumerGroupClaimMockRecorder) HighWaterMarkOffset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HighWaterMarkOffset", reflect.TypeOf((*MockConsumerGroupClaim)(nil).HighWaterMarkOffset))
}

// InitialOffset mocks base method.
func (m *MockConsumerGroupClaim) InitialOffset() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitialOffset")
	ret0, _ := ret[0].(int64)
	return ret0
}

// InitialOffset indicates an expected call of InitialOffset.
func (mr *MockConsumerGroupClaimMockRecorder) InitialOffset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitialOffset", reflect.TypeOf((*MockConsumerGroupClaim)(nil).InitialOffset))
}

// Messages mocks base method.
func (m *MockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Messages")
	ret0, _ := ret[0].(<-chan *sarama.ConsumerMessage)
	return ret0
}

// Messages indicates an expected call of Messages.
func (mr *MockConsumerGroupClaimMockRecorder) Messages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Messages", reflect.TypeOf((*MockConsumerGroupClaim)(nil).Messages))
}

// Partition mocks base method.
func (m *MockConsumerGroupClaim) Partition() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Partition")
	ret0, _ := ret[0].(int32)
	return ret0
}

// Partition indicates an expected call of Partition.
func (mr *MockConsumerGroupClaimMockRecorder) Partition() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Partition", reflect.TypeOf((*MockConsumerGroupClaim)(nil).Partition))
}

// Topic mocks base method.
func (m *MockConsumerGroupClaim) Topic() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Topic")
	ret0, _ := ret[0].(string)
	return ret0
}

// Topic indicates an expected call of Topic.
func (mr *MockConsumerGroupClaimMockRecorder) Topic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Topic", reflect.TypeOf((*MockConsumerGroupClaim)(nil).Topic))
}
