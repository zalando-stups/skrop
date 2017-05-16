package filters

import (
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"github.com/stretchr/testify/assert"
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
	options, _ := quality.CreateOptions(image)

	assert.Equal(t, 75, options.Quality)
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
		Args: []interface{}{80.0, 100.0},
		Err:  true,
	}})
}
