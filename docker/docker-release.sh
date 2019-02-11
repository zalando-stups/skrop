#!/usr/bin/env bash

readonly IMAGE_NAME=${1-"skrop/skrop"}

. docker/version.sh

./docker/skrop-build.sh

cp skrop ./docker
cp ./eskip/sample.eskip ./docker

./docker/docker-build.sh

echo >&2 "Tagging Docker images…"

docker tag "$IMAGE_NAME" "${IMAGE_NAME}:latest"
docker tag "$IMAGE_NAME" "${IMAGE_NAME}:${NEXT_PATCH_VERSION}"

echo >&2 "Logging in to Docker Hub…"

echo "$DOCKER_USERNAME"
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin

echo >&2 "Pushing Docker images…"

docker push "${IMAGE_NAME}:latest"
docker push "${IMAGE_NAME}:${NEXT_PATCH_VERSION}"

echo >&2 "Success!"
