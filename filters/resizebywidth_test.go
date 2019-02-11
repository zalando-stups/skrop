package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
	"testing"
)

func TestNewResizeByWidth(t *testing.T) {
	name := NewResizeByWidth().Name()
	assert.Equal(t, "width", name)
}

func TestResizeByWidth_Name(t *testing.T) {
	c := resizeByWidth{}
	assert.Equal(t, "width", c.Name())
}

func TestResizeByWidth_CreateOptions_Shrink_EnlargeAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Width / 2
	resizeByWidth := resizeByWidth{width: newSize, enlarge: true}
	options, _ := resizeByWidth.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newSize, options.Width)
}

func TestResizeByWidth_CreateOptions_Enlarge_EnlargeAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Width * 2
	resizeByWidth := resizeByWidth{width: newSize, enlarge: true}
	options, _ := resizeByWidth.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newSize, options.Width)
}

func TestResizeByWidth_CreateOptions_Shrink_EnlargeNotAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Width / 2
	resizeByWidth := resizeByWidth{width: newSize, enlarge: false}
	options, _ := resizeByWidth.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newSize, options.Width)
}

func TestResizeByWidth_CreateOptions_Enlarge_EnlargeNotAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Width * 2
	resizeByWidth := resizeByWidth{width: newSize, enlarge: false}
	options, _ := resizeByWidth.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, 0, options.Width)
}

func TestResizeByWidth_CanBeMerged_True(t *testing.T) {
	s := resizeByWidth{}
	opt := &bimg.Options{}
	self := &bimg.Options{Width: 265}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestResizeByWidth_CanBeMerged_False(t *testing.T) {
	s := resizeByWidth{}
	opt := &bimg.Options{Width: 265}
	self := &bimg.Options{Width: 144}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestResizeByWidth_Merge(t *testing.T) {
	s := resizeByWidth{}
	self := &bimg.Options{Width: 265}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Width, opt.Width)
}

func TestResizeByWidth_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewResizeByWidth, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "one arg",
		Args: []interface{}{256.0},
		Err:  false,
	}, {
		Msg:  "two args",
		Args: []interface{}{100.0, "DO_NOT_ENLARGE"},
		Err:  false,
	}, {
		Msg:  "wrong 2nd args",
		Args: []interface{}{1000.0, 100.0},
		Err:  true,
	}, {
		Msg:  "more than two args",
		Args: []interface{}{2000.0, "DO_NOT_ENLARGE", 1.0},
		Err:  true,
	}})
}
