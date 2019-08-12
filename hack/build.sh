#!/usr/bin/env bash
set -e

OUTFILE=${1:-/usr/local/bin/irisctl}
BUILD_DATE=$(date)
BUILD_VERSION=$(cat VERSION)
GIT_TAG=$(git describe --tags)

LDFLAGS="
 -X 'github.com/olegsu/iris/pkg/util.GitTag=${GIT_TAG}'
 -X 'github.com/olegsu/iris/pkg/util.BuildVersion=${BUILD_VERSION}'
 -X 'github.com/olegsu/iris/pkg/util.BuildDate=${BUILD_DATE}'
"
# echo $LDFLAGS
go build -ldflags "$LDFLAGS" -o "$OUTFILE" main.go
