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
	ListResourceRecordSets(context.Context, *route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

// AssertHostedZoneExists asserts whether or not the Route53 zone name
// it's passed is found amongst those reported by the AWS API.
func AssertHostedZoneExists(t *testing.T, ctx context.Context, client Route53Client, zoneName string) {
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

// AssertRecordExistsInHostedZone asserts whether or not the Route53 record
// name it's passed exists amongst those associated with the the Route53 zone whose
// name it's passed.
func AssertRecordExistsInHostedZone(t *testing.T, ctx context.Context, client Route53Client, recordInput AssertRecordInput) {
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
		if strings.ToLower(*r.Name) == strings.ToLower(recordName) && &recordInput.RecordType != nil && r.Type == recordInput.RecordType {
			recordFound = true
			break
		}
	}

	assert.True(t, recordFound, fmt.Sprintf("record '%s' not found", recordName))
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
