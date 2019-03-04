FROM golang:1.11-alpine3.8 as builder

# Add basic tools
RUN apk add --no-cache --update curl bash make git

# upload coverage reports to Codecov.io: pass CODECOV_TOKEN as build-arg
ARG CODECOV_TOKEN
# default codecov bash uploader
ARG CODECOV_BASH_URL=https://codecov.io/bash
# set Codecov expected env
ARG VCS_COMMIT_ID
ARG VCS_BRANCH_NAME
ARG VCS_SLUG
ARG CI_BUILD_URL
ARG CI_BUILD_ID

RUN go get github.com/stretchr/testify/mock
RUN go get github.com/stretchr/testify/assert
RUN mkdir /iris
COPY . /iris
WORKDIR /iris

# Run tests
RUN sh ./hack/test.sh
# Report coverage
RUN if [ "$CODECOV_TOKEN" != "" ]; then curl -s $CODECOV_BASH_URL | bash -s; fi

# Build binary
RUN CGO_ENABLED=0 go build -mod=vendor -v -o "./dist/bin/iris" *.go


FROM alpine:3.8

RUN apk add --no-cache ca-certificates

COPY --from=builder /iris/dist/bin/iris /usr/bin/iris
COPY --from=builder /iris/hack/docker_entrypoint.sh /entrypoint.sh
COPY --from=builder /iris/VERSION /VERSION

ENV PATH $PATH:/usr/bin/iris
ENTRYPOINT ["sh", "entrypoint.sh"]

CMD ["--help"]