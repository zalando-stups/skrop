package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/tools"
)

const CropByHeightName = "cropByHeight"

type cropByHeight struct {
	height   int
	cropType string
}

func NewCropByHeight() filters.Spec {
	return &cropByHeight{}
}

func (c *cropByHeight) Name() string {
	return CropByHeightName
}

func (c *cropByHeight) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for crop by height ", c)

	imageSize, err := image.Size()

	if err != nil {
		return nil, err
	}

	return &bimg.Options{
		Width:   imageSize.Width,
		Height:  c.height,
		Gravity: cropTypeToGravity[c.cropType],
		Crop:    true}, nil
}

func (c *cropByHeight) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 1 || len(args) > 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &cropByHeight{cropType: Center}

	f.height, err = tools.ParseEskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	if len(args) == 2 {
		if cropType, ok := args[1].(string); ok && cropTypes[cropType] {
			f.cropType = cropType
		} else {
			return nil, filters.ErrInvalidFilterParameters
		}
	}

	return f, nil
}

func (c *cropByHeight) Request(ctx filters.FilterContext) {}

func (c *cropByHeight) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, c)
}
