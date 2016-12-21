package filters

import (
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
	"testing"
)

func TestNewLongerEdgeResize(t *testing.T) {
	if NewLongerEdgeResize().Name() != "longerEdgeResize" {
		t.Error("New longer edge resize name incorrect")
	}
}

func TestLongerEdgeResize_Name(t *testing.T) {
	c := longerEdgeResize{}
	if c.Name() != "longerEdgeResize" {
		t.Error("Longer edge resize name incorrect")
	}
}

func TestLongerEdgeResize_CreateOptions(t *testing.T) {
	resize := longerEdgeResize{size: 800}
	image := imagefiltertest.LandscapeImage()
	options, _ := resize.CreateOptions(image)
	if (*options != bimg.Options{Width: 800}) {
		t.Error("Create options for longer edge resize didn't return a correct value for a landscape image")
	}

	image = imagefiltertest.PortraitImage()
	options, _ = resize.CreateOptions(image)
	if (*options != bimg.Options{Height: 800}) {
		t.Error("Create options for longer edge resize didn't return a correct value for a portrait image")
	}
}

func TestLongerEdgeResize_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewLongerEdgeResize, []imagefiltertest.CreateTestItem{{
		"no args",
		nil,
		true,
	}, {
		"one arg",
		[]interface{}{800.0},
		false,
	}, {
		"more than 1 args",
		[]interface{}{800.0, 600.0},
		true,
	}})
}
