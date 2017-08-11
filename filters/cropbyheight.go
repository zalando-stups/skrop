package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-incubator/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
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

func (s *cropByHeight) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.Width == 0 && other.Height == 0 && !other.Crop) ||
		(other.Width == self.Width && other.Height == self.Height && other.Crop == self.Crop)
}

func (s *cropByHeight) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	other.Gravity = self.Gravity
	other.Crop = self.Crop
	return other
}

func (c *cropByHeight) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 1 || len(args) > 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &cropByHeight{cropType: Center}

	f.height, err = parse.EskipIntArg(args[0])

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
