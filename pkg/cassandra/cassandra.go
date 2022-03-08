// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package cassandra

import (
	"context"
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SessionInterface interface {
	//Query(stmt string, values ...interface{}) *Query
	Query(string, ...interface{}) *gocql.Query
}

type QueryInterface interface {
	//Bind(...interface{}) QueryInterface
	Exec() error
}

type Session struct {
	session *gocql.Session
}

type GetSessionInput struct {
	InstanceIp string
	Keyspace   string
	Username   string
	Password   string
}

type AssertCassandraQueryInput struct {
	Query *gocql.Query
}

func GetSession(t *testing.T, ctx context.Context, input GetSessionInput) (session SessionInterface) {
	cluster := gocql.NewCluster(input.InstanceIp)
	cluster.Consistency = gocql.One
	cluster.Keyspace = input.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: input.Username,
		Password: input.Password,
	}
	session, err := cluster.CreateSession()
	require.Nil(t, err)
	return session
}

func AssertCassandraQuerySucceeds(t *testing.T, ctx context.Context, query QueryInterface) {
	err := query.Exec()

	assert.Nil(t, err)
}
