package aws

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/golang/mock/gomock"
	"github.com/hbocodelabs/infratest/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

func TestGetEKSAuthE(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockEKSClient(ctrl)
	clusterName := "my-cluster"
	clusterEndpoint := "my-cluster.eks.amazonaws.com"
	clusterCAData := "cadata"
	clusterCADataBytes := []byte(clusterCAData)
	clusterCADataEncoded := base64.StdEncoding.EncodeToString([]byte(clusterCAData))
	describeClusterInput := &eks.DescribeClusterInput{
		Name: &clusterName,
	}
	ctx := context.Background()
	describeClusterOutput := &eks.DescribeClusterOutput{
		Cluster: &types.Cluster{
			Endpoint: &clusterEndpoint,
			CertificateAuthority: &types.Certificate{
				Data: &clusterCADataEncoded,
			},
		},
	}
	input := &GetEKSTokenInput{
		ClusterName: clusterName,
	}
	client.EXPECT().DescribeCluster(ctx, describeClusterInput).Times(1).Return(describeClusterOutput, nil)
	fakeTest := &testing.T{}

	output, err := GetEKSAuthE(fakeTest, ctx, client, input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, clusterEndpoint, output.Endpoint)
	assert.Equal(t, clusterCADataBytes, output.CAData)
}

func TestGetEKSClientset(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGenerator := mock.NewMockGenerator(ctrl)
	mockKubernetes := mock.NewMockKubernetes(ctrl)
	clusterName := "my-cluster"
	clusterEndpoint := "my-cluster.eks.amazonaws.com"
	clusterCAData := "cadata"
	clusterCADataBytes := []byte(clusterCAData)
	tokenData := "token"
	tokenObj := &token.Token{
		Token: tokenData,
	}
	getTokenOpts := &token.GetTokenOptions{
		ClusterID: clusterName,
	}
	restConfig := &rest.Config{
		Host:        clusterEndpoint,
		BearerToken: tokenData,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: clusterCADataBytes,
		},
	}
	clientset := &kubernetes.Clientset{}
	//clientset, err := kubernetes.NewForConfig(restConfig)
	//require.Nil(t, err)
	fakeTest := &testing.T{}
	ctx := context.Background()
	getEKSClientEInput := &GetEKSClientsetInput{
		ClusterName:     clusterName,
		ClusterEndpoint: clusterEndpoint,
		ClusterCAData:   clusterCADataBytes,
		GetWithOptions:  mockGenerator.GetWithOptions,
		NewConfigFunc:   mockKubernetes.NewForConfig,
	}
	mockGenerator.EXPECT().GetWithOptions(getTokenOpts).Times(1).Return(tokenObj, nil)
	mockKubernetes.EXPECT().NewForConfig(restConfig).Times(1).Return(clientset, nil)

	actualClientSet, err := GetEKSClientsetE(fakeTest, ctx, getEKSClientEInput)

	require.Nil(t, err)
	require.NotNil(t, actualClientSet)
	assert.Equal(t, clientset, actualClientSet)
}
