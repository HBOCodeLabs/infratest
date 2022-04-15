// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package aws

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
)

// Route53Client is an AWS Route53 API client.
// Typically, it's a [Route53](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/route53#Client).
type Route53Client interface {
	ListHostedZonesByName(context.Context, *route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesOutput, error)
	ListHostedZonesByVPC(context.Context, *route53.ListHostedZonesByVPCInput) (*route53.ListHostedZonesByVPCOutput, error)
	ListResourceRecordSets(context.Context, *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

// AssertRoute53HostedZoneExists asserts whether or not the Route53 zone name
// it's passed is found amongst those reported by the AWS API.
func AssertRoute53HostedZoneExists(t *testing.T, ctx context.Context, client Route53Client, zoneName string) {
	_, found, err := findZoneE(ctx, client, zoneName)

	assert.Nil(t, err)
	assert.True(t, found, fmt.Sprintf("'%s' not found", zoneName))
}

// AssertRecordInput is used as an input to the AssertRecordExistsInHostedZone method.
type AssertRecordInput struct {
	// The record name.
	RecordName string

	// The record type.
	RecordType types.RRType

	// The zone name.
	ZoneName string
}

// AssertRoute53RecordExistsInHostedZone asserts whether or not the Route53 record
// name it's passed exists amongst those associated with the the Route53 zone whose
// name it's passed.
func AssertRoute53RecordExistsInHostedZone(t *testing.T, ctx context.Context, client Route53Client, recordInput AssertRecordInput) {
	recordFound := false
	zoneName := recordInput.ZoneName
	recordName := recordInput.RecordName

	z, zoneFound, err := findZoneE(ctx, client, zoneName)

	assert.Nil(t, err)
	assert.True(t, zoneFound, fmt.Sprintf("zone '%s' not found", zoneName))

	if !zoneFound {
		return
	}

	recs, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		StartRecordName: &recordName,
		HostedZoneId:    z.Id,
	})
	assert.Nil(t, err)

	for _, r := range recs.ResourceRecordSets {
		if recordInput.RecordType == "" || r.Type == recordInput.RecordType {
			recordFound = strings.EqualFold(*r.Name, recordName)

			if recordFound {
				break
			}
		}
	}

	assert.True(t, recordFound, fmt.Sprintf("record '%s' not found", recordName))
}

// AssertRoute53ZoneIsAssociatedVPCInput is used as input to the AssertRoute53ZoneIsAssociatedWithVPC method.
type AssertRoute53ZoneIsAssociatedWithVPCInput struct {
	// The ID of the VPC to check for zone association (required).
	VPCID string

	// The region of the VPC to check for zone association (required).
	VPCRegion types.VPCRegion

	// The name of the zone to check for VPC association (required).
	ZoneName string
}

// AssertRoute53ZoneIsAssociatedWithVPC asserts whether or not the Route53 zone
// is associated with the given VPC.
func AssertRoute53ZoneIsAssociatedWithVPC(t *testing.T, ctx context.Context, client Route53Client, associationInput AssertRoute53ZoneIsAssociatedWithVPCInput) {
	input := route53.ListHostedZonesByVPCInput{
		VPCId:     &associationInput.VPCID,
		VPCRegion: associationInput.VPCRegion,
	}
	zones := make([]string, 0)

	for {
		output, err := client.ListHostedZonesByVPC(ctx, &input)
		assert.Nil(t, err)
		for _, zone := range output.HostedZoneSummaries {
			zones = append(zones, *zone.Name)
		}

		input.NextToken = output.NextToken
		if input.NextToken == nil {
			break
		}
	}

	// AWS returns zone names with a trailing period, i.e. "myzone.com." instead of
	// "myzone.com". We need to add the period if it's missing.
	zoneName := associationInput.ZoneName
	if !strings.HasSuffix(zoneName, ".") {
		zoneName += "."
	}
	assert.Contains(t, zones, zoneName)
}

func findZoneE(ctx context.Context, client Route53Client, zoneName string) (*types.HostedZone, bool, error) {
	zones, err := client.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{
		DNSName: &zoneName,
	})
	if err != nil {
		return nil, false, err
	}

	var zone *types.HostedZone
	for _, n := range zones.HostedZones {
		if zoneName == *n.Name {
			zone = &n
			break
		}
	}

	if zone != nil {
		return zone, true, nil
	}

	return nil, false, nil
}
