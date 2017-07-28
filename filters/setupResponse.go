package filters

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
)

// This filter is the default filter for every configuration of skrop
// It
type setupResponse struct{}

const (
	SetupResponseName = "setupResponse"
	SkropImage        = "skImage"
	SkropOptions      = "skOptions"
)

func NewSetupResponse() filters.Spec {
	return &setupResponse{}
}

func (s *setupResponse) Name() string {
	return SetupResponseName
}

func (s *setupResponse) CreateFilter(args []interface{}) (filters.Filter, error) {
	if len(args) != 0 {
		return nil, filters.ErrInvalidFilterParameters
	}

	return &addBackground{}, nil
}

func (s *setupResponse) Request(ctx filters.FilterContext) {}

//initialize the statebag where the options created by each filter and the image are stored
func (s *setupResponse) Response(ctx filters.FilterContext) {
	rsp := ctx.Response()

	defer rsp.Body.Close()

	rsp.Header.Del("Content-Length")

	buf, err := ioutil.ReadAll(rsp.Body)
	imageBytesLength := len(buf)

	log.Debug("Image bytes length: ", imageBytesLength)

	if err != nil {
		log.Error("Failed to process image ", err.Error())
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		})
		return
	}

	if imageBytesLength == 0 {
		log.Error("Original image is empty. Nothing to process")
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		})
		return
	}

	ctx.StateBag()[SkropImage] = bimg.NewImage(buf)
	ctx.StateBag()[SkropOptions] = &bimg.Options{}
}
