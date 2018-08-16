#!/bin/bash

for pkg in `go list ./... | grep -v /vendor/`; do
    go test -v --cover $pkg
done