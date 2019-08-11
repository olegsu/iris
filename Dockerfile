FROM golang as builder
RUN mkdir -p /build/dist
ADD . /build/
WORKDIR /build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/iris main.go

FROM alpine:latest
RUN apk --no-cache add --update ca-certificates
COPY --from=builder /build/dist/iris /usr/local/bin/iris
COPY hack/docker_entrypoint.sh /entrypoint.sh
COPY VERSION /VERSION

ENTRYPOINT ["sh", "entrypoint.sh"]
CMD [ "--help" ]
