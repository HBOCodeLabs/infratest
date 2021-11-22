// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
)

// Route53Client is an AWS Route53 API client.
// Typically, it's a [Route53](https://docs.aws.amazon.com/sdk-for-go/api/service/route53/#Route53).
type Route53Client interface {
	ListHostedZonesByNameInput(*route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesOutput, error)
	ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
}

// AssertRoute53HostedZoneExists asserts whether or not the Route53 zone name
// it's passed is found amongst those reported by the AWS API.
func AssertRoute53HostedZoneExists(t *testing.T, client Route53Client, zoneName string) {
	_, found, err := findZone(client, zoneName)

	assert.Nil(t, err)
	assert.True(t, found, fmt.Sprintf("'%s' not found", zoneName))
}

// AssertRoute53RecordExistsInHostedZone asserts whether or not the Route53 record
// name it's passed exists amongst those associated with the the Route53 zone whose
// name it's passed.
func AssertRoute53RecordExistsInHostedZone(t *testing.T, client Route53Client, recordName string, zoneName string) {
	recordFound := false

	z, zoneFound, err := findZone(client, zoneName)

	assert.Nil(t, err)
	assert.True(t, zoneFound, fmt.Sprintf("zone '%s' not found", zoneName))

	if !zoneFound {
		return
	}

	recs, err := client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		StartRecordName: &recordName,
		HostedZoneId:    z.Id,
	})
	assert.Nil(t, err)

	for _, r := range recs.ResourceRecordSets {
		if strings.ToLower(*r.Name) == strings.ToLower(recordName) {
			recordFound = true
			break
		}
	}

	assert.True(t, recordFound, fmt.Sprintf("record '%s' not found", recordName))
}

func findZone(client Route53Client, zoneName string) (*types.HostedZone, bool, error) {
	zones, err := client.ListHostedZonesByNameInput(&route53.ListHostedZonesByNameInput{
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
