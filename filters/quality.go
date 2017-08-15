package filters

import (
	log "github.com/sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/parse"
)

const (
	QualityName = "quality"
)

type quality struct {
	percentage int
}

func NewQuality() filters.Spec {
	return &quality{}
}

func (r *quality) Name() string {
	return QualityName
}

func (r *quality) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for quality ", r)

	return &bimg.Options{
		Quality: r.percentage}, nil
}

func (s *quality) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return other.Quality == 0 || other.Quality == self.Quality
}

func (s *quality) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Quality = self.Quality
	return other
}

func (r *quality) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &quality{}

	f.percentage, err = parse.EskipIntArg(args[0])

	if err != nil {
		return nil, err
	} else if f.percentage > 100 {
		return nil, filters.ErrInvalidFilterParameters
	}

	return f, nil
}

func (r *quality) Request(ctx filters.FilterContext) {}

func (r *quality) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
