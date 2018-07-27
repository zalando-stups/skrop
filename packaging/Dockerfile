FROM golang:1.10.3-alpine

RUN apk add --update git bash

ADD install-vips.sh install-vips.sh
RUN ./install-vips.sh

RUN go get github.com/Masterminds/glide
