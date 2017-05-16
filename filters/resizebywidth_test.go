package filters

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.bus.zalan.do/hollywood/walk-of-fame/filters/imagefiltertest"
)

func TestNewResizeByWidth(t *testing.T) {
	name := NewResizeByWidth().Name()
	assert.Equal(t, "width", name)
}

func TestResizeByWidth_Name(t *testing.T) {
	c:= resizeByWidth{}
	assert.Equal(t, "width", c.Name())
}

func TestResizeByWidth_CreateOptions(t *testing.T) {
	resizeByWidth := resizeByWidth{width: 256}
	image := imagefiltertest.LandscapeImage()
	options, _ := resizeByWidth.CreateOptions(image)

	assert.Equal(t, 256, options.Width)
}

func TestResizeByWidth_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewResizeByWidth, []imagefiltertest.CreateTestItem{{
		Msg: "no args",
		Args: nil,
		Err: true,
	},{
		Msg: "one arg",
		Args: []interface{}{256.0},
		Err: false,
	},{
		Msg: "more than one args",
		Args: []interface{}{256.0, 100.0},
		Err: true,
	}})
}