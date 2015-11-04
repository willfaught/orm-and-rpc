#!/bin/sh

# Generated

set -e -x

go get github.com/gocql/gocql github.com/motemen/gore

CASSANDRA=${CASSANDRA:-"127.0.0.1"}

gore <<EOF
:import github.com/gocql/gocql
:import github.com/motemen/gore
var cluster = gocql.NewCluster($CASSANDRA)
cluster.Keyspace = "company"
cluster.Port = 9042
cluster.ProtoVersion = 1
var session1, err = cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session1.Close()
cluster.Keyspace = "company"
session2, err := cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session2.Close()
metadata1, err := session1.KeyspaceMetadata("offer")
if err != nil {
    log.Fatalln("Error:", err)
}
metadata2, err := session2.KeyspaceMetadata("offer2")
if err != nil {
    log.Fatalln("Error:", err)
}
if _, ok = metadata1.Tables["offer"]; !ok {
    log.Fatalln("Error: Table company.offer does not exist")
}
if _, ok := metadata2.Tables["offer2"]; !ok {
    log.Fatalln("Error: Table company.offer2 does not exist")
}
var all = session1.Query("select * from offer")
if err := all.Exec(); err != nil {
    log.Fatal("Error:", err)
}
var iter = all.Iter()
for _, m1 := range iter.SliceMap() {
    if session2.Query("select * from offer2 where ID = ?", m1["ID"]).Scan() == gocql.ErrNotFound {
        if err := session2.Query("insert into offer2 (Added, Created, Deleted, ID, Updated) values (?, ?, ?, ?, ?)", "TBD", m1["Created"], m1["Deleted"], m1["ID"], m1["Updated"]).Exec(); err != nil {
            log.Fatalln("Error:", err)
        }
    }
    time.Sleep(time.Second / 2)
}
if err := iter.Close(); err != nil {
    log.Fatalln("Error:", err)
}
if err := session1.Query("drop table offer"); err != nil {
    log.Fatalln("Error:", err)
}
EOF
