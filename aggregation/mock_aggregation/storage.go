// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/micky/Documents/playground/go/iot/internal/aggregation/storage.go

// Package mock_aggregation is a generated GoMock package.
package mock_aggregation

import (
	aggregation "github.com/Michaellqa/iot/aggregation"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Write mocks base method
func (m *MockStorage) Write(record aggregation.Record) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", record)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockStorageMockRecorder) Write(record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockStorage)(nil).Write), record)
}

// MockAsyncStorage is a mock of AsyncStorage interface
type MockAsyncStorage struct {
	ctrl     *gomock.Controller
	recorder *MockAsyncStorageMockRecorder
}

// MockAsyncStorageMockRecorder is the mock recorder for MockAsyncStorage
type MockAsyncStorageMockRecorder struct {
	mock *MockAsyncStorage
}

// NewMockAsyncStorage creates a new mock instance
func NewMockAsyncStorage(ctrl *gomock.Controller) *MockAsyncStorage {
	mock := &MockAsyncStorage{ctrl: ctrl}
	mock.recorder = &MockAsyncStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAsyncStorage) EXPECT() *MockAsyncStorageMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockAsyncStorage) Add(r aggregation.Record) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", r)
}

// Add indicates an expected call of Add
func (mr *MockAsyncStorageMockRecorder) Add(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockAsyncStorage)(nil).Add), r)
}

// Wait mocks base method
func (m *MockAsyncStorage) Wait() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Wait")
}

// Wait indicates an expected call of Wait
func (mr *MockAsyncStorageMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockAsyncStorage)(nil).Wait))
}

// Close mocks base method
func (m *MockAsyncStorage) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockAsyncStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAsyncStorage)(nil).Close))
}
