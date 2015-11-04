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
var session4, err = cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session4.Close()
cluster.Keyspace = "company"
session5, err := cluster.CreateSession()
if err != nil {
    log.Fatalln("Error:", err)
}
defer session5.Close()
metadata4, err := session4.KeyspaceMetadata("offer4")
if err != nil {
    log.Fatalln("Error:", err)
}
metadata5, err := session5.KeyspaceMetadata("offer5")
if err != nil {
    log.Fatalln("Error:", err)
}
if _, ok = metadata4.Tables["offer4"]; !ok {
    log.Fatalln("Error: Table company.offer4 does not exist")
}
if _, ok := metadata5.Tables["offer5"]; !ok {
    log.Fatalln("Error: Table company.offer5 does not exist")
}
var all = session4.Query("select * from offer4")
if err := all.Exec(); err != nil {
    log.Fatal("Error:", err)
}
var iter = all.Iter()
for _, m4 := range iter.SliceMap() {
    if session5.Query("select * from offer5 where ID = ?", m4["ID"]).Scan() == gocql.ErrNotFound {
        if err := session5.Query("insert into offer5 (Created, Deleted, ID, Name, Updated) values (?, ?, ?, ?, ?)", m4["Created"], m4["Deleted"], m4["ID"], "TBD", m4["Updated"]).Exec(); err != nil {
            log.Fatalln("Error:", err)
        }
    }
    time.Sleep(time.Second / 2)
}
if err := iter.Close(); err != nil {
    log.Fatalln("Error:", err)
}
if err := session4.Query("drop table offer4"); err != nil {
    log.Fatalln("Error:", err)
}
EOF
