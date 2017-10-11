package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	LongerEdgeResizeName = "longerEdgeResize"
)

type longerEdgeResize struct {
	size int
}

func NewLongerEdgeResize() filters.Spec {
	return &longerEdgeResize{}
}

func (r *longerEdgeResize) Name() string {
	return LongerEdgeResizeName
}

func (r *longerEdgeResize) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for longer edge resize ", r)

	imageSize, err := image.Size()

	if err != nil {
		return nil, err
	}

	if imageSize.Width > imageSize.Height {
		return &bimg.Options{
			Width: r.size}, nil
	} else {
		return &bimg.Options{
			Height: r.size}, nil
	}
}

func (s *longerEdgeResize) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	if self.Width != 0 {
		return other.Width == 0 || other.Width == self.Width
	}

	//if Height was set
	return other.Height == 0 || other.Height == self.Height

}

func (s *longerEdgeResize) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {

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

func (r *longerEdgeResize) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &longerEdgeResize{}

	f.size, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *longerEdgeResize) Request(ctx filters.FilterContext) {}

func (r *longerEdgeResize) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
