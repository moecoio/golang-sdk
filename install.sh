#!/bin/bash

sudo apt-get install golang

export GOPATH=`pwd`:`pwd`/vendor
go install -v github.com/mihalicyn/gatt
go install -v github.com/mattn/go-sqlite3
go install -v github.com/sirupsen/logrus
go install -v github.com/pkg/errors
