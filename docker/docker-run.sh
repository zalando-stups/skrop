#!/usr/bin/env bash

rm -rf "$$(pwd)"/mylocalfilecache
mkdir "$$(pwd)"/mylocalfilecache
docker run --rm -v "$$(pwd)"/images:/images -v "$$(pwd)"/mylocalfilecache:/mylocalfilecache -e STRIP_METADATA='TRUE' -p 9090:9090 skrop/skrop -verbose
