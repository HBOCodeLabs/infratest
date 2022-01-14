// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package cassandra

import (
	"context"
	"log"
	"testing"

	"github.com/gocql/gocql"
)

type SessionInterface interface {
	//Query(stmt string, values ...interface{}) *Query
	Query(stmt string, values ...interface{}) QueryInterface
}

type QueryInterface interface {
	Bind(...interface{}) QueryInterface
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
	//Query string
	Query *gocql.Query
}

func NewSession(session *gocql.Session) SessionInterface {
	return &Session{
		session,
	}
}

func GetSession(t *testing.T, ctx context.Context, input GetSessionInput) (cassandra.SessionInterface, err error) {
	cluster := gocql.NewCluster(input.InstanceIp)
	cluster.Consistency = gocql.One
	cluster.Keyspace = input.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: input.Username,
		Password: input.Password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	return cassandra.NewSession(session), err
}

func AssertCassandraQuerySucceedsE(t *testing.T, ctx context.Context, session Session, queryinput AssertCassandraQueryInput) (assertion bool, err error) {
	assertion = false

	connectionStatus = false
	result, err := session.Query(queryinput.Query).Exec()
	if err != nil {
		return assertion, err
	}
	connectionStatus - true

	queryStatus := false
	if result != nil {
		queryStatus = true
	}

	assertion = queryStatus && connectionStatus
	return

}
