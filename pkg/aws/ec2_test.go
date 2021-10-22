// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"

	"testing"
)

type EC2ClientMock struct {
	DescribeInstancesOutput *ec2.DescribeInstancesOutput
	DescribeVolumesOutput   *ec2.DescribeVolumesOutput
	DescribeTagsOutput      *ec2.DescribeTagsOutput
	Test                    *testing.T
}

func (c EC2ClientMock) DescribeInstances(ctx context.Context, input *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return c.DescribeInstancesOutput, nil
}

func (c EC2ClientMock) DescribeVolumes(ctx context.Context, input *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	return c.DescribeVolumesOutput, nil
}

func (c EC2ClientMock) DescribeTags(ctx context.Context, input *ec2.DescribeTagsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeTagsOutput, error) {
	hasInstanceResourceTypeFilter := false
	resourceTypeFilterName := "resource-type"
	for _, filter := range input.Filters {
		filterName := filter.Name
		if *filterName == resourceTypeFilterName {
			for _, value := range filter.Values {
				if value == "instance" {
					hasInstanceResourceTypeFilter = true
				}
			}
		}
	}
	assert.True(c.Test, hasInstanceResourceTypeFilter, "DescribeTags was called without specifying a resource type filter.")
	return c.DescribeTagsOutput, nil
}

func TestAssertEC2VolumeEncryptedE_Match(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	encrypted := true
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}

	// Execute
	result, err := AssertEC2VolumeEncryptedE(context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
	})

	// Assert
	assert.True(t, result)
	assert.Nil(t, err)

}

func TestAssertEC2VolumeEncrypted_Match(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	encrypted := true
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}
	fakeTest := &testing.T{}

	// Execute
	AssertEC2VolumeEncrypted(fakeTest, context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
	})

	// Assert
	assert.False(t, fakeTest.Failed())
}


func TestAssertEC2VolumeEncryptedE_NoMatch(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	encrypted := false
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}

	// Execute
	result, err := AssertEC2VolumeEncryptedE(context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
	})

	// Assert
	assert.False(t, result)
	assert.Nil(t, err)

}

func TestAssertEC2VolumeEncrypted_NoMatch(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	encrypted := false
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}
	fakeTest := &testing.T{}

	// Execute
	AssertEC2VolumeEncrypted(fakeTest, context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
	})

	// Assert
	assert.True(t, fakeTest.Failed())
}

func TestAssertEC2VolumeEncryptedE_MatchWithKMSKeyID(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	encrypted := true
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}

	// Execute
	result, err := AssertEC2VolumeEncryptedE(context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
		KMSKeyID:   kmsKeyID,
	})

	// Assert
	assert.True(t, result)
	assert.Nil(t, err)
}

func TestAssertEC2VolumeEncrypted_MatchWithKMSKeyID(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	encrypted := true
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}
	fakeTest := &testing.T{}

	// Execute
	AssertEC2VolumeEncrypted(fakeTest, context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
		KMSKeyID:   kmsKeyID,
	})

	// Assert
	assert.False(t, fakeTest.Failed())
}

func TestAssertEC2VolumeEncryptedE_NoMatchWithKMSKeyID(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	kmsKeyID2 := "/key/id2"
	encrypted := true
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID2,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}

	// Execute
	result, err := AssertEC2VolumeEncryptedE(context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
		KMSKeyID:   kmsKeyID,
	})

	// Assert
	assert.False(t, result)
	assert.Nil(t, err)
}

