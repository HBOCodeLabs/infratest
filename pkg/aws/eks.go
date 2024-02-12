package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/eks"

	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type EKSClient interface {
	DescribeCluster(context.Context, *eks.DescribeClusterInput, ...func(*eks.Options)) (*eks.DescribeClusterOutput, error)
}

// GetEKSClusterEOptions is a struct for use with functional options for the GetEKSClusterE method.
type GetEKSClusterEOptions struct {
	// Options that are passed to the underlying DescribeCluster method.
	EKSOptions []func(*eks.Options)
}

// GetEKSClusterEOptionsFunc is a type used for functional options for the GetEKSClusterE method.
type GetEKSClusterEOptionsFunc func(GetEKSClusterEOptions) error

type GetEKSClusterOutput struct {
	Endpoint string
	CAData   []byte
}

/*
GetEKSClusterE returns some metadata about the specified EKS cluster, such as the endpoint and the CA certificate information.
It must be passed an AWS SDK v2 [EKS client object](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/eks#Client).
*/
func GetEKSClusterE(ctx context.Context, client EKSClient, clusterName string, optFns ...GetEKSClusterEOptionsFunc) (output *GetEKSClusterOutput, err error) {
	describeClusterInput := &eks.DescribeClusterInput{
		Name: &clusterName,
	}
	opts := GetEKSClusterEOptions{
		EKSOptions: []func(*eks.Options){},
	}

	for _, f := range optFns {
		err := f(opts)
		if err != nil {
			return nil, err
		}
	}

	describeOutput, err := client.DescribeCluster(ctx, describeClusterInput, opts.EKSOptions...)
	if err != nil {
		return
	}
	output = &GetEKSClusterOutput{}
	output.Endpoint = *describeOutput.Cluster.Endpoint
	output.CAData, err = base64.StdEncoding.DecodeString(*describeOutput.Cluster.CertificateAuthority.Data)
	return
}

// generator is an interface used for mocking the [generator interface](https://pkg.go.dev/sigs.k8s.io/aws-iam-authenticator@v0.5.3/pkg/token#generator)
// from the `aws-iam-authenticator/token` package.
type generator interface {
	GetWithOptions(*token.GetTokenOptions) (token.Token, error)
}

type GetEKSTokenEOptions struct {
	// The object used for generating the token. Generally this should only be specified in the context of tests.
	Generator generator
	// The input object passed to the GetWithOptions method.
	GetTokenOptions *token.GetTokenOptions
}

// GetEKSTokenEOptionsFunc is a type for the functional options of the GetEKSTokenE method.
type GetEKSTokenEOptionsFunc func(*GetEKSTokenEOptions) error

// GetEKSTokenE generates a new bearer token for authenticating with EKS clusters.
// It assumes you have AWS credentials configured in your environment in accordance with
// the [`aws-iam-authenticator` guidelines](https://pkg.go.dev/sigs.k8s.io/aws-iam-authenticator@v0.5.3#readme-specifying-credentials-using-aws-profiles).
// You can alter that configuring by passing in functional options that modify the GetTokenOptions object.
func GetEKSTokenE(ctx context.Context, clusterName string, opts ...func(*GetEKSTokenEOptions) error) (tkn token.Token, err error) {
	getTokenOpts := &token.GetTokenOptions{
		ClusterID: clusterName,
	}
	generator, err := token.NewGenerator(true, false)
	if err != nil {
		return
	}
	getEKSTokenEOptions := &GetEKSTokenEOptions{
		Generator:       generator,
		GetTokenOptions: getTokenOpts,
	}

	for _, fn := range opts {
		err = fn(getEKSTokenEOptions)
		if err != nil {
			return
		}
	}

	tkn, err = getEKSTokenEOptions.Generator.GetWithOptions(getEKSTokenEOptions.GetTokenOptions)
	return
}
