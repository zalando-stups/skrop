#!/usr/bin/env bash

readonly DOCKER_IMAGE_VERSION=${1-"latest"}

function make_build_image() {
  docker build -t zalando-incubator/skrop-build packaging/ \
  && return 0
}

function make_binary() {
  docker run --rm \
    -v ${PWD}:/go/src/github.com/zalando-incubator/skrop \
    -e "GOPATH=/go" \
    -e "GOOS=linux" \
    -w /go/src/github.com/zalando-incubator/skrop \
    zalando-incubator/skrop-build sh \
    -c 'godep restore && go build --ldflags="-s" ./cmd/skrop' \
  && return 0
}

function make_production_image() {
  docker build -t zalando-incubator/skrop:${DOCKER_IMAGE_VERSION} . \
  && return 0
}

echo >&2 "Building build Docker image…"

make_build_image
if [ "$?" -ne 0 ]; then
  echo >&2 "Aborting!"
  exit 1
fi

echo >&2 "Building binary…"

make_binary
if [ "$?" -ne 0 ]; then
  echo >&2 "Aborting!"
  exit 1
fi

echo >&2 "Building Docker image…"

make_production_image
if [ "$?" -ne 0 ]; then
  echo >&2 "Aborting!"
  exit 1
fi

echo >&2 "Success!"
