package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const CropName = "crop"

type crop struct {
	width    int
	height   int
	cropType string
}

func NewCrop() filters.Spec {
	return &crop{}
}

func (c *crop) Name() string {
	return CropName
}

func (c *crop) CreateOptions(_ *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for crop ", c)

	return &bimg.Options{
		Width:   c.width,
		Height:  c.height,
		Gravity: cropTypeToGravity[c.cropType],
		Crop:    true}, nil
}

func (c *crop) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 2 || len(args) > 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &crop{cropType: Center}

	f.width, err = parseEskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	f.height, err = parseEskipIntArg(args[1])

	if err != nil {
		return nil, err
	}

	if len(args) == 3 {
		if cropType, ok := args[2].(string); ok && cropTypes[cropType] {
			f.cropType = cropType
		} else {
			return nil, filters.ErrInvalidFilterParameters
		}
	}

	return f, nil
}

func (c *crop) Request(ctx filters.FilterContext) {}

func (c *crop) Response(ctx filters.FilterContext) {
	handleResponse(ctx, c)
}
