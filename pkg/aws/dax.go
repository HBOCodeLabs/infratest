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

//AssertDAXClusterSubnetGroup asserts that a DAX cluster has a given subnet group associated to it.
func AssertDAXClusterSubnetGroup(t *testing.T, ctx context.Context, client DAXClient, name string, subnetGroupName string) {
	input := &dax.DescribeClustersInput{
		ClusterNames: []string{name},
	}
	output, err := client.DescribeClusters(ctx, input)
	assert.Nil(t, err)
	assert.Equal(t, subnetGroupName, *output.Clusters[0].SubnetGroup)
}
