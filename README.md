# SKROP

A skipper based media service using the vips library.

## Usage

```
go run cmd/skrop/main.go -routes-file eskip/sample.eskip -verbose
```

## Packaging

In order to package skrop for production, you're going to need [Docker](https://docs.docker.com).
To build a Docker image, just run the build script: (the version is optional)

Eg:
```
make docker version=1.0.0
```

Now you can start skrop in a Docker container and pass arguments to it:

```
make docker-run
```
