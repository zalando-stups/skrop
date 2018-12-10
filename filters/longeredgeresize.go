package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
)

// LongerEdgeResizeName is the name of the filter
const LongerEdgeResizeName = "longerEdgeResize"

type longerEdgeResize struct {
	size int
}

// NewLongerEdgeResize creates a new filter of this type
func NewLongerEdgeResize() filters.Spec {
	return &longerEdgeResize{}
}

func (f *longerEdgeResize) Name() string {
	return LongerEdgeResizeName
}

func (f *longerEdgeResize) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for longer edge resize ", f)

	imageSize, err := imageContext.Image.Size()

	if err != nil {
		return nil, err
	}

	if imageSize.Width > imageSize.Height {
		return &bimg.Options{
			Width: f.size}, nil
	}

	return &bimg.Options{
		Height: f.size}, nil
}

func (f *longerEdgeResize) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	if self.Width != 0 {
		return other.Width == 0 || other.Width == self.Width
	}

	//if Height was set
	return other.Height == 0 || other.Height == self.Height

}

func (f *longerEdgeResize) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {

	//if Width was set
	if self.Width != 0 {
		other.Width = self.Width
	}

	//if Height was set
	if self.Height != 0 {
		other.Height = self.Height
	}

	return other
}

func (f *longerEdgeResize) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &longerEdgeResize{}

	c.size, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (f *longerEdgeResize) Request(ctx filters.FilterContext) {}

func (f *longerEdgeResize) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
