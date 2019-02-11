package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
	"testing"
)

func TestNewOverlayImage(t *testing.T) {
	name := NewOverlayImage().Name()
	assert.Equal(t, "overlayImage", name)
}

func TestOverlay_Name(t *testing.T) {
	c := &overlay{}
	assert.Equal(t, "overlayImage", c.Name())
}

func TestOverlay_CreateOptions_SE(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	overArr, _ := readImage("../images/star.png")
	overImage := bimg.NewImage(overArr)
	overSize, _ := overImage.Size()
	overlay := &overlay{file: "../images/star.png",
		opacity:           0.9,
		horizontalGravity: bimg.GravityEast,
		verticalGravity:   bimg.GravitySouth,
		leftMargin:        10,
		rightMargin:       20,
		topMargin:         30,
		bottomMargin:      40,
	}

	options, _ := overlay.CreateOptions(buildParameters(nil, image))
	over := options.WatermarkImage

	assert.Equal(t, overArr, over.Buf)
	assert.Equal(t, float32(0.9), over.Opacity)
	assert.Equal(t, size.Height-overSize.Height-40, over.Top)
	assert.Equal(t, size.Width-overSize.Width-20, over.Left)
}

func TestOverlay_CreateOptions_NW(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	overArr, _ := readImage("../images/star.png")
	overlay := &overlay{file: "../images/star.png",
		opacity:           0.9,
		horizontalGravity: bimg.GravityWest,
		verticalGravity:   bimg.GravityNorth,
		leftMargin:        10,
		rightMargin:       20,
		topMargin:         30,
		bottomMargin:      40,
	}

	options, _ := overlay.CreateOptions(buildParameters(nil, image))
	over := options.WatermarkImage

	assert.Equal(t, overArr, over.Buf)
	assert.Equal(t, float32(0.9), over.Opacity)
	assert.Equal(t, 30, over.Top)
	assert.Equal(t, 10, over.Left)
}

func TestOverlay_CreateOptions_CC(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	overArr, _ := readImage("../images/star.png")
	overImage := bimg.NewImage(overArr)
	overSize, _ := overImage.Size()
	overlay := &overlay{file: "../images/star.png",
		opacity:           0.9,
		horizontalGravity: bimg.GravityCentre,
		verticalGravity:   bimg.GravityCentre,
		leftMargin:        0,
		rightMargin:       0,
		topMargin:         0,
		bottomMargin:      0,
	}

	options, _ := overlay.CreateOptions(buildParameters(nil, image))
	over := options.WatermarkImage

	assert.Equal(t, overArr, over.Buf)
	assert.Equal(t, float32(0.9), over.Opacity)
	assert.Equal(t, int(size.Height/2)-int(overSize.Height/2), over.Top)
	assert.Equal(t, int(size.Width/2)-int(overSize.Width/2), over.Left)
}

func TestOverlay_CanBeMerged_True(t *testing.T) {
	s := overlay{}
	opt := &bimg.Options{}
	self := &bimg.Options{WatermarkImage: bimg.WatermarkImage{Opacity: 3.4, Left: 10, Top: 20}}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestOverlay_CanBeMerged_False(t *testing.T) {
	s := overlay{}
	opt := &bimg.Options{WatermarkImage: bimg.WatermarkImage{Opacity: 3.7, Left: 20, Top: 30}}
	self := &bimg.Options{WatermarkImage: bimg.WatermarkImage{Opacity: 3.4, Left: 10, Top: 20}}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestOverlay_Merge(t *testing.T) {
	s := overlay{}
	self := &bimg.Options{WatermarkImage: bimg.WatermarkImage{Opacity: 3.4, Left: 10, Top: 20}}

	opt := s.Merge(&bimg.Options{}, self)

	assert.True(t, equals(opt.WatermarkImage, self.WatermarkImage))
}

func TestOverlay_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewOverlayImage, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "two args",
		Args: []interface{}{25.0, 35.0},
		Err:  true,
	}, {
		Msg:  "three args",
		Args: []interface{}{"abc", 2.6, "NE"},
		Err:  false,
	}, {
		Msg:  "wrong type args",
		Args: []interface{}{"abc", -2.6, true},
		Err:  true,
	}, {
		Msg:  "wrong type args",
		Args: []interface{}{1.0, 2.6, "NE"},
		Err:  true,
	}, {
		Msg:  "wrong type args",
		Args: []interface{}{"abc", "", "NE"},
		Err:  true,
	}, {
		Msg:  "five args error",
		Args: []interface{}{"abc", 2.6, "NE", 1.0, 2.0},
		Err:  true,
	}, {
		Msg:  "seven args",
		Args: []interface{}{"abc", 2.6, "NE", 1.0, 2.0, 3.0, 4.0},
		Err:  false,
	}, {
		Msg:  "wrong type args",
		Args: []interface{}{"abc", 2.6, "NE", false, 2.0, 3.0, 4.0},
		Err:  true,
	}, {
		Msg:  "wrong type args",
		Args: []interface{}{"abc", 2.6, "NE", 1.0, "2.0", 3.0, 4.0},
		Err:  true,
	}, {
		Msg:  "wrong type args",
		Args: []interface{}{"abc", 2.6, "NE", 1.0, 2.0, false, 4.0},
		Err:  true,
	}, {
		Msg:  "wrong type args",
		Args: []interface{}{"abc", 2.6, "NE", 1.0, 2.0, 3.0, ""},
		Err:  true,
	}, {
		Msg:  "gravity error",
		Args: []interface{}{"abc", 2.6, "NA", 1.0, 2.0, 3.0, 4.0},
		Err:  true,
	}})
}
