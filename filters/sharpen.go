package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

// For infomations about the parameters meanings and default values have a look here:
// http://www.vips.ecs.soton.ac.uk/supported/current/doc/html/libvips/libvips-convolution.html#vips-sharpen

const (
	SharpenName = "sharpen"
)

type sharpen bimg.Sharpen

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

	f.Radius, err = parseEskipIntArg(args[0])
	f.X1, err = parseEskipFloatArg(args[1])
	f.Y2, err = parseEskipFloatArg(args[2])
	f.Y3, err = parseEskipFloatArg(args[3])
	f.M1, err = parseEskipFloatArg(args[4])
	f.M2, err = parseEskipFloatArg(args[5])

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *sharpen) Request(ctx filters.FilterContext) {}

func (r *sharpen) Response(ctx filters.FilterContext) {
	handleResponse(ctx, r)
}
