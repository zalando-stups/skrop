package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
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

func (r *crop) Name() string {
	return CropName
}

func (r *crop) CreateOptions(_ *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for crop ", c)

	return &bimg.Options{
		Width:   c.width,
		Height:  c.height,
		Gravity: cropTypeToGravity[c.cropType],
		Crop:    true}, nil
}

func (r *crop) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.Width == 0 && other.Height == 0 && !other.Crop) ||
		(other.Width == self.Width && other.Height == self.Height && other.Crop == self.Crop)
}

func (r *crop) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	other.Gravity = self.Gravity
	other.Crop = self.Crop
	return other
}

func (r *crop) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 2 || len(args) > 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &crop{cropType: Center}

	f.width, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	f.height, err = parse.EskipIntArg(args[1])

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

func (r *crop) Request(ctx filters.FilterContext) {}

func (r *crop) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, c)
}
