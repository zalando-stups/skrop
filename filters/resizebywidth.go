package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando-incubator/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	ResizeByWidthName = "width"
)

type resizeByWidth struct {
	width   int
	enlarge bool
}

func NewResizeByWidth() filters.Spec {
	return &resizeByWidth{}
}

func (r *resizeByWidth) Name() string {
	return ResizeByWidthName
}

func (r *resizeByWidth) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for resize by width ", r)

	if !r.enlarge {
		size, err := image.Size()
		if err != nil {
			return nil, err
		}

		// enlargement not allowed here
		if size.Width <= r.width {
			return &bimg.Options{}, nil
		}
	}

	return &bimg.Options{
		Width: r.width}, nil
}

func (s *resizeByWidth) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return other.Width == 0 || other.Width == self.Width
}

func (s *resizeByWidth) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	return other
}

func (r *resizeByWidth) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 && len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resizeByWidth{}

	f.width, err = parse.EskipIntArg(args[0])
	if err != nil {
		return nil, err
	}

	f.enlarge = true

	if len(args) == 2 {
		cons, err := parse.EskipStringArg(args[1])
		if err != nil {
			return nil, err
		}

		f.enlarge = !(cons == doNotEnlarge)
	}

	return f, nil
}

func (r *resizeByWidth) Request(ctx filters.FilterContext) {}

func (r *resizeByWidth) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
