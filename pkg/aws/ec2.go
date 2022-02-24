// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	resourceIDFilterName            string = "resource-id"
	resourceTypeFilterName          string = "resource-type"
	resourceTypeFilterValueInstance string = "instance"
)

type EC2Client interface {
	DescribeInstances(context.Context, *ec2.DescribeInstancesInput, ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)

	DescribeVolumes(context.Context, *ec2.DescribeVolumesInput, ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error)

	DescribeTags(context.Context, *ec2.DescribeTagsInput, ...func(*ec2.Options)) (*ec2.DescribeTagsOutput, error)

	DescribeSecurityGroups(context.Context, *ec2.DescribeSecurityGroupsInput, ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
}

// AssertEC2VolumeEncryptedInput is used as an input to the AssertEC2VolumeEncryptedE and AssertEC2VolumeEncrypted methods.
type AssertEC2VolumeEncryptedInput struct {
	// The device ID that the volume is mapped to on the instance.
	DeviceID string
	// The KMS key ID that must be used to encrypt the volume. If left blank, then matching on this attribute will not be performed.
	KMSKeyID string
	// The Instance ID the volume must be attached to.
	InstanceID string
}

// AssertEC2TagValueEInput is used as an input to the AssertEC2TagValueE method. This is deprecated.
type AssertEC2TagValueEInput struct {
	// The name of the tag to assert exists.
	TagName string
	// The value of the tag to assert.
	Value string
	// The Instance ID that the tag mustbe set on.
	InstanceID string
}

// AssertVolumeAttributesInput is used as an input to the AssertEC2VolumeType,AssertEC2VolumeIops,AssertEC2VolumeThroughput  methods.
type AssertVolumeAttributesInput struct {
	// The Instance ID that is used to get devices associated to it.
	InstanceID string
	// The device ID that the volume is mapped to on the instance.
	// Used for informational purpose
	DeviceID string
	// The Volume Type for each volume
	VolumeType string
	// The Volume IOPS for each volume
	VolumeIOPS *int32
	// The Volume throughput for each volume
	VolumeThroughput *int32
}

// AssertEC2TagValueInput is used as an input to the AssertEC2TagValue method.
type AssertEC2TagValueInput struct {
	// The name of the tag to assert exists.
	TagName string
	// The value of the tag to assert.
	Value string
	// The Instance ID that the method will assert has a tag with the specified tag name and the specified value.
	InstanceID string
}

// AssertEC2VolumeEncryptedE asserts that a volume attached to an EC2 instance is encrypted and (optionally) done so using a specified KMS Key.
// This function is deprecated in favor of the AssertEC2VolumeEncrypted method.
func AssertEC2VolumeEncryptedE(ctx context.Context, client EC2Client, input AssertEC2VolumeEncryptedInput) (assertion bool, err error) {
	assertion = false

	instance, err := getEC2InstanceByInstanceIDE(ctx, client, input.InstanceID)
	if err != nil {
		return assertion, err
	}

	deviceFound := false
	deviceEncrypted := false
	kmsKeyMatches := false
	var trueVal bool = true
	for _, v := range instance.BlockDeviceMappings {
		if *v.DeviceName == input.DeviceID {
			deviceFound = true
			volume, err := getEC2VolumeByVolumeIDE(ctx, client, *v.Ebs.VolumeId)
			if err != nil {
				return false, err
			}
			if *volume.Encrypted == trueVal {
				deviceEncrypted = true
			}
			if input.KMSKeyID == "" || *volume.KmsKeyId == input.KMSKeyID {
				kmsKeyMatches = true
			}
		}
	}

	assertion = deviceEncrypted && deviceFound && kmsKeyMatches
	return
}

// AssertEC2VolumeEncrypted asserts that an EBS volume is encrypted, optionally using a specified KMS key.
func AssertEC2VolumeEncrypted(t *testing.T, ctx context.Context, client EC2Client, input AssertEC2VolumeEncryptedInput) {

	instance, err := getEC2InstanceByInstanceIDE(ctx, client, input.InstanceID)
	assert.Nil(t, err)

	deviceFound := false
	deviceEncrypted := false
	var trueVal bool = true
	for _, v := range instance.BlockDeviceMappings {
		if *v.DeviceName == input.DeviceID {
			deviceFound = true
			volume, err := getEC2VolumeByVolumeIDE(ctx, client, *v.Ebs.VolumeId)
			assert.Nil(t, err)
			if *volume.Encrypted == trueVal {
				deviceEncrypted = true
			}
			if input.KMSKeyID != "" {
				assert.Equal(t, input.KMSKeyID, *volume.KmsKeyId, "Volume with device ID '%s' for instance '%s' was not encrypted using the correct KMS Key ID.", input.DeviceID, input.InstanceID)
			}
		}
	}

	assert.True(t, deviceFound, "Volume with device ID '%s' was not found for instance '%s'.", input.DeviceID, input.InstanceID)
	assert.True(t, deviceEncrypted, "Volume with device ID '%s' for instance '%s' was not encrypted.", input.DeviceID, input.InstanceID)
}

// AssertVolumeType asserts the right volume type
func AssertEC2VolumeType(t *testing.T, ctx context.Context, client EC2Client, input AssertVolumeAttributesInput) {

	instance, err := getEC2InstanceByInstanceIDE(ctx, client, input.InstanceID)

	require.NoError(t, err)

	for _, v := range instance.BlockDeviceMappings {
		if *v.DeviceName == input.DeviceID {
			volume, err := getEC2VolumeByVolumeIDE(ctx, client, *v.Ebs.VolumeId)
			require.NoError(t, err)
			volumeType := fmt.Sprintf("%v", volume.VolumeType)
			assert.Equal(t, input.VolumeType, volumeType, "Volume with device ID '%s' does not have the right volume type.", input.DeviceID)
		}
	}
}

// AssertVolumeThroughput & IOPs asserts associated throughput for given volume type
func AssertEC2VolumeThroughput(t *testing.T, ctx context.Context, client EC2Client, input AssertVolumeAttributesInput) {

	instance, err := getEC2InstanceByInstanceIDE(ctx, client, input.InstanceID)
	require.NoError(t, err)

	for _, v := range instance.BlockDeviceMappings {
		if *v.DeviceName == input.DeviceID {
			volume, err := getEC2VolumeByVolumeIDE(ctx, client, *v.Ebs.VolumeId)
			require.NoError(t, err)
			if input.VolumeType != "gp2" {
				assert.Equal(t, input.VolumeThroughput, volume.Throughput, "Volume with device ID '%s' does not have the right throughput associated to volume.", input.DeviceID)
			} else {
				t.Logf("This test is ignored since it is not gp3 volume type : %s", input.VolumeType)
			}
		}
	}
}

// AssertVolumeIops asserts associated Iops for given volume type
func AssertEC2VolumeIOPS(t *testing.T, ctx context.Context, client EC2Client, input AssertVolumeAttributesInput) {

	instance, err := getEC2InstanceByInstanceIDE(ctx, client, input.InstanceID)
	require.NoError(t, err)

	for _, v := range instance.BlockDeviceMappings {
		if *v.DeviceName == input.DeviceID {
			volume, err := getEC2VolumeByVolumeIDE(ctx, client, *v.Ebs.VolumeId)
			require.NoError(t, err)
			if input.VolumeType != "gp2" {
				assert.Equal(t, input.VolumeIOPS, volume.Iops, "Volume with device ID '%s' does not have the right IOPS value associated to volume.", input.DeviceID)
			} else {
				t.Logf("This test is ignored since it is not gp3 volume type : %s", input.VolumeType)
			}
		}
	}
}

// AssertEC2TagValue asserts that an EC2 instance has a tag with the given value.
func AssertEC2TagValue(t *testing.T, ctx context.Context, client EC2Client, input AssertEC2TagValueInput) {
	resourceTypeFilterName := resourceTypeFilterName
	resourceTypeFilterValue := resourceTypeFilterValueInstance
	resourceIDFilterName := resourceIDFilterName
	describeTagsInput := &ec2.DescribeTagsInput{
		Filters: []types.Filter{
			{
				Name:   &resourceTypeFilterName,
				Values: []string{resourceTypeFilterValue},
			},
			{
				Name:   &resourceIDFilterName,
				Values: []string{input.InstanceID},
			},
		},
	}
	describeTagsOutput, err := client.DescribeTags(ctx, describeTagsInput)
	assert.Nil(t, err)
	hasMatch := false
	for _, tag := range describeTagsOutput.Tags {
		tagKey := *tag.Key
		tagValue := *tag.Value
		if tagKey == input.TagName {
			hasMatch = true
			assert.Equal(t, input.Value, tagValue, "Tag with key '%s' does not match expected value.", tagKey)
		}
	}
	assert.True(t, hasMatch, "Tag with key '%s' does not exist.", input.TagName)
}

func getEC2InstanceByInstanceIDE(ctx context.Context, client EC2Client, InstanceID string) (types.Instance, error) {
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{InstanceID},
	}
	describeInstancesOutput, err := client.DescribeInstances(ctx, describeInstancesInput)
	if err != nil {
		return types.Instance{}, err
	}
	if len(describeInstancesOutput.Reservations) == 0 {
		err = fmt.Errorf("instance with ID '%s' was not found", InstanceID)
		return types.Instance{}, err
	}
	instance := describeInstancesOutput.Reservations[0].Instances[0]
	return instance, nil
}

