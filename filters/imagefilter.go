package filters

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"bytes"
	"io"
	"net/http"
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

	defer rsp.Body.Close()

	rsp.Header.Del("Content-Length")

	r, err := handleImageTransform(rsp.Body, f)

	if err != nil {
		log.Error("failed to process image ", err.Error())
		ctx.Serve(&http.Response{
			StatusCode:http.StatusInternalServerError,
			Body: ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		})
		return
	}

	rsp.Body = ioutil.NopCloser(r)
}

func handleImageTransform(r io.Reader, f ImageFilter) (io.Reader, error) {
	var err error
	imageBytes, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	imageBytesLength := len(imageBytes)

	log.Debug("Image bytes length: ", imageBytesLength)

	if imageBytesLength == 0 {
		return nil, errors.New("original image is empty. nothing to process")
	}

	originalImage := bimg.NewImage(imageBytes)

	options, err := f.CreateOptions(originalImage)

	if err != nil {
		return nil, err
	}

	transBytes, err := transformImage(originalImage, options)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(transBytes), nil
}

func transformImage(image *bimg.Image, opts *bimg.Options) ([]byte, error) {
	defOpt := applyDefaults(opts)

	transformedImageBytes, err := image.Process(*defOpt)

	if err != nil {
		return nil, err
	}

	return transformedImageBytes, nil
}

func applyDefaults(o *bimg.Options) *bimg.Options {
	if o.Quality == 0 {
		o.Quality = Quality
	}
	if o.Background == bimg.ColorBlack {
		o.Background = bimg.Color{R: 255, G: 255, B: 255}
	}
	return o
}
