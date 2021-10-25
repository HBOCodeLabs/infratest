package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
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

func (c Route53ClientMock) ListHostedZonesByNameInput(listHostedZonesByNameInput *route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesOutput, error) {
	return c.listHostedZonesOutput, c.err
}

func TestAssertRoute53HostedZoneExists_NotFound(t *testing.T) {
	fakeTest := &testing.T{}
	client := NewMockClient(&route53.ListHostedZonesOutput{}, nil)
	AssertRoute53HostedZoneExists(fakeTest, client, "bar.com")

	assert.True(t, fakeTest.Failed(), "expected AssertRoute53HostedZoneExists to fail")
}

func TestAssertRoute53HostedZoneExists_Error(t *testing.T) {
	fakeTest := &testing.T{}
	client := NewMockClient(&route53.ListHostedZonesOutput{}, errors.New("some error"))
	AssertRoute53HostedZoneExists(fakeTest, client, "foo.com")

	assert.True(t, fakeTest.Failed(), "expected AssertRoute53HostedZoneExists to fail")
}

func TestAssertRoute53HostedZoneExists_Found(t *testing.T) {
	fakeTest := &testing.T{}
	name := "foo.com"
	client := NewMockClient(&route53.ListHostedZonesOutput{
		HostedZones: []types.HostedZone{
			types.HostedZone{
				Name: &name,
			},
		},
	}, nil)
	AssertRoute53HostedZoneExists(fakeTest, client, name)

	assert.False(t, fakeTest.Failed(), "expected AssertRoute53HostedZoneExists to pass")
}
