// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.

package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/dax/types"
	"github.com/stretchr/testify/assert"
)

// DAXClient serves as a stub client interface for the AWS SDK [DAX client](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dax#Client).
type DAXClient interface {
	DescribeClusters(context.Context, *dax.DescribeClustersInput, ...func(*dax.Options)) (*dax.DescribeClustersOutput, error)
}

// AssertDAXClusterEncrypted asserts that a DAX cluster has server side encryption enabled.
func AssertDAXClusterEncrypted(t *testing.T, ctx context.Context, client DAXClient, name string) {
	input := &dax.DescribeClustersInput{
		ClusterNames: []string{name},
	}
	output, err := client.DescribeClusters(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, types.SSEStatusEnabled, output.Clusters[0].SSEDescription.Status)
}

// AssertDAXClusterSubnetGroup asserts that a DAX cluster has a given subnet group associated to it.
func AssertDAXClusterSubnetGroup(t *testing.T, ctx context.Context, client DAXClient, name string, subnetGroupName string) {
	input := &dax.DescribeClustersInput{
		ClusterNames: []string{name},
	}
	output, err := client.DescribeClusters(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, subnetGroupName, *output.Clusters[0].SubnetGroup)
}

// AssertDAXClusterSecurityGroup asserts that a DAX cluster is associated with a given security group. It does not assert
// that the group provided is the _only_ security group associated with the cluster.
func AssertDAXClusterSecurityGroup(t *testing.T, ctx context.Context, client DAXClient, ec2client EC2Client, name string, securityGroupName string) {
	securityGroupOutput, err := GetEC2SecurityGroupByName(ctx, ec2client, securityGroupName)
	assert.Nil(t, err, "An error occurred while retrieving the named security group.")
	assert.NotNil(t, securityGroupOutput, "A security group with the specified name does not exist.")
	expectedSecurityGroupID := *securityGroupOutput.GroupId

	input := &dax.DescribeClustersInput{
		ClusterNames: []string{name},
	}
	output, err := client.DescribeClusters(ctx, input)
	assert.Nil(t, err)

	securityGroupMatchFound := false
	for _, securityGroupAttachment := range output.Clusters[0].SecurityGroups {
		if *securityGroupAttachment.SecurityGroupIdentifier == expectedSecurityGroupID {
			securityGroupMatchFound = true
		}
	}
	assert.True(t, securityGroupMatchFound, "A security group with the name specified is not associated with the DAX cluster.")
}
