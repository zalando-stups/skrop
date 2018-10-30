package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"github.com/danpersa/bimg"
	"strconv"
)

// CropByFocalAreaName is the name of the filter
const CropByFocalAreaName = "cropByFocalArea"

type cropByFocalArea struct {
}

// NewCropByFocalArea creates a new filter of this type
func NewCropByFocalArea() filters.Spec {
	return &cropByFocalArea{}
}

func (f *cropByFocalArea) Name() string {
	return CropByFocalAreaName
}

func (f *cropByFocalArea) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for crop by focal area ", f)

	imageSize, err := imageContext.Image.Size()

	if err != nil {
		return nil, err
	}

	desired_width, err := strconv.Atoi(imageContext.PathParam("desired_width"))

	if err != nil {
		return nil, err
	}

	desired_height, err := strconv.Atoi(imageContext.PathParam("desired_height"))

	if err != nil {
		return nil, err
	}

	focalPointX, err := strconv.Atoi(imageContext.PathParam("focalPointX"))

	if err != nil {
		return nil, err
	}

	focalPointY, err  := strconv.Atoi(imageContext.PathParam("focalPointY"))

	if err != nil {
		return nil, err
	}

	if desired_width > imageSize.Width {
		return nil, filters.ErrInvalidFilterParameters
	}

	if desired_height > imageSize.Height {
		return nil, filters.ErrInvalidFilterParameters
	}

	if float64(focalPointX) + 0.5*float64(desired_width) > float64(imageSize.Width) {
		focalPointX = focalPointX - (imageSize.Width - (focalPointX + int(0.5*float64(desired_width))))
	}

	if float64(focalPointY) + 0.5*float64(desired_height) > float64(imageSize.Height) {
		focalPointY = focalPointY - (imageSize.Height - (focalPointY + int(0.5*float64(desired_height))))
	}

	return &bimg.Options{
		AreaWidth:  desired_width,
		AreaHeight: desired_height,
		Top:    1,
		Left:    0,
		}, nil
}

func (f *cropByFocalArea) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.AreaWidth == 0 || other.AreaWidth == self.AreaWidth) &&
			(other.AreaHeight == 0 || other.AreaHeight == self.AreaHeight) &&
			(other.Top == 0 || other.Top == self.Top) &&
			(other.Left == 0 || other.Left == self.Left) &&
			(other.Width == 0) &&
			(other.Height == 0)
}

func (f *cropByFocalArea) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.AreaWidth = self.AreaWidth
	other.AreaHeight = self.AreaHeight
	other.Top = self.Top
	other.Left = self.Left
	return other
}

func (f *cropByFocalArea) CreateFilter(args []interface{}) (filters.Filter, error) {

	c := &cropByFocalArea{}

	return c, nil
}

func (f *cropByFocalArea) Request(ctx filters.FilterContext) {}

func (f *cropByFocalArea) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
