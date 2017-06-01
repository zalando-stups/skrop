package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/tools"
)

const (
	ResizeByHeightName = "height"
)

type resizeByHeight struct {
	height int
}

func NewResizeByHeight() filters.Spec {
	return &resizeByHeight{}
}

func (r *resizeByHeight) Name() string {
	return ResizeByHeightName
}

func (r *resizeByHeight) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for resize by width ", r)

	return &bimg.Options{
		Height: r.height}, nil
}

func (r *resizeByHeight) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resizeByHeight{}

	f.height, err = tools.ParseEskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *resizeByHeight) Request(ctx filters.FilterContext) {}

func (r *resizeByHeight) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
