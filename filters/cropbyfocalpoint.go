package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
	"strconv"
)

// CropByFocalPointName is the name of the filter
const CropByFocalPointName = "cropByFocalPoint"

type cropByFocalPoint struct {
	targetX     float64
	targetY     float64
	aspectRatio float64
	minWidth	int
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// NewCropByFocalPoint creates a new filter of this type
func NewCropByFocalPoint() filters.Spec {
	return &cropByFocalPoint{}
}

func (f *cropByFocalPoint) Name() string {
	return CropByFocalPointName
}

func (f *cropByFocalPoint) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for crop by focal point ", f)

	imageSize, err := imageContext.Image.Size()

	if err != nil {
		return nil, err
	}

	focalPointX := imageContext.PathParam("focalPointX")
	focalPointY := imageContext.PathParam("focalPointY")

	if focalPointX == "" || focalPointY == "" {
		return nil, filters.ErrInvalidFilterParameters
	}

	sourceX, err := strconv.Atoi(focalPointX)

	if err != nil {
		return nil, err
	}

	sourceY, err := strconv.Atoi(focalPointY)

	if err != nil {
		return nil, err
	}

	x := sourceX
	y := sourceY
	if f.minWidth != -1 {
		minHeight := int(f.aspectRatio * float64(f.minWidth))

		minX := int(float64(Min(f.minWidth, imageSize.Width)) * f.targetX)
		maxX := imageSize.Width - int(float64(Min(f.minWidth, imageSize.Width)) * (1 - f.targetX))
		minY := int(float64(Min(minHeight, imageSize.Height)) * f.targetY)
		maxY := imageSize.Height - int(float64(Min(minHeight, imageSize.Height)) * (1 - f.targetY))

		if x < minX {
			x = minX
		}
		if x > maxX {
			x = maxX
		}
		if y < minY {
			y = minY
		}
		if y > maxY {
			y = maxY
		}
	}

	right := imageSize.Width - x
	bottom := imageSize.Height - y

	cropLeftWidth := int(float64(x) / f.targetX)
	cropRightWidth := int(float64(right) / (float64(1) - f.targetX))

	width := cropRightWidth

	if cropLeftWidth < cropRightWidth {
		width = cropLeftWidth
	}

	cropTopHeight := int(float64(y) / f.targetY)
	cropBottomHeight := int(float64(bottom) / (float64(1) - f.targetY))

	height := cropBottomHeight

	if cropTopHeight < cropBottomHeight {
		height = int(float64(y) / f.targetY)
	}

	ratio := float64(height) / float64(width)

	if ratio > f.aspectRatio {
		height = int(float64(width) * f.aspectRatio)
	} else {
		width = int(float64(height) / f.aspectRatio)
	}

	return &bimg.Options{
		AreaWidth:  width,
		AreaHeight: height,
		Top:    y - int(float64(height) * f.targetY),
		Left:   x - int(float64(width) * f.targetX)}, nil
}

func (f *cropByFocalPoint) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.AreaWidth == 0 || other.AreaWidth == self.AreaWidth) &&
			(other.AreaHeight == 0 || other.AreaHeight == self.AreaHeight) &&
			(other.Top == 0 || other.Top == self.Top) &&
			(other.Left == 0 || other.Left == self.Left) &&
			(other.Width == 0) &&
			(other.Height == 0)
}

func (f *cropByFocalPoint) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.AreaWidth = self.AreaWidth
	other.AreaHeight = self.AreaHeight
	other.Top = self.Top
	other.Left = self.Left
	return other
}

func (f *cropByFocalPoint) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) < 3 || len(args) > 4 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &cropByFocalPoint{}

	c.targetX, err = parse.EskipFloatArg(args[0])

	if err != nil {
		return nil, err
	}

	c.targetY, err = parse.EskipFloatArg(args[1])

	if err != nil {
		return nil, err
	}

	c.aspectRatio, err = parse.EskipFloatArg(args[2])

	if err != nil {
		return nil, err
	}

	if len(args)  == 4 {
		c.minWidth, err = parse.EskipIntArg(args[3])

		if err != nil {
			return nil, err
		}
	} else {
		c.minWidth = -1
	}

	return c, nil
}

func (f *cropByFocalPoint) Request(ctx filters.FilterContext) {}

func (f *cropByFocalPoint) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
