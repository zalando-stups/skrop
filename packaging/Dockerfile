FROM golang:1.9.1-alpine

RUN apk add --update git bash

ADD install-vips.sh install-vips.sh
RUN ./install-vips.sh

RUN go get github.com/tools/godep
