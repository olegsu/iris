FROM alpine:3.8

RUN apk add --update ca-certificates

COPY dist/linux_386/iris /usr/local/bin/

ENTRYPOINT [ "iris" ]

CMD [ "--help" ]