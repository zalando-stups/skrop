package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
	"math"
)

const (
	North  = "north"
	South  = "south"
	East   = "east"
	West   = "west"
	Center = "center"
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

func handleResponse(ctx filters.FilterContext, f ImageFilter) {
	rsp := ctx.Response()

	rsp.Header.Del("Content-Length")

	in := rsp.Body
	r, w := io.Pipe()
	rsp.Body = r

	go handleImageTransform(w, in, f)
}

func handleImageTransform(out *io.PipeWriter, in io.ReadCloser, f ImageFilter) error {
	defer func() {
		in.Close()
	}()

	imageBytes, err := ioutil.ReadAll(in)

	if err != nil {
		return err
	}

	log.Debug("Image bytes length: ", len(imageBytes))

	responseImage := bimg.NewImage(imageBytes)

	options, err := f.CreateOptions(responseImage)

	if err != nil {
		return err
	}

	return transformImage(out, responseImage, options)
}

func transformImage(out *io.PipeWriter, image *bimg.Image, opts *bimg.Options) error {
	var err error

	defer func() {
		if err == nil {
			err = io.EOF
		}
		out.CloseWithError(err)
	}()

	newImage, err := image.Process(*opts)

	if err != nil {
		return err
	}

	_, err = out.Write(newImage)

	return err
}

func parseEskipIntArg(arg interface{}) (int, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return int(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func parseEskipUint8Arg(arg interface{}) (uint8, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return uint8(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}

func parseEskipStringArg(arg interface{}) (string, error) {
	if str, ok := arg.(string); ok {
		return string(str), nil
	} else {
		return "", filters.ErrInvalidFilterParameters
	}
}