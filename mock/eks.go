// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/aws/eks.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	eks "github.com/aws/aws-sdk-go-v2/service/eks"
	gomock "github.com/golang/mock/gomock"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
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

// MockGenerator is a mock of Generator interface.
type MockGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockGeneratorMockRecorder
}

// MockGeneratorMockRecorder is the mock recorder for MockGenerator.
type MockGeneratorMockRecorder struct {
	mock *MockGenerator
}

// NewMockGenerator creates a new mock instance.
func NewMockGenerator(ctrl *gomock.Controller) *MockGenerator {
	mock := &MockGenerator{ctrl: ctrl}
	mock.recorder = &MockGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenerator) EXPECT() *MockGeneratorMockRecorder {
	return m.recorder
}

// GetWithOptions mocks base method.
func (m *MockGenerator) GetWithOptions(arg0 *token.GetTokenOptions) (token.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithOptions", arg0)
	ret0, _ := ret[0].(token.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithOptions indicates an expected call of GetWithOptions.
func (mr *MockGeneratorMockRecorder) GetWithOptions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithOptions", reflect.TypeOf((*MockGenerator)(nil).GetWithOptions), arg0)
}

// MockKubernetes is a mock of Kubernetes interface.
type MockKubernetes struct {
	ctrl     *gomock.Controller
	recorder *MockKubernetesMockRecorder
}

// MockKubernetesMockRecorder is the mock recorder for MockKubernetes.
type MockKubernetesMockRecorder struct {
	mock *MockKubernetes
}

// NewMockKubernetes creates a new mock instance.
func NewMockKubernetes(ctrl *gomock.Controller) *MockKubernetes {
	mock := &MockKubernetes{ctrl: ctrl}
	mock.recorder = &MockKubernetesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubernetes) EXPECT() *MockKubernetesMockRecorder {
	return m.recorder
}

// NewForConfig mocks base method.
func (m *MockKubernetes) NewForConfig(arg0 *rest.Config) (*kubernetes.Clientset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewForConfig", arg0)
	ret0, _ := ret[0].(*kubernetes.Clientset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewForConfig indicates an expected call of NewForConfig.
func (mr *MockKubernetesMockRecorder) NewForConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewForConfig", reflect.TypeOf((*MockKubernetes)(nil).NewForConfig), arg0)
}
