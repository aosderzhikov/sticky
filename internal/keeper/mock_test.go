// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -source=handler.go -destination=mock_test.go -package=keeper
//

// Package keeper is a generated GoMock package.
package keeper

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockService) Delete(key string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", key)
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), key)
}

// Get mocks base method.
func (m *MockService) Get(key string) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockServiceMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), key)
}

// Set mocks base method.
func (m *MockService) Set(key string, value []byte, ttl time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, value, ttl)
}

// Set indicates an expected call of Set.
func (mr *MockServiceMockRecorder) Set(key, value, ttl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockService)(nil).Set), key, value, ttl)
}
