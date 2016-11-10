#!/bin/sh
go run `ls cmd/nejireco-moody/*.go|grep _test.go -v` $@
