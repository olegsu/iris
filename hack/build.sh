#!/usr/bin/env bash
set -e

OUTFILE=${1:-/usr/local/bin/irisctl}
BUILD_DATE=$(date)
GIT_TAG=$(git describe --tags)
GIT_COMMIT=$(git rev-parse --short HEAD)

LDFLAGS="-s -w
 -X 'github.com/olegsu/iris/pkg/util.BuildVersion=${GIT_TAG}'
 -X 'github.com/olegsu/iris/pkg/util.BuildDate=${BUILD_DATE}'
 -X 'github.com/olegsu/iris/pkg/util.BuildCommit=${GIT_COMMIT}'
 -X 'github.com/olegsu/iris/pkg/util.BuildBy=Makefile'
"

go build -ldflags "$LDFLAGS" -o "$OUTFILE" main.go
