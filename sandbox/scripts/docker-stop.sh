#!/usr/bin/env bash
set -e

docker rm $(docker stop $(docker ps -a -q --filter="label=artifactId=skrop-sandbox"))
echo 'container removed'