func getEC2VolumeByVolumeIDE(ctx context.Context, client EC2Client, VolumeID string) (types.Volume, error) {
	describeVolumesInput := &ec2.DescribeVolumesInput{
		VolumeIds: []string{VolumeID},
	}
	describeVolumesOutput, err := client.DescribeVolumes(ctx, describeVolumesInput)
	if err != nil {
		return types.Volume{}, err
	}
	if len(describeVolumesOutput.Volumes) == 0 {
		err = fmt.Errorf("volume with ID '%s' was not found", VolumeID)
		return types.Volume{}, err
	}
	volume := describeVolumesOutput.Volumes[0]
	return volume, nil
}

func getEC2InstancesByTagE(ctx context.Context, client EC2Client, tags map[string][]string) (instances []types.Instance, err error) {
	var filters []types.Filter

	for tagName, tagValues := range tags {
		tagKeyName := fmt.Sprintf("tag:%s", tagName)
		filter := types.Filter{
			Name:   &tagKeyName,
			Values: tagValues,
		}
		filters = append(filters, filter)
	}
	describeInstancesInput := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	output, err := client.DescribeInstances(ctx, describeInstancesInput)
	if err != nil {
		return nil, err
	} else if output == nil {
		return
	}
	reservations := output.Reservations
	for _, reservation := range reservations {
		instances = append(instances, reservation.Instances...)
	}
	return
}

