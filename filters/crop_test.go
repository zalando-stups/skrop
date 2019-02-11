package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
	"testing"
)

func TestNewCrop(t *testing.T) {
	if NewCrop().Name() != "crop" {
		t.Error("New crop name incorrect")
	}
}

func TestCrop_Name(t *testing.T) {
	c := crop{}
	if c.Name() != "crop" {
		t.Error("Crop name incorrect")
	}
}

func TestCrop_CreateOptions(t *testing.T) {
	c := crop{width: 800, height: 600, cropType: North}
	options, _ := c.CreateOptions(nil)

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 600, options.Height)
	assert.Equal(t, true, options.Crop)
	assert.Equal(t, bimg.GravityNorth, options.Gravity)
}

func TestCrop_CanBeMerged_True(t *testing.T) {
	s := crop{}
	opt := &bimg.Options{}
	self := &bimg.Options{Width: 100, Height: 350, Gravity: 2, Crop: true}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestCrop_CanBeMerged_False(t *testing.T) {
	s := crop{}
	opt := &bimg.Options{Width: 100, Height: 350, Gravity: 2, Crop: true}
	self := &bimg.Options{Width: 225, Height: 365, Gravity: 2, Crop: true}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestCrop_Merge(t *testing.T) {
	s := crop{}
	self := &bimg.Options{Width: 100, Height: 350, Gravity: 2, Crop: true}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Width, opt.Width)
	assert.Equal(t, self.Height, opt.Height)
	assert.Equal(t, self.Gravity, opt.Gravity)
	assert.Equal(t, self.Crop, opt.Crop)
}

func TestCrop_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewCrop, []imagefiltertest.CreateTestItem{{
		"no args",
		nil,
		true,
	}, {
		"two args",
		[]interface{}{800.0, 600.0},
		false,
	}, {
		"three args",
		[]interface{}{800.0, 600.0, North},
		false,
	}, {
		"more than 3 args",
		[]interface{}{800.0, 600.0, North, "whaaat?"},
		true,
	}, {
		"less than 2 args",
		[]interface{}{800.0},
		true,
	}})
}
