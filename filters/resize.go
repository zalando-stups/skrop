package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando-incubator/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"math"
)

const (
	ResizeName           = "resize"
	ignoreAspectRatioStr = "ignoreAspectRatio"
)

type resize struct {
	width           int
	height          int
	keepAspectRatio bool
}

func NewResize() filters.Spec {
	return &resize{}
}

func (r *resize) Name() string {
	return ResizeName
}

func (r *resize) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for resize ", r)

	if !r.keepAspectRatio {
		return &bimg.Options{
			Width:  r.width,
			Height: r.height,
			Force:  true}, nil
	}

	size, err := image.Size()
	if err != nil {
		return nil, err
	}

	// calculate height keeping width
	ht := int(math.Floor(float64(size.Height*r.width) / float64(size.Width)))

	// if height is less or equal than desired, return transform by width
	if ht <= r.height {
		return &bimg.Options{
			Width: r.width}, nil
	} else {
		// otherwise transform by height
		return &bimg.Options{
			Height: r.height}, nil
	}

}

func (r *resize) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 2 && len(args) != 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resize{}

	f.width, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	f.height, err = parse.EskipIntArg(args[1])

	if err != nil {
		return nil, err
	}

	if len(args) == 3 {
		ratio, err := parse.EskipStringArg(args[2])
		if err != nil {
			return nil, err
		}

		f.keepAspectRatio = !(ratio == ignoreAspectRatioStr)

	} else {
		f.keepAspectRatio = true
	}

	return f, nil
}

func (r *resize) Request(ctx filters.FilterContext) {}

func (r *resize) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
