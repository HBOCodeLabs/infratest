package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Client interface {
	DescribeInstances(context.Context, *ec2.DescribeInstancesInput, ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)

	DescribeVolumes(context.Context, *ec2.DescribeVolumesInput, ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error)

	DescribeTags(context.Context, *ec2.DescribeTagsInput, ...func(*ec2.Options)) (*ec2.DescribeTagsOutput, error)
}

// AssertEC2VolumeEncryptedInput is used as an input to the AssertEC2VolumeEncryptedE method.
type AssertEC2VolumeEncryptedInput struct {
	// The device ID that the volume is mapped to on the instance.
	DeviceID string
	// The KMS key ID that must be used to encrypt the volume. If left blank, then matching on this attribute will not be performed.
	KMSKeyID string
	// The Instance ID the volume must be attached to.
	InstanceID string
}

// AssertEC2TagValueEInput is used as an input to the AssertEC2TagValueE method.
type AssertEC2TagValueEInput struct {
	// The name of the tag to assert exists.
	TagName string
	// The value of the tag to assert.
	Value string
	// The Instance ID that the tag mustbe set on.
	InstanceID string
}

// AssertEC2VolumeEncryptedE asserts that a volume attached to an EC2 instance is encrypted and (optionally) done so using a specified KMS Key.
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

// AssertEC2TagValueE asserts that an EC2 instance has a given tag with the given value.
func AssertEC2TagValueE(ctx context.Context, client EC2Client, input AssertEC2TagValueEInput) (assertion bool, err error) {
	resourceTypeFilterName := "resource-type"
	resourceTypeFilterValue := "instance"
	describeTagsInput := &ec2.DescribeTagsInput{
		Filters: []types.Filter{
			{
				Name:   &resourceTypeFilterName,
				Values: []string{resourceTypeFilterValue},
			},
		},
	}
	describeTagsOutput, err := client.DescribeTags(ctx, describeTagsInput)
	if err != nil {
		return false, err
	}
	hasTagMatch := false
	for _, tag := range describeTagsOutput.Tags {
		tagKey := tag.Key
		tagValue := tag.Value
		if *tagKey == input.TagName {
			if *tagValue == input.Value {
				hasTagMatch = true
			}
		}
	}
	return hasTagMatch, nil
}

func getEC2InstanceByInstanceIDE(ctx context.Context, client EC2Client, InstanceID string) (ec2types.Instance, error) {
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{InstanceID},
	}
	describeInstancesOutput, err := client.DescribeInstances(ctx, describeInstancesInput)
	if err != nil {
		return types.Instance{}, err
	}
	if len(describeInstancesOutput.Reservations) == 0 {
		err = fmt.Errorf("Instance with ID '%s' was not found.", InstanceID)
		return types.Instance{}, err
	}
	instance := describeInstancesOutput.Reservations[0].Instances[0]
	return instance, nil
}

func getEC2VolumeByVolumeIDE(ctx context.Context, client EC2Client, VolumeID string) (ec2types.Volume, error) {
	describeVolumesInput := &ec2.DescribeVolumesInput{
		VolumeIds: []string{VolumeID},
	}
	describeVolumesOutput, err := client.DescribeVolumes(ctx, describeVolumesInput)
	if err != nil {
		return types.Volume{}, err
	}
	if len(describeVolumesOutput.Volumes) == 0 {
		err = fmt.Errorf("Volume with ID '%s' was not found.", VolumeID)
		return types.Volume{}, err
	}
	volume := describeVolumesOutput.Volumes[0]
	return volume, nil
}
