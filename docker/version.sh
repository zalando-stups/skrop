#!/usr/bin/env bash

export CURRENT_VERSION="$(git describe --tags --always)"
echo "CURRENT_VERSION=${CURRENT_VERSION}"
export NEXT_PATCH_VERSION="$(go run docker/version/version.go patch ${CURRENT_VERSION})"
echo "NEXT_PATCH_VERSION=${NEXT_PATCH_VERSION}"
