package k8s

import (
	"context"
	"os"
	"path/filepath"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClientsetOptionsE is used for passing functional options to the GetEKSClientset method.
type GetClientsetOptionsE struct {
	// The method used to get the clientset object. Generally this should only be specified in the context of tests.
	NewForConfig func(*rest.Config) (*k8s.Clientset, error)
	// The input object passed to the NewForConfig method when generating the clientset.
	RESTConfig rest.Config
}

// Kubernetes is an interface used for mocking the Kubernetes client-go package.
type kubernetes interface {
	NewForConfig(*rest.Config) (*k8s.Clientset, error)
}

func getDefaultKubeconfigPathE() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	path = filepath.Join(home, ".kube", "config")

	return
}

func getKubeconfigPathE(kubeconfigPath string) (outPath string, err error) {
	if kubeconfigPath == "" {
		outPath, err = getDefaultKubeconfigPathE()
	} else {
		outPath = kubeconfigPath
	}
	return
}

// GetClientsetEOptionsFunc is the type for the functional options arguments for the GetClientsetE method.
type GetClientsetEOptionsFunc func(*GetClientsetOptionsE) error

// WithGetClientsetEHost sets a host name when invoking the GetClientsetE method.
func WithGetClientsetEHost(host string) (f GetClientsetEOptionsFunc) {
	f = func(opts *GetClientsetOptionsE) error {
		opts.RESTConfig.Host = host
		return nil
	}
	return
}

// WithGetClientsetEToken sets a token when invoking the GetClientsetE method.
func WithGetClientsetEToken(token string) (f GetClientsetEOptionsFunc) {
	f = func(opts *GetClientsetOptionsE) error {
		opts.RESTConfig.BearerToken = token
		return nil
	}
	return
}

// WithGetClientsetETLSCAData sets the expected CA certificate data when invoking the GetClientsetE method.
func WithGetClientsetETLSCAData(caData []byte) (f GetClientsetEOptionsFunc) {
	f = func(gco *GetClientsetOptionsE) error {
		gco.RESTConfig.TLSClientConfig.CAData = caData
		return nil
	}
	return
}

// WithGetClientsetEKubeconfigPath sets the GetClientsetE method to configure from a Kubeconfig file at a particular path.
// This should almost always be called only by itself, not with other `WithGetClientsetE` methods.
func WithGetClientsetEKubeconfigPath(path string) (f GetClientsetEOptionsFunc) {
	f = func(gco *GetClientsetOptionsE) error {
		config, err := clientcmd.BuildConfigFromFlags("", path)
		if err != nil {
			return err
		}
		gco.RESTConfig = *config
		return nil
	}
	return
}

/* GetClientsetE returns a Kuberenets client-go Clientset object with a friendly interface.
 */
func GetClientsetE(ctx context.Context, opts ...GetClientsetEOptionsFunc) (clientset *k8s.Clientset, err error) {
	restConfig := rest.Config{}
	getClientsetEOptions := &GetClientsetOptionsE{
		NewForConfig: k8s.NewForConfig,
		RESTConfig:   restConfig,
	}

	for _, fn := range opts {
		err = fn(getClientsetEOptions)
		if err != nil {
			return
		}
	}

	clientset, err = getClientsetEOptions.NewForConfig(&getClientsetEOptions.RESTConfig)
	return
}
