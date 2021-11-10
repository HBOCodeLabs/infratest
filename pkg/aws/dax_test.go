// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.

package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dax"
	"github.com/aws/aws-sdk-go-v2/service/dax/types"

	gomock "github.com/golang/mock/gomock"

	"github.com/hbocodelabs/infratest/mock"

	"github.com/stretchr/testify/assert"
)

func TestAssertDAXClusterEncrypted_Pass(t *testing.T) {
	// Setup
	t.Parallel()

	fakeTest := &testing.T{}

	ctrl := gomock.NewController(t)
	client := mock.NewMockDAXClient(ctrl)

	clusterName := "daxcluster"
	expectedInput := &dax.DescribeClustersInput{
		ClusterNames: []string{clusterName},
	}
	output := &dax.DescribeClustersOutput{
		Clusters: []types.Cluster{
			{
				ClusterName:    &clusterName,
				SSEDescription: &types.SSEDescription{Status: types.SSEStatusEnabled},
			},
		},
	}
	ctx := context.Background()
	client.EXPECT().
		DescribeClusters(ctx, expectedInput).
		Times(1).
		DoAndReturn(
			func(context.Context, *dax.DescribeClustersInput, ...func(dax.Options)) (*dax.DescribeClustersOutput, error) {
				return output, nil
			},
		)

	// Execute
	AssertDAXClusterEncrypted(fakeTest, ctx, client, clusterName)

	// Assert
	ctrl.Finish()
	assert.False(t, fakeTest.Failed())
}

func TestAssertDAXClusterEncrypted_Fail(t *testing.T) {
	// Setup
	t.Parallel()

	fakeTest := &testing.T{}

	ctrl := gomock.NewController(t)
	client := mock.NewMockDAXClient(ctrl)

	clusterName := "daxcluster"
	expectedInput := &dax.DescribeClustersInput{
		ClusterNames: []string{clusterName},
	}
	output := &dax.DescribeClustersOutput{
		Clusters: []types.Cluster{
			{
				ClusterName:    &clusterName,
				SSEDescription: &types.SSEDescription{Status: types.SSEStatusDisabled},
			},
		},
	}
	ctx := context.Background()
	client.EXPECT().
		DescribeClusters(ctx, expectedInput).
		Times(1).
		DoAndReturn(
			func(context.Context, *dax.DescribeClustersInput, ...func(dax.Options)) (*dax.DescribeClustersOutput, error) {
				return output, nil
			},
		)

	// Execute
	AssertDAXClusterEncrypted(fakeTest, ctx, client, clusterName)

	// Assert
	ctrl.Finish()
	assert.True(t, fakeTest.Failed())
}
