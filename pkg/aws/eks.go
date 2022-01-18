package aws

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/stretchr/testify/require"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type EKSClient interface {
	DescribeCluster(context.Context, *eks.DescribeClusterInput, ...*eks.Options) (*eks.DescribeClusterOutput, error)
}

type GetEKSTokenInput struct {
	ClusterName string
}

type GetEKSTokenOutput struct {
	Token    string
	Endpoint string
	CAData   []byte
}

func GetEKSClusterE(t *testing.T, ctx context.Context, client EKSClient, input *GetEKSTokenInput) (output *GetEKSTokenOutput, err error) {
	describeClusterInput := &eks.DescribeClusterInput{
		Name: &input.ClusterName,
	}
	describeOutput, err := client.DescribeCluster(ctx, describeClusterInput)
	require.Nil(t, err)
	require.NotNil(t, describeOutput)
	output = &GetEKSTokenOutput{}
	output.Endpoint = *describeOutput.Cluster.Endpoint
	output.CAData, err = base64.StdEncoding.DecodeString(*describeOutput.Cluster.CertificateAuthority.Data)
	return
}

type GetEKSClientsetInput struct {
	ClusterName     string
	ClusterEndpoint string
	GetWithOptions  func(*token.GetTokenOptions) (*token.Token, error)
	NewConfigFunc   func(*rest.Config) (*kubernetes.Clientset, error)
	ClusterCAData   []byte
}

type Generator interface {
	GetWithOptions(*token.GetTokenOptions) (*token.Token, error)
}

type Kubernetes interface {
	NewForConfig(*rest.Config) (*kubernetes.Clientset, error)
}

func GetEKSClientsetE(t *testing.T, ctx context.Context, input *GetEKSClientsetInput) (clientset *kubernetes.Clientset, err error) {
	opts := &token.GetTokenOptions{
		ClusterID: input.ClusterName,
	}
	token, err := input.GetWithOptions(opts)
	require.Nil(t, err)

	clientset, err = input.NewConfigFunc(&rest.Config{
		Host:        input.ClusterEndpoint,
		BearerToken: token.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: input.ClusterCAData,
		},
	})
	require.Nil(t, err)
	return
}
