FROM alpine:3.8

RUN apk --no-cache add --update ca-certificates

COPY dist/iris_linux_386/iris /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/iris"]
CMD [ "--help" ]
