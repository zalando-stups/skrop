package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
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
	options, _ := resizeByWidth.CreateOptions(image)

	assert.Equal(t, newSize, options.Width)
}

func TestResizeByWidth_CreateOptions_Enlarge_EnlargeAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Width * 2
	resizeByWidth := resizeByWidth{width: newSize, enlarge: true}
	options, _ := resizeByWidth.CreateOptions(image)

	assert.Equal(t, newSize, options.Width)
}

func TestResizeByWidth_CreateOptions_Shrink_EnlargeNotAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Width / 2
	resizeByWidth := resizeByWidth{width: newSize, enlarge: false}
	options, _ := resizeByWidth.CreateOptions(image)

	assert.Equal(t, newSize, options.Width)
}

func TestResizeByWidth_CreateOptions_Enlarge_EnlargeNotAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Width * 2
	resizeByWidth := resizeByWidth{width: newSize, enlarge: false}
	options, _ := resizeByWidth.CreateOptions(image)

	assert.Equal(t, 0, options.Width)
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
