package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"testing"
)

func TestNewtransformFromQueryParams(t *testing.T) {
	name := NewTransformFromQueryParams().Name()
	assert.Equal(t, TransformByQueryParamsName, name)
}

func TestTransformFromQueryParams_CanBeMerged(t *testing.T) {
	ea := transformFromQueryParams{}
	assert.Equal(t, false, ea.CanBeMerged(nil, nil))
}

func TestTransformFromQueryParams_CreateOptions(t *testing.T) {
	ea := transformFromQueryParams{}
	img := imagefiltertest.LandscapeImage()
	options := make(map[string][]string)
	options[cropParameters] = []string{"10,10,100,100"}
	ctx := &ImageFilterContext{
		Image:      img,
		Parameters: options,
	}
	opts, err := ea.CreateOptions(ctx)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 10, opts.Top)
	assert.Equal(t, 10, opts.Left)
	assert.Equal(t, 100, opts.AreaHeight)
	assert.Equal(t, 100, opts.AreaWidth)
}
func TestTransformFromQueryParams_CreateOptions_Defaults(t *testing.T) {
	ea := transformFromQueryParams{}
	options := make(map[string][]string)
	img := imagefiltertest.LandscapeImage()
	ctx := &ImageFilterContext{
		Image:      img,
		Parameters: options,
	}
	//set defaults if not a valid value
	options[cropParameters] = []string{"a,b, c, d"}
	imgSize, _ := img.Size()
	opts, err := ea.CreateOptions(ctx)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, 0, opts.Left)
	assert.Equal(t, 0, opts.Top)
	assert.Equal(t, imgSize.Width, opts.AreaWidth)
	assert.Equal(t, imgSize.Height, opts.AreaHeight)

}
func TestTransformFromQueryParams_CreateOptions_ExtremeValues(t *testing.T) {
	//When given values exceed image size
	ea := transformFromQueryParams{}
	img := imagefiltertest.LandscapeImage()
	imgSize, _ := img.Size()
	options := make(map[string][]string)
	options[cropParameters] = []string{"100,100,10000,10000"}
	ctx := &ImageFilterContext{
		Image:      img,
		Parameters: options,
	}
	opts, err := ea.CreateOptions(ctx)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, imgSize.Height, opts.AreaHeight+opts.Top)
	assert.Equal(t, imgSize.Width, opts.AreaWidth+opts.Left)

}
