package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	ResizeName = "resize"
)

type resize struct {
	width  int
	height int
}

func NewResize() filters.Spec {
	return &resize{}
}

func (r *resize) Name() string {
	return ResizeName
}

func (r *resize) CreateOptions(_ *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for resize ", r)

	return &bimg.Options{
		Width:  r.width,
		Height: r.height}, nil
}

func (r *resize) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resize{}

	f.width, err = parseEskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	f.height, err = parseEskipIntArg(args[1])

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *resize) Request(ctx filters.FilterContext) {}

func (r *resize) Response(ctx filters.FilterContext) {
	handleResponse(ctx, r)
}
