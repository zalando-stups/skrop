package filters

import (
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
	"io/ioutil"
	"os"
)

const (
	// OverlayImageName is the name of the filter
	OverlayImageName = "overlayImage"
	// NE North East
	NE = "NE"
	// NC North Center
	NC = "NC"
	// NW North West
	NW = "NW"
	// CE Centre East
	CE = "CE"
	// CC Centre Center
	CC = "CC"
	// CW Centre West
	CW = "CW"
	// SE South East
	SE = "SE"
	// SC South Center
	SC = "SC"
	// SW South West
	SW = "SW"
)

var (
	gravityType = map[string]bool{
		NE: true,
		NC: true,
		NW: true,
		CE: true,
		CC: true,
		CW: true,
		SE: true,
		SC: true,
		SW: true,
	}
	verticalGravity = map[string]bimg.Gravity{
		NE: bimg.GravityNorth,
		NC: bimg.GravityNorth,
		NW: bimg.GravityNorth,
		CE: bimg.GravityCentre,
		CC: bimg.GravityCentre,
		CW: bimg.GravityCentre,
		SE: bimg.GravitySouth,
		SC: bimg.GravitySouth,
		SW: bimg.GravitySouth,
	}
	horizontalGravity = map[string]bimg.Gravity{
		NE: bimg.GravityEast,
		NC: bimg.GravityCentre,
		NW: bimg.GravityWest,
		CE: bimg.GravityEast,
		CC: bimg.GravityCentre,
		CW: bimg.GravityWest,
		SE: bimg.GravityEast,
		SC: bimg.GravityCentre,
		SW: bimg.GravityWest,
	}
)

type overlay struct {
	file              string
	opacity           float64
	verticalGravity   bimg.Gravity
	horizontalGravity bimg.Gravity
	rightMargin       int
	leftMargin        int
	topMargin         int
	bottomMargin      int
}

// NewOverlayImage creates a new filter of this type
func NewOverlayImage() filters.Spec {
	return &overlay{}
}

func (f *overlay) Name() string {
	return OverlayImageName
}

func (f *overlay) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	origSize, err := imageContext.Image.Size()
	if err != nil {
		return nil, err
	}

	overArr, err := readImage(f.file)
	if err != nil {
		return nil, err
	}

	overImage := bimg.NewImage(overArr)
	overSize, err := overImage.Size()
	if err != nil {
		return nil, err
	}

	var x, y int
	switch f.verticalGravity {
	case bimg.GravityNorth:
		y = f.topMargin
	case bimg.GravityCentre:
		y = f.topMargin + int(float64(origSize.Height-f.topMargin-f.bottomMargin)/2) - int(
			float64(overSize.Height)/2)
	case bimg.GravitySouth:
		y = origSize.Height - f.bottomMargin - overSize.Height
	}

	switch f.horizontalGravity {
	case bimg.GravityWest:
		x = f.leftMargin
	case bimg.GravityCentre:
		x = f.leftMargin + int(float64(origSize.Width-f.leftMargin-f.rightMargin)/2) - int(float64(overSize.Width)/2)
	case bimg.GravityEast:
		x = origSize.Width - f.rightMargin - overSize.Width
	}

	return &bimg.Options{WatermarkImage: bimg.WatermarkImage{Buf: overArr,
		Opacity: float32(f.opacity),
		Left:    x,
		Top:     y,
	}}, nil
}

func (f *overlay) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	zero := bimg.WatermarkImage{}

	//it can be merged if the background was not set (in options or in self) or if they are set to the same value
	return other.Width == 0 && other.Height == 0 && (equals(other.WatermarkImage, zero) || equals(other.WatermarkImage, self.WatermarkImage))
}

func equals(one bimg.WatermarkImage, two bimg.WatermarkImage) bool {
	return one.Opacity == two.Opacity &&
		one.Top == two.Top &&
		one.Left == two.Left &&
		len(one.Buf) == len(two.Buf)
}

func (f *overlay) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.WatermarkImage = self.WatermarkImage
	return other
}

func readImage(file string) ([]byte, error) {
	img, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	buf, err := ioutil.ReadAll(img)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (f *overlay) CreateFilter(args []interface{}) (filters.Filter, error) {
	//imageOverlay(<filename>, <opacity>, <gravity>, <right_margin>, <left_margin>, <top_margin>, <bottom_margin>)
	//imageOverlay("filename", 1.0, NE, 0, 0, 0, 0)
	//imageOverlay("filename", 1.0, NE)
	var err error

	if len(args) != 3 && len(args) != 7 {
		return nil, filters.ErrInvalidFilterParameters
	}

	o := &overlay{}

	o.file, err = parse.EskipStringArg(args[0])
	if err != nil {
		return nil, err
	}

	o.opacity, err = parse.EskipFloatArg(args[1])
	if err != nil {
		return nil, err
	}
	if o.opacity < 0 {
		o.opacity = 0
	} else if o.opacity > 1.0 {
		o.opacity = 1
	}

	gravity, err := parse.EskipStringArg(args[2])
	if err != nil {
		return nil, err
	}
	if !gravityType[gravity] {
		return nil, filters.ErrInvalidFilterParameters
	}

	o.verticalGravity = verticalGravity[gravity]
	o.horizontalGravity = horizontalGravity[gravity]

	if len(args) == 3 {
		return o, nil
	}

	o.topMargin, err = parse.EskipIntArg(args[3])
	if err != nil {
		return nil, err
	}

	o.rightMargin, err = parse.EskipIntArg(args[4])
	if err != nil {
		return nil, err
	}

	o.bottomMargin, err = parse.EskipIntArg(args[5])
	if err != nil {
		return nil, err
	}

	o.leftMargin, err = parse.EskipIntArg(args[6])
	if err != nil {
		return nil, err
	}

	return o, nil

}

func (f *overlay) Request(ctx filters.FilterContext) {}

func (f *overlay) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
