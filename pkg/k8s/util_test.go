package k8s

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hbocodelabs/infratest/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func TestGetKubeconfigPathE_NoPath(t *testing.T) {
	t.Parallel()
	homeDir, err := os.UserHomeDir()
	require.Nil(t, err)
	expectedPath := filepath.Join(homeDir, ".kube", "config")

	actualPath, err := getKubeconfigPathE("")

	require.Nil(t, err)
	require.Equal(t, expectedPath, actualPath)
}

func TestGetKubeconfigPathE_Path(t *testing.T) {
	t.Parallel()
	expectedPath := filepath.Join("/tmp", ".kube", "config")

	actualPath, err := getKubeconfigPathE(expectedPath)

	require.Nil(t, err)
	require.Equal(t, expectedPath, actualPath)
}

func withGetClientsetMock(mock kubernetes) (f GetClientsetEOptionsFunc) {
	f = func(opts *GetClientsetOptionsE) error {
		opts.NewForConfig = mock.NewForConfig
		return nil
	}
	return
}

func TestGetEKSClientset_WithTokenHost(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	mockKubernetes := mock.NewMockkubernetes(ctrl)
	clusterEndpoint := "my-cluster.eks.amazonaws.com"
	clusterCAData := "cadata"
	clusterCADataBytes := []byte(clusterCAData)
	tokenData := "token"
	restConfig := &rest.Config{
		Host:        clusterEndpoint,
		BearerToken: tokenData,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: clusterCADataBytes,
		},
	}
	clientset := &k8s.Clientset{}
	ctx := context.Background()
	mockKubernetes.EXPECT().NewForConfig(restConfig).Times(1).Return(clientset, nil)

	actualClientSet, err := GetClientsetE(
		ctx,
		WithGetClientsetEHost(clusterEndpoint),
		withGetClientsetMock(mockKubernetes),
		WithGetClientsetEToken(tokenData),
		WithGetClientsetETLSCAData(clusterCADataBytes),
	)

	require.Nil(t, err)
	require.NotNil(t, actualClientSet)
	assert.Equal(t, clientset, actualClientSet)
}

func TestWithGetClientsetHost(t *testing.T) {
	t.Parallel()
	expectedHostName := "host"
	restConfig := &rest.Config{}
	getClientsetEOptions := &GetClientsetOptionsE{
		RESTConfig: restConfig,
	}

	f := WithGetClientsetEHost(expectedHostName)
	f(getClientsetEOptions)

	assert.Equal(t, expectedHostName, restConfig.Host)

}

func TestWithClientsetToken(t *testing.T) {
	t.Parallel()
	expectedTokenData := "token"
	restConfig := &rest.Config{}
	getClientsetEOptions := &GetClientsetOptionsE{
		RESTConfig: restConfig,
	}

	f := WithGetClientsetEToken(expectedTokenData)
	f(getClientsetEOptions)

	assert.Equal(t, expectedTokenData, restConfig.BearerToken)
}

func TestWithClientsetCA(t *testing.T) {
	t.Parallel()
	expectedCAData := "CA"
	expectedCADataBytes := []byte(expectedCAData)
	restConfig := &rest.Config{}
	getClientsetEOptions := &GetClientsetOptionsE{
		RESTConfig: restConfig,
	}

	f := WithGetClientsetETLSCAData(expectedCADataBytes)
	f(getClientsetEOptions)

	assert.Equal(t, expectedCADataBytes, restConfig.TLSClientConfig.CAData)
}
