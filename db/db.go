package db

import (
	"log"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitDB() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "personal_cloud"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to ScyllaDB:", err)
	}
	Session = session
}
