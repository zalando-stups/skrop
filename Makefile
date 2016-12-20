.DEFAULT_GOAL := all

godep:
	godep restore

build: godep
	go build ./cmd/skrop

docker:
	./packaging/build.sh $(version)

docker-run:
	docker run --rm -w /app -p 9090:9090 -v $(PWD):/app zalando-incubator/skrop -verbose -routes-file=eskip/sample.eskip

test:
	go test ./...

all: build test
