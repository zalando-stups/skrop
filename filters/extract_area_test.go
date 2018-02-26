package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"testing"
)

func TestNewExtractArea(t *testing.T) {
	name := NewExtractArea().Name()
	assert.Equal(t, ExtractArea, name)
}

func TestExtractArea_CanBeMerged(t *testing.T) {
	ea := extractArea{}
	assert.Equal(t, false, ea.CanBeMerged(nil, nil))
}

func TestExtractArea_CreateOptions(t *testing.T) {
	ea := extractArea{}
	img := imagefiltertest.LandscapeImage()
	options := make(map[string][]string)
	options[cropFrom] = []string{"10,10"}
	options[cropHeight] = []string{"100"}
	options[cropWidth] = []string{"100"}
	_, err := ea.CreateOptions(img, options)
	assert.Nil(t, err, "error should be nil")

	options[cropFrom] = []string{"a,10"}
	options[cropHeight] = []string{"100"}
	options[cropWidth] = []string{"100"}
	_, err = ea.CreateOptions(img, options)
	assert.NotNil(t, err, "error should not be nil when point has invalid value")

	options[cropFrom] = []string{"a,10"}
	options[cropHeight] = []string{"a"}
	options[cropWidth] = []string{"100"}
	_, err = ea.CreateOptions(img, options)
	assert.NotNil(t, err, "error should not be nil when height has invalid value")

	options[cropFrom] = []string{"a,10"}
	options[cropHeight] = []string{"100"}
	options[cropWidth] = []string{"a"}
	_, err = ea.CreateOptions(img, options)
	assert.NotNil(t, err, "error should not be nil when width has invalid value")

}
