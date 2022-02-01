package integration_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"

	infrak8s "github.com/hbocodelabs/infratest/pkg/k8s"

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
	t.Log("Creating KIND cluster")
	err := createKINDCluster(clusterName, k8sVersion, kubeConfigPath)
	require.Nil(t, err)
	defer func() {
		t.Log("Deleting KIND cluster")
		deleteKINDCluster(clusterName, kubeConfigPath)
	}()

	testCases := []struct {
		name             string
		image            string
		command          []string
		arguments        []string
		timeoutPeriod    time.Duration
		testFailExpected bool
	}{
		{
			name:             "Job succeeds",
			image:            "ubuntu:20.04",
			command:          []string{"/bin/bash", "-c", "--"},
			arguments:        []string{"sleep 5; exit 0;"},
			timeoutPeriod:    5 * time.Minute,
			testFailExpected: false,
		},
		{
			name:             "Job fails",
			image:            "ubuntu:20.04",
			command:          []string{"/bin/bash", "-c", "--"},
			arguments:        []string{"sleep 5; exit 1;"},
			timeoutPeriod:    5 * time.Minute,
			testFailExpected: true,
		},
		{
			name:             "Context timeout expired",
			image:            "ubuntu:20.04",
			command:          []string{"/bin/bash", "-c", "--"},
			arguments:        []string{"sleep 60; exit 1;"},
			timeoutPeriod:    30 * time.Second,
			testFailExpected: true,
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
	ctx := context.Background()
	clientset, err := infrak8s.GetClientsetE(ctx, infrak8s.WithGetClientsetEKubeconfigPath(kubeConfigPath))
	require.Nil(t, err)
	jobClient := clientset.BatchV1().Jobs(namespace)
	require.Nil(t, err)

	t.Run("TestCases", func(t *testing.T) {
		for _, testCase := range testCases {
			// Yes, this is necessary (and works!). See the example at https://go.dev/blog/subtests.
			testCase := testCase

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

				ctx, cancelFunc := context.WithTimeout(context.Background(), testCase.timeoutPeriod)
				defer cancelFunc()
				fakeTest := &testing.T{}
				infrak8s.AssertJobSucceeds(fakeTest, ctx, jobClient, jobSpec)
				assert.Equal(t, testCase.testFailExpected, fakeTest.Failed())
			})

		}
	})
}
