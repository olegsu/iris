#!/bin/bash

set -e
set -o pipefail

rm -rf .cover/ .test/
mkdir .cover/ .test/
trap "rm -rf .test/" EXIT
for pkg in `go list ./... | grep -v /vendor/`; do
    go test -v -covermode=atomic \
        -coverprofile=".cover/$(echo $pkg | sed 's/\//_/g').cover.out" $pkg
done

echo "mode: set" > .cover/cover.out && cat .cover/*.cover.out | grep -v mode: | sort -r | \
   awk '{if($1 != last) {print $0;last=$1}}' >> .cover/cover.out

go tool cover -html=.cover/cover.out -o=.cover/coverage.html

CODECOV_BASH_URL=https://codecov.io/bash
if [ "$CODECOV_TOKEN" != "" ]; then curl -s $CODECOV_BASH_URL | bash -s; fi