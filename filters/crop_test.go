package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
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
