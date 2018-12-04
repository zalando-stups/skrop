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

RUN go build ./cmd/skrop

# final stage
FROM danpersa/alpine-mozjpeg-vips:3.3.1-8.7.0

ARG ROUTES_FILE

COPY --from=0 /go/src/github.com/zalando-stups/skrop/skrop /usr/local/bin/
ADD $ROUTES_FILE skrop.eskip

ENTRYPOINT skrop -routes-file skrop.eskip ${SKROP_ARGS}
