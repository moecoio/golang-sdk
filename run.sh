#!/bin/bash
export GOPATH=`pwd`:`pwd`/vendor
nohup go run cmd/moecosdk.go &

