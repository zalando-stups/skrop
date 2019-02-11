package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
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

func TestAddBackground_CanBeMerged_True(t *testing.T) {
	s := addBackground{}
	opt := &bimg.Options{}
	self := &bimg.Options{Background: bimg.Color{R: 240, G: 0, B: 200}}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestAddBackground_CanBeMerged_False(t *testing.T) {
	s := addBackground{}
	opt := &bimg.Options{Background: bimg.Color{R: 10, G: 153, B: 200}}
	self := &bimg.Options{Background: bimg.Color{R: 240, G: 0, B: 200}}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestAddBackground_Merge(t *testing.T) {
	s := addBackground{}
	self := &bimg.Options{Background: bimg.Color{R: 240, G: 0, B: 200}}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Background, opt.Background)
}

func TestAddBackground_CreateOptionsPNG(t *testing.T) {
	image := imagefiltertest.PNGImage()
	addBackground := addBackground{R: 1, G: 2, B: 3}
	options, _ := addBackground.CreateOptions(buildParameters(nil, image))

	background := options.Background

	assert.Equal(t, uint8(1), background.R)
	assert.Equal(t, uint8(2), background.G)
	assert.Equal(t, uint8(3), background.B)
}

func TestAddBackground_CreateOptionsJPG(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	addBackground := addBackground{R: 1, G: 2, B: 3}
	options, _ := addBackground.CreateOptions(buildParameters(nil, image))

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
