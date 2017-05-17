package filters

import (
	"image/color"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
)

const (
	AddBackgroundName = "addBackground"
)

type addBackground color.RGBA

func NewAddBackground() filters.Spec {
	return &addBackground{}
}

func (r *addBackground) Name() string {
	return AddBackgroundName
}

func (r *addBackground) CreateOptions(image *bimg.Image) (*bimg.Options, error) {

}

func (r *addBackground) CreateFilter(args []interface{}) (filters.Filter, error) {

}

func (r *addBackground) Request(ctx filters.FilterContext) {}

func (r *addBackground) Response(ctx filters.FilterContext) {
	handleResponse(ctx, r)
}