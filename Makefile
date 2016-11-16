VERBOSE_FLAG = $(if $(VERBOSE),-v)

VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

build: deps
	go build -ldflags "$(LDFLAGS)" $(VERBOSE_FLAG) -o bin/nrec-moody ./cmd/nrec-moody

test: deps
	go test $(VERBOSE_FLAG) ./...

deps:
	go get -d $(VERBOSE_FLAG)

.PHONY: build test deps