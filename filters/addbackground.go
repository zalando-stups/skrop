package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/tools"
)

const (
	AddBackgroundName = "addBackground"
)

type addBackground bimg.Color

func NewAddBackground() filters.Spec {
	return &addBackground{}
}

func (r *addBackground) Name() string {
	return AddBackgroundName
}

func (r *addBackground) CreateOptions(image *bimg.Image) (*bimg.Options, error) {
	log.Debug("Create options for adding background ", r)

	if image.Type() != "png" {
		return &bimg.Options{}, nil
	}

	backgroundColor := bimg.Color{R: r.R, G: r.G, B: r.B}

	return &bimg.Options{
		Background: backgroundColor}, nil

}

func (r *addBackground) CreateFilter(args []interface{}) (filters.Filter, error) {
	if len(args) != 3 {
		return nil, filters.ErrInvalidFilterParameters
	}

	var err error
	f:= &addBackground{}

	f.R, err = tools.ParseEskipUint8Arg(args[0])
	if err != nil {
		return nil, err
	}

	f.G, err = tools.ParseEskipUint8Arg(args[1])
	if err != nil {
		return nil, err
	}

	f.B, err = tools.ParseEskipUint8Arg(args[2])
	if err != nil {
		return nil, err
	}


	return f, nil
}

func (r *addBackground) Request(ctx filters.FilterContext) {}

func (r *addBackground) Response(ctx filters.FilterContext) {
	HandleImageResponse(ctx, r)
}
