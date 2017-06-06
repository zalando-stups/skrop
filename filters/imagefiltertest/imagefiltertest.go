package imagefiltertest

import (
	"testing"

	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	LandscapeImageFile = "../images/lisbon-tram.jpg"
	PortraitImageFile  = "../images/big-ben.jpg"
	PNGImageFile       = "../images/bag.png"
)

type FakeImageFilter bimg.Options

type CreateTestItem struct {
	Msg  string
	Args []interface{}
	Err  bool
}

func (h *FakeImageFilter) CreateOptions(_ *bimg.Image) (*bimg.Options, error) {
	options := bimg.Options(*h)
	return &options, nil
}

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

func LandscapeImage() *bimg.Image {
	buffer, _ := bimg.Read(LandscapeImageFile)
	return bimg.NewImage(buffer)
}

func PortraitImage() *bimg.Image {
	buffer, _ := bimg.Read(PortraitImageFile)
	return bimg.NewImage(buffer)
}

func PNGImage() *bimg.Image {
	buffer, _ := bimg.Read(PNGImageFile)
	return bimg.NewImage(buffer)
}
