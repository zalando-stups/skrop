package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
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

	imgSize, err := image.Size()

	if err != nil {
		return nil, err
	}

	// resize only if the original width is bigger than the required width
	if imgSize.Width > r.width {
		return &bimg.Options{
			Width: r.width}, nil
	} else {
		return &bimg.Options{}, nil
	}
}

func (r *resizeByWidth) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resizeByWidth{}

	f.width, err = parseEskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *resizeByWidth) Request(ctx filters.FilterContext) {}

func (r *resizeByWidth) Response(ctx filters.FilterContext) {
	handleResponse(ctx, r)
}
