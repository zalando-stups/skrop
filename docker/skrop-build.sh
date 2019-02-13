#!/usr/bin/env bash

PROJECT_DIR="/go/src/github.com/zalando-stups/skrop"

docker pull skrop/skrop-build
docker run -t -v $(pwd):${PROJECT_DIR} skrop/skrop-build sh -c "cd ${PROJECT_DIR} && go build ./cmd/skrop"

if [ "$?" -ne 0 ]; then
  echo >&2 "Error building Skrop. Aborting!"
  exit 1
fi

echo >&2 "Skrop executable was successfully built!"