type AssertEC2InstancesSubnetBalancedInput struct {
	// A list of instances
	Instances []types.Instance

	// A list of subnets
	Subnets []types.Subnet
}

// AssertEC2InstancesBalancedInSubnets asserts that EC2 instances in a list are spread evenly throughout a list of subnets,
// such that instance number 'x' in the list should be placed in the subnet with an index of 'x modulus the length of the
// subnet list'.
func AssertEC2InstancesBalancedInSubnets(t *testing.T, ctx context.Context, input AssertEC2InstancesSubnetBalancedInput) {
	subnetListLength := len(input.Subnets)
	assert.Greater(t, subnetListLength, 0, "The provided subnet list does not contain any elements.")
	assert.Greater(t, len(input.Instances), 0, "The provided instance list does not contain any elements.")
	for instanceIndex, instance := range input.Instances {
		acutalSubnetID := instance.SubnetId
		expectedSubnetID := input.Subnets[instanceIndex%subnetListLength].SubnetId
		assert.Equal(t, expectedSubnetID, acutalSubnetID, "Instance with ID '%s' is not in expected subnet.")
	}
}

// CreateFiltersFromMap is a utility method that creates a Filter object from a map of strings.
// It's designed to make creating filter objects easier without worrying about pointers and the like.
func CreateFiltersFromMap(input map[string][]string) (output []types.Filter) {
	for filterKey, filterValues := range input {
		// This is required since the value of the filterKey variable changes, and we have to pass a pointer.
		filterKeyCopy := filterKey
		filter := types.Filter{
			Name:   &filterKeyCopy,
			Values: filterValues,
		}
		output = append(output, filter)
	}
	return
}

// GetEC2SecurityGroupByName returns a security group object based on the name provided. If no matching group
// is found, it will return a nil value.
func GetEC2SecurityGroupByName(ctx context.Context, client EC2Client, name string) (securityGroup *types.SecurityGroup, err error) {
	filterKey := "group-name"
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   &filterKey,
				Values: []string{name},
			},
		},
	}
	output, err := client.DescribeSecurityGroups(ctx, input)
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, err
	}
	securityGroup = &output.SecurityGroups[0]
	return
}
