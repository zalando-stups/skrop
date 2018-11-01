package filters

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/zalando/skipper/filters"
)

func TestNewCropByFocalArea(t *testing.T) {
	name := NewCropByFocalArea().Name()
	assert.Equal(t, "cropByFocalArea", name)
}

func TestCropByFocalArea_Name(t *testing.T) {
	c := cropByFocalArea{}
	assert.Equal(t, "cropByFocalArea", c.Name())
}

func TestCropByFocalArea_CreateOptions(t *testing.T) {
	// Landscape image is 1000x668
	c := cropByFocalArea{}
	image := imagefiltertest.LandscapeImage()
	fc := createDefaultContext(t, "doesnotmatter.com")
	fc.FParams = make(map[string]string)
	fc.FParams["focalPointX"] = "50";
	fc.FParams["focalPointY"] = "50";
	fc.FParams["cropExpectedWidth"] = "200";
	fc.FParams["cropExpectedHeight"] = "100";

	options, _ := c.CreateOptions(buildParameters(fc, image))

	assert.Equal(t, 200, options.AreaWidth)
	assert.Equal(t, 100, options.AreaHeight)
	assert.Equal(t, 0, options.Top)
	assert.Equal(t, 0, options.Left)

	c = cropByFocalArea{}
	image = imagefiltertest.LandscapeImage()
	fc = createDefaultContext(t, "doesnotmatter.com")
	fc.FParams = make(map[string]string)
	fc.FParams["focalPointX"] = "300";
	fc.FParams["focalPointY"] = "300";
	fc.FParams["cropExpectedWidth"] = "200";
	fc.FParams["cropExpectedHeight"] = "100";

	options, _ = c.CreateOptions(buildParameters(fc, image))

	assert.Equal(t, 200, options.AreaWidth)
	assert.Equal(t, 100, options.AreaHeight)
	assert.Equal(t, 250, options.Top)
	assert.Equal(t, 200, options.Left)

	c = cropByFocalArea{}
	image = imagefiltertest.LandscapeImage()
	fc = createDefaultContext(t, "doesnotmatter.com")
	fc.FParams = make(map[string]string)
	fc.FParams["focalPointX"] = "1000";
	fc.FParams["focalPointY"] = "1000";
	fc.FParams["cropExpectedWidth"] = "200";
	fc.FParams["cropExpectedHeight"] = "100";

	options, _ = c.CreateOptions(buildParameters(fc, image))

	assert.Equal(t, 200, options.AreaWidth)
	assert.Equal(t, 100, options.AreaHeight)
	assert.Equal(t, 568, options.Top)
	assert.Equal(t, 900, options.Left)
}

func TestCropByFocalArea_CreateOptions_MissingPathParam(t *testing.T) {
	c := cropByFocalArea{}
	image := imagefiltertest.LandscapeImage()
	fc := createDefaultContext(t, "doesnotmatter.com")
	fc.FParams = make(map[string]string)
	fc.FParams["focalPointY"] = "334";

	options, err := c.CreateOptions(buildParameters(fc, image))

	assert.Nil(t, options)
	assert.Equal(t, filters.ErrInvalidFilterParameters, err)

	fc = createDefaultContext(t, "doesnotmatter.com")
	fc.FParams = make(map[string]string)
	fc.FParams["cropExpectedWidth"] = "334";

	options, err = c.CreateOptions(buildParameters(fc, image))

	assert.Nil(t, options)
	assert.Equal(t, filters.ErrInvalidFilterParameters, err)
}

func TestCropByFocalArea_CreateOptions_InvalidPathParam(t *testing.T) {
	c := cropByFocalArea{targetX: 0.5, targetY: 0.5, aspectRatio: 1.5}
	image := imagefiltertest.LandscapeImage()
	fc := createDefaultContext(t, "doesnotmatter.com")
	fc.FParams = make(map[string]string)
	fc.FParams["focalPointX"] = "xyz";
	fc.FParams["focalPointY"] = "abc";
	fc.FParams["cropExpectedWidth"] = "200";
	fc.FParams["cropExpectedHeight"] = "100";

	options, err := c.CreateOptions(buildParameters(fc, image))

	assert.Nil(t, options)
	assert.NotNil(t, err)

	fc.FParams["focalPointX"] = "100";
	fc.FParams["focalPointY"] = "abc";
	fc.FParams["cropExpectedWidth"] = "200";
	fc.FParams["cropExpectedHeight"] = "100";

	options, err = c.CreateOptions(buildParameters(fc, image))

	assert.Nil(t, options)
	assert.NotNil(t, err)
}

func TestCropByFocalArea_CanBeMerged(t *testing.T) {
	ea := transformFromQueryParams{}
	assert.Equal(t, false, ea.CanBeMerged(nil, nil))
}