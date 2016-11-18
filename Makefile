VERBOSE_FLAG = $(if $(VERBOSE),-v)

VERSION := ${shell gobump show ./cmd/nrec-moody|sed -e 's/{"version":"\(.*\)"}/\1/g'}
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

version:
	@echo $(VERSION)

build: deps
	go build -ldflags "$(LDFLAGS)" $(VERBOSE_FLAG) ./cmd/nrec-moody

test: deps
	go test $(VERBOSE_FLAG) ./...

deps:
	go get github.com/motemen/gobump/cmd/gobump
	go get -d $(VERBOSE_FLAG)

install: deps
	go install -ldflags "$(LDFLAGS)" $(VERBOSE_FLAG) ./cmd/nrec-moody

.PHONY: build test deps install version