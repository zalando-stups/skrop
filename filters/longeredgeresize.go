package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/parse"
)

const (
	LongerEdgeResizeName = "longerEdgeResize"
)

type longerEdgeResize struct {
	size int
}

func NewLongerEdgeResize() filters.Spec {
	return &longerEdgeResize{}
}

func (r *longerEdgeResize) Name() string {
	return LongerEdgeResizeName
}

func (r *longerEdgeResize) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for longer edge resize ", r)

	imageSize, err := image.Size()

	if err != nil {
		return nil, err
	}

	if imageSize.Width > imageSize.Height {
		return &bimg.Options{
			Width: r.size}, nil
	} else {
		return &bimg.Options{
			Height: r.size}, nil
	}
}

func (r *longerEdgeResize) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &longerEdgeResize{}

	f.size, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *longerEdgeResize) Request(ctx filters.FilterContext) {}

func (r *longerEdgeResize) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
