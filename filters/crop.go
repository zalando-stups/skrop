package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	CropName = "crop"
)

type crop struct {
	width  int
	height int
}

func NewCrop() filters.Spec {
	return &crop{}
}

func (c *crop) Name() string {
	return CropName
}

func (c *crop) createOptions() *bimg.Options {
	log.Debug("Create options for crop ", c)

	return &bimg.Options{
		Width:   c.width,
		Height:  c.height,
		Crop:    true,
		Gravity: bimg.GravityCentre}
}

func (c *crop) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &crop{}

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

func (c *crop) Request(ctx filters.FilterContext) {}

func (c *crop) Response(ctx filters.FilterContext) {
	handleResponse(ctx, c)
}
