package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
	"math"
)

const (
	// ResizeName is the name of the filter
	ResizeName           = "resize"
	ignoreAspectRatioStr = "ignoreAspectRatio"
)

type resize struct {
	width           int
	height          int
	keepAspectRatio bool
}

// NewResize creates a new filter of this type
func NewResize() filters.Spec {
	return &resize{}
}

func (f *resize) Name() string {
	return ResizeName
}

func (f *resize) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for resize ", f)

	if !f.keepAspectRatio {
		return &bimg.Options{
			Width:  f.width,
			Height: f.height,
			Force:  true}, nil
	}

	size, err := imageContext.Image.Size()
	if err != nil {
		return nil, err
	}

	// calculate height keeping width
	ht := int(math.Floor(float64(size.Height*f.width) / float64(size.Width)))

	// if height is less or equal than desired, return transform by width
	if ht <= f.height {
		return &bimg.Options{
			Width: f.width}, nil
	}
	// otherwise transform by height
	return &bimg.Options{
		Height: f.height}, nil

}

func (f *resize) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.AreaWidth == 0 && other.AreaHeight == 0 ) && ((other.Width == 0 && other.Height == 0) ||
		(self.Width == other.Width && self.Height == other.Height))
}

func (f *resize) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	return other
}

func (f *resize) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 2 && len(args) != 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &resize{}

	c.width, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	c.height, err = parse.EskipIntArg(args[1])

	if err != nil {
		return nil, err
	}

	if len(args) == 3 {
		ratio, err := parse.EskipStringArg(args[2])
		if err != nil {
			return nil, err
		}

		c.keepAspectRatio = !(ratio == ignoreAspectRatioStr)

	} else {
		c.keepAspectRatio = true
	}

	return c, nil
}

func (f *resize) Request(ctx filters.FilterContext) {}

func (f *resize) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
