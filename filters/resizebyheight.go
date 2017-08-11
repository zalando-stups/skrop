package filters

import (
	log "github.com/sirupsen/logrus"
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

		// enlargement not allowed here
		if size.Height <= r.height {
			return &bimg.Options{}, nil
		}
	}

	return &bimg.Options{
		Height: r.height}, nil
}

func (s *resizeByHeight) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return other.Height == 0 || other.Height == self.Height
}

func (s *resizeByHeight) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Height = self.Height
	return other
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

func (r *resizeByHeight) Request(ctx filters.FilterContext) {}

func (r *resizeByHeight) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
