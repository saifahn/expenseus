#!/bin/bash

# print commands
set -x

SCRIPTS="$(dirname $0)"
ROOT="$(dirname "$SCRIPTS")"

# source the env
source $ROOT/.env

rm -rf $ROOT/webserver
go build $ROOT/cmd/webserver
./webserver

# turn off printing commands
set +x
