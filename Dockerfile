# build stage
FROM golang:1.11.1-alpine AS build-env

ENV GOPATH /go
ENV GOOS linux

RUN apk add --update git bash

ADD packaging/install-vips.sh install-vips.sh
RUN ./install-vips.sh

ADD . /go/src/github.com/zalando-stups/skrop
WORKDIR /go/src/github.com/zalando-stups/skrop

RUN go get github.com/Masterminds/glide

RUN glide install

RUN go build ./cmd/skrop

# final stage
FROM skrop/alpine-mozjpeg-vips:3.3.1-8.7.0

ARG ROUTES_FILE

# Build-time metadata as defined at http://label-schema.org
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="Skrop" \
      org.label-schema.description="Image transformation service using libvips, based on Skipper." \
      org.label-schema.url="https://github.com/zalando-stups/skrop" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/zalando-stups/skrop" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0"

COPY --from=0 /go/src/github.com/zalando-stups/skrop/skrop /usr/local/bin/
ADD $ROUTES_FILE skrop.eskip

ENTRYPOINT skrop -routes-file skrop.eskip ${SKROP_ARGS}
