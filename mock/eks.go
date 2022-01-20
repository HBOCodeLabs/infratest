// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/aws/eks.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	eks "github.com/aws/aws-sdk-go-v2/service/eks"
	gomock "github.com/golang/mock/gomock"
	token "sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

// MockEKSClient is a mock of EKSClient interface.
type MockEKSClient struct {
	ctrl     *gomock.Controller
	recorder *MockEKSClientMockRecorder
}

// MockEKSClientMockRecorder is the mock recorder for MockEKSClient.
type MockEKSClientMockRecorder struct {
	mock *MockEKSClient
}

// NewMockEKSClient creates a new mock instance.
func NewMockEKSClient(ctrl *gomock.Controller) *MockEKSClient {
	mock := &MockEKSClient{ctrl: ctrl}
	mock.recorder = &MockEKSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEKSClient) EXPECT() *MockEKSClientMockRecorder {
	return m.recorder
}

// DescribeCluster mocks base method.
func (m *MockEKSClient) DescribeCluster(arg0 context.Context, arg1 *eks.DescribeClusterInput, arg2 ...*eks.Options) (*eks.DescribeClusterOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeCluster", varargs...)
	ret0, _ := ret[0].(*eks.DescribeClusterOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeCluster indicates an expected call of DescribeCluster.
func (mr *MockEKSClientMockRecorder) DescribeCluster(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeCluster", reflect.TypeOf((*MockEKSClient)(nil).DescribeCluster), varargs...)
}

// Mockgenerator is a mock of generator interface.
type Mockgenerator struct {
	ctrl     *gomock.Controller
	recorder *MockgeneratorMockRecorder
}

// MockgeneratorMockRecorder is the mock recorder for Mockgenerator.
type MockgeneratorMockRecorder struct {
	mock *Mockgenerator
}

// NewMockgenerator creates a new mock instance.
func NewMockgenerator(ctrl *gomock.Controller) *Mockgenerator {
	mock := &Mockgenerator{ctrl: ctrl}
	mock.recorder = &MockgeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockgenerator) EXPECT() *MockgeneratorMockRecorder {
	return m.recorder
}

// GetWithOptions mocks base method.
func (m *Mockgenerator) GetWithOptions(arg0 *token.GetTokenOptions) (token.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithOptions", arg0)
	ret0, _ := ret[0].(token.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithOptions indicates an expected call of GetWithOptions.
func (mr *MockgeneratorMockRecorder) GetWithOptions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithOptions", reflect.TypeOf((*Mockgenerator)(nil).GetWithOptions), arg0)
}
