package filters

import (
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
	"math"
)

type ImageFilter interface {
	CreateOptions() *bimg.Options
}

func handleResponse(ctx filters.FilterContext, f ImageFilter) {
	rsp := ctx.Response()

	rsp.Header.Del("Content-Length")

	in := rsp.Body
	r, w := io.Pipe()
	rsp.Body = r

	options := f.CreateOptions()

	go transformImage(w, in, options)
}

func transformImage(out *io.PipeWriter, in io.ReadCloser, opts *bimg.Options) error {
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

	return err
}

func parseEskipIntArg(arg interface{}) (int, error) {
	if number, ok := arg.(float64); ok && math.Trunc(number) == number {
		return int(number), nil
	} else {
		return 0, filters.ErrInvalidFilterParameters
	}
}
