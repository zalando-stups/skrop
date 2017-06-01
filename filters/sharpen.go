package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/tools"
)

// For infomations about the parameters meanings and default values have a look here:
// http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/libvips-convolution.html#vips-sharpen

const (
	SharpenName = "sharpen"
)

type sharpen struct {
	Radius int
	X1     float64
	Y2     float64
	Y3     float64
	M1     float64
	M2     float64
}

func NewSharpen() filters.Spec {
	return &sharpen{}
}

func (r *sharpen) Name() string {
	return SharpenName
}

func (r *sharpen) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for sharpen ", r)

	sha := bimg.Sharpen{Radius: r.Radius, X1: r.X1, Y2: r.Y2, Y3: r.Y3, M1: r.M1, M2: r.M2}

	return &bimg.Options{
		Sharpen: sha}, nil
}

func (r *sharpen) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 6 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &sharpen{}

	f.Radius, err = tools.ParseEskipIntArg(args[0])
	if err != nil {
		return nil, err
	}

	f.X1, err = tools.ParseEskipFloatArg(args[1])
	if err != nil {
		return nil, err
	}

	f.Y2, err = tools.ParseEskipFloatArg(args[2])
	if err != nil {
		return nil, err
	}

	f.Y3, err = tools.ParseEskipFloatArg(args[3])
	if err != nil {
		return nil, err
	}

	f.M1, err = tools.ParseEskipFloatArg(args[4])
	if err != nil {
		return nil, err
	}

	f.M2, err = tools.ParseEskipFloatArg(args[5])
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *sharpen) Request(ctx filters.FilterContext) {}

func (r *sharpen) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
