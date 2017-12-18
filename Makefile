.DEFAULT_GOAL := all

version ?= latest
routes_file ?= ./eskip/sample.eskip
docker_tag ?= zalando-stups/skrop

glide:
	glide install

build: glide
	go build ./cmd/skrop

docker:
	./packaging/build.sh $(version) $(routes_file) $(docker_tag)

docker-run:
	docker run --rm -p 9090:9090 zalando-stups/skrop -verbose

test:
	go test ./...

update-deps:
	glide update

all: build test
