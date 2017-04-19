FROM alpine:3.4

RUN apk update && \
  apk add \
    ca-certificates \
    mailcap
    build-base \
    gcc \
    abuild \
    binutils && \
  rm -rf /var/cache/apk/*

ADD ims /bin/
ENTRYPOINT ["/bin/ims"]
