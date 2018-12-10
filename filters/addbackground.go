package filters

import (
	"github.com/h2non/bimg"
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
)

const (
	// AddBackgroundName is the name of the filter
	AddBackgroundName = "addBackground"
)

type addBackground bimg.Color

// NewAddBackground creates a new add backgroud filter
func NewAddBackground() filters.Spec {
	return &addBackground{}
}

func (s *addBackground) Name() string {
	return AddBackgroundName
}

func (s *addBackground) CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for adding background ", s)

	if imageContext.Image.Type() != "png" {
		return &bimg.Options{}, nil
	}

	backgroundColor := bimg.Color{R: s.R, G: s.G, B: s.B}

	return &bimg.Options{
		Background: backgroundColor}, nil

}

func (s *addBackground) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	zero := bimg.Color{}

	//it can be merged if the background was not set (in options or in self) or if they are set to the same value
	return other.Background == zero || other.Background == self.Background
}

func (s *addBackground) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Background = self.Background
	return other
}

func (s *addBackground) CreateFilter(args []interface{}) (filters.Filter, error) {
	if len(args) != 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	var err error
	f := &addBackground{}

	f.R, err = parse.EskipUint8Arg(args[0])
	if err != nil {
		return nil, err
	}

	f.G, err = parse.EskipUint8Arg(args[1])
	if err != nil {
		return nil, err
	}

	f.B, err = parse.EskipUint8Arg(args[2])
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *addBackground) Request(ctx filters.FilterContext) {}

func (s *addBackground) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, s)
}
