package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
)

// ResizeByWidthName is the name of the filter
const ResizeByWidthName = "width"

type resizeByWidth struct {
	width   int
	enlarge bool
}

// NewResizeByWidth creates a new filter of this type
func NewResizeByWidth() filters.Spec {
	return &resizeByWidth{}
}

func (f *resizeByWidth) Name() string {
	return ResizeByWidthName
}

func (f *resizeByWidth) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for resize by width ", f)

	if !f.enlarge {
		size, err := imageContext.Image.Size()
		if err != nil {
			return nil, err
		}

		// enlargement not allowed here
		if size.Width <= f.width {
			return &bimg.Options{}, nil
		}
	}

	return &bimg.Options{
		Width: f.width}, nil
}

func (f *resizeByWidth) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return other.Width == 0 || other.Width == self.Width
}

func (f *resizeByWidth) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	return other
}

func (f *resizeByWidth) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 && len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	r := &resizeByWidth{}

	r.width, err = parse.EskipIntArg(args[0])
	if err != nil {
		return nil, err
	}

	r.enlarge = true

	if len(args) == 2 {
		cons, err := parse.EskipStringArg(args[1])
		if err != nil {
			return nil, err
		}

		r.enlarge = !(cons == doNotEnlarge)
	}

	return r, nil
}

func (f *resizeByWidth) Request(ctx filters.FilterContext) {}

func (f *resizeByWidth) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
