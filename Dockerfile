FROM alpine:3.8

RUN apk --no-cache add --update ca-certificates

COPY dist/iris_linux_386/iris /usr/local/bin/
COPY hack/docker_entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "entrypoint.sh"]
CMD [ "--help" ]
