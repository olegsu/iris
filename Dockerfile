FROM golang as builder
RUN mkdir -p /build/dist
ADD . /build/
WORKDIR /build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ./hack/build.sh dist/iris

FROM alpine:latest
RUN apk --no-cache add --update ca-certificates
COPY hack/docker_entrypoint.sh /entrypoint.sh
COPY --from=builder /build/dist/iris /usr/local/bin/iris

ENTRYPOINT ["sh", "entrypoint.sh"]
CMD [ "--help" ]
