package filters

import (
	"github.com/h2non/bimg"
	log "github.com/sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"strconv"
	"strings"
)

const (
	TransformByQueryParamsName = "transformByQueryParams"
	cropParameters             = "crop"
)

type transformFromQueryParams struct{}

func NewTransformFromQueryParams() filters.Spec {
	return &transformFromQueryParams{}
}

func (t *transformFromQueryParams) Name() string {
	return TransformByQueryParamsName
}

func (t *transformFromQueryParams) CreateFilter(config []interface{}) (filters.Filter, error) {
	return t, nil
}

func (t *transformFromQueryParams) CreateOptions(ctx *ImageFilterContext) (*bimg.Options, error) {
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

	log.Debugf("Crop options x=%d y=%d width=%d height=%d\n", x, y, width, height)

	return &bimg.Options{
		Top:        y,
		Left:       x,
		AreaHeight: height,
		AreaWidth:  width,
	}, nil
}

func (t *transformFromQueryParams) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return false
}

func (t *transformFromQueryParams) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	return nil
}

func (t *transformFromQueryParams) Request(ctx filters.FilterContext) {
}

func (e *transformFromQueryParams) Response(ctx filters.FilterContext) {
	log.Debugf("Response %s\n", TransformByQueryParamsName)
	HandleImageResponse(ctx, e)
	FinalizeResponse(ctx)
}
