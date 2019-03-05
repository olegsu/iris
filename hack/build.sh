#!/bin/bash
set -e
OUTFILE=/usr/local/bin/irisctl
go build -o $OUTFILE main.go

chmod +x $OUTFILE