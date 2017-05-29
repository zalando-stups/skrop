package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"testing"
)

func TestNewResizeByHeight(t *testing.T) {
	name := NewResizeByHeight().Name()
	assert.Equal(t, "height", name)
}

func TestResizeByHeight_Name(t *testing.T) {
	c := resizeByHeight{}
	assert.Equal(t, "height", c.Name())
}

func TestResizeByHeight_CreateOptions(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	resizeByHeight := resizeByHeight{height: 150}
	options, _ := resizeByHeight.CreateOptions(image)

	assert.Equal(t, 150, options.Width)
}

func TestResizeByHeight_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewResizeByHeight, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "one arg",
		Args: []interface{}{256.0},
		Err:  false,
	}, {
		Msg:  "more than one args",
		Args: []interface{}{256.0, 100.0},
		Err:  true,
	}})
}
