package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
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
	options, _ := c.CreateOptions(image)

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 668, options.Height)
	assert.Equal(t, true, options.Crop)
	assert.Equal(t, bimg.GravityNorth, options.Gravity)
}

func TestCropByWidth_CreateFilter(t *testing.T) {
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
