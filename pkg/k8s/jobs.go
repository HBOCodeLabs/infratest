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

// AssertJobSucceedsOptions is a struct used for functional options for the AssertJobSucceeds method.
type AssertJobSucceedsOptions struct {
	// The job to run.
	JobSpec apiv1.Job
	// The CreateOptions used when creating the job.
	CreateOptions metav1.CreateOptions
	// The GetOptions used when retrieving the job for subsequent checks.
	GetOptions metav1.GetOptions
}

type AssertJobSucceedsOptsFunc func(AssertJobSucceedsOptions) error

/*
AssertJobSucceeds will start a Kubernetes Job using the provided client and spec, then after it completes will either
fail the passed test (if the job fails) or pass the test if the job succeeds. It should be passed a
JobClient object and a Job object.
*/
func AssertJobSucceeds(t *testing.T, ctx context.Context, jobClient JobClient, jobSpec *apiv1.Job, optFns ...AssertJobSucceedsOptsFunc) {
	createOpts := metav1.CreateOptions{}
	getOpts := metav1.GetOptions{}
	opts := AssertJobSucceedsOptions{
		JobSpec:       *jobSpec,
		CreateOptions: createOpts,
		GetOptions:    getOpts,
	}

	for _, f := range optFns {
		err := f(opts)
		if err != nil {
			t.Error(err)
			return
		}
	}

	job, err := jobClient.Create(ctx, &opts.JobSpec, opts.CreateOptions)
	if err != nil {
		t.Error(err)
		return
	}

	for !isJobCompleted(job) {
		t.Logf("Job %s is still running", job.Name)
		time.Sleep(5 * time.Second)
		job, err = jobClient.Get(ctx, job.Name, opts.GetOptions)
		if err != nil {
			t.Error(err)
			return
		}
	}
	if !k8s.IsJobSucceeded(job) {
		t.Errorf("Job with name '%s' did not complete successfully.", job.Name)
	}
}

// isJobCompleted returns a boolean value indicating if a job has completed (whether successfully or not).
func isJobCompleted(job *apiv1.Job) (isCompleted bool) {
	for _, condition := range job.Status.Conditions {
		if (condition.Type == apiv1.JobComplete || condition.Type == apiv1.JobFailed) && condition.Status == v1.ConditionTrue {
			isCompleted = true
		}
	}
	return
}
