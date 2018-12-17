#!/usr/bin/env bash
set -e

# go to scripts folder
cd $(dirname $0)

project_root="${PWD}/../"

sh docker-clean.sh

# build image
latest_version=$(git log --pretty=format:'%h' -n 1 || echo '0')
echo "new version: ${latest_version}"

docker build -t skrop-sandbox:${latest_version} ${project_root}
