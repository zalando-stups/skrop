package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const CropByWidthName = "cropByWidth"

type cropByWidth struct {
	width    int
	cropType string
}

func NewCropByWidth() filters.Spec {
	return &cropByWidth{}
}

func (c *cropByWidth) Name() string {
	return CropByWidthName
}

func (c *cropByWidth) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for crop by width ", c)

	imageSize, err := image.Size()

	if err != nil {
		return nil, err
	}

	return &bimg.Options{
		Width:   c.width,
		Height:  imageSize.Height,
		Gravity: cropTypeToGravity[c.cropType],
		Crop:    true}, nil
}

func (s *cropByWidth) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.Width == 0 && other.Height == 0 && !other.Crop) ||
		(other.Width == self.Width && other.Height == self.Height && other.Crop == self.Crop)
}

func (s *cropByWidth) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	other.Gravity = self.Gravity
	other.Crop = self.Crop
	return other
}

func (c *cropByWidth) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 1 || len(args) > 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &cropByWidth{cropType: Center}

	f.width, err = parse.EskipIntArg(args[0])

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

func (c *cropByWidth) Request(ctx filters.FilterContext) {}

func (c *cropByWidth) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, c)
}
