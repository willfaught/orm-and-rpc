#!/bin/sh

# Generated

set -e -x

cqlsh <<EOF
use company;
create table offer(Created int, Deleted int, ID text, Name text, Updated int, primary key (ID));
EOF
