#!/bin/sh
go run `ls cmd/nejireco-pubsub/*.go|grep _test.go -v` $@
