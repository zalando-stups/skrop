package filters

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
	"bytes"
)

const (
	North        = "north"
	South        = "south"
	East         = "east"
	West         = "west"
	Center       = "center"
	Quality      = 100
	doNotEnlarge = "DO_NOT_ENLARGE"
)

var (
	cropTypeToGravity map[string]bimg.Gravity
	cropTypes         map[string]bool
)

func init() {
	cropTypes = map[string]bool{
		North:  true,
		South:  true,
		East:   true,
		West:   true,
		Center: true}
	cropTypeToGravity = map[string]bimg.Gravity{
		North:  bimg.GravityNorth,
		South:  bimg.GravitySouth,
		East:   bimg.GravityEast,
		West:   bimg.GravityWest,
		Center: bimg.GravityCentre}
}

type ImageFilter interface {
	CreateOptions(image *bimg.Image) (*bimg.Options, error)
}

func HandleImageResponse(ctx filters.FilterContext, f ImageFilter) {
	rsp := ctx.Response()

	rsp.Header.Del("Content-Length")

	handleImageTransform(rsp, f)
}

func handleImageTransform(rsp *http.Response, f ImageFilter) error {
	var err error
	imageBytes, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return err
	}

	imageBytesLength := len(imageBytes)

	log.Debug("Image bytes length: ", imageBytesLength)

	if imageBytesLength == 0 {
		return errors.New("original image is empty. nothing to process")
	}

	originalImage := bimg.NewImage(imageBytes)

	options, err := f.CreateOptions(originalImage)

	if err != nil {
		return err
	}

	err = transformImage(rsp, originalImage, options)
	return err
}

func transformImage(rsp *http.Response, image *bimg.Image, opts *bimg.Options) error {
	defOpt := applyDefaults(opts)
	transformedImageBytes, err := image.Process(*defOpt)

	if err != nil {
		return err
	}

	rsp.Body = ioutil.NopCloser(bytes.NewReader(transformedImageBytes))

	return nil
}

func applyDefaults(o *bimg.Options) *bimg.Options {
	if o.Quality == 0 {
		o.Quality = Quality
	}
	return o
}
