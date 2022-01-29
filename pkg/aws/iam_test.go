// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package aws

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	gomock "github.com/golang/mock/gomock"
	"github.com/hbocodelabs/infratest/mock"

	"github.com/stretchr/testify/assert"
)

func TestParseIAMPolicyField_ReturnsMultiArray(t *testing.T) {
	resourceString := []interface{}{
		"someresource",
		"someotherresource",
	}
	resources := parseIAMPolicyField(resourceString)

	assert.Len(t, resources, 2)
	assert.Equal(t, "someresource", resources[0])
	assert.Equal(t, "someotherresource", resources[1])
}

func TestParseIAMPolicyField_ReturnsSingleArray(t *testing.T) {
	resourceString := []interface{}{
		"someresource",
	}
	resources := parseIAMPolicyField(resourceString)

	assert.Len(t, resources, 1)
	assert.Equal(t, "someresource", resources[0])
}

func TestFindIAMPolicyResource_SingleEntry(t *testing.T) {
	resourceName := "someresource"
	statements := []StatementEntry{
		{
			Effect: "Allow",
			Action: []interface{}{"s3:*"},
			Resource: []interface{}{
				resourceName,
			},
		},
	}

	foundIndex := findIAMPolicyResource(statements, resourceName, 0)
	assert.Equal(t, 0, foundIndex)

	notFoundIndex := findIAMPolicyResource(statements, resourceName, 1)
	assert.Equal(t, -1, notFoundIndex)

}

func TestFindIAMPolicyResource_MultiEntry(t *testing.T) {
	resourceName := "someresource"
	statements := []StatementEntry{
		{
			Effect: "Allow",
			Action: []interface{}{"s3:*"},
			Resource: []interface{}{
				resourceName,
			},
		},
		{
			Effect: "Allow",
			Action: []interface{}{"ec2:*"},
			Resource: []interface{}{
				"someOtherResource",
			},
		},
		{
			Effect: "Allow",
			Action: []interface{}{"eks:*"},
			Resource: []interface{}{
				resourceName,
			},
		},
	}

	firstFoundIndex := findIAMPolicyResource(statements, resourceName, 0)
	assert.Equal(t, 0, firstFoundIndex)

	secondFoundIndex := findIAMPolicyResource(statements, resourceName, firstFoundIndex+1)
	assert.Equal(t, 2, secondFoundIndex)

	notFoundIndex := findIAMPolicyResource(statements, resourceName, secondFoundIndex+1)
	assert.Equal(t, -1, notFoundIndex)
}

func TestFindIAMPolicyAction_Found(t *testing.T) {
	actionName := "s3:*"
	effect := "Allow"
	statement := StatementEntry{
		Effect:   "Allow",
		Action:   []interface{}{"ec2:*", actionName},
		Resource: []interface{}{"*"},
	}

	firstFoundIndex := findIamPolicyAction(statement, actionName, effect)
	assert.Equal(t, 1, firstFoundIndex)
}

func TestFindIAMPolicyAction_NotFound_Action(t *testing.T) {
	actionName := "s3:*"
	effect := "Allow"
	statement := StatementEntry{
		Effect:   "Allow",
		Action:   []interface{}{"ec2:*", "s3:GetObject"},
		Resource: []interface{}{"*"},
	}

	notFoundIndex := findIamPolicyAction(statement, actionName, effect)
	assert.Equal(t, -1, notFoundIndex)
}

func TestFindIAMPolicyAction_NotFound_Effect(t *testing.T) {
	actionName := "s3:*"
	effect := "Allow"
	statement := StatementEntry{
		Effect:   "Deny`",
		Action:   []interface{}{"ec2:*", actionName},
		Resource: []interface{}{"*"},
	}

	notFoundIndex := findIamPolicyAction(statement, actionName, effect)
	assert.Equal(t, -1, notFoundIndex)
}

