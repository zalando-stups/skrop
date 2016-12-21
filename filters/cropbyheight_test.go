package filters

import (
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
	"testing"
)

func TestNewCropByHeight(t *testing.T) {
	if NewCropByHeight().Name() != "cropByHeight" {
		t.Error("New crop by height name incorrect")
	}
}

func TestCropByHeight_Name(t *testing.T) {
	c := cropByHeight{}
	if c.Name() != "cropByHeight" {
		t.Error("Crop by height name incorrect")
	}
}

func TestCropByHeight_CreateOptions(t *testing.T) {
	c := cropByHeight{height: 400, cropType: North}
	image := imagefiltertest.LandscapeImage()
	options, _ := c.CreateOptions(image)
	if (*options != bimg.Options{Width: 1000, Height: 400, Crop: true, Gravity: bimg.GravityNorth}) {
		t.Error("Create options for crop didn't return a correct value, ", *options)
	}
}

func TestCropByHeight_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewCropByHeight, []imagefiltertest.CreateTestItem{{
		"no args",
		nil,
		true,
	}, {
		"one arg",
		[]interface{}{400.0},
		false,
	}, {
		"two args",
		[]interface{}{400.0, North},
		false,
	}, {
		"more than 2 args",
		[]interface{}{400.0, 200.0, North},
		true,
	}})
}
