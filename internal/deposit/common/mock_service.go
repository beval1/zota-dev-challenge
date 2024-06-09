// Code generated by MockGen. DO NOT EDIT.
// Source: internal/deposit/service.go

// Package deposit is a generated GoMock package.
package common

import (
	reflect "reflect"
	"zota-dev-challenge/internal/deposit/shared"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceInterface is a mock of ServiceInterface interface.
type MockServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockServiceInterfaceMockRecorder
}

// MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.
type MockServiceInterfaceMockRecorder struct {
	mock *MockServiceInterface
}

// NewMockServiceInterface creates a new mock instance.
func NewMockServiceInterface(ctrl *gomock.Controller) *MockServiceInterface {
	mock := &MockServiceInterface{ctrl: ctrl}
	mock.recorder = &MockServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceInterface) EXPECT() *MockServiceInterfaceMockRecorder {
	return m.recorder
}

// ProcessDeposit mocks base method.
func (m *MockServiceInterface) ProcessDeposit(r *shared.ClientRequest) (*shared.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessDeposit", r)
	ret0, _ := ret[0].(*shared.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProcessDeposit indicates an expected call of ProcessDeposit.
func (mr *MockServiceInterfaceMockRecorder) ProcessDeposit(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessDeposit", reflect.TypeOf((*MockServiceInterface)(nil).ProcessDeposit), r)
}
