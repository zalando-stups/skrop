#!/usr/bin/env bash

readonly DOCKER_IMAGE_VERSION=${1-"latest"}
readonly ROUTES_FILE=${2-"eskip/sample.eskip"}

readonly PROJECT_NAME="zalando-incubator/skrop"
readonly GO_PROJECT_NAME="github.com/$PROJECT_NAME"

function make_build_image() {
  docker build -t "$PROJECT_NAME-build" packaging/ \
  && return 0
}

function make_binary() {
  docker run --rm \
    -v ${PWD}:/go/src/${GO_PROJECT_NAME} \
    -e "GOPATH=/go" \
    -e "GOOS=linux" \
    -w /go/src/${GO_PROJECT_NAME} \
    "$PROJECT_NAME-build" sh \
    -c 'godep restore && go build ./cmd/skrop' \
  && return 0
}

function make_production_image() {
  docker build \
  --build-arg ROUTES_FILE=${ROUTES_FILE} \
  -t ${PROJECT_NAME}:${DOCKER_IMAGE_VERSION} . \
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