func TestAssertEC2VolumeEncrypted_NoMatchWithKMSKeyID(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	volumeId := "v123dfasd92"
	deviceName := "/dev/sdc"
	kmsKeyID := "/key/id"
	kmsKeyID2 := "/key/id2"
	encrypted := true
	instanceOutput := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			types.Reservation{
				Instances: []types.Instance{
					types.Instance{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							types.InstanceBlockDeviceMapping{
								DeviceName: &deviceName,
								Ebs: &types.EbsInstanceBlockDevice{
									VolumeId: &volumeId,
								},
							},
						},
					},
				},
			},
		},
	}
	volumeOutput := &ec2.DescribeVolumesOutput{
		Volumes: []types.Volume{
			types.Volume{
				Encrypted: &encrypted,
				KmsKeyId:  &kmsKeyID2,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesOutput: instanceOutput,
		DescribeVolumesOutput:   volumeOutput,
	}
	fakeTest := &testing.T{}

	// Execute
	AssertEC2VolumeEncrypted(fakeTest, context.Background(), clientMock, AssertEC2VolumeEncryptedInput{
		DeviceID:   deviceName,
		InstanceID: instanceID,
		KMSKeyID:   kmsKeyID,
	})

	// Assert
	assert.True(t, fakeTest.Failed())
}

func TestAssertEC2TagValue_NoMatch(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	tagName := "MyTag"
	tagValue := "TagValue"
	wrongTagValue := "OtherValue"
	nextToken := ""
	describeTagsOutput := &ec2.DescribeTagsOutput{
		NextToken: &nextToken,
		Tags: []types.TagDescription{
			{
				ResourceId: &instanceID,
				ResourceType: types.ResourceTypeInstance,
				Key: &tagName,
				Value: &wrongTagValue,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeTagsOutput: describeTagsOutput,
		Test: t,
	}
	describeTagsInput := AssertEC2TagValueInput{
		TagName: tagName,
		Value: tagValue,
		InstanceID: instanceID,
	}
	ctx := context.Background()
	fakeTest := &testing.T{}

	// Test
	AssertEC2TagValue(fakeTest, ctx, clientMock, describeTagsInput)
	assert.True(t, fakeTest.Failed(), "AssertEC2TagValue did not fail the test when the tag value did not match.")
}

func TestAssertEC2TagValueE_NoMatch(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	tagName := "MyTag"
	tagValue := "TagValue"
	wrongTagValue := "OtherValue"
	nextToken := ""
	describeTagsOutput := &ec2.DescribeTagsOutput{
		NextToken: &nextToken,
		Tags: []types.TagDescription{
			{
				ResourceId:   &instanceID,
				ResourceType: types.ResourceTypeInstance,
				Key:          &tagName,
				Value:        &wrongTagValue,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeTagsOutput: describeTagsOutput,
		Test:               t,
	}
	describeTagsInput := AssertEC2TagValueEInput{
		TagName:    tagName,
		Value:      tagValue,
		InstanceID: instanceID,
	}
	ctx := context.Background()

	// Test
	result, err := AssertEC2TagValueE(ctx, clientMock, describeTagsInput)
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, result, "AssertEC2TagValueE returned 'true' when tag value did not match.")
}

func TestAssertEC2TagValueE_Match(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	tagName := "MyTag"
	tagValue := "TagValue"
	nextToken := ""
	describeTagsOutput := &ec2.DescribeTagsOutput{
		NextToken: &nextToken,
		Tags: []types.TagDescription{
			{
				ResourceId:   &instanceID,
				ResourceType: types.ResourceTypeInstance,
				Key:          &tagName,
				Value:        &tagValue,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeTagsOutput: describeTagsOutput,
		Test:               t,
	}
	describeTagsInput := AssertEC2TagValueEInput{
		TagName:    tagName,
		Value:      tagValue,
		InstanceID: instanceID,
	}
	ctx := context.Background()

	// Test
	result, err := AssertEC2TagValueE(ctx, clientMock, describeTagsInput)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, result, "AssertEC2TagValueE returned 'false' when tag value matched.")
}

func TestAssertEC2TagValue_Match(t *testing.T) {
	// Setup
	instanceID := "i546acas321sd"
	tagName := "MyTag"
	tagValue := "TagValue"
	nextToken := ""
	describeTagsOutput := &ec2.DescribeTagsOutput{
		NextToken: &nextToken,
		Tags: []types.TagDescription{
			{
				ResourceId: &instanceID,
				ResourceType: types.ResourceTypeInstance,
				Key: &tagName,
				Value: &tagValue,
			},
		},
	}
	clientMock := &EC2ClientMock{
		DescribeTagsOutput: describeTagsOutput,
		Test: t,
	}
	describeTagsInput := AssertEC2TagValueInput{
		TagName: tagName,
		Value: tagValue,
		InstanceID: instanceID,
	}
	ctx := context.Background()
	fakeTest := &testing.T{}

	// Test
	AssertEC2TagValue(fakeTest, ctx, clientMock, describeTagsInput)
	assert.False(t, fakeTest.Failed(), "AssertEC2TagValue failed the test when tag value matched.")
}
