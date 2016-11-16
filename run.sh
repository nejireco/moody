#!/bin/sh
go run `ls cmd/nrec-moody/*.go|grep _test.go -v` $@
