package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getDefaultKubeconfigPath() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	path = filepath.Join(home, ".kube", "config")

	return
}

func getClient(kubeconfigPath string) (client *kubernetes.Clientset, err error) {
	if kubeconfigPath == "" {
		kubeconfigPath, err = getDefaultKubeconfigPath()
		if err != nil {
			return
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return
	}

	client, err = kubernetes.NewForConfig(config)
	return
}
