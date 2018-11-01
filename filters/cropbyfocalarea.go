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

	focalPointXString := imageContext.PathParam("focalPointX")
	focalPointYString := imageContext.PathParam("focalPointY")
	desiredWidthString := imageContext.PathParam("desiredWidth")
	desiredHeightString := imageContext.PathParam("desiredHeight")

	if focalPointXString == "" || focalPointYString == "" || desiredWidthString == "" || desiredHeightString == "" {
		return nil, filters.ErrInvalidFilterParameters
	}

	focalPointX, err := strconv.Atoi(focalPointXString)

	if err != nil {
		return nil, err
	}

	focalPointY, err  := strconv.Atoi(focalPointYString)

	if err != nil {
		return nil, err
	}

	desiredWidth, err := strconv.Atoi(desiredWidthString)

	if err != nil {
		return nil, err
	}

	desiredHeight, err := strconv.Atoi(desiredHeightString)

	if err != nil {
		return nil, err
	}

	if desiredWidth > imageSize.Width {
		return nil, filters.ErrInvalidFilterParameters
	}

	if desiredHeight > imageSize.Height {
		return nil, filters.ErrInvalidFilterParameters
	}

	if float64(focalPointX) + 0.5*float64(desiredWidth) > float64(imageSize.Width) {
		focalPointX = imageSize.Width - int(0.5*float64(desiredWidth))
	} else if float64(focalPointX) - 0.5*float64(desiredWidth) < 0 {
		focalPointX = desiredWidth/2
	}

	if float64(focalPointY) + 0.5*float64(desiredHeight) > float64(imageSize.Height) {
		focalPointY = imageSize.Height - int(0.5*float64(desiredHeight))
	} else if float64(focalPointY) - 0.5*float64(desiredHeight) < 0 {
		focalPointY = desiredHeight/2
	}

	return &bimg.Options{
		AreaWidth:  desiredWidth,
		AreaHeight: desiredHeight,
		Top:    focalPointY - desiredHeight/2,
		Left:    focalPointX - desiredWidth/2,
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
