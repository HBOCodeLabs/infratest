package k8s

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hbocodelabs/infratest/mock"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAssertJobSucceeds_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	jobClient := mock.NewMockJobClient(ctrl)
	createOpts := metav1.CreateOptions{}
	getOpts := metav1.GetOptions{}
	ctx := context.Background()
	jobName := "job"
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Status: batchv1.JobStatus{
			Succeeded: 0,
			Failed:    0,
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
	fakeTest := &testing.T{}

	jobClient.EXPECT().Create(ctx, job, createOpts).Return(job, nil)
	jobClient.EXPECT().Get(ctx, jobName, getOpts).DoAndReturn(func(context.Context, string, metav1.GetOptions) (*batchv1.Job, error) {
		returnJob := job.DeepCopy()
		returnJob.Status.Succeeded = 1
		return returnJob, nil
	})

	err := AssertJobSucceeds(fakeTest, context.TODO(), jobClient, job)

	assert.Nil(t, err)
	assert.False(t, fakeTest.Failed())
}

func TestAssertJobSucceeds_Fails(t *testing.T) {
	ctrl := gomock.NewController(t)
	jobClient := mock.NewMockJobClient(ctrl)
	createOpts := metav1.CreateOptions{}
	getOpts := metav1.GetOptions{}
	ctx := context.Background()
	jobName := "job"
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Status: batchv1.JobStatus{
			Succeeded: 0,
			Failed:    0,
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
	fakeTest := &testing.T{}

	jobClient.EXPECT().Create(ctx, job, createOpts).Return(job, nil)
	jobClient.EXPECT().Get(ctx, jobName, getOpts).DoAndReturn(func(context.Context, string, metav1.GetOptions) (*batchv1.Job, error) {
		returnJob := job.DeepCopy()
		returnJob.Status.Failed = 1
		return returnJob, nil
	})

	err := AssertJobSucceeds(fakeTest, context.TODO(), jobClient, job)

	assert.Nil(t, err)
	assert.True(t, fakeTest.Failed())
}
