#!/bin/bash
set -e

VERSION=$(git describe --tags --abbrev=0)
REVISION=$(git rev-parse --short HEAD)
LDFLAGS="-X 'main.version=${VERSION}' -X 'main.revision=${REVISION}'"

XC_ARCH=${XC_ARCH:-386 amd64}
XC_OS=${XC_OS:-darwin linux windows}

rm -rf pkg/
gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -output "pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" \
    -ldflags "${LDFLAGS}" \
    ./cmd/nrec-moody