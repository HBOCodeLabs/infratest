package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/route53"
)

// Route53Client is an AWS Route53 API client.
// Typically, it's a [Route53](https://docs.aws.amazon.com/sdk-for-go/api/service/route53/#Route53).
type Route53Client interface {
	ListHostedZonesByName(*route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesOutput, error)
}

// AssertRoute53HostedZoneExists reports whether or not the Route53 zone name
// it's passed is found amongst those reported by the AWS API.
func AssertRoute53HostedZoneExists(client Route53Client, zoneName string) (bool, error) {
	zones, err := client.ListHostedZonesByName(&route53.ListHostedZonesByNameInput{
		DNSName: &zoneName,
	})
	if err != nil {
		return false, err
	}

	for _, n := range zones.HostedZones {
		if zoneName == *n.Name {
			return true, nil
		}
	}

	return false, fmt.Errorf("'%s' not found", zoneName)
}
