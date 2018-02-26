package filters

import (
	"errors"
	"fmt"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"strconv"
	"strings"
)

const (
	ExtractArea = "extractArea"
	cropFrom    = "point"
	cropWidth   = "width"
	cropHeight  = "height"
)

type extractArea struct{}

func NewExtractArea() filters.Spec {
	return &extractArea{}
}

func (t *extractArea) Name() string {
	return ExtractArea
}

func (t *extractArea) CreateFilter(config []interface{}) (filters.Filter, error) {
	return t, nil
}

func (t *extractArea) CreateOptions(img *bimg.Image, options map[string][]string) (*bimg.Options, error) {

	// GET starting position to crop from
	point, ok := options[cropFrom]
	if !ok {
		return nil, errors.New(fmt.Sprintf("could not find information about where to start the crop from."))
	}
	point = strings.Split(point[0], ",")
	if len(point) != 2 {
		return nil, errors.New(fmt.Sprintf("received %d number of values for the field %s. Expected is 2", len(point), cropFrom))
	}

	x, err := strconv.Atoi(point[0])
	if err != nil {
		return nil, err
	}

	y, err := strconv.Atoi(point[1])
	if err != nil {
		return nil, err
	}

	// Get height of the area to be cropped
	vals, ok := options[cropHeight]
	if !ok {
		return nil, errors.New(fmt.Sprintf("could not find information about height of the area to crop."))
	}
	if len(vals) != 1 {
		return nil, errors.New(fmt.Sprintf("received invalid number of values for the field %s, %d ", cropHeight, vals))
	}

	height, err := strconv.Atoi(vals[0])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("received invalid value for the field %s with value %s", cropHeight, vals[0]))
	}

	// Get width of the area to be cropped
	vals, ok = options[cropWidth]
	if !ok {
		return nil, errors.New(fmt.Sprintf("could not find information about height of the area to crop."))
	}
	if len(vals) != 1 {
		return nil, errors.New(fmt.Sprintf("received invalid number of values for the field %s, %d ", cropWidth, vals))
	}

	width, err := strconv.Atoi(vals[0])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("received invalid value for the field %s with value %s", cropWidth, vals[0]))
	}

	imgSize, err := img.Size()
	if err != nil {
		return nil, err
	}

	return &bimg.Options{
		Top:        y,
		Left:       x,
		AreaHeight: min(imgSize.Height, height),
		AreaWidth:  min(imgSize.Width, width),
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (t *extractArea) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return false
}

func (t *extractArea) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.AreaWidth = self.AreaWidth
	other.AreaHeight = self.AreaHeight
	other.Top = self.Top
	other.Left = self.Left
	return other
}

func (t *extractArea) Request(ctx filters.FilterContext) {

}

func (e *extractArea) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, e)
}
