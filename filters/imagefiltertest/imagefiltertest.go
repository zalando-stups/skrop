package imagefiltertest

import (
	"testing"

	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
)

const (
	// LandscapeImageFile  is the path to the sample image landscape
	LandscapeImageFile = "../images/lisbon-tram.jpg"
	// PortraitImageFile   is the path to the sample image portrait
	PortraitImageFile = "../images/big-ben.jpg"
	// PNGImageFile        is the path to the sample PNG image
	PNGImageFile = "../images/bag.png"
)

// CreateTestItem is a utility to test arguments in eskip files
type CreateTestItem struct {
	Msg  string
	Args []interface{}
	Err  bool
}

// TestCreate autonmatically creates a test for the CreateTestItem
func TestCreate(t *testing.T, spec func() filters.Spec, items []CreateTestItem) {
	for _, ti := range items {
		func() {
			f, err := spec().CreateFilter(ti.Args)
			switch {
			case ti.Err && err == nil:
				t.Error(ti.Msg, "failed to fail")
			case !ti.Err && err != nil:
				t.Error(ti.Msg, err)
			case err == nil && f == nil:
				t.Error(ti.Msg, "failed to create filter")
			}
		}()
	}
}

// LandscapeImage returns a landscape test image
func LandscapeImage() *bimg.Image {
	buffer, _ := bimg.Read(LandscapeImageFile)
	return bimg.NewImage(buffer)
}

// PortraitImage returns a portrait test image
func PortraitImage() *bimg.Image {
	buffer, _ := bimg.Read(PortraitImageFile)
	return bimg.NewImage(buffer)
}

// PNGImage returns a PNG test image
func PNGImage() *bimg.Image {
	buffer, _ := bimg.Read(PNGImageFile)
	return bimg.NewImage(buffer)
}
