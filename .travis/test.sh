#!/usr/bin/env bash

set -e

function cover_all {
  for d in $(go list ./... | grep -v vendor); do
    local name="$(echo "$d" | grep -o "[^/]*$")"
    echo >&2 "Testing ${d}…"
    go test -coverprofile="$name.coverage.txt" -covermode=atomic ${d}
  done
}

cover_all
gocovmerge *.coverage.txt > coverage.txt
