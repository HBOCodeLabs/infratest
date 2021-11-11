// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
	"github.com/hbocodelabs/infratest/mock"
	"github.com/stretchr/testify/assert"

	"testing"
)

type EC2ClientMock struct {
	DescribeInstancesInput  *ec2.DescribeInstancesInput
	DescribeInstancesOutput *ec2.DescribeInstancesOutput
	DescribeVolumesOutput   *ec2.DescribeVolumesOutput
	DescribeTagsOutput      *ec2.DescribeTagsOutput
	Test                    *testing.T
}

func (c EC2ClientMock) DescribeInstances(ctx context.Context, input *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	if c.Test == nil && c.DescribeInstancesInput != nil {
		return nil, fmt.Errorf("Mock object not set up with a test object.")
	}
	if c.DescribeInstancesInput != nil {
		assert.Equal(c.Test, c.DescribeInstancesInput, input)
	}
	return c.DescribeInstancesOutput, nil
}

func (c EC2ClientMock) DescribeVolumes(ctx context.Context, input *ec2.DescribeVolumesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	return c.DescribeVolumesOutput, nil
}

func (c EC2ClientMock) DescribeTags(ctx context.Context, input *ec2.DescribeTagsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeTagsOutput, error) {
	if c.Test == nil {
		return nil, fmt.Errorf("Mock object not set up with a test object.")
	}
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

// This is a stub function; tests for this will use the new Mock object.
func (c EC2ClientMock) DescribeSecurityGroups(ctx context.Context, input *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
	return nil, nil
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
						BlockDeviceMappings: []types.InstanceBlockDeviceMapping{
							{
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
			{
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
	describeTagsInput := AssertEC2TagValueInput{
		TagName:    tagName,
		Value:      tagValue,
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
	describeTagsInput := AssertEC2TagValueInput{
		TagName:    tagName,
		Value:      tagValue,
		InstanceID: instanceID,
	}
	ctx := context.Background()
	fakeTest := &testing.T{}

	// Test
	AssertEC2TagValue(fakeTest, ctx, clientMock, describeTagsInput)
	assert.False(t, fakeTest.Failed(), "AssertEC2TagValue failed the test when tag value matched.")
}

func TestGetEC2InstancesByTag(t *testing.T) {
	// Setup
	tagName := "myTag"
	tagKeyName := fmt.Sprintf("tag:%s", tagName)
	tagValues := []string{"myValue1", "myValuy2"}
	tags := map[string][]string{
		tagName: tagValues,
	}
	filters := []types.Filter{
		{
			Name:   &tagKeyName,
			Values: tagValues,
		},
	}
	expectedInput := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	instanceID := "abc123456"
	output := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId: &instanceID,
					},
				},
			},
		},
		NextToken: nil,
	}
	expectedOutput := []types.Instance{
		{
			InstanceId: &instanceID,
		},
	}
	clientMock := &EC2ClientMock{
		DescribeInstancesInput:  expectedInput,
		DescribeInstancesOutput: output,
		Test:                    t,
	}
	ctx := context.Background()

	// Execute
	actualOutput, err := getEC2InstancesByTagE(ctx, clientMock, tags)

	// Assert
	assert.Nil(t, err, "getEC2InstancesByTagE returned an unexpected error")
	assert.ElementsMatch(t, expectedOutput, actualOutput, "getEC2InstancesByTagE did not return the expected results")
}

func TestAssertEC2InstancesSubnetBalanced_Matched(t *testing.T) {
	subnetID1 := "s123456"
	subnetID2 := "s7891011"
	subnets := []types.Subnet{
		{
			SubnetId: &subnetID1,
		},
		{
			SubnetId: &subnetID2,
		},
	}
	instanceID1 := "a123456"
	instanceID2 := "b123456"
	instanceID3 := "c123456"
	instances := []types.Instance{
		{
			InstanceId: &instanceID1,
			SubnetId:   &subnetID1,
		},
		{
			InstanceId: &instanceID2,
			SubnetId:   &subnetID2,
		},
		{
			InstanceId: &instanceID3,
			SubnetId:   &subnetID1,
		},
	}
	input := AssertEC2InstancesSubnetBalancedInput{
		Instances: instances,
		Subnets:   subnets,
	}
	fakeTest := &testing.T{}
	ctx := context.Background()

	AssertEC2InstancesBalancedInSubnets(fakeTest, ctx, input)

	assert.False(t, fakeTest.Failed())
}

func TestCreateFiltersFromMap(t *testing.T) {
	filterKey := "key"
	filterKey2 := "otherkey"
	filterValues := []string{
		"hello",
		"there",
	}
	filterValues2 := []string{
		"something",
		"else",
	}
	inputMap := map[string][]string{
		filterKey:  filterValues,
		filterKey2: filterValues2,
	}
	expectedOutput := []types.Filter{
		{
			Name:   &filterKey,
			Values: filterValues,
		},
		{
			Name:   &filterKey2,
			Values: filterValues2,
		},
	}

	actualOutput := CreateFiltersFromMap(inputMap)

	assert.ElementsMatch(t, expectedOutput, actualOutput)
}

func TestGetEC2SecurityGroupByNameE(t *testing.T) {
	// Setup
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockEC2Client(ctrl)

	securityGroupName := "security-group"
	expectedFilterName := "group-name"
	expectedInput := &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   &expectedFilterName,
				Values: []string{securityGroupName},
			},
		},
	}
	output := &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []types.SecurityGroup{
			{
				GroupName: &securityGroupName,
			},
		},
	}
	expectedOutput := &types.SecurityGroup{
		GroupName: &securityGroupName,
	}
	mockClient.EXPECT().DescribeSecurityGroups(ctx, expectedInput).
		Times(1).
		DoAndReturn(
			func(context.Context, *ec2.DescribeSecurityGroupsInput, ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
				return output, nil
			},
		)

	// Execute
	actualOutput, err := GetEC2SecurityGroupByName(ctx, mockClient, securityGroupName)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, actualOutput)
}
