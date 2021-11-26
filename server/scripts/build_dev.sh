#!/bin/bash

rm ./webserver

export MODE="production"

go build ./cmd/webserver
