package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/tools"
)

const (
	QualityName = "quality"
)

type quality struct {
	percentage int
}

func NewQuality() filters.Spec {
	return &quality{}
}

func (r *quality) Name() string {
	return QualityName
}

func (r *quality) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for quality ", r)

	return &bimg.Options{
		Quality: r.percentage}, nil
}

func (r *quality) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &quality{}

	f.percentage, err = tools.ParseEskipIntArg(args[0])

	if err != nil {
		return nil, err
	} else if f.percentage > 100 {
		return nil, filters.ErrInvalidFilterParameters
	}

	return f, nil
}

func (r *quality) Request(ctx filters.FilterContext) {}

func (r *quality) Response(ctx filters.FilterContext) {
	HandleResponse(ctx, r)
}
