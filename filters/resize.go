package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
	"math"
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

func (h *resize) Name() string {
	return ResizeName
}

func (h *resize) CreateFilter(args []interface{}) (filters.Filter, error) {
	if len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	f := &resize{}

	if width, ok := args[0].(float64); ok && math.Trunc(width) == width {
		f.width = int(width)
	} else {
		log.Error("Width not correct")
		return nil, filters.ErrInvalidFilterParameters
	}

	if height, ok := args[1].(float64); ok && math.Trunc(height) == height {
		f.width = int(height)
	} else {
		log.Error("Height not correct")
		return nil, filters.ErrInvalidFilterParameters
	}

	return f, nil
}

func (h *resize) Request(ctx filters.FilterContext) {}

func resizeImage(out *io.PipeWriter, in io.ReadCloser, opts *bimg.Options) error {
	var err error

	defer func() {
		if err == nil {
			err = io.EOF
		}
		out.CloseWithError(err)
		in.Close()
	}()

	responseImage, err := ioutil.ReadAll(in)

	if err != nil {
		return err
	}

	newImage, err := bimg.NewImage(responseImage).Process(*opts)

	if err != nil {
		return err
	}

	_, err = out.Write(newImage)
	if err != nil {
		return err
	}

	return nil
}

func (h *resize) Response(ctx filters.FilterContext) {
	rsp := ctx.Response()

	rsp.Header.Del("Content-Length")

	in := rsp.Body
	r, w := io.Pipe()
	rsp.Body = r

	options := &bimg.Options{
		Width:  h.width,
		Height: h.height}

	go resizeImage(w, in, options)
}
