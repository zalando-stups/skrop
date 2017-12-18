package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

// For infomations about the parameters meanings and default values have a look here:
// http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/libvips-convolution.html#vips-sharpen

// SharpenName is the name of the filter
const SharpenName = "sharpen"

type sharpen struct {
	Radius int
	X1     float64
	Y2     float64
	Y3     float64
	M1     float64
	M2     float64
}

// NewSharpen creates a new filter of this type
func NewSharpen() filters.Spec {
	return &sharpen{}
}

func (f *sharpen) Name() string {
	return SharpenName
}

func (f *sharpen) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for sharpen ", r)

	sha := bimg.Sharpen{Radius: r.Radius, X1: r.X1, Y2: r.Y2, Y3: r.Y3, M1: r.M1, M2: r.M2}

	return &bimg.Options{
		Sharpen: sha}, nil
}

func (f *sharpen) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	zero := bimg.Sharpen{}

	//it can be merged if the background was not set (in options or in self) or if they are set to the same value
	return other.Sharpen == zero || other.Sharpen == self.Sharpen
}

func (f *sharpen) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Sharpen = self.Sharpen
	return other
}

func (f *sharpen) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 6 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &sharpen{}

	f.Radius, err = parse.EskipIntArg(args[0])
	if err != nil {
		return nil, err
	}

	f.X1, err = parse.EskipFloatArg(args[1])
	if err != nil {
		return nil, err
	}

	f.Y2, err = parse.EskipFloatArg(args[2])
	if err != nil {
		return nil, err
	}

	f.Y3, err = parse.EskipFloatArg(args[3])
	if err != nil {
		return nil, err
	}

	f.M1, err = parse.EskipFloatArg(args[4])
	if err != nil {
		return nil, err
	}

	f.M2, err = parse.EskipFloatArg(args[5])
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *sharpen) Request(ctx filters.FilterContext) {}

func (f *sharpen) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
