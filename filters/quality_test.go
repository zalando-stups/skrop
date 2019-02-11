package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
	"testing"
)

func TestNewQuality(t *testing.T) {
	name := NewQuality().Name()
	assert.Equal(t, "quality", name)
}

func TestNewQuality_Name(t *testing.T) {
	c := quality{}
	assert.Equal(t, "quality", c.Name())
}

func TestNewQuality_CreateOptions(t *testing.T) {
	quality := quality{percentage: 75}
	image := imagefiltertest.LandscapeImage()
	options, _ := quality.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, 75, options.Quality)
}

func TestNewQuality_CanBeMerged_True(t *testing.T) {
	s := quality{}
	opt := &bimg.Options{}
	self := &bimg.Options{Quality: 85}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestNewQuality_CanBeMerged_False(t *testing.T) {
	s := quality{}
	opt := &bimg.Options{Quality: 27}
	self := &bimg.Options{Quality: 85}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestNewQuality_Merge(t *testing.T) {
	s := quality{}
	self := &bimg.Options{Quality: 85}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Quality, opt.Quality)
}

func TestNewQuality_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewQuality, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "one arg",
		Args: []interface{}{80.0},
		Err:  false,
	}, {
		Msg:  "too high value",
		Args: []interface{}{110.0},
		Err:  true,
	}, {
		Msg:  "more than one args",
		Args: []interface{}{80.0, 90.0},
		Err:  true,
	}})
}
