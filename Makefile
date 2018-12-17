.DEFAULT_GOAL := all

routes_file ?= ./eskip/sample.eskip
docker_tag ?= skrop/skrop

build:
	./docker/skrop-build.sh

docker:
	./docker/docker-build.sh

docker-run:
	./docker/docker-run.sh

test:
	go test ./...

build-docker-vips:
	docker build -f Dockerfile-Vips --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` --build-arg VCS_REF=`git rev-parse --short HEAD` -t skrop/alpine-mozjpeg-vips:3.3.1-8.7.0 .
	docker push skrop/alpine-mozjpeg-vips:3.3.1-8.7.0
