package filters

import (
	"image/color"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	AddBackgroundName = "addBackground"
)

type addBackground color.RGBA

func NewAddBackground() filters.Spec {
	return &addBackground{}
}

func (r *addBackground) Name() string {
	return AddBackgroundName
}

func (r *addBackground) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	
}

func (r *addBackground) CreateFilter(args []interface{}) (filters.Filter, error) {
	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	col, err := parseEskipRGBAArg(args[0])
	f := addBackground(col)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *addBackground) Request(ctx filters.FilterContext) {}

func (r *addBackground) Response(ctx filters.FilterContext) {
	handleResponse(ctx, r)
}