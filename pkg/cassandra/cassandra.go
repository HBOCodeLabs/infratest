// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package cassandra

import (
	"context"
	"testing"
	"github.com/gocql/gocql"
)

type goCqlClient interface {
	cluster *gocql..ClusterConfig
	session *gocql..Session

}

type DBConnection struct {
	cluster *gocql..ClusterConfig
	session *gocql..Session
}

type ClusterConfigInput struct {
	ip string
	keyspace string
}
var connection DBConnection



func SetupDBConnection(input ClusterConfigInput) {
	connection.cluster - gocql.NewCluster(input.ip)
	connection.cluster.Consistency = gocql.One
	connection.cluster.keyspace = input.keyspace
	connection.session,  _ = connection.cluster.CreateSession()


}

func AssertCassandraQuerySucceeds(ClusterConfigInput) {
	assertion = false

	query := "select key from system.local"
	if err := connection.session.Query(query).Exec();

	err != nil {
		return err
	}


}
