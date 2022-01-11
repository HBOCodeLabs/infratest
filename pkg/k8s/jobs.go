package k8s

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	apiv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobClient interface {
	Create(context.Context, *apiv1.Job, metav1.CreateOptions) (*apiv1.Job, error)
	Get(context.Context, string, metav1.GetOptions) (*apiv1.Job, error)
}

// GetJobClient returns a JobClient object given a path to a kubeconfig file and
// a namespace name. If a blank string is passed as the kubeconfig path, then
// the method will use the path $HOME/.kube/config, where $HOME is the user's home
// directory as determined by the OS.
func GetJobClient(kubeconfigPath string, namespace string) (client JobClient, err error) {
	clientset, err := getClientset(kubeconfigPath)
	if err != nil {
		return
	}
	client = clientset.BatchV1().Jobs(namespace)
	return
}

/*
AssertJobSucceeds will start a Kubernetes Job using the provided client and spec, then after it completes will either
fail the passed test (if the job fails) or pass the test if the job succeeds. It should be passed a
JobClient object and a Job object.

Example:
```
// Gets a client for the default Kubeconfig path and the 'default' namespace
client := GetJobClient("", "default")
ctx := context.Background()
jobName := "job"
job := &batchv1.Job{
	ObjectMeta: metav1.ObjectMeta{
		Name: jobName,
	},
	Spec: batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "my-container",
						Image: "my-image",
					},
				},
			},
		},
	},
}
err := AssertJobSucceeds(t, ctx, client, job)
```
*/
func AssertJobSucceeds(t *testing.T, ctx context.Context, jobClient JobClient, job *apiv1.Job) (err error) {
	createOpts := metav1.CreateOptions{}
	getOpts := metav1.GetOptions{}
	job, err = jobClient.Create(ctx, job, createOpts)
	require.Nil(t, err)

	for job.Status.Active > 0 || job.Status.StartTime == nil {
		t.Logf("Job is still running")
		time.Sleep(5 * time.Second)
		job, err = jobClient.Get(ctx, job.Name, getOpts)
		require.Nil(t, err)
	}
	if job.Status.Failed > 0 {
		t.Fail()
	}
	return
}
