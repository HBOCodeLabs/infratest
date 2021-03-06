// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/aws/iam.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	iam "github.com/aws/aws-sdk-go-v2/service/iam"
	gomock "github.com/golang/mock/gomock"
)

// MockIAMClient is a mock of IAMClient interface.
type MockIAMClient struct {
	ctrl     *gomock.Controller
	recorder *MockIAMClientMockRecorder
}

// MockIAMClientMockRecorder is the mock recorder for MockIAMClient.
type MockIAMClientMockRecorder struct {
	mock *MockIAMClient
}

// NewMockIAMClient creates a new mock instance.
func NewMockIAMClient(ctrl *gomock.Controller) *MockIAMClient {
	mock := &MockIAMClient{ctrl: ctrl}
	mock.recorder = &MockIAMClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAMClient) EXPECT() *MockIAMClientMockRecorder {
	return m.recorder
}

// GetRole mocks base method.
func (m *MockIAMClient) GetRole(arg0 context.Context, arg1 *iam.GetRoleInput, arg2 ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRole", varargs...)
	ret0, _ := ret[0].(*iam.GetRoleOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockIAMClientMockRecorder) GetRole(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockIAMClient)(nil).GetRole), varargs...)
}
