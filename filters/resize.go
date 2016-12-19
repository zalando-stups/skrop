package filters

import (
	"github.com/zalando/skipper/filters"
	"net/http"
)

const (
	ResizeName = "resize"
)

type resize struct {
	width  int
	height int
}

func NewResize() filters.Spec {
	return &resize{}
}

func (h *resize) Name() string { return ResizeName }

func (h *resize) CreateFilter(_ []interface{}) (filters.Filter, error) { return h, nil }
func (h *resize) Request(ctx filters.FilterContext)                    {}
func (h *resize) Response(ctx filters.FilterContext)                   { ctx.Response().StatusCode = http.StatusOK }
