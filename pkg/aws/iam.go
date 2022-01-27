// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package aws

import (
	"context"
	"net/url"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2/json"

	"testing"
)

type IAMClient interface {
	GetRole(context.Context, *iam.GetRoleInput, ...func(*iam.Options)) (*iam.GetRoleOutput, error)
}

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect   string
	Action   interface{}
	Resource interface{}
}

// AssertIAMPolicyDocumentContainsResourceAction will assert the an IAM Policy Document provided contains a Statement with the given Resource, Action, and Effect.
// If such a Statement does not exist within the Policy the test will immediately fail.
func AssertIAMPolicyDocumentContainsResourceAction(t *testing.T, resource string, action string, effect string, policyDocument PolicyDocument) {

	if !checkIAMPolicyDocumentContainsResourceAction(resource, action, effect, policyDocument) {
		t.Logf("Could not locate combination of resource '%s', action '%s', effect '%s' in provided policy document.", resource, action, effect)
		t.Fail()
	}
}

// AssertIAMPolicyDocumentsContainsResourceAction will assert the _at least one_ IAM Policy Document in a provided set contains a Statement with the given Resource, Action, and Effect.
// If such a Statement does not exist within the provided Policies the test will immediately fail.
func AssertIAMPolicyDocumentsContainResourceAction(t *testing.T, resource string, action string, effect string, policyDocuments []PolicyDocument) {
	for i := 0; i < len(policyDocuments); i++ {
		if checkIAMPolicyDocumentContainsResourceAction(resource, action, effect, policyDocuments[i]) {
			return
		}
	}

	t.Logf("Could not locate combination of resource '%s', action '%s', effect '%s' in provided policy documents.", resource, action, effect)
	t.Fail()
}

// checkIAMPolicyDocumentContainsResourceAction checks if a provided Policy Document contains a given Resource, Action, and Effect combination.
func checkIAMPolicyDocumentContainsResourceAction(resource string, action string, effect string, policyDocument PolicyDocument) bool {
	resourceIndex := 0
	statements := policyDocument.Statement
	for resourceIndex != -1 {
		resourceIndex = findIAMPolicyResource(statements, resource, resourceIndex)
		if resourceIndex != -1 {
			actionIndex := findIamPolicyAction(statements[resourceIndex], action, effect)
			if actionIndex != -1 {
				return true
			}
			resourceIndex++
		}
	}

	return false
}

// getIAMPolicyDefaultVersionE gets the current version of the IAM Policy for a given ARN.
func getIAMPolicyDefaultVersionE(context context.Context, policyArn string, client *iam.Client) (PolicyDocument, error) {
	IAMGetPolicyInput := &iam.GetPolicyInput{
		PolicyArn: &policyArn,
	}
	IAMPolicyOutput, err := client.GetPolicy(context, IAMGetPolicyInput)
	if err != nil {
		return PolicyDocument{}, err
	}

	iamPolicyDefaultVersionID := IAMPolicyOutput.Policy.DefaultVersionId
	IAMPolicyDefaultVersion, err := client.GetPolicyVersion(context, &iam.GetPolicyVersionInput{
		PolicyArn: &policyArn,
		VersionId: iamPolicyDefaultVersionID,
	})
	if err != nil {
		return PolicyDocument{}, err
	}

	IAMPolicyDocument, err := unMarshallPolicyDocument(*IAMPolicyDefaultVersion.PolicyVersion.Document)
	return *IAMPolicyDocument, err
}

// getIAMRolePolicyNamesE returns an array of strings containing the names of all inline policies attached to a Role.
func getIAMRolePolicyNamesE(context context.Context, client iam.ListRolePoliciesAPIClient, roleName string) ([]string, error) {
	listRolePoliciesInput := &iam.ListRolePoliciesInput{
		RoleName: &roleName,
	}
	paginator := iam.NewListRolePoliciesPaginator(client, listRolePoliciesInput)

	policyNames := []string{}

	hasMoreItems := true
	for hasMoreItems {
		listRolePolicyOutput, err := paginator.NextPage(context)
		if err != nil {
			return []string{}, err
		}
		policyNames = append(policyNames, listRolePolicyOutput.PolicyNames...)
		hasMoreItems = listRolePolicyOutput.IsTruncated
	}
	return policyNames, nil
}

