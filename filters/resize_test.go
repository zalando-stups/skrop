package filters

import (
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
	"testing"
)

func TestNewResize(t *testing.T) {
	if NewResize().Name() != "resize" {
		t.Error("New resize name incorrect")
	}
}

func TestResize_Name(t *testing.T) {
	c := resize{}
	if c.Name() != "resize" {
		t.Error("Resize name incorrect")
	}
}

func TestResize_CreateOptions(t *testing.T) {
	resize := resize{width: 800, height: 600}
	options, _ := resize.CreateOptions(nil)
	if (*options != bimg.Options{Width: 800, Height: 600}) {
		t.Error("Create options for resize didn't return a correct value")
	}
}

func TestResize_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewResize, []imagefiltertest.CreateTestItem{{
		"no args",
		nil,
		true,
	}, {
		"two args",
		[]interface{}{800.0, 600.0},
		false,
	}, {
		"more than 2 args",
		[]interface{}{800.0, 600.0, "whaaat?"},
		true,
	}, {
		"less than 2 args",
		[]interface{}{800.0},
		true,
	}})
}
