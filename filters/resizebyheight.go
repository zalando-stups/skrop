package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando-incubator/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	ResizeByHeightName = "height"
)

type resizeByHeight struct {
	height  int
	enlarge bool
}

func NewResizeByHeight() filters.Spec {
	return &resizeByHeight{}
}

func (r *resizeByHeight) Name() string {
	return ResizeByHeightName
}

func (r *resizeByHeight) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for resize by width ", r)

	if !r.enlarge {
		size, err := image.Size()
		if err != nil {
			return nil, err
		}

		// in case the image is bigger than the requested, do not change it
		if size.Height <= r.height {
			return &bimg.Options{}, nil
		}
	}

	return &bimg.Options{
		Height: r.height}, nil
}

func (r *resizeByHeight) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 && len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resizeByHeight{}

	f.height, err = parse.EskipIntArg(args[0])
	if err != nil {
		return nil, err
	}

	if len(args) == 2 {
		var cons string
		cons, err = parse.EskipStringArg(args[1])
		if err != nil {
			return nil, err
		}

		f.enlarge = !(cons == doNotEnlarge)
	} else {
		f.enlarge = true
	}

	return f, nil
}

func (r *resizeByHeight) Request(ctx filters.FilterContext) {}

func (r *resizeByHeight) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
