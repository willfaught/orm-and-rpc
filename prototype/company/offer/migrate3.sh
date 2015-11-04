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
var session3, err = cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session3.Close()
cluster.Keyspace = "company"
session4, err := cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session4.Close()
metadata3, err := session3.KeyspaceMetadata("offer3")
if err != nil {
    log.Fatalln("Error:", err)
}
metadata4, err := session4.KeyspaceMetadata("offer4")
if err != nil {
    log.Fatalln("Error:", err)
}
if _, ok = metadata3.Tables["offer3"]; !ok {
    log.Fatalln("Error: Table company.offer3 does not exist")
}
if _, ok := metadata4.Tables["offer4"]; !ok {
    log.Fatalln("Error: Table company.offer4 does not exist")
}
var all = session3.Query("select * from offer3")
if err := all.Exec(); err != nil {
    log.Fatal("Error:", err)
}
var iter = all.Iter()
for _, m3 := range iter.SliceMap() {
    if session4.Query("select * from offer4 where ID = ?", m3["ID"]).Scan() == gocql.ErrNotFound {
        if err := session4.Query("insert into offer4 (Created, Deleted, ID, Updated) values (?, ?, ?, ?)", m3["Created"], m3["Deleted"], m3["ID"], m3["Updated"]).Exec(); err != nil {
            log.Fatalln("Error:", err)
        }
    }
    time.Sleep(time.Second / 2)
}
if err := iter.Close(); err != nil {
    log.Fatalln("Error:", err)
}
if err := session3.Query("drop table offer3"); err != nil {
    log.Fatalln("Error:", err)
}
EOF
