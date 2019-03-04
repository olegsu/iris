FROM alpine:3.8

RUN apk add --update ca-certificates

COPY dist/linux_386/iris /usr/local/bin/
COPY hack/docker_entrypoint.sh /entrypoint.sh
COPY VERSION /VERSION

ENTRYPOINT ["sh", "entrypoint.sh"]
CMD [ "--help" ]