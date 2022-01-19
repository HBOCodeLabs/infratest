package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/eks"

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

func GetEKSClusterE(ctx context.Context, client EKSClient, input *GetEKSTokenInput) (output *GetEKSTokenOutput, err error) {
	describeClusterInput := &eks.DescribeClusterInput{
		Name: &input.ClusterName,
	}
	describeOutput, err := client.DescribeCluster(ctx, describeClusterInput)
	if err != nil {
		return
	}
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

// GetEKSClientsetE returns a Kuberenets client-go Clientset object that is set up for connectivity to
// a specified EKS cluster name.
func GetEKSClientsetE(ctx context.Context, input *GetEKSClientsetInput) (clientset *kubernetes.Clientset, err error) {
	opts := &token.GetTokenOptions{
		ClusterID: input.ClusterName,
	}
	token, err := input.GetWithOptions(opts)
	if err != nil {
		return
	}

	clientset, err = input.NewConfigFunc(&rest.Config{
		Host:        input.ClusterEndpoint,
		BearerToken: token.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: input.ClusterCAData,
		},
	})
	return
}
