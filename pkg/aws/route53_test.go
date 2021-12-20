// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package aws

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
)

type Route53ClientMock struct {
	listHostedZonesOutput *route53.ListHostedZonesOutput
	listHostedZonesErr    error

	listResourceRecordSetsOutput *route53.ListResourceRecordSetsOutput
	listResourceRecordSetsErr    error
}

func (c Route53ClientMock) ListHostedZonesByNameInput(input *route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesOutput, error) {
	return c.listHostedZonesOutput, c.listHostedZonesErr
}

func (c Route53ClientMock) ListResourceRecordSets(ctx context.Context, input *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
	return c.listResourceRecordSetsOutput, c.listResourceRecordSetsErr
}

func TestAssertRoute53HostedZoneExists_NotFound(t *testing.T) {
	fakeTest := &testing.T{}
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{},
		listHostedZonesErr:    nil,
	}
	AssertRoute53HostedZoneExists(fakeTest, client, "bar.com")

	assert.True(t, fakeTest.Failed(), "expected AssertRoute53HostedZoneExists to fail")
}

func TestAssertRoute53HostedZoneExists_Error(t *testing.T) {
	fakeTest := &testing.T{}
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{},
		listHostedZonesErr:    errors.New("some error"),
	}
	AssertRoute53HostedZoneExists(fakeTest, client, "foo.com")

	assert.True(t, fakeTest.Failed(), "expected AssertRoute53HostedZoneExists to fail")
}

func TestAssertRoute53HostedZoneExists_Found(t *testing.T) {
	fakeTest := &testing.T{}
	name := "foo.com"
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{
				types.HostedZone{
					Name: &name,
				},
			},
		},
		listHostedZonesErr: nil,
	}
	AssertRoute53HostedZoneExists(fakeTest, client, name)

	assert.False(t, fakeTest.Failed(), "expected AssertRoute53HostedZoneExists to pass")
}

func TestAssertRecordExistsInHostedZone_Found(t *testing.T) {
	fakeTest := &testing.T{}
	zoneName := "foo.com"
	recordName := fmt.Sprintf("foo.%s", zoneName)
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{
				types.HostedZone{
					Name: &zoneName,
				},
			},
		},
		listHostedZonesErr: nil,
		listResourceRecordSetsOutput: &route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []types.ResourceRecordSet{
				types.ResourceRecordSet{
					Name: &recordName,
				},
			},
		},
		listResourceRecordSetsErr: nil,
	}

	AssertRecordExistsInHostedZone(fakeTest, context.Background(), client, AssertRecordInput{
		RecordName: recordName,
		ZoneName:   zoneName,
	})

	assert.False(t, fakeTest.Failed(), "expected AssertRecordExistsInZone to pass")
}

func TestAssertRecordExistsInHostedZone_RecordNotFound(t *testing.T) {
	fakeTest := &testing.T{}
	zoneName := "foo.com"
	recordName := fmt.Sprintf("foo.%s", zoneName)
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{
				types.HostedZone{
					Name: &zoneName,
				},
			},
		},
		listHostedZonesErr: nil,
		listResourceRecordSetsOutput: &route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []types.ResourceRecordSet{},
		},
		listResourceRecordSetsErr: nil,
	}

	AssertRecordExistsInHostedZone(fakeTest, context.Background(), client, AssertRecordInput{
		RecordName: recordName,
		ZoneName:   zoneName,
	})

	assert.True(t, fakeTest.Failed(), "expected AssertRecordExistsInZone to fail")
}

func TestAssertRecordExistsInHostedZone_RecordTypeNotFound(t *testing.T) {
	fakeTest := &testing.T{}
	zoneName := "foo.com"
	recordName := fmt.Sprintf("foo.%s", zoneName)
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{
				types.HostedZone{
					Name: &zoneName,
				},
			},
		},
		listHostedZonesErr: nil,
		listResourceRecordSetsOutput: &route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []types.ResourceRecordSet{
				types.ResourceRecordSet{
					Name: &recordName,
					Type: types.RRTypeA,
				},
			},
		},
		listResourceRecordSetsErr: nil,
	}

	AssertRecordExistsInHostedZone(fakeTest, context.Background(), client, AssertRecordInput{
		RecordName: recordName,
		RecordType: types.RRTypeSoa,
		ZoneName:   zoneName,
	})

	assert.True(t, fakeTest.Failed(), "expected AssertRecordExistsInZone to fail")
}

func TestAssertRecordExistsInHostedZone_ListResourceRecordSets_Error(t *testing.T) {
	fakeTest := &testing.T{}
	zoneName := "foo.com"
	recordName := fmt.Sprintf("foo.%s", zoneName)
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{
				types.HostedZone{
					Name: &zoneName,
				},
			},
		},
		listHostedZonesErr: nil,
		listResourceRecordSetsOutput: &route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []types.ResourceRecordSet{},
		},
		listResourceRecordSetsErr: errors.New("some error"),
	}

	AssertRecordExistsInHostedZone(fakeTest, context.Background(), client, AssertRecordInput{
		RecordName: recordName,
		ZoneName:   zoneName,
	})

	assert.True(t, fakeTest.Failed(), "expected AssertRecordExistsInZone to fail")
}

func TestAssertRecordExistsInHostedZone_ZoneNotFound(t *testing.T) {
	fakeTest := &testing.T{}
	zoneName := "foo.com"
	recordName := fmt.Sprintf("foo.%s", zoneName)
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{},
		},
		listHostedZonesErr: nil,
		listResourceRecordSetsOutput: &route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []types.ResourceRecordSet{
				types.ResourceRecordSet{
					Name: &recordName,
				},
			},
		},
		listResourceRecordSetsErr: nil,
	}

	AssertRecordExistsInHostedZone(fakeTest, context.Background(), client, AssertRecordInput{
		RecordName: recordName,
		ZoneName:   zoneName,
	})

	assert.True(t, fakeTest.Failed(), "expected AssertRecordExistsInZone to fail")
}

func TestAssertRecordExistsInHostedZone_ListHostedZonesByNameInput_Error(t *testing.T) {
	fakeTest := &testing.T{}
	zoneName := "foo.com"
	recordName := fmt.Sprintf("foo.%s", zoneName)
	client := Route53ClientMock{
		listHostedZonesOutput: &route53.ListHostedZonesOutput{
			HostedZones: []types.HostedZone{
				types.HostedZone{
					Name: &zoneName,
				},
			},
		},
		listHostedZonesErr: errors.New("some error"),
		listResourceRecordSetsOutput: &route53.ListResourceRecordSetsOutput{
			ResourceRecordSets: []types.ResourceRecordSet{
				types.ResourceRecordSet{
					Name: &recordName,
				},
			},
		},
		listResourceRecordSetsErr: nil,
	}

	AssertRecordExistsInHostedZone(fakeTest, context.Background(), client, AssertRecordInput{
		RecordName: recordName,
		ZoneName:   zoneName,
	})

	assert.True(t, fakeTest.Failed(), "expected AssertRecordExistsInZone to fail")
}
