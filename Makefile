# Change this and commit to create new release
VERSION=0.0.1

install:
	@glide install

build:
	@CGO_ENABLED=0 go build -v -o "./dist/bin/iris" *.go