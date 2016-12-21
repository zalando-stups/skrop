.DEFAULT_GOAL := all

godep:
	godep restore

build: godep
	go build ./cmd/skrop

docker:
	./packaging/build.sh $(version) $(routes_file)

docker-run:
	docker run --rm -p 9090:9090 zalando-incubator/skrop -verbose

test:
	go test ./...

all: build test
