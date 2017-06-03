package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/parse"
)

const (
	ResizeByWidthName = "width"
)

type resizeByWidth struct {
	width int
}

func NewResizeByWidth() filters.Spec {
	return &resizeByWidth{}
}

func (r *resizeByWidth) Name() string {
	return ResizeByWidthName
}

func (r *resizeByWidth) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for resize by width ", r)

	return &bimg.Options{
		Width: r.width}, nil
}

func (r *resizeByWidth) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resizeByWidth{}

	f.width, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *resizeByWidth) Request(ctx filters.FilterContext) {}

func (r *resizeByWidth) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
