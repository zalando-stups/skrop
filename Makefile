.DEFAULT_GOAL := all

routes_file ?= ./eskip/sample.eskip
docker_tag ?= skrop/skrop

CURRENT_VERSION    = $(shell git describe --tags --always --dirty)
VERSION           ?= $(CURRENT_VERSION)
NEXT_PATCH         = $(shell go run packaging/version/version.go patch $(CURRENT_VERSION))
COMMIT_HASH        = $(shell git rev-parse --short HEAD)


build:
	go build ./cmd/skrop

docker:
	./packaging/build.sh $(VERSION) $(routes_file) $(docker_tag)

docker-run:
	rm -rf "$$(pwd)"/mylocalfilecache
	mkdir "$$(pwd)"/mylocalfilecache
	docker run --rm -v "$$(pwd)"/images:/images -v "$$(pwd)"/mylocalfilecache:/mylocalfilecache -e STRIP_METADATA='TRUE' -p 9090:9090 skrop/skrop -verbose

test: build test-only

test-only:
	go test ./...

init-deps:
	./packaging/build.sh
	go get ./cmd/skrop/

all: init-deps build test

tag:
	echo "Creating tag for version: $(VERSION)"
	git tag $(VERSION) -a -m "Generated tag from TravisCI for build $(TRAVIS_BUILD_NUMBER)"

push-tags:
	git push -q --tags https://$(GITHUB_AUTH)@github.com/zalando-stups/skrop

release-patch:
	echo "Incrementing patch version"
	make VERSION=$(NEXT_PATCH) tag push-tags

ci-user:
	git config --global user.email "builds@travis-ci.com"
	git config --global user.name "Travis CI"

ci-release-patch: ci-user release-patch

ci-test:
	./.travis/test.sh

ci-trigger: ci-test
ifeq ($(TRAVIS_BRANCH)_$(TRAVIS_PULL_REQUEST), master_false)
	echo "Merge to 'master'. Tagging patch version up."
	make ci-release-patch
else
	echo "Not a merge to 'master'. Not versionning this merge."
endif

build-docker-vips:
	docker build -f Dockerfile-Vips --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` --build-arg VCS_REF=`git rev-parse --short HEAD` -t skrop/alpine-mozjpeg-vips:3.3.1-8.7.0 .
	docker push skrop/alpine-mozjpeg-vips:3.3.1-8.7.0
