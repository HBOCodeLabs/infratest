package integration_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"

	infrak8s "github.com/hbocodelabs/infratest/pkg/k8s"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/kind/pkg/cluster"
)

func createKINDCluster(name string, version string, kubeconfigPath string) (err error) {
	nodeImage := fmt.Sprintf("kindest/node:v%s", version)
	provider := cluster.NewProvider()
	err = provider.Create(
		name,
		cluster.CreateWithNodeImage(nodeImage),
		cluster.CreateWithKubeconfigPath(kubeconfigPath),
	)
	return
}

func deleteKINDCluster(name string, kubeconfigPath string) (err error) {
	provider := cluster.NewProvider()
	err = provider.Delete(name, kubeconfigPath)
	return
}

func TestAssertJobSucceeds(t *testing.T) {
	clusterName := strings.ToLower(random.UniqueId())
	namespace := strings.ToLower(random.UniqueId())
	kubeConfigPath := filepath.Join(os.TempDir(), clusterName)
	k8sVersion := os.Getenv("K8S_VERSION")
	if k8sVersion == "" {
		k8sVersion = "1.21.1"
	}
	err := createKINDCluster(clusterName, k8sVersion, kubeConfigPath)
	require.Nil(t, err)
	defer deleteKINDCluster(clusterName, kubeConfigPath)

	testCases := []struct {
		name       string
		image      string
		command    []string
		arguments  []string
		testResult bool
	}{
		{
			name:       "Job finishes successfully",
			image:      "ubuntu:20.04",
			command:    []string{"/bin/bash", "-c", "--"},
			arguments:  []string{"sleep 5; exit 0;"},
			testResult: false,
		},
		{
			name:       "Job fails",
			image:      "ubuntu:20.04",
			command:    []string{"/bin/bash", "-c", "--"},
			arguments:  []string{"sleep 5; exit 1;"},
			testResult: true,
		},
	}

	kubectlOptions := &k8s.KubectlOptions{
		ConfigPath: kubeConfigPath,
	}
	err = k8s.CreateNamespaceE(t, kubectlOptions, namespace)
	require.Nil(t, err)
	defer func() {
		err = k8s.DeleteNamespaceE(t, kubectlOptions, namespace)
		require.Nil(t, err)
	}()
	jobClient, err := infrak8s.GetJobClientE(kubeConfigPath, namespace)
	require.Nil(t, err)

	t.Run("TestCases", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				jobName := strings.ToLower(random.UniqueId())
				backoffLimit := int32(1)
				jobSpec := &batchv1.Job{
					ObjectMeta: metav1.ObjectMeta{
						Name:      jobName,
						Namespace: namespace,
					},
					Spec: batchv1.JobSpec{
						BackoffLimit: &backoffLimit,
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								RestartPolicy: corev1.RestartPolicyNever,
								Containers: []corev1.Container{
									corev1.Container{
										Name:    jobName,
										Image:   testCase.image,
										Command: testCase.command,
										Args:    testCase.arguments,
									},
								},
							},
						},
					},
				}
				input := infrak8s.AssertJobSucceedsInput{
					JobSpec: jobSpec,
				}
				ctx := context.Background()
				fakeTest := &testing.T{}
				infrak8s.AssertJobSucceeds(fakeTest, ctx, jobClient, input)
				assert.Equal(t, testCase.testResult, fakeTest.Failed())
			})

		}
	})
}
