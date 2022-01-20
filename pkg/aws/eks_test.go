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
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

func TestGetEKSClusterE(t *testing.T) {
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
	input := &GetEKSClusterInput{
		ClusterName: clusterName,
	}
	client.EXPECT().DescribeCluster(ctx, describeClusterInput).Times(1).Return(describeClusterOutput, nil)

	output, err := GetEKSClusterE(ctx, client, input)

	require.Nil(t, err)
	require.NotNil(t, output)
	assert.Equal(t, clusterEndpoint, output.Endpoint)
	assert.Equal(t, clusterCADataBytes, output.CAData)
}

func withGetEKSTokenEMock(mock generator) (f GetEKSTokenEOptionsFunc) {
	f = func(gee *GetEKSTokenEOptions) error {
		gee.Generator = mock
		return nil
	}
	return
}

func TestGetEKSToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGenerator := mock.NewMockgenerator(ctrl)
	clusterName := "my-cluster"
	tokenData := "token"
	tokenObj := token.Token{
		Token: tokenData,
	}
	getTokenOpts := &token.GetTokenOptions{
		ClusterID: clusterName,
	}
	ctx := context.Background()
	mockGenerator.EXPECT().GetWithOptions(getTokenOpts).Times(1).Return(tokenObj, nil)

	actualToken, err := GetEKSTokenE(ctx, clusterName, withGetEKSTokenEMock(mockGenerator))

	require.Nil(t, err)
	require.NotNil(t, actualToken)
	assert.Equal(t, tokenData, actualToken.Token)
}
