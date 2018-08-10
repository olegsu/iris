FROM golang:1.10-alpine3.8 as builder

RUN mkdir -p /go/src/github.com/olegsu/iris
WORKDIR /go/src/github.com/olegsu/iris

RUN apk add --update make

COPY . .

# Run tests
RUN make test

# Build binary
RUN make build


FROM alpine:3.6

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/olegsu/iris/dist/bin/iris /usr/bin/iris
COPY --from=builder /go/src/github.com/olegsu/iris/hack/docker_entrypoint.sh /entrypoint.sh
COPY --from=builder /go/src/github.com/olegsu/iris/VERSION /VERSION

ENV PATH $PATH:/usr/bin/iris
ENTRYPOINT ["sh", "entrypoint.sh"]

CMD ["--help"]