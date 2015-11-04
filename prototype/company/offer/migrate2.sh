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
var session2, err = cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session2.Close()
cluster.Keyspace = "company"
session3, err := cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session3.Close()
metadata2, err := session2.KeyspaceMetadata("offer2")
if err != nil {
    log.Fatalln("Error:", err)
}
metadata3, err := session3.KeyspaceMetadata("offer3")
if err != nil {
    log.Fatalln("Error:", err)
}
if _, ok = metadata2.Tables["offer2"]; !ok {
    log.Fatalln("Error: Table company.offer2 does not exist")
}
if _, ok := metadata3.Tables["offer3"]; !ok {
    log.Fatalln("Error: Table company.offer3 does not exist")
}
var all = session2.Query("select * from offer2")
if err := all.Exec(); err != nil {
    log.Fatal("Error:", err)
}
var iter = all.Iter()
for _, m2 := range iter.SliceMap() {
    if session3.Query("select * from offer2 where ID = ?", m2["ID"]).Scan() == gocql.ErrNotFound {
        if err := session3.Query("insert into offer3 (Created, Deleted, ID, Renamed, Updated) values (?, ?, ?, ?, ?)", m2["Created"], m2["Deleted"], m2["ID"], m2["Added"], m2["Updated"]).Exec(); err != nil {
            log.Fatalln("Error:", err)
        }
    }
    time.Sleep(time.Second / 2)
}
if err := iter.Close(); err != nil {
    log.Fatalln("Error:", err)
}
if err := session2.Query("drop table offer2"); err != nil {
    log.Fatalln("Error:", err)
}
EOF
