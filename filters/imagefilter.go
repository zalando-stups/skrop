package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"image/color"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"strconv"
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

func parseEskipRGBAArg(arg interface{}) (color.RGBA, error) {
	str, ok := arg.(string)
	if !ok {
		return color.RGBA{}, filters.ErrInvalidFilterParameters
	}
	str = strings.TrimPrefix(str, "#")
	r, g, b, a, err := parseRGBA(str)

	if err != nil {
		return color.RGBA{}, filters.ErrInvalidFilterParameters
	}
	return color.RGBA{R: r, G: g, B: b, A: a}, nil
}

func parseRGBA(str string) (uint8, uint8, uint8, uint8, error) {
	if len(str) < 6 {
		return 0, 0, 0, 0, filters.ErrInvalidFilterParameters
	}

	var red, green, blue, alpha uint8

	r, err1 := strconv.ParseUint(str[0:2], 16, 8)
	g, err2 := strconv.ParseUint(str[2:4], 16, 8)
	b, err3 := strconv.ParseUint(str[4:6], 16, 8)

	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, 0, filters.ErrInvalidFilterParameters
	}

	red = uint8(r)
	green = uint8(g)
	blue = uint8(b)

	if len(str) == 8 {
		a, err4 := strconv.ParseUint(str[6:8], 16, 8)
		if err4 != nil {
			return 0, 0, 0, 0, filters.ErrInvalidFilterParameters
		}

		alpha = uint8(a)
	} else {
		alpha = 255
	}

	return red, green, blue, alpha, nil

}
