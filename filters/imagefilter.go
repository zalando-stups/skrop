package filters

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/messages"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
)

const (
	// North Gravity
	North = "north"
	// South Gravity
	South = "south"
	// East Gravity
	East = "east"
	// West Gravity
	West = "west"
	// Center Gravity
	Center = "center"
	// Quality used by default if not specified
	Quality      = 100
	doNotEnlarge = "DO_NOT_ENLARGE"
	skropImage   = "skImage"
	skropOptions = "skOptions"
	skropInit    = "skInit"
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

// ImageFilter defines what a filter should implement
type ImageFilter interface {
	CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error)
	CanBeMerged(other *bimg.Options, self *bimg.Options) bool
	Merge(other *bimg.Options, self *bimg.Options) *bimg.Options
}

type ImageFilterContext struct {
	Image      *bimg.Image
	Parameters map[string][]string
}

func errorResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewBufferString(messages.Error500)),
	}
}

func buildParameters(ctx filters.FilterContext, image *bimg.Image) *ImageFilterContext {
	parameters := map[string][]string(nil)
	if ctx != nil {
		parameters = ctx.Request().URL.Query()
	}
	return &ImageFilterContext{
		Image:      image,
		Parameters: parameters,
	}
}

// HandleImageResponse should be called by the Response of every filter. It transforms the image
func HandleImageResponse(ctx filters.FilterContext, f ImageFilter) error {

	//in case the response had an error from the backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return fmt.Errorf("processing skipped, as the backend/filter reported %d status code", ctx.Response().StatusCode)
	}

	//executed while processing the first filter
	if _, ok := ctx.StateBag()[skropInit]; !ok {
		initResponse(ctx)
		ctx.StateBag()[skropInit] = true
	}

	image, ok := ctx.StateBag()[skropImage].(*bimg.Image)
	if !ok {
		//call the init func
		log.Error("context state bag does not contains the key ", skropImage)
		ctx.Serve(errorResponse())
		return errors.New("processing failed, image not exists in the state bag")
	}

	opt, err := f.CreateOptions(buildParameters(ctx, image))
	if err != nil {
		log.Error("Failed to create options ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	opts, ok := ctx.StateBag()[skropOptions].(*bimg.Options)
	if !ok {
		log.Error("context state bag does not contains the key ", skropImage)
		ctx.Serve(errorResponse())
		return errors.New("processing failed, initialization of options not successful")
	}

	if f.CanBeMerged(opts, opt) {
		ctx.StateBag()[skropOptions] = f.Merge(opts, opt)
		log.Debug("Filter ", f, " merged in ", ctx.StateBag()[skropOptions])
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
	newOption, err := f.CreateOptions(buildParameters(ctx, newImage))
	if err != nil {
		log.Error("Failed to create new options ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	ctx.StateBag()[skropImage] = newImage
	ctx.StateBag()[skropOptions] = newOption
	return nil
}

// FinalizeResponse is called at the end of the transformations on an image to empty the queue of
// operations to perform
func FinalizeResponse(ctx filters.FilterContext) {
	//in case the response had an error fromt he backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return
	}

	if _, ok := ctx.StateBag()[skropInit]; !ok {
		return
	}

	image := ctx.StateBag()[skropImage].(*bimg.Image)
	opts := ctx.StateBag()[skropOptions].(*bimg.Options)

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

	ctx.StateBag()[skropImage] = bimg.NewImage(buf)
	ctx.StateBag()[skropOptions] = &bimg.Options{}
}
