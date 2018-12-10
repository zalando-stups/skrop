package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
)

// CropName is the name of the filter
const CropName = "crop"

type crop struct {
	width    int
	height   int
	cropType string
}

// NewCrop creates a new filter of this type
func NewCrop() filters.Spec {
	return &crop{}
}

func (f *crop) Name() string {
	return CropName
}

func (f *crop) CreateOptions(_ *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for crop ", f)

	return &bimg.Options{
		Width:   f.width,
		Height:  f.height,
		Gravity: cropTypeToGravity[f.cropType],
		Crop:    true}, nil
}

func (f *crop) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.Width == 0 && other.Height == 0 && !other.Crop) ||
		(other.Width == self.Width && other.Height == self.Height && other.Crop == self.Crop)
}

func (f *crop) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	other.Gravity = self.Gravity
	other.Crop = self.Crop
	return other
}

func (f *crop) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 2 || len(args) > 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &crop{cropType: Center}

	c.width, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	c.height, err = parse.EskipIntArg(args[1])

	if err != nil {
		return nil, err
	}

	if len(args) == 3 {
		if cropType, ok := args[2].(string); ok && cropTypes[cropType] {
			c.cropType = cropType
		} else {
			return nil, filters.ErrInvalidFilterParameters
		}
	}

	return c, nil
}

func (f *crop) Request(ctx filters.FilterContext) {}

func (f *crop) Response(ctx filters.FilterContext) {
	log.Debugf("Response %s\n", CropName)
	HandleImageResponse(ctx, f)
}
