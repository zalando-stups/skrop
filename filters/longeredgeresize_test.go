package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
)

func TestNewLongerEdgeResize(t *testing.T) {
	name := NewLongerEdgeResize().Name()
	assert.Equal(t, "longerEdgeResize", name)
}

func TestLongerEdgeResize_Name(t *testing.T) {
	c := longerEdgeResize{}
	assert.Equal(t, "longerEdgeResize", c.Name())
}

func TestLongerEdgeResize_CreateOptions_Landscape(t *testing.T) {
	resize := longerEdgeResize{size: 800}
	image := imagefiltertest.LandscapeImage()
	options, _ := resize.CreateOptions(image)

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 0, options.Height)
}

func TestLongerEdgeResize_CreateOptions_Portrait(t *testing.T) {
	resize := longerEdgeResize{size: 800}
	image := imagefiltertest.PortraitImage()
	options, _ := resize.CreateOptions(image)

	assert.Equal(t, 0, options.Width)
	assert.Equal(t, 800, options.Height)
}

func TestLongerEdgeResize_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewLongerEdgeResize, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "one arg",
		Args: []interface{}{800.0},
		Err:  false,
	}, {
		Msg:  "more than 1 args",
		Args: []interface{}{800.0, 600.0},
		Err:  true,
	}})
}
