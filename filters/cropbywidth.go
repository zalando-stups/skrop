package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
)

// CropByWidthName is the name of the filter
const CropByWidthName = "cropByWidth"

type cropByWidth struct {
	width    int
	cropType string
}

// NewCropByWidth creates a new filter of this type
func NewCropByWidth() filters.Spec {
	return &cropByWidth{}
}

func (f *cropByWidth) Name() string {
	return CropByWidthName
}

func (f *cropByWidth) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for crop by width ", f)

	imageSize, err := imageContext.Image.Size()

	if err != nil {
		return nil, err
	}

	return &bimg.Options{
		Width:   f.width,
		Height:  imageSize.Height,
		Gravity: cropTypeToGravity[f.cropType],
		Crop:    true}, nil
}

func (f *cropByWidth) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.Width == 0 && other.Height == 0 && !other.Crop) ||
		(other.Width == self.Width && other.Height == self.Height && other.Crop == self.Crop)
}

func (f *cropByWidth) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	other.Gravity = self.Gravity
	other.Crop = self.Crop
	return other
}

func (f *cropByWidth) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 1 || len(args) > 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &cropByWidth{cropType: Center}

	c.width, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	if len(args) == 2 {
		if cropType, ok := args[1].(string); ok && cropTypes[cropType] {
			c.cropType = cropType
		} else {
			return nil, filters.ErrInvalidFilterParameters
		}
	}

	return c, nil
}

func (f *cropByWidth) Request(ctx filters.FilterContext) {}

func (f *cropByWidth) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
