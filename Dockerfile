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
ENV PATH $PATH:/usr/bin/iris
ENTRYPOINT ["iris"]

CMD ["--help"]