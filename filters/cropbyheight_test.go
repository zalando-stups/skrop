package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
)

func TestNewCropByHeight(t *testing.T) {
	name := NewCropByHeight().Name()
	assert.Equal(t, "cropByHeight", name)
}

func TestCropByHeight_Name(t *testing.T) {
	c := cropByHeight{}
	assert.Equal(t, "cropByHeight", c.Name())
}

func TestCropByHeight_CreateOptions(t *testing.T) {
	c := cropByHeight{height: 400, cropType: North}
	image := imagefiltertest.LandscapeImage()
	options, _ := c.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, 1000, options.Width)
	assert.Equal(t, 400, options.Height)
	assert.Equal(t, true, options.Crop)
	assert.Equal(t, bimg.GravityNorth, options.Gravity)
}

func TestCropByHeight_CanBeMerged_True(t *testing.T) {
	s := cropByHeight{}
	opt := &bimg.Options{}
	self := &bimg.Options{Height: 350, Gravity: 2, Crop: true}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestCropByHeight_CanBeMerged_False(t *testing.T) {
	s := cropByHeight{}
	opt := &bimg.Options{Height: 350, Gravity: 2, Crop: true}
	self := &bimg.Options{Height: 365, Gravity: 2, Crop: true}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestCropByHeight_Merge(t *testing.T) {
	s := cropByHeight{}
	self := &bimg.Options{Height: 350, Gravity: 2, Crop: true}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Height, opt.Height)
	assert.Equal(t, self.Gravity, opt.Gravity)
	assert.Equal(t, self.Crop, opt.Crop)
}

func TestCropByHeight_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewCropByHeight, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "one arg",
		Args: []interface{}{400.0},
		Err:  false,
	}, {
		Msg:  "two args",
		Args: []interface{}{400.0, North},
		Err:  false,
	}, {
		Msg:  "more than 2 args",
		Args: []interface{}{400.0, 200.0, North},
		Err:  true,
	}})
}
