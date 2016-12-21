# Skrop &nbsp; [![Build Status](https://travis-ci.org/zalando-incubator/skrop.svg?branch=master)](https://travis-ci.org/zalando-incubator/skrop) [![codecov](https://codecov.io/gh/zalando-incubator/skrop/branch/master/graph/badge.svg)](https://codecov.io/gh/zalando-incubator/skrop) [![Go Report Card](https://goreportcard.com/badge/github.com/zalando-incubator/skrop)](https://goreportcard.com/report/github.com/zalando-incubator/skrop) [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/zalando-incubator/skrop/master/LICENSE)

Skrop is a media service based on [Skipper](https://github.com/zalando/skipper) and the [vips](https://github.com/jcupitt/libvips) library.

## Usage

```
go run cmd/skrop/main.go -routes-file eskip/sample.eskip -verbose
```

## Packaging

In order to package skrop for production, you're going to need [Docker](https://docs.docker.com).
To build a Docker image, just run the build script (the `version` and `routes_file` arguments are optional):

```
make docker version=1.0.0 routes_file=eskip/sample.eskip
```

Now you can start Skrop in a Docker container:

```
make docker-run
```
