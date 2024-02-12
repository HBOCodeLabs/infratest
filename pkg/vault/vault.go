// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.

/*
Package vault provides methods for asserting that objects in Vault exist and match a given specification.

# Basics

Methods in the vault package make use of the functional options pattern as a way to easily evolve interfaces
without breaking backwards compatibility. Many methods require the use of the Vault LogicalClient struct,
which is the object used to do most read / write / non-admin functions in the Go Vault client. For these,
you can use the `WithLogicalClient` functional option.

# Examples

Execute a method using a Vault client configured using standard Vault environment variables

	import (
		"os"
		"testing"
		"github.com/hbocodelabs/infratest/pkg/vault"
	)

	func TestVaultThing(t *testing.T) {
		rootToken = os.GetEnv("VAULT_TOKEN")
		vaultAddress = os.GetEnv("VAULT_ADDR")
		ctx := context.Background()

		clientConfig := &api.Config{
			Address:    vaultAddress,
			MaxRetries: 100,
		}
		client, err := api.NewClient(clientConfig)
		require.Nil(t, err, "Vault NewClient method returned an unexpected error.")
		client.SetToken(rootToken)

		logicalClient := client.Logical()

		expectedPath := "secret/data/hello"
		expectedSecretData := map[string]interface{}{
			"data": map[string]interface{}{
				"username": "myname",
				"password": "password",
			},
		}

		// Do something that is supposed to create the expected secret here

		vault.AssertVaultSecretExists(t, ctx, vault.WithLogicalClient(logicalClient), vault.WithPath(expectedPath), vault.WithKey("username"), vault.WithValue("myname"))
		vault.AssertVaultSecretExists(t, ctx, vault.WithLogicalClient(logicalClient), vault.WithPath(expectedPath), vault.WithKey("password"), vault.WithValue("password"))
	}
*/
package vault

import (
	"context"

	"github.com/hashicorp/vault/api"
	"github.com/hbocodelabs/infratest/pkg/test"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=../../mock/vault.go -package=mock github.com/hbocodelabs/infratest/pkg/vault LogicalClient

const (
	constFailureMissingKey      = "Secret did not contain the desired key."
	constFailureKeyNotSpecified = "A Key must be set by one or more passed functional options."
	constFailureValueNotMatch   = "The specified value does not match the actual value."
)

// LogicalClient is an interface that matches with the Vault Logical Client (https://pkg.go.dev/github.com/hashicorp/vault/api#Logical) API type.
type LogicalClient interface {
	Read(string) (*api.Secret, error)
	Delete(string) (*api.Secret, error)
}

// AssertVaultOptions is a struct that is used for passing options to methods used for asserting
// against Vault objects. It should never be used directly; instead create functional option methods
// which modify it, according to the AssertVaultOptsFunc interface.
type AssertVaultOptions struct {
	LogicalClient LogicalClient
	// The path at which the object (secret, role, etc) should exist
	Path string
	// The desired key for a Secret object
	Key string
	// The desired value for an object's key
	Value interface{}
}

// AssertVaultOptsFunc defines the interface used for all functional options used by Vault related
// methods.
type AssertVaultOptsFunc func(opts *AssertVaultOptions) error

// WithLogicalClient allows you to pass an existing Vault Logical Client object for use with
// the assertion methods. This object is returned by the `Logical` method of the Vault Client object.
// See https://pkg.go.dev/github.com/hashicorp/vault/api#Client.Logical.
func WithLogicalClient(client LogicalClient) AssertVaultOptsFunc {
	return func(opts *AssertVaultOptions) error {
		opts.LogicalClient = client
		return nil
	}
}

// WithPath sets the path at which various assertion methods will test for
// a Vault object's existence.
func WithPath(path string) AssertVaultOptsFunc {
	return func(opts *AssertVaultOptions) error {
		opts.Path = path
		return nil
	}
}

// WithKey sets the key to assert exists when using objects that have keys, such as
// Secrets.
func WithKey(key string) AssertVaultOptsFunc {
	return func(opts *AssertVaultOptions) error {
		opts.Key = key
		return nil
	}
}

// WithValue sets the value to assert that an object at a given path (and key in some cases)
// is equal to.
func WithValue(value string) AssertVaultOptsFunc {
	return func(opts *AssertVaultOptions) error {
		opts.Value = value
		return nil
	}
}

/*
AssertSecretExists asserts that a key/value secret exists at a given path and has a given
key present in the secret data. If the "WithValue" functional option is used, it wil also assert
that the value of given secret key is the same as the passed in value.

# Examples

Assert that a secret and key exists, using a passed in client, ignoring the value of the secret key.

	expectedPath := "path"
	expectedKey := "key"
	AssertSecretExists(
		t,
		ctx,
		WithClient(client),
		WithPath(expectedPath),
		WithKey(expectedKey),
	)

Assert that a secret and key exists, using a passed in client, with a particular value.

	expectedPath := "path"
	expectedKey := "key"
	expectedValue := "value"
	AssertSecretExists(
		t,
		ctx,
		WithClient(client),
		WithPath(expectedPath),
		WithKey(expectedKey),
		WithValue(expectedValue),
	)
*/
func AssertSecretExists(ctx context.Context, t test.T, optFns ...AssertVaultOptsFunc) {
	opts := &AssertVaultOptions{}
	for _, optFn := range optFns {
		err := optFn(opts)
		require.Nil(t, err, "Optional function threw an unexpected error.")
	}

	require.NotNil(t, opts.LogicalClient, "A Client must be set by one or more passed functional options.")
	require.NotEmpty(t, opts.Path, "A Path must be set by one or more passed functional options.")
	require.NotEmpty(t, opts.Key, constFailureKeyNotSpecified)

	secret, err := opts.LogicalClient.Read(opts.Path)
	require.Nil(t, err, "Vault client read returned an unexpected error.")
	secretData := secret.Data["data"].(map[string]interface{})

	require.Containsf(t, secretData, opts.Key, constFailureMissingKey)

	if opts.Value != nil {
		actualValue := secretData[opts.Key]
		require.Equal(t, opts.Value, actualValue, constFailureValueNotMatch)
	}

}
