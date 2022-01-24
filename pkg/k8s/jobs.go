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

/*
AssertJobSucceeds will start a Kubernetes Job using the provided client and spec, then after it completes will either
fail the passed test (if the job fails) or pass the test if the job succeeds. It should be passed a
JobClient object and a Job object.
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
