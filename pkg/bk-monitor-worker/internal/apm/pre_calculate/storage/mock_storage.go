// Code generated by MockGen. DO NOT EDIT.
// Source: internal/apm/pre_calculate/storage/storage.go

// Package store is a generated GoMock package.
package storage

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// Exist mocks base method.
func (m *MockBackend) Exist(req ExistRequest) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exist", req)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exist indicates an expected call of Exist.
func (mr *MockBackendMockRecorder) Exist(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exist", reflect.TypeOf((*MockBackend)(nil).Exist), req)
}

// Query mocks base method.
func (m *MockBackend) Query(queryRequest QueryRequest) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", queryRequest)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockBackendMockRecorder) Query(queryRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockBackend)(nil).Query), queryRequest)
}

// ReceiveSaveRequest mocks base method.
func (m *MockBackend) ReceiveSaveRequest(errorReceiveChan chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReceiveSaveRequest", errorReceiveChan)
}

// ReceiveSaveRequest indicates an expected call of ReceiveSaveRequest.
func (mr *MockBackendMockRecorder) ReceiveSaveRequest(errorReceiveChan interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveSaveRequest", reflect.TypeOf((*MockBackend)(nil).ReceiveSaveRequest), errorReceiveChan)
}

// Run mocks base method.
func (m *MockBackend) Run(errorReceiveChan chan<- error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", errorReceiveChan)
}

// Run indicates an expected call of Run.
func (mr *MockBackendMockRecorder) Run(errorReceiveChan interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockBackend)(nil).Run), errorReceiveChan)
}

// SaveRequest mocks base method.
func (m *MockBackend) SaveRequest() chan<- SaveRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRequest")
	ret0, _ := ret[0].(chan<- SaveRequest)
	return ret0
}

// SaveRequest indicates an expected call of SaveRequest.
func (mr *MockBackendMockRecorder) SaveRequest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRequest", reflect.TypeOf((*MockBackend)(nil).SaveRequest))
}