func TestCheckIAMPolicyContainsResourceAction_Single(t *testing.T) {
	action := "s3:*"
	effect := "Allow"
	resource := "*"
	statement := StatementEntry{
		Effect:   effect,
		Action:   []interface{}{action},
		Resource: []interface{}{resource},
	}
	policyDocument := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			statement,
		},
	}

	found := checkIAMPolicyDocumentContainsResourceAction(resource, action, effect, policyDocument)
	assert.True(t, found)
}

func TestCheckIAMPolicyContainsResourceAction_Multi(t *testing.T) {
	action := "s3:*"
	effect := "Allow"
	resource := "*"
	statement := StatementEntry{
		Effect:   effect,
		Action:   action,
		Resource: "arn:aws:s3:::something",
	}
	statement2 := StatementEntry{
		Effect:   effect,
		Action:   []interface{}{action},
		Resource: []interface{}{resource},
	}
	policyDocument := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			statement,
			statement2,
		},
	}

	found := checkIAMPolicyDocumentContainsResourceAction(resource, action, effect, policyDocument)
	assert.True(t, found)
}

func TestParseIAMPolicyField_SingleResource(t *testing.T) {
	resource := "*"

	result := parseIAMPolicyField(resource)
	assert.Len(t, result, 1)
	assert.Equal(t, resource, result[0])
}

func TestParseIAMPolicyField_MultiResource(t *testing.T) {
	firstResource := "aws:iam:something"
	secondResource := "aws:iam:somethingElse"
	resource := []interface{}{
		firstResource,
		secondResource,
	}

	result := parseIAMPolicyField(resource)
	assert.Len(t, result, 2)
	assert.Equal(t, firstResource, result[0])
	assert.Equal(t, secondResource, resource[1])
}

func TestUnMarshallPolicyDocument_DecodeJson(t *testing.T) {
	policyDocument := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect:   "Allow",
				Action:   "s3:*",
				Resource: "arn:aws:s3:::somebucket",
			},
		},
	}

	policyDocumentJson, err := json.Marshal(policyDocument)
	if err != nil {
		t.Fatal(err)
	}
	policyDocumentEncoded := url.QueryEscape(string(policyDocumentJson))

	policyDocumentResult, err := unMarshallPolicyDocument(policyDocumentEncoded)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, policyDocument, *policyDocumentResult)
}

func TestAssertIamRoleComponent_MaxDuration_Success(t *testing.T) {
	fakeTest := &testing.T{}

	ctrl := gomock.NewController(t)
	client := mock.NewMockIAMClient(ctrl)

	var roleMaxDuration int32 = 5200
	roleName := "testIam"
	input := &iam.GetRoleInput{
		RoleName: &roleName,
	}

	role := types.Role{
		RoleName:           &roleName,
		MaxSessionDuration: &roleMaxDuration,
	}
	output := &iam.GetRoleOutput{
		Role: &role,
	}
	ctx := context.Background()
	client.EXPECT().
		GetRole(ctx, input).
		Times(1).
		Return(output, nil)

	AssertIamRoleMaxSessionDuration(fakeTest, ctx, client, roleName, roleMaxDuration)
	ctrl.Finish()
	assert.False(t, fakeTest.Failed())

}

func TestAssertIamRoleComponent_MaxDuration_Fail(t *testing.T) {
	fakeTest := &testing.T{}

	ctrl := gomock.NewController(t)
	client := mock.NewMockIAMClient(ctrl)

	var actualRoleMaxDuration int32 = 5200
	var expectedRoleMaxDuration int32 = 2000
	roleName := "testIam"
	input := &iam.GetRoleInput{
		RoleName: &roleName,
	}

	role := types.Role{
		RoleName:           &roleName,
		MaxSessionDuration: &actualRoleMaxDuration,
	}
	output := &iam.GetRoleOutput{
		Role: &role,
	}
	ctx := context.Background()
	client.EXPECT().
		GetRole(ctx, input).
		Times(1).
		Return(output, nil)

	AssertIamRoleMaxSessionDuration(fakeTest, ctx, client, roleName, expectedRoleMaxDuration)
	ctrl.Finish()
	assert.True(t, fakeTest.Failed())

}
