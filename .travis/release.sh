#!/usr/bin/env bash

./docker/docker-release.sh

./.travis/git-tag.sh
