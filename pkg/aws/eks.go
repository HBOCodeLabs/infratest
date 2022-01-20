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

type GetEKSClusterInput struct {
	ClusterName string
}

type GetEKSClusterOutput struct {
	Token    string
	Endpoint string
	CAData   []byte
}

/*	GetEKSClusterE returns some metadata about the specified EKS cluster, such as the endpoint and the CA certificate information.
		It must be passed an AWS SDK v2 [EKS client object](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/eks#Client).
   	This is meant to be used in tandem with the GetEKSClientset method, such as:
   	```go
	 	getClusterInput := &GetEKSClusterInput{
	 		ClusterName: clusterName,
	 	}
	 	cluster, err := GetEKSClusterE(ctx, client, getClusterInput)
		if err != nil {
			return nil, err
		}
		generator, err := token.NewGenerator
		getClientsetInput := &GetEKSClientsetInput{
			ClusterName: clusterName,
			ClusterEndpoint: output.Endpoint,
			ClusterCAData: output.CAData,
			GetWithOptions: token.GetWithOptions,
			NewForConfig: kubernetes.NewForConfig,
		}
		clientset, err := GetEKSClientsetE(ctx, getClientsetInput)
		```
*/
func GetEKSClusterE(ctx context.Context, client EKSClient, input *GetEKSClusterInput) (output *GetEKSClusterOutput, err error) {
	describeClusterInput := &eks.DescribeClusterInput{
		Name: &input.ClusterName,
	}
	describeOutput, err := client.DescribeCluster(ctx, describeClusterInput)
	if err != nil {
		return
	}
	output = &GetEKSClusterOutput{}
	output.Endpoint = *describeOutput.Cluster.Endpoint
	output.CAData, err = base64.StdEncoding.DecodeString(*describeOutput.Cluster.CertificateAuthority.Data)
	return
}

/*	GetEKSClientsetInput is used as an input object for the GetEKSClientsetE method.
 */
type GetEKSClientsetInput struct {
	// The name of the EKS cluster.
	ClusterName string
	// The API endpoint of the cluster.
	ClusterEndpoint string
	// The (not base64 encoded) string data for the CA certificate used by the cluster.
	ClusterCAData []byte
}

// GetEKSClientsetOptions is used for passing functional options to the GetEKSClientset method.
type GetEKSClientsetOptions struct {
	// The object used for generating the token. Generally this should only be specified in the context of tests.
	Generator generator
	// The method used to get the clientset object. Generally this should only be specified in the context of tests.
	NewForConfig func(*rest.Config) (*kubernetes.Clientset, error)
	// The input object passed to the GetWithOptions method.
	GetTokenOptions *token.GetTokenOptions
	// The input object passed to the NewForConfig method when generating the clientset.
	RESTConfig *rest.Config
}

// generator is an interface used for mocking the [generator interface](https://pkg.go.dev/sigs.k8s.io/aws-iam-authenticator@v0.5.3/pkg/token#generator)
// from the `aws-iam-authenticator/token` package.
type generator interface {
	GetWithOptions(*token.GetTokenOptions) (token.Token, error)
}

// Kubernetes is an interface used for mocking the Kubernetes client-go package.
type Kubernetes interface {
	NewForConfig(*rest.Config) (*kubernetes.Clientset, error)
}

/* GetEKSClientsetE returns a Kuberenets client-go Clientset object that is set up for connectivity to
   a specified EKS cluster name.
*/
func GetEKSClientsetE(ctx context.Context, input *GetEKSClientsetInput, opts ...func(*GetEKSClientsetOptions) error) (clientset *kubernetes.Clientset, err error) {
	getTokenOpts := &token.GetTokenOptions{
		ClusterID: input.ClusterName,
	}
	generator, err := token.NewGenerator(true, false)
	if err != nil {
		return
	}
	getEKSClientsetOptions := &GetEKSClientsetOptions{
		Generator:    generator,
		NewForConfig: kubernetes.NewForConfig,
	}

	for _, fn := range opts {
		err = fn(getEKSClientsetOptions)
		if err != nil {
			return
		}
	}

	token, err := getEKSClientsetOptions.Generator.GetWithOptions(getTokenOpts)
	if err != nil {
		return
	}

	restConfig := &rest.Config{
		Host:        input.ClusterEndpoint,
		BearerToken: token.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: input.ClusterCAData,
		},
	}
	// I'm calling this again since it potentially may alter the Config object.
	// There is probably a better way to do this.
	for _, fn := range opts {
		err = fn(getEKSClientsetOptions)
		if err != nil {
			return
		}
	}

	clientset, err = getEKSClientsetOptions.NewForConfig(restConfig)
	return
}

type GetEKSTokenEOptions struct {
	// The object used for generating the token. Generally this should only be specified in the context of tests.
	Generator generator
	// The input object passed to the GetWithOptions method.
	GetTokenOptions *token.GetTokenOptions
}

// GetEKSTokenE generates a new bearer token for authenticating with EKS clusters.
// It assumes you have AWS credentials configured in your environment in accordance with
// the [`aws-iam-authenticator` guidelines](https://pkg.go.dev/sigs.k8s.io/aws-iam-authenticator@v0.5.3#readme-specifying-credentials-using-aws-profiles).
// You can alter that configuring by passing in functional options that modify the GetTokenOptions object.
func GetEKSTokenE(clusterName string, opts ...func(*GetEKSTokenEOptions) error) (tkn token.Token, err error) {
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

	tkn, err = getEKSTokenEOptions.Generator.GetWithOptions(getTokenOpts)
	return
}
