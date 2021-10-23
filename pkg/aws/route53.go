package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/stretchr/testify/assert"
)

// Route53Client is an AWS Route53 API client.
// Typically, it's a [Route53](https://docs.aws.amazon.com/sdk-for-go/api/service/route53/#Route53).
type Route53Client interface {
	ListHostedZonesByNameInput(*route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesOutput, error)
}

// AssertRoute53HostedZoneExists asserts whether or not the Route53 zone name
// it's passed is found amongst those reported by the AWS API.
func AssertRoute53HostedZoneExists(t *testing.T, client Route53Client, zoneName string) {
	found := false
	zones, err := client.ListHostedZonesByNameInput(&route53.ListHostedZonesByNameInput{
		DNSName: &zoneName,
	})
	assert.Nil(t, err)

	for _, n := range zones.HostedZones {
		if zoneName == *n.Name {
			found = true
			break
		}
	}

	assert.True(t, found, fmt.Sprintf("'%s' not found", zoneName))
}
