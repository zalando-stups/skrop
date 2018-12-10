package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
)

// ResizeByHeightName is the name of the filter
const ResizeByHeightName = "height"

type resizeByHeight struct {
	height  int
	enlarge bool
}

// NewResizeByHeight creates a new filter of this type
func NewResizeByHeight() filters.Spec {
	return &resizeByHeight{}
}

func (f *resizeByHeight) Name() string {
	return ResizeByHeightName
}

func (f *resizeByHeight) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for resize by width ", f)

	if !f.enlarge {
		size, err := imageContext.Image.Size()
		if err != nil {
			return nil, err
		}

		// enlargement not allowed here
		if size.Height <= f.height {
			return &bimg.Options{}, nil
		}
	}

	return &bimg.Options{
		Height: f.height}, nil
}

func (f *resizeByHeight) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return other.Height == 0 || other.Height == self.Height
}

func (f *resizeByHeight) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Height = self.Height
	return other
}

func (f *resizeByHeight) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 && len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &resizeByHeight{}

	c.height, err = parse.EskipIntArg(args[0])
	if err != nil {
		return nil, err
	}

	c.enlarge = true

	if len(args) == 2 {
		cons, err := parse.EskipStringArg(args[1])
		if err != nil {
			return nil, err
		}

		c.enlarge = !(cons == doNotEnlarge)
	}

	return c, nil
}

func (f *resizeByHeight) Request(ctx filters.FilterContext) {}

func (f *resizeByHeight) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
