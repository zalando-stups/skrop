#!/usr/bin/env bash
set -e

# remove containers
containers=$(docker ps -a -f "label=artifactId=skrop-sandbox" -q)
if [ ! -z "$containers" ]; then
    docker rm -f ${containers}
fi

# remove images
images=$(docker images -f "label=artifactId=skrop-sandbox" -q)
if [ ! -z "$images" ]; then
    docker rmi -f ${images}
fi
