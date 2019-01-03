# Getting Started

## Creating a routes file

In order to safely use different image transformations, you need to create a routes file.
The routes file uses the [eskip format](WIP)

Here is how a route from the [sample.eskip](../eskip/sample.eskip) file looks like:

```
small: Path("/images/S/:image")
  -> modPath("^/images/S", "/images")
  -> finalizeResponse()
  -> longerEdgeResize(800)
  -> "http://localhost:9090";
```

What this means, is that when somebody does a `GET` to `http://skrop.url/images/S/myimage.jpg`,
Skrop will call `http://localhost:9090/images/myimage.jpg` to retrieve
the image from there, resize it, so that its longer edge is 800px and return
the resized image back as response.

## Start skrop using docker

As soon as the routes file is ready, you can start skrop using docker.

```bash
WIP
```
