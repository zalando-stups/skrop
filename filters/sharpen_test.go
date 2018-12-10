package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
	"testing"
)

func TestNewSharpen(t *testing.T) {
	name := NewSharpen().Name()
	assert.Equal(t, "sharpen", name)
}

func TestSharpen_Name(t *testing.T) {
	c := sharpen{}
	assert.Equal(t, "sharpen", c.Name())
}

func TestSharpen_CreateOptions(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	sharpen := sharpen{Radius: 1, X1: 2, Y2: 3, Y3: 4, M1: 5, M2: 6}
	options, _ := sharpen.CreateOptions(buildParameters(nil, image))

	sha := options.Sharpen

	assert.Equal(t, int(1), sha.Radius)
	assert.Equal(t, float64(2), sha.X1)
	assert.Equal(t, float64(3), sha.Y2)
	assert.Equal(t, float64(4), sha.Y3)
	assert.Equal(t, float64(5), sha.M1)
	assert.Equal(t, float64(6), sha.M2)
}

func TestSharpen_CanBeMerged_True(t *testing.T) {
	s := sharpen{}
	opt := &bimg.Options{}
	self := &bimg.Options{Sharpen: bimg.Sharpen{Radius: 1, X1: 2, Y2: 3, Y3: 4, M1: 5, M2: 6}}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestSharpen_CanBeMerged_False(t *testing.T) {
	s := sharpen{}
	opt := &bimg.Options{Sharpen: bimg.Sharpen{Radius: 1, X1: 2, Y2: 3, Y3: 4, M1: 5, M2: 6}}
	self := &bimg.Options{Sharpen: bimg.Sharpen{Radius: 9, X1: 8, Y2: 7, Y3: 6, M1: 5, M2: 4}}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestSharpen_Merge(t *testing.T) {
	s := sharpen{}
	self := &bimg.Options{Sharpen: bimg.Sharpen{Radius: 1, X1: 2, Y2: 3, Y3: 4, M1: 5, M2: 6}}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Sharpen, opt.Sharpen)
}

func TestSharpen_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewSharpen, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "three args",
		Args: []interface{}{25.0, 35.0, 103.0},
		Err:  true,
	}, {
		Msg:  "type error",
		Args: []interface{}{25.0, 35.0, 103.0, "abc", 1.0, 2.6},
		Err:  true,
	}, {
		Msg:  "six args",
		Args: []interface{}{1.0, 3.0, 2.0, 1.0, 2.5, 3.7},
		Err:  false,
	}})
}
