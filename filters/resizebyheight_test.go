package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
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

func TestResizeByHeight_CreateOptions_Shrink_EnlargeAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Height / 2
	resizeByHeight := resizeByHeight{height: newSize, enlarge: true}
	options, _ := resizeByHeight.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newSize, options.Height)
}

func TestResizeByHeight_CreateOptions_Enlarge_EnlargeAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Height * 2
	resizeByHeight := resizeByHeight{height: newSize, enlarge: true}
	options, _ := resizeByHeight.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newSize, options.Height)
}

func TestResizeByHeight_CreateOptions_Shrink_EnlargeNotAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Height / 2
	resizeByHeight := resizeByHeight{height: newSize, enlarge: false}
	options, _ := resizeByHeight.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, newSize, options.Height)
}

func TestResizeByHeight_CreateOptions_Enlarge_EnlargeNotAllowed(t *testing.T) {
	image := imagefiltertest.LandscapeImage()
	size, _ := image.Size()
	newSize := size.Height * 2
	resizeByHeight := resizeByHeight{height: newSize, enlarge: false}
	options, _ := resizeByHeight.CreateOptions(buildParameters(nil, image))

	assert.Equal(t, 0, options.Height)
}

func TestResizeByHeight_CanBeMerged_True(t *testing.T) {
	s := resizeByHeight{}
	opt := &bimg.Options{}
	self := &bimg.Options{Height: 365}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestResizeByHeight_CanBeMerged_False(t *testing.T) {
	s := resizeByHeight{}
	opt := &bimg.Options{Height: 365}
	self := &bimg.Options{Height: 963}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestResizeByHeight_Merge(t *testing.T) {
	s := resizeByHeight{}
	self := &bimg.Options{Height: 365}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Height, opt.Height)
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