// getIAMRolePolicyDocuments returns an array of PolicyDocument structs representing all the inline IAM Policies attached to a role.
func getIAMRolePolicyDocuments(context context.Context, client *iam.Client, roleName string, policyNames []string) ([]PolicyDocument, error) {
	rolePolicyDocuments := []PolicyDocument{}

	for i := 0; i < len(policyNames); i++ {
		getIAMRolePolicyInput := &iam.GetRolePolicyInput{
			PolicyName: &policyNames[i],
			RoleName:   &roleName,
		}
		getRolePolicyOutput, err := client.GetRolePolicy(context, getIAMRolePolicyInput)
		if err != nil {
			return rolePolicyDocuments, err
		}
		policyDocument, err := unMarshallPolicyDocument(*getRolePolicyOutput.PolicyDocument)
		if err != nil {
			return rolePolicyDocuments, err
		}
		rolePolicyDocuments = append(rolePolicyDocuments, *policyDocument)
	}

	return rolePolicyDocuments, nil
}

// findIamPolicyAction returns the index of a particular Action in an IAM Policy Document Statement. If the Action is not found, it will
// return -1.
func findIamPolicyAction(statement StatementEntry, action string, effect string) int {
	actions := parseIAMPolicyField(statement.Action)
	for i := 0; i < len(actions); i++ {
		if actions[i] == action && statement.Effect == effect {
			return i
		}
	}
	return -1
}

// findIAMPolicyResource returns the index (zero based) of the Statement that contains a given Resource within an IAM Policy's Statement collection, starting at the index given.
// If the resource is not found, it will return a value of -1.
func findIAMPolicyResource(statements []StatementEntry, resource string, startIndex int) int {
	for i := startIndex; i < len(statements); i++ {
		statement := statements[i]
		resources := parseIAMPolicyField(statement.Resource)
		for r := 0; r < len(resources); r++ {
			if resources[r] == resource {
				return i
			}
		}
	}
	return -1
}

// ParseIAMPolicyField takes an input and returns a single element array (if the passed input is a string)
// or an array of strings (if the passed input is an array). This is because the Resource and Action fields of an IAM Policy document statement
// can either be a single string or an array of strings.
func parseIAMPolicyField(field interface{}) []string {
	var array []string

	if reflect.TypeOf(field).String() != "string" {
		intArray := field.([]interface{})
		strArray := make([]string, len(intArray))
		for i := 0; i < len(intArray); i++ {
			strArray[i] = intArray[i].(string)
		}
		array = append(array, strArray...)
	} else {
		array = append(array, field.(string))
	}
	return array
}

func unMarshallPolicyDocument(document string) (*PolicyDocument, error) {
	var policyDocument PolicyDocument
	documentUnencoded, err := url.QueryUnescape(document)
	if err != nil {
		return &PolicyDocument{}, err
	}
	err = json.Unmarshal([]byte(documentUnencoded), &policyDocument)
	if err != nil {
		return &PolicyDocument{}, err
	}
	return &policyDocument, nil
}

func getIAMRole(context context.Context, client IAMClient, roleName string) (output *iam.GetRoleOutput, err error) {

	getIamRoleInput := &iam.GetRoleInput{
		RoleName: &roleName,
	}

	IAMRoleOutput, err := client.GetRole(context, getIamRoleInput)
	if err != nil {
		return IAMRoleOutput, err
	}

	return IAMRoleOutput, err
}

type AssertIamRoleComponentInput struct {
	RoleName       string
	AssertionKey   string
	AssertionValue interface{}
}

func AssertIamRoleComponent(t *testing.T, ctx context.Context, client IAMClient, input AssertIamRoleComponentInput) {

	getIamRoleOutput, err := getIAMRole(ctx, client, input.RoleName)
	assert.Nil(t, err)

	assertionObject := getIamRoleOutput[input.AssertionKey]

	assert.True(t, assertionObject, input.AssertionValue)

}
