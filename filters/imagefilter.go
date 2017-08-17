package filters

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zalando-incubator/skrop/messages"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
)

const (
	North        = "north"
	South        = "south"
	East         = "east"
	West         = "west"
	Center       = "center"
	Quality      = 100
	doNotEnlarge = "DO_NOT_ENLARGE"
	SkropImage   = "skImage"
	SkropOptions = "skOptions"
	SkropInit    = "skInit"
)

var (
	cropTypeToGravity map[string]bimg.Gravity
	cropTypes         map[string]bool
)

func init() {
	cropTypes = map[string]bool{
		North:  true,
		South:  true,
		East:   true,
		West:   true,
		Center: true}
	cropTypeToGravity = map[string]bimg.Gravity{
		North:  bimg.GravityNorth,
		South:  bimg.GravitySouth,
		East:   bimg.GravityEast,
		West:   bimg.GravityWest,
		Center: bimg.GravityCentre}
}

type ImageFilter interface {
	CreateOptions(image *bimg.Image) (*bimg.Options, error)
	CanBeMerged(other *bimg.Options, self *bimg.Options) bool
	Merge(other *bimg.Options, self *bimg.Options) *bimg.Options
}

func errorResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewBufferString(messages.Error500)),
	}
}

func HandleImageResponse(ctx filters.FilterContext, f ImageFilter) error {

	//in case the response had an error from the backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return errors.New(fmt.Sprintf("processing skipped, as the backend/filter reported %d status code", ctx.Response().StatusCode))
	}

	//executed while processing the first filter
	if _, ok := ctx.StateBag()[SkropInit]; !ok {
		initResponse(ctx)
		ctx.StateBag()[SkropInit] = true
	}

	image, ok := ctx.StateBag()[SkropImage].(*bimg.Image)
	if !ok {
		//call the init func
		log.Error("context state bag does not contains the key ", SkropImage)
		ctx.Serve(errorResponse())
		return errors.New("processing failed, image not exists in the state bag")
	}

	opt, err := f.CreateOptions(image)
	if err != nil {
		log.Error("Failed to create options ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	opts, ok := ctx.StateBag()[SkropOptions].(*bimg.Options)
	if !ok {
		log.Error("context state bag does not contains the key ", SkropImage)
		ctx.Serve(errorResponse())
		return errors.New("processing failed, initialization of options not successful.")
	}

	if f.CanBeMerged(opts, opt) {
		ctx.StateBag()[SkropOptions] = f.Merge(opts, opt)
		log.Debug("Filter ", f, " merged in ", ctx.StateBag()[SkropOptions])
		return nil
	}

	//transform the image
	buf, err := transformImage(image, opts)
	if err != nil {
		log.Error("Failed to process image ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	// set opt in the stateBag
	newImage := bimg.NewImage(buf)
	newOption, err := f.CreateOptions(newImage)
	if err != nil {
		log.Error("Failed to create new options ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	ctx.StateBag()[SkropImage] = newImage
	ctx.StateBag()[SkropOptions] = newOption
	return nil
}

func FinalizeResponse(ctx filters.FilterContext) {
	//in case the response had an error fromt he backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return
	}

	if _, ok := ctx.StateBag()[SkropInit]; !ok {
		return
	}

	image := ctx.StateBag()[SkropImage].(*bimg.Image)
	opts := ctx.StateBag()[SkropOptions].(*bimg.Options)

	buf, err := transformImage(image, opts)
	if err != nil {
		log.Error("failed to process image ", err.Error())
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(messages.Error500)),
		})
		return
	}

	rsp := ctx.Response()

	defer rsp.Body.Close()

	rsp.Body = ioutil.NopCloser(bytes.NewReader(buf))
}

func transformImage(image *bimg.Image, opts *bimg.Options) ([]byte, error) {
	defOpt := applyDefaults(opts)

	log.Debug("successfully applied the following options on the image: ", opts)

	transformedImageBytes, err := image.Process(*defOpt)

	if err != nil {
		return nil, err
	}

	return transformedImageBytes, nil
}

func applyDefaults(o *bimg.Options) *bimg.Options {
	if o.Quality == 0 {
		o.Quality = Quality
	}
	if o.Background == bimg.ColorBlack {
		o.Background = bimg.Color{R: 255, G: 255, B: 255}
	}
	return o
}

func initResponse(ctx filters.FilterContext) {
	rsp := ctx.Response()

	defer rsp.Body.Close()

	rsp.Header.Del("Content-Length")

	buf, err := ioutil.ReadAll(rsp.Body)
	imageBytesLength := len(buf)

	log.Debug("Image bytes length: ", imageBytesLength)

	if err != nil {
		log.Error("failed to process image ", err.Error())
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(messages.Error500)),
		})
		return
	}

	if imageBytesLength == 0 {
		log.Error("original image is empty")
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(messages.Error404)),
		})
		return
	}

	ctx.StateBag()[SkropImage] = bimg.NewImage(buf)
	ctx.StateBag()[SkropOptions] = &bimg.Options{}
}
