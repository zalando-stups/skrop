#!/usr/bin/env bash

readonly ROUTES_FILE=${1-"./sample.eskip"}
readonly IMAGE_NAME=${2-"skrop/skrop"}

. docker/version.sh

function build_docker_image() {
  docker pull "$IMAGE_NAME" || true
  docker build --pull --cache-from "$IMAGE_NAME" \
  --build-arg ROUTES_FILE=${ROUTES_FILE} \
  --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
  --build-arg VCS_REF=`git rev-parse --short HEAD` \
  --build-arg VERSION=${NEXT_PATCH_VERSION} \
  -t ${IMAGE_NAME} ./docker \
  && return 0
}

echo >&2 "Building Docker image…"
build_docker_image
if [ "$?" -ne 0 ]; then
  echo >&2 "Aborting!"
  exit 1
fi
echo >&2 "Docker image build successfully…"
