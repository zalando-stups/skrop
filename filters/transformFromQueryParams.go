package filters

import (
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"strconv"
	"strings"
)

const (
	ExtractArea    = "transformFromQueryParams"
	cropParameters = "crop"
)

type transformFromQueryParams struct{}

func NewTransformFromQueryParams() filters.Spec {
	return &transformFromQueryParams{}
}

func (t *transformFromQueryParams) Name() string {
	return ExtractArea
}

func (t *transformFromQueryParams) CreateFilter(config []interface{}) (filters.Filter, error) {
	return t, nil
}

func (t *transformFromQueryParams) CreateOptions(ctx *ImageFilterContext, _ filters.FilterContext) (*bimg.Options, error) {
	// Get crop prams from the request
	params, ok := ctx.Parameters[cropParameters]
	if !ok {
		return &bimg.Options{}, nil
	}

	params = strings.Split(params[0], ",")
	if len(params) != 4 {
		return &bimg.Options{}, nil
	}

	imgSize, err := ctx.Image.Size()
	if err != nil {
		return nil, err
	}

	// Get x
	x, err := strconv.Atoi(strings.TrimSpace(params[0]))
	if err != nil {
		x = 0
	}

	// Get y
	y, err := strconv.Atoi(strings.TrimSpace(params[1]))
	if err != nil {
		y = 0
	}

	// Get height
	height, err := strconv.Atoi(strings.TrimSpace(params[2]))
	if err != nil {
		height = imgSize.Height
	}

	// Get width
	width, err := strconv.Atoi(strings.TrimSpace(params[3]))
	if err != nil {
		width = imgSize.Width
	}

	if y+height > imgSize.Height {
		height = imgSize.Height - y
	}

	if x+width > imgSize.Width {
		width = imgSize.Width - x
	}

	return &bimg.Options{
		Top:        y,
		Left:       x,
		AreaHeight: height,
		AreaWidth:  width,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (t *transformFromQueryParams) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return false
}

func (t *transformFromQueryParams) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.AreaWidth = self.AreaWidth
	other.AreaHeight = self.AreaHeight
	other.Top = self.Top
	other.Left = self.Left
	return other
}

func (t *transformFromQueryParams) Request(ctx filters.FilterContext) {

}

func (e *transformFromQueryParams) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, e)
}
