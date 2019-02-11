package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
)

func TestNewResize(t *testing.T) {
	name := NewResize().Name()
	assert.Equal(t, "resize", name)
}

func TestResize_Name(t *testing.T) {
	r := resize{}
	assert.Equal(t, "resize", r.Name())
}

func TestResize_CreateOptions_IgnoreProportions_Explicit(t *testing.T) {
	r := resize{width: 800, height: 600, keepAspectRatio: false}
	options, _ := r.CreateOptions(buildParameters(nil, imagefiltertest.LandscapeImage()))

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 600, options.Height)
	assert.False(t, r.keepAspectRatio)
}

func TestResize_CreateOptions_IgnoreProportions_Implicit(t *testing.T) {
	r := resize{width: 800, height: 600}

	options, _ := r.CreateOptions(buildParameters(nil, imagefiltertest.LandscapeImage()))

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 600, options.Height)
	assert.False(t, r.keepAspectRatio)
}

func TestResize_CreateOptions_WithProportions1(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()

	newHeight := size.Height - 10

	r := resize{width: size.Width, height: newHeight, keepAspectRatio: true}
	options, _ := r.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newHeight, options.Height)
	assert.Zero(t, options.Width)
}

func TestResize_CreateOptions_WithProportions2(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()

	newWidth := size.Width - 10

	r := resize{width: newWidth, height: size.Height, keepAspectRatio: true}
	options, _ := r.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newWidth, options.Width)
	assert.Zero(t, options.Height)
}

func TestResize_CanBeMerged_True(t *testing.T) {
	s := resize{}
	opt := &bimg.Options{}
	self := &bimg.Options{Width: 265, Height: 365}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestResize_CanBeMerged_False(t *testing.T) {
	s := resize{}
	opt := &bimg.Options{Width: 265, Height: 365}
	self := &bimg.Options{Width: 144, Height: 963}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestResize_Merge(t *testing.T) {
	s := resize{}
	self := &bimg.Options{Width: 265, Height: 365}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Width, opt.Width)
	assert.Equal(t, self.Height, opt.Height)
}

func TestResize_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewResize, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "two args",
		Args: []interface{}{800.0, 600.0},
		Err:  false,
	}, {
		Msg:  "three args",
		Args: []interface{}{800.0, 600.0, "ignoreAspectRatio"},
		Err:  false,
	}, {
		Msg:  "more than 3 args",
		Args: []interface{}{800.0, 200.0, "test", "Whaaat!"},
		Err:  true,
	}})
}

func TestResize_CreateFilter_IgnoreAspectRatio(t *testing.T) {

	f := &resize{}

	nf, _ := f.CreateFilter([]interface{}{100.0, 50.0, "ignoreAspectRatio"})
	i := nf.(*resize)

	assert.Equal(t, false, i.keepAspectRatio)
}
func TestResize_CreateFilter_KeepAspectRatio(t *testing.T) {

	f := &resize{}

	nf, _ := f.CreateFilter([]interface{}{100.0, 50.0, "anything"})
	i := nf.(*resize)

	assert.Equal(t, true, i.keepAspectRatio)
}
