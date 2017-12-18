package imagefiltertest

import (
	"testing"

	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	// LandscapeImageFile  is the path to the sample image landscape
	LandscapeImageFile = "../images/lisbon-tram.jpg"
	// PortraitImageFile   is the path to the sample image portrait
	PortraitImageFile  = "../images/big-ben.jpg"
	// PNGImageFile        is the path to the sample PNG image
	PNGImageFile       = "../images/bag.png"
)

type FakeImageFilter bimg.Options

// CreateTestItem is a utility to test arguments in eskip files
type CreateTestItem struct {
	Msg  string
	Args []interface{}
	Err  bool
}

func (f *FakeImageFilter) CreateOptions(_ *bimg.Image) (*bimg.Options, error) {
	options := bimg.Options(*s)
	return &options, nil
}

func (f *FakeImageFilter) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.Width == 0 && other.Height == 0) ||
		(other.Width == self.Width && other.Height == self.Height)
}

func (f *FakeImageFilter) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	other.Quality = self.Quality
	other.Background = self.Background
	return other
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
