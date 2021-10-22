package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/stretchr/testify/assert"
)

type Route53ClientMock struct {
	listHostedZonesOutput *route53.ListHostedZonesOutput
	err                   error
}

func NewMockClient(listHostedZonesOutput *route53.ListHostedZonesOutput, err error) Route53ClientMock {
	return Route53ClientMock{
		listHostedZonesOutput,
		err,
	}
}

func (c Route53ClientMock) ListHostedZonesByName(listHostedZonesByNameInput *route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesOutput, error) {
	return c.listHostedZonesOutput, c.err
}

func TestAssertRoute53HostedZoneExists_NotFound(t *testing.T) {
	client := NewMockClient(&route53.ListHostedZonesOutput{}, nil)
	exists, err := AssertRoute53HostedZoneExists(client, "foo.com")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "'foo.com' not found")
	assert.False(t, exists)
}

func TestAssertRoute53HostedZoneExists_Error(t *testing.T) {
	client := NewMockClient(&route53.ListHostedZonesOutput{}, errors.New("some error"))
	exists, err := AssertRoute53HostedZoneExists(client, "foo.com")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "some error")
	assert.False(t, exists)
}

func TestAssertRoute53HostedZoneExists_Found(t *testing.T) {
	name := "foo.com"
	client := NewMockClient(&route53.ListHostedZonesOutput{
		HostedZones: []*route53.HostedZone{
			&route53.HostedZone{
				Name: &name,
			},
		},
	}, nil)
	exists, err := AssertRoute53HostedZoneExists(client, name)

	assert.Nil(t, err)
	assert.True(t, exists)
}
