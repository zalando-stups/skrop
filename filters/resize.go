package filters

import (
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
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

func (h *resize) Request(ctx filters.FilterContext) {
	log.Debug("We start doing the resize request")
}

func resizeImage(out *io.PipeWriter, in io.ReadCloser) {
	var (
		err error
	)

	defer func() {

		if err == nil {
			err = io.EOF
		}

		out.CloseWithError(err)
		in.Close()
	}()

	options := bimg.Options{
		Width: 600,
	}

	responseImage, err := ioutil.ReadAll(in)

	newImage, err := bimg.NewImage(responseImage).Process(options)

	out.Write(newImage)
}

func (h *resize) Response(ctx filters.FilterContext) {
	rsp := ctx.Response()
	in := rsp.Body
	r, w := io.Pipe()
	rsp.Body = r
	go resizeImage(w, in)
}

func (h *resize) CreateFilter(args []interface{}) (filters.Filter, error) {
	if len(args) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	var width int
	switch w := args[0].(type) {
	case float64:
		width = int(w)
	default:
		log.Error("Width not correct")
		return nil, filters.ErrInvalidFilterParameters
	}

	var height int
	switch h := args[0].(type) {
	case float64:
		height = int(h)
	default:
		log.Error("Height not correct")
		return nil, filters.ErrInvalidFilterParameters
	}

	return &resize{width: width, height: height}, nil
}
