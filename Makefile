.DEFAULT_GOAL := all

version ?= latest
routes_file ?= ./eskip/sample.eskip
docker_tag ?= zalando-incubator/skrop

godep:
	godep restore

build: godep
	go build ./cmd/skrop

docker:
	./packaging/build.sh $(version) $(routes_file) $(docker_tag)

docker-run:
	docker run --rm -p 9090:9090 zalando-incubator/skrop -verbose

test:
	go test ./...

all: build test
