FROM golang:latest as builder
RUN mkdir -p /go/src/github.com/olegsu/iris
WORKDIR /go/src/github.com/olegsu/iris
COPY . .
RUN "./hack/build.sh"


FROM alpine:3.6

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/olegsu/iris/dist/bin/iris /usr/bin/iris
ENV PATH $PATH:/usr/bin/iris
ENTRYPOINT ["iris"]

CMD ["--help"]