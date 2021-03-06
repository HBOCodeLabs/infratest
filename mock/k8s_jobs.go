// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/k8s/jobs.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/batch/v1"
	v10 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockJobClient is a mock of JobClient interface.
type MockJobClient struct {
	ctrl     *gomock.Controller
	recorder *MockJobClientMockRecorder
}

// MockJobClientMockRecorder is the mock recorder for MockJobClient.
type MockJobClientMockRecorder struct {
	mock *MockJobClient
}

// NewMockJobClient creates a new mock instance.
func NewMockJobClient(ctrl *gomock.Controller) *MockJobClient {
	mock := &MockJobClient{ctrl: ctrl}
	mock.recorder = &MockJobClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJobClient) EXPECT() *MockJobClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockJobClient) Create(arg0 context.Context, arg1 *v1.Job, arg2 v10.CreateOptions) (*v1.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockJobClientMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockJobClient)(nil).Create), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockJobClient) Get(arg0 context.Context, arg1 string, arg2 v10.GetOptions) (*v1.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockJobClientMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockJobClient)(nil).Get), arg0, arg1, arg2)
}
