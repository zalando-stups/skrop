package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
	"testing"
)

func TestNewBlur(t *testing.T) {
	name := NewBlur().Name()
	assert.Equal(t, "blur", name)
}

func TestBlur_Name(t *testing.T) {
	c := blur{}
	assert.Equal(t, "blur", c.Name())
}

func TestBlur_CreateOptions_ExplicitParam(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	blur := blur{Sigma: 19, MinAmpl: 21}
	options, _ := blur.CreateOptions(image, make(map[string][]string))

	blu := options.GaussianBlur

	assert.Equal(t, float64(19), blu.Sigma)
	assert.Equal(t, float64(21), blu.MinAmpl)
}

func TestBlur_CreateOptions_ImplicitParam(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	blur := blur{Sigma: 19}
	options, _ := blur.CreateOptions(image, make(map[string][]string))

	blu := options.GaussianBlur

	assert.Equal(t, float64(19), blu.Sigma)
	assert.Equal(t, float64(0), blu.MinAmpl)
}

func TestBlur_CanBeMerged_True(t *testing.T) {
	s := blur{}
	opt := &bimg.Options{}
	self := &bimg.Options{GaussianBlur: bimg.GaussianBlur{Sigma: 0.5, MinAmpl: 1.7}}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestBlur_CanBeMerged_False(t *testing.T) {
	s := blur{}
	opt := &bimg.Options{GaussianBlur: bimg.GaussianBlur{Sigma: 0.7, MinAmpl: 7.3}}
	self := &bimg.Options{GaussianBlur: bimg.GaussianBlur{Sigma: 0.5, MinAmpl: 1.7}}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestBlur_Merge(t *testing.T) {
	s := blur{}
	self := &bimg.Options{GaussianBlur: bimg.GaussianBlur{Sigma: 0.5, MinAmpl: 1.7}}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.GaussianBlur, opt.GaussianBlur)
}

func TestBlur_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewBlur, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "two args",
		Args: []interface{}{25.0, 35.0},
		Err:  false,
	}, {
		Msg:  "type error",
		Args: []interface{}{"abc", 2.6},
		Err:  true,
	}, {
		Msg:  "more args",
		Args: []interface{}{1.0, 3.0, 2.0, 1.0, 2.5, 3.7},
		Err:  true,
	}})
}
