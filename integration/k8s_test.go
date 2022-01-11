package integration_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"

	infrak8s "github.com/hbocodelabs/infratest/pkg/k8s"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/kind/pkg/cluster"
)

func TestAssertJobSucceeds(t *testing.T) {
	clusterName := strings.ToLower(random.UniqueId())
	namespace := strings.ToLower(random.UniqueId())
	jobName := strings.ToLower(random.UniqueId())
	kubeConfigPath, err := k8s.GetKubeConfigPathE(t)
	require.Nil(t, err)

	provider := cluster.NewProvider()
	err = provider.Create(clusterName)
	require.Nil(t, err)
	defer provider.Delete(clusterName, kubeConfigPath)

	kubectlOptions := &k8s.KubectlOptions{
		ConfigPath: kubeConfigPath,
	}
	err = k8s.CreateNamespaceE(t, kubectlOptions, namespace)
	require.Nil(t, err)
	defer func() {
		err = k8s.DeleteNamespaceE(t, kubectlOptions, namespace)
		require.Nil(t, err)
	}()
	jobClient, err := infrak8s.GetJobClient(kubeConfigPath, namespace)
	require.Nil(t, err)
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						corev1.Container{
							Name:  jobName,
							Image: "ubuntu:latest",
							Command: []string{
								"/bin/bash",
								"-c",
								"--",
							},
							Args: []string{
								"sleep 5; exit 0;",
							},
						},
					},
				},
			},
		},
	}
	ctx := context.Background()

	err = infrak8s.AssertJobSucceeds(t, ctx, jobClient, jobSpec)
	require.Nil(t, err)

}
