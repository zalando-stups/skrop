package filters

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
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
	options, _ := sharpen.CreateOptions(image)

	sha := options.Sharpen

	assert.Equal(t, int(1), sha.Radius)
	assert.Equal(t, float64(2), sha.X1)
	assert.Equal(t, float64(3), sha.Y2)
	assert.Equal(t, float64(4), sha.Y3)
	assert.Equal(t, float64(5), sha.M1)
	assert.Equal(t, float64(6), sha.M2)
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
		Msg:  "six args",
		Args: []interface{}{1.0, 3.0, 2.0, 1.0, 2.5, 3.7},
		Err:  false,
	}})
}