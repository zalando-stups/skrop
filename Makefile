.DEFAULT_GOAL := all

routes_file ?= ./eskip/sample.eskip
docker_tag ?= zalando-stups/skrop

CURRENT_VERSION    = $(shell git describe --tags --always --dirty)
VERSION           ?= $(CURRENT_VERSION)
NEXT_MAJOR         = $(shell go run packaging/version/version.go major $(CURRENT_VERSION))
NEXT_MINOR         = $(shell go run packaging/version/version.go minor $(CURRENT_VERSION))
NEXT_PATCH         = $(shell go run packaging/version/version.go patch $(CURRENT_VERSION))
COMMIT_HASH        = $(shell git rev-parse --short HEAD)

glide:
	glide install

build: glide
	go build ./cmd/skrop

docker:
	./packaging/build.sh $(VERSION) $(routes_file) $(docker_tag)

docker-run:
	docker run --rm -v "$$(pwd)"/images:/images -p 9090:9090 zalando-stups/skrop -verbose

test: build test-only

test-only:
	go test ./...

init-deps:
	./packaging/build.sh
	go get github.com/Masterminds/glide
	go get ./cmd/skrop/

update-deps:
	glide update

all: init-deps build test

tag:
	git tag $(VERSION)

push-tags:
	git push --tags https://$(GITHUB_AUTH)@github.com/zalando-stups/skrop

release-major:
	echo Incrementing major version
	make VERSION=$(NEXT_MAJOR) tag push-tags

release-minor:
	echo Incrementing minor version
	make VERSION=$(NEXT_MINOR) tag push-tags

release-patch:
	echo Incrementing patch version
	make VERSION=$(NEXT_PATCH) tag push-tags

ci-user:
	git config --global user.email "builds@travis-ci.com"
	git config --global user.name "Travis CI"

ci-release-major: ci-user init-deps release-major
ci-release-minor: ci-user init-deps release-minor
ci-release-patch: ci-user init-deps release-patch

ci-test:
	./.travis/test.sh

ci-trigger: ci-test
ifeq ($(TRAVIS_BRANCH)_$(TRAVIS_PULL_REQUEST)_$(findstring major-release,$(TRAVIS_COMMIT_MESSAGE)), master_false_major-release)
	make ci-release-major
else ifeq ($(TRAVIS_BRANCH)_$(TRAVIS_PULL_REQUEST)_$(findstring minor-release,$(TRAVIS_COMMIT_MESSAGE)), master_false_minor-release)
	make ci-release-minor
else ifeq ($(TRAVIS_BRANCH)_$(TRAVIS_PULL_REQUEST), master_false)
	make ci-release-patch
else
	echo Not a merge to 'master'. Not versionning this merge.
endif
