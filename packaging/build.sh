#!/usr/bin/env bash

readonly DOCKER_IMAGE_VERSION=${1-"latest"}

echo >&2 "Building build Docker image…"

docker build -t zalando-incubator/skrop-build packaging/

echo >&2 "Building binary…"

docker run --rm \
  -v ${PWD}:/go/src/github.com/zalando-incubator/skrop \
  -e "GOPATH=/go" \
  -e "GOOS=linux" \
  -w /go/src/github.com/zalando-incubator/skrop \
  zalando-incubator/skrop-build sh \
  -c 'godep restore && go build --ldflags="-s" ./cmd/skrop'

echo >&2 "Building Docker image…"

docker build -t zalando-incubator/skrop:${DOCKER_IMAGE_VERSION} .
