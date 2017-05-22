package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"testing"
)

func TestNewAddBackground(t *testing.T) {
	name := NewAddBackground().Name()
	assert.Equal(t, "addBackground", name)
}

func TestAddBackground_Name(t *testing.T) {
	c := addBackground{}
	assert.Equal(t, "addBackground", c.Name())
}

func TestAddBackground_CreateOptionsPNG(t *testing.T) {
	image := imagefiltertest.PNGImage()
	addBackground := addBackground{R: 1, G: 2, B: 3}
	options, _ := addBackground.CreateOptions(image)

	background := options.Background

	assert.Equal(t, uint8(1), background.R)
	assert.Equal(t, uint8(2), background.G)
	assert.Equal(t, uint8(3), background.B)
}

func TestAddBackground_CreateOptionsJPG(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	addBackground := addBackground{R: 1, G: 2, B: 3}
	options, _ := addBackground.CreateOptions(image)

	background := options.Background

	assert.Equal(t, uint8(0), background.R)
	assert.Equal(t, uint8(0), background.G)
	assert.Equal(t, uint8(0), background.B)
}

func TestAddBackground_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewAddBackground, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "three args",
		Args: []interface{}{25.0, 35.0, 103.0},
		Err:  false,
	}, {
		Msg:  "more than three args",
		Args: []interface{}{25.0, 100.0, 25.0, 100.0},
		Err:  true,
	}})
}
