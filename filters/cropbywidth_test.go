package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
)

func TestNewCropByWidth(t *testing.T) {
	name := NewCropByWidth().Name()
	assert.Equal(t, "cropByWidth", name)
}

func TestCropByWidth_Name(t *testing.T) {
	c := cropByWidth{}
	assert.Equal(t, "cropByWidth", c.Name())
}

func TestCropByWidth_CreateOptions(t *testing.T) {
	c := cropByWidth{width: 800, cropType: North}
	image := imagefiltertest.LandscapeImage()
	options, _ := c.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 668, options.Height)
	assert.Equal(t, true, options.Crop)
	assert.Equal(t, bimg.GravityNorth, options.Gravity)
}

func TestCropByWidth_CanBeMerged_True(t *testing.T) {
	s := cropByWidth{}
	opt := &bimg.Options{}
	self := &bimg.Options{Width: 100, Gravity: 2, Crop: true}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestCropByWidth_CanBeMerged_False(t *testing.T) {
	s := cropByWidth{}
	opt := &bimg.Options{Width: 100, Gravity: 2, Crop: true}
	self := &bimg.Options{Width: 225, Gravity: 2, Crop: true}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestCropByWidth_Merge(t *testing.T) {
	s := cropByWidth{}
	self := &bimg.Options{Width: 100, Gravity: 2, Crop: true}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Width, opt.Width)
	assert.Equal(t, self.Gravity, opt.Gravity)
	assert.Equal(t, self.Crop, opt.Crop)
}

func TestCropByWidth_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewCropByWidth, []imagefiltertest.CreateTestItem{{
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
