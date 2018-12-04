#!/usr/bin/env bash

readonly DOCKER_IMAGE_VERSION=${1-"latest"}
readonly ROUTES_FILE=${2-"eskip/sample.eskip"}
readonly DOCKER_IMAGE_NAME=${3-"zalando-stups/skrop"}

readonly GO_PROJECT_NAME="github.com/zalando-stups/skrop"

function make_production_image() {
  docker build \
  --build-arg ROUTES_FILE=${ROUTES_FILE} \
  -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION} . \
  && return 0
}

echo >&2 "Building Docker image…"

make_production_image
if [ "$?" -ne 0 ]; then
  echo >&2 "Aborting!"
  exit 1
fi

echo >&2 "Success!"
