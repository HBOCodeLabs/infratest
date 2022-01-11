package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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

func getClientset(kubeconfigPath string) (client *kubernetes.Clientset, err error) {
	kubeconfigPath, err = getKubeconfigPathE(kubeconfigPath)
	if err != nil {
		return
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return
	}

	client, err = kubernetes.NewForConfig(config)
	return
}
