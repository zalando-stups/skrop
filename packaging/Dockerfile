FROM golang:1.7.4-alpine

RUN apk add --update git bash

ADD install-vips.sh install-vips.sh
RUN ./install-vips.sh

RUN go get github.com/tools/godep
