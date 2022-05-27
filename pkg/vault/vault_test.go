// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package vault

import (
	"context"
	"fmt"
	"testing"

	matchers "github.com/Storytel/gomock-matchers"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/vault/api"
	"github.com/hbocodelabs/infratest/mock"
	"github.com/stretchr/testify/assert"
)

const (
	constExpectedKey   = "key"
	constExpectedPath  = "path"
	constExpectedValue = "value"
)

func assertPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Log("Test did not invoke panic.")
		t.Fail()
	}
}

func mockPanic() func() {
	return func() {
		panic("I'm panicking!")
	}
}

func TestWithClient(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockLogicalClient(ctrl)
	opts := &AssertVaultOptions{}

	method := WithLogicalClient(mockClient)
	err := method(opts)

	assert.Nil(t, err, "WithClient returned an unexpected error.")
	assert.Equal(t, mockClient, opts.LogicalClient, "WithClient did not set the Client property as expected.")
}

func TestWithPath(t *testing.T) {
	t.Parallel()
	opts := &AssertVaultOptions{}

	method := WithPath(constExpectedPath)
	err := method(opts)

	assert.Nil(t, err, "WithPath returned an unexpected error.")
	assert.Equal(t, constExpectedPath, opts.Path, "WithPath did not set the Path property as expected.")
}

func TestWithKey(t *testing.T) {
	t.Parallel()
	opts := &AssertVaultOptions{}

	method := WithKey(constExpectedKey)
	err := method(opts)

	assert.Nil(t, err, "WithKey returned an unexpected error.")
	assert.Equal(t, constExpectedKey, opts.Key, "WithKey did not set the Key property as expected.")
}

func TestWithValue(t *testing.T) {
	t.Parallel()
	opts := &AssertVaultOptions{}

	method := WithValue(constExpectedValue)
	err := method(opts)

	assert.Nil(t, err, "WithValue returned an unexpected error.")
	assert.Equal(t, constExpectedValue, opts.Value, "WithValue did not set the Value property as expected.")
}

func TestAssertVaultSecretExists_NoClient(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	fakeTest.EXPECT().FailNow().Times(1).DoAndReturn(mockPanic())
	m := matchers.Regexp("A Client must be set by one or more passed functional options")
	fakeTest.EXPECT().Errorf(gomock.Any(), m).Times(1)
	fakePath := "path"
	defer assertPanic(t)

	AssertSecretExists(ctx, fakeTest, WithPath(fakePath))
}

func TestAssertVaultSecretExists_NoPath(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	client := mock.NewMockLogicalClient(ctrl)
	fakeTest.EXPECT().FailNow().Times(1).DoAndReturn(mockPanic())
	m := matchers.Regexp("A Path must be set by one or more passed functional options")
	fakeTest.EXPECT().Errorf(gomock.Any(), m).Times(1)
	defer assertPanic(t)

	AssertSecretExists(ctx, fakeTest, WithLogicalClient(client))
}

func TestAssertVaultSecretExists_ReadError(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	client := mock.NewMockLogicalClient(ctrl)
	client.EXPECT().Read(gomock.Any()).Return(nil, fmt.Errorf("an error"))
	fakeTest.EXPECT().FailNow().Times(1).DoAndReturn(mockPanic())
	m := matchers.Regexp("Vault client read returned an unexpected error")
	fakeTest.EXPECT().Errorf(gomock.Any(), m).Times(1)
	defer assertPanic(t)

	AssertSecretExists(ctx, fakeTest, WithLogicalClient(client), WithPath(constExpectedPath), WithKey(constExpectedKey))
}

func TestAssertVaultSecretExists_NoKey(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	client := mock.NewMockLogicalClient(ctrl)
	fakeTest.EXPECT().FailNow().Times(1).DoAndReturn(mockPanic())
	m := matchers.Regexp(constFailureKeyNotSpecified)
	fakeTest.EXPECT().Errorf(gomock.Any(), m).Times(1)
	defer assertPanic(t)

	AssertSecretExists(ctx, fakeTest, WithLogicalClient(client), WithPath(constExpectedPath))
}

func TestAssertVaultSecretExists_NoMatchingKey(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	client := mock.NewMockLogicalClient(ctrl)
	mockSecret := &api.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"somekey": "someData",
			},
		},
	}
	client.EXPECT().Read(constExpectedPath).Times(1).Return(mockSecret, nil)
	m := matchers.Regexp(constFailureMissingKey)
	fakeTest.EXPECT().Errorf(gomock.Any(), m).Times(1)
	fakeTest.EXPECT().FailNow().Times(1).DoAndReturn(mockPanic())
	defer assertPanic(t)

	AssertSecretExists(ctx, fakeTest, WithLogicalClient(client), WithPath(constExpectedPath), WithKey(constExpectedKey))
}

func TestAssertVaultSecretExists_NoMatchingValue(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	client := mock.NewMockLogicalClient(ctrl)
	mockSecret := &api.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				constExpectedKey: "someData",
			},
		},
	}
	client.EXPECT().Read(constExpectedPath).Times(1).Return(mockSecret, nil)
	m := matchers.Regexp(constFailureValueNotMatch)
	fakeTest.EXPECT().Errorf(gomock.Any(), m).Times(1)
	fakeTest.EXPECT().FailNow().Times(1).DoAndReturn(mockPanic())
	defer assertPanic(t)

	AssertSecretExists(
		ctx,
		fakeTest,
		WithLogicalClient(client),
		WithPath(constExpectedPath),
		WithKey(constExpectedKey),
		WithValue(constExpectedValue),
	)
}

func TestAssertVaultSecretExists_MatchingKey(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	client := mock.NewMockLogicalClient(ctrl)
	mockSecret := &api.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				constExpectedKey: "someData",
			},
		},
	}
	client.EXPECT().Read(constExpectedPath).Times(1).Return(mockSecret, nil)

	AssertSecretExists(ctx, fakeTest, WithLogicalClient(client), WithPath(constExpectedPath), WithKey(constExpectedKey))
}

func TestAssertVaultSecretExists_MatchingValue(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	fakeTest := mock.NewMockT(ctrl)
	client := mock.NewMockLogicalClient(ctrl)
	mockSecret := &api.Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				constExpectedKey: constExpectedValue,
			},
		},
	}
	client.EXPECT().Read(constExpectedPath).Times(1).Return(mockSecret, nil)

	AssertSecretExists(
		ctx,
		fakeTest,
		WithLogicalClient(client),
		WithPath(constExpectedPath),
		WithKey(constExpectedKey),
		WithValue(constExpectedValue),
	)
}
