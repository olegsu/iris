FROM golang:1.10-alpine3.8 as builder

# Add basic tools
RUN apk add --no-cache --update curl bash make git

# upload coverage reports to Codecov.io: pass CODECOV_TOKEN as build-arg
ARG CODECOV_TOKEN
# default codecov bash uploader
ARG CODECOV_BASH_URL=https://codecov.io/bash
# set Codecov expected env
ARG VCS_COMMIT_ID
ARG VCS_BRANCH_NAME

RUN mkdir -p /go/src/github.com/olegsu/iris
WORKDIR /go/src/github.com/olegsu/iris

COPY . .

# Run tests
RUN make test
# Report coverage
RUN if [ "$CODECOV_TOKEN" != "" ]; then curl -s $CODECOV_BASH_URL | bash -s; fi

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