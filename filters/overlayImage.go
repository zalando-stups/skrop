package filters

import (
	"errors"
	"github.com/zalando-incubator/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"os"
)

const (
	OverlayImageName = "overlayImage"
	NE               = "NE"
	NC               = "NC"
	NW               = "NW"
	CE               = "CE"
	CC               = "CC"
	CW               = "CW"
	SE               = "SE"
	SC               = "SC"
	SW               = "SW"
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

func NewOverlayImage() filters.Spec {
	return &overlay{}
}

func (r *overlay) Name() string {
	return OverlayImageName
}

func (r *overlay) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	origSize, err := image.Size()
	if err != nil {
		return nil, err
	}

	overArr := readImage(r.file)
	overImage := bimg.NewImage(overArr)
	overSize, err := overImage.Size()
	if err != nil {
		return nil, err
	}

	var x, y int
	switch r.verticalGravity {
	case bimg.GravityNorth:
		y = r.topMargin
	case bimg.GravityCentre:
		y = r.topMargin + int(float64(origSize.Height-r.topMargin-r.bottomMargin)/2) - int(float64(overSize.Height)/2)
	case bimg.GravitySouth:
		y = origSize.Height - r.bottomMargin - overSize.Height
	}

	// in case y overflows the image
	if y < 0 || y+overSize.Height > origSize.Height {
		return nil, errors.New("Error: the overlay image in placed outside the image area on the y axe")
	}

	switch r.horizontalGravity {
	case bimg.GravityWest:
		x = r.leftMargin
	case bimg.GravityCentre:
		x = r.leftMargin + int(float64(origSize.Width-r.leftMargin-r.rightMargin)/2) - int(float64(overSize.Width)/2)
	case bimg.GravityEast:
		x = origSize.Width - r.rightMargin - overSize.Width
	}

	// in case y overflows the image
	if x < 0 || x+overSize.Width > origSize.Width {
		return nil, errors.New("Error: the overlay image in placed outside the image area on the x axe")
	}

	return &bimg.Options{WatermarkImage: bimg.WatermarkImage{Buf: overArr,
		Opacity: float32(r.opacity),
		Left:    x,
		Top:     y,
	}}, nil
}

func readImage(file string) []byte {
	img, _ := os.Open(file)
	buf, _ := ioutil.ReadAll(img)
	defer img.Close()
	return buf
}

func (r *overlay) CreateFilter(args []interface{}) (filters.Filter, error) {
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

func (r *overlay) Request(ctx filters.FilterContext) {}

func (r *overlay) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
