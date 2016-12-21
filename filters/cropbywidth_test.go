package filters

import (
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
	"testing"
)

func TestNewCropByWidth(t *testing.T) {
	if NewCropByWidth().Name() != "cropByWidth" {
		t.Error("New crop by width name incorrect")
	}
}

func TestCropByWidth_Name(t *testing.T) {
	c := cropByWidth{}
	if c.Name() != "cropByWidth" {
		t.Error("Crop by width name incorrect")
	}
}

func TestCropByWidth_CreateOptions(t *testing.T) {
	c := cropByWidth{width: 800, cropType: North}
	image := imagefiltertest.LandscapeImage()
	options, _ := c.CreateOptions(image)
	if (*options != bimg.Options{Width: 800, Height: 668, Crop: true, Gravity: bimg.GravityNorth}) {
		t.Error("Create options for crop didn't return a correct value, ", *options)
	}
}

func TestCropByWidth_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewCropByWidth, []imagefiltertest.CreateTestItem{{
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
