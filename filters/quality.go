package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"github.com/h2non/bimg"
)

// QualityName is the name of the filter
const QualityName = "quality"

type quality struct {
	percentage int
}

// NewQuality creates a new filter of this type
func NewQuality() filters.Spec {
	return &quality{}
}

func (f *quality) Name() string {
	return QualityName
}

func (f *quality) CreateOptions(_ *ImageFilterContext) (*bimg.Options, error) {
	log.Debug("Create options for quality ", f)

	return &bimg.Options{
		Quality: f.percentage}, nil
}

func (f *quality) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return other.Quality == 0 || other.Quality == self.Quality
}

func (f *quality) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Quality = self.Quality
	return other
}

func (f *quality) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	c := &quality{}

	c.percentage, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	} else if c.percentage > 100 {
		return nil, filters.ErrInvalidFilterParameters
	}

	return c, nil
}

func (f *quality) Request(ctx filters.FilterContext) {}

func (f *quality) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, f)
}
