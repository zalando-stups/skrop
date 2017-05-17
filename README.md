# Skrop &nbsp; [![Build Status](https://travis-ci.org/zalando-incubator/skrop.svg?branch=master)](https://travis-ci.org/zalando-incubator/skrop) [![codecov](https://codecov.io/gh/zalando-incubator/skrop/branch/master/graph/badge.svg)](https://codecov.io/gh/zalando-incubator/skrop) [![Go Report Card](https://goreportcard.com/badge/github.com/zalando-incubator/skrop)](https://goreportcard.com/report/github.com/zalando-incubator/skrop) [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://raw.githubusercontent.com/zalando-incubator/skrop/master/LICENSE)

Skrop is a media service based on [Skipper](https://github.com/zalando/skipper) and the [vips](https://github.com/jcupitt/libvips) library.

## Usage

In order to be able to use Skrop, you have to be familiar with how
[Skipper](https://github.com/zalando/skipper) works.

### Install dependencies:
```
./packaging/build.sh
go get github.com/tools/godep
go get ./cmd/skrop/
```
### Run Skrop
```
go run cmd/skrop/main.go -routes-file eskip/sample.eskip -verbose
```
### Test
To test if everything is configured correctly you should open in your browser
```
http://localhost:9090/images/big-ben.jpg
```
and the resized version
```
http://localhost:9090/images/S/big-ben.jpg
```

Here is how a route from the `sample.eskip` file would look like:
 
```
small: Path("/images/S/:image")
  -> modPath("^/images/S", "/images")
  -> longerEdgeResize(800)
  -> "http://localhost:9090";
```

What it means, is that when somebody does a `GET` to `http://skrop.url/images/S/myimage.jpg`,
Skrop will call `http://localhost:9090/images/myimage.jpg` to retrieve
the image from there, resize it, so that its longer edge is 800px and return 
the resized image back in the response.


## Filters

Skrop provides a set of filters, which you can use within the routes:

```
longerEdgeResize(size) — resizes the image to have the longer edge as specified,
                         while at the same time preserving the aspect ratio 
crop(width, height, type) — crops the image to have the specified width and height
                            the type can be "north", "south", "east" and "west"
cropByHeight(height, type) — crops the image to have the specified height
cropByWidth(width, type) — crops the image to have the specified width
resize(width, height) — resizes an image without preserving its aspect ratio
```

## Packaging

In order to package skrop for production, you're going to need [Docker](https://docs.docker.com).
To build a Docker image, just run the build script (the arguments are optional):


```
make docker version=1.0.0 routes_file=eskip/sample.eskip docker_tag=zalando-incubator/skrop
```

Now you can start Skrop in a Docker container:

```
make docker-run
```
