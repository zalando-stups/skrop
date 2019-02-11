## Filters
Skrop provides a set of filters, which you can use within the routes:

* **longerEdgeResize(size)** — resizes the image to have the longer edge as specified, while at the same time preserving the aspect ratio
* **crop(width, height, type)** — crops the image to have the specified width and height the type can be "north", "south", "east" and "west"
* **cropByHeight(height, type)** — crops the image to have the specified height
* **cropByWidth(width, type)** — crops the image to have the specified width
* **resize(width, height, opt-keep-aspect-ratio)** — resizes an image. Third parameter is optional: "ignoreAspectRatio" to ignore the aspect ratio, anything else to keep it
* **addBackground(R, G, B)** — adds the background to a PNG image with transparency
* **convertImageType(type)** — converts between different formats (for the list of supported types see [here](https://github.com/h2non/bimg/blob/master/type.go)
* **sharpen(radius, X1, Y2, Y3, M1, M2)** — sharpens the image (for info about the meaning of the parameters and the suggested values see [here](http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/libvips-convolution.html#vips-sharpen))
* **width(size, opt-enlarge)** — resizes the image to the specified width keeping the ratio. If the second arg is specified and it is equals to "DO_NOT_ENLARGE", the image will not be enlarged
* **height(size, opt-enlarge)** — resizes the image to the specified height keeping the ratio. If the second arg is specified and it is equals to "DO_NOT_ENLARGE", the image will not be enlarged
* **blur(sigma, min_ampl)** — blurs the image (for info see [here](http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/libvips-convolution.html#vips-gaussblur))
* **imageOverlay(filename, opacity, gravity, opt-top-margin, opt-right-margin, opt-bottom-margin, opt-left-margin)** — puts an image onverlay over the required image
* **transformByQueryParams()** - transforms the image based on the request query parameters (supports only crop for now) e.g: localhost:9090/images/S/big-ben.jpg?crop=120,300,500,300.
* **cropByFocalPoint(targetX, targetY, aspectRatio, minWidth)** — crops the image based on a focal point on both the source as well as on the target and desired aspect ratio of the target. TargetX and TargetY are the definition of the target image focal point defined as relative values for both width and height, i.e. if the focal point of the target image should be right in the center it would be 0.5 and 0.5. This filter expects two PathParams named **focalPointX** and **focalPointY** which are absolute X and Y coordinates of the focal point in the source image. The fourth parameter is optional; when given the filter will ensure that the resulting image has at least the specified minimum width if not it will crop the biggest possible part based on the focal point.

_Note:_ As Skrop is built on top of Skipper, it supports all filters supported by [Skipper](https://www.github.com/zalando/skipper) as well.

### About filters
The eskip file defines a list of configurations. Every configuration is composed by a route and a list of filters to
apply for that specific route. The filters are applied starting with the last one, to the first one.

The `finalizeResponse()` filter needs to be added at the end of the pipeline (beginning of the route),
because it triggers the last transformation of the image.

Because of performance, most of the filters don't not trigger a transformation of the image, but if possible it is
merged with the result of the previous filters. The image is actually transformed every time the filter cannot be
merged with the previous one e.g. both edit the same attribute and also at the end of the filter chain by the 
`finalizeResponse()` filter.

## Metadata
By default metadata are kept in the processed images. If you are not interested in metadata and 
you want them stripped from all the images that are processed, you can add the following 
environment variable to the running system:

```
STRIP_METADATA=TRUE
``` 
