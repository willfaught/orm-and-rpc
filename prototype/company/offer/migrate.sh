#!/bin/sh

# Generated

set -e -x

path=$(dirname $0)

eval "$path/migrate1.sh"
eval "$path/migrate2.sh"
eval "$path/migrate3.sh"
eval "$path/migrate4.sh"
