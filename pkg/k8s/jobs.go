package k8s

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	apiv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobClient is an interface that partially implements the [JobInterface object](https://pkg.go.dev/k8s.io/client-go@v0.23.1/kubernetes/typed/batch/v1#JobInterface).
type JobClient interface {
	Create(context.Context, *apiv1.Job, metav1.CreateOptions) (*apiv1.Job, error)
	Get(context.Context, string, metav1.GetOptions) (*apiv1.Job, error)
}

// AssertJobSucceedsInput is a struct used as an input to the AssertJobSucceeds method.
type AssertJobSucceedsInput struct {
	// The job to run
	JobSpec *apiv1.Job
}

// GetJobClientE returns a JobClient object given a path to a kubeconfig file and
// a namespace name. If a blank string is passed as the kubeconfig path, then
// the method will use the path $HOME/.kube/config, where $HOME is the user's home
// directory as determined by the OS.
func GetJobClientE(kubeconfigPath string, namespace string) (client JobClient, err error) {
	ctx := context.Background()
	clientset, err := GetClientsetE(ctx, WithGetClientsetEKubeconfigPath(kubeconfigPath))
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
// Or, get a clientset for an EKS cluster
client := aws.
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
func AssertJobSucceeds(t *testing.T, ctx context.Context, jobClient JobClient, input AssertJobSucceedsInput) {
	createOpts := metav1.CreateOptions{}
	getOpts := metav1.GetOptions{}
	job, err := jobClient.Create(ctx, input.JobSpec, createOpts)
	if err != nil {
		t.Error(err)
		return
	}

	for !IsJobCompleted(job) {
		t.Logf("Job is still running")
		time.Sleep(5 * time.Second)
		job, err = jobClient.Get(ctx, job.Name, getOpts)
		if err != nil {
			t.Error(err)
			return
		}
	}
	if !k8s.IsJobSucceeded(job) {
		t.Fail()
	}
}

func IsJobCompleted(job *apiv1.Job) (isCompleted bool) {
	for _, condition := range job.Status.Conditions {
		if (condition.Type == apiv1.JobComplete || condition.Type == apiv1.JobFailed) && condition.Status == v1.ConditionTrue {
			isCompleted = true
		}
	}
	return
}
