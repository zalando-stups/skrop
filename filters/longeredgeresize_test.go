package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
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
	options, _ := resize.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 0, options.Height)
}

func TestLongerEdgeResize_CreateOptions_Portrait(t *testing.T) {
	resize := longerEdgeResize{size: 800}
	image := imagefiltertest.PortraitImage()
	options, _ := resize.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, 0, options.Width)
	assert.Equal(t, 800, options.Height)
}

func TestLongerEdgeResize_CanBeMerged_True(t *testing.T) {
	s := longerEdgeResize{}
	opt := &bimg.Options{}
	self := &bimg.Options{Width: 100}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestLongerEdgeResize_CanBeMerged_False(t *testing.T) {
	s := longerEdgeResize{}
	opt := &bimg.Options{Height: 350}
	self := &bimg.Options{Height: 365}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestLongerEdgeResize_Merge(t *testing.T) {
	s := longerEdgeResize{}
	self := &bimg.Options{Width: 100}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Width, opt.Width)
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
