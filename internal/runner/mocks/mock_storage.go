// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/runner/storage.go

// Package mock_runner is a generated GoMock package.
package mock_runner

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/sonyamoonglade/delivery-service/internal/entity"
	dto "github.com/sonyamoonglade/delivery-service/internal/runner/transport/dto"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockStorage) All(ctx context.Context) ([]*entity.Runner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", ctx)
	ret0, _ := ret[0].([]*entity.Runner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockStorageMockRecorder) All(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockStorage)(nil).All), ctx)
}

// Ban mocks base method.
func (m *MockStorage) Ban(phoneNumber string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ban", phoneNumber)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ban indicates an expected call of Ban.
func (mr *MockStorageMockRecorder) Ban(phoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ban", reflect.TypeOf((*MockStorage)(nil).Ban), phoneNumber)
}

// BeginWork mocks base method.
func (m *MockStorage) BeginWork(dto dto.RunnerBeginWorkDto) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginWork", dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// BeginWork indicates an expected call of BeginWork.
func (mr *MockStorageMockRecorder) BeginWork(dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginWork", reflect.TypeOf((*MockStorage)(nil).BeginWork), dto)
}

// GetByTelegramId mocks base method.
func (m *MockStorage) GetByTelegramId(tgUsrID int64) (*entity.Runner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTelegramId", tgUsrID)
	ret0, _ := ret[0].(*entity.Runner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTelegramId indicates an expected call of GetByTelegramId.
func (mr *MockStorageMockRecorder) GetByTelegramId(tgUsrID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTelegramId", reflect.TypeOf((*MockStorage)(nil).GetByTelegramId), tgUsrID)
}

// IsKnownByTelegramId mocks base method.
func (m *MockStorage) IsKnownByTelegramId(usrID int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsKnownByTelegramId", usrID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsKnownByTelegramId indicates an expected call of IsKnownByTelegramId.
func (mr *MockStorageMockRecorder) IsKnownByTelegramId(usrID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsKnownByTelegramId", reflect.TypeOf((*MockStorage)(nil).IsKnownByTelegramId), usrID)
}

// IsRunner mocks base method.
func (m *MockStorage) IsRunner(usrPhoneNumber string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunner", usrPhoneNumber)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRunner indicates an expected call of IsRunner.
func (mr *MockStorageMockRecorder) IsRunner(usrPhoneNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunner", reflect.TypeOf((*MockStorage)(nil).IsRunner), usrPhoneNumber)
}

// Register mocks base method.
func (m *MockStorage) Register(dto dto.RegisterRunnerDto) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", dto)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockStorageMockRecorder) Register(dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockStorage)(nil).Register), dto)
}
