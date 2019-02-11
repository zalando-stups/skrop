package filters

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/h2non/bimg"
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/messages"
	"github.com/zalando/skipper/filters"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	Quality          = 100
	doNotEnlarge     = "DO_NOT_ENLARGE"
	skropImage       = "skImage"
	hasMergedFilters = "hasMergedFilters"
	skropOptions     = "skOptions"
	skropInit        = "skInit"
)

var (
	cropTypeToGravity map[string]bimg.Gravity
	cropTypes         map[string]bool
	stripMetadata     bool
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

	val, exists := os.LookupEnv("STRIP_METADATA")
	if exists && strings.ToUpper(val) == "TRUE" {
		stripMetadata = true
	}
}

// ImageFilter defines what a filter should implement
type ImageFilter interface {
	CreateOptions(imageContext *ImageFilterContext) (*bimg.Options, error)
	CanBeMerged(other *bimg.Options, self *bimg.Options) bool
	Merge(other *bimg.Options, self *bimg.Options) *bimg.Options
}

type ImageFilterContext struct {
	Image         *bimg.Image
	Parameters    map[string][]string
	filterContext *filters.FilterContext
}

func (c *ImageFilterContext) PathParam(key string) string { return (*c.filterContext).PathParam(key) }

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
		Image:         image,
		Parameters:    parameters,
		filterContext: &ctx,
	}
}

// HandleImageResponse should be called by the Response of every filter. It transforms the image
func HandleImageResponse(ctx filters.FilterContext, f ImageFilter) error {

	log.Debug("Handle Image Response")

	//in case the response had an error from the backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return fmt.Errorf("processing skipped, as the backend/filter reported %d status code", ctx.Response().StatusCode)
	}

	//executed while processing the first filter
	if _, ok := ctx.StateBag()[skropInit]; !ok {
		initResponse(ctx)
		ctx.StateBag()[skropInit] = true
		ctx.StateBag()[hasMergedFilters] = false
	}

	image, ok := ctx.StateBag()[skropImage].(*bimg.Image)
	if !ok {
		//call the init func
		log.Error("context state bag does not contains the key ", skropImage)
		ctx.Serve(errorResponse())
		return errors.New("processing failed, image not exists in the state bag")
	}

	optionsFromRequest, err := f.CreateOptions(buildParameters(ctx, image))
	if err != nil {
		log.Error("Failed to create options ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	optionsFromStateBag, ok := ctx.StateBag()[skropOptions].(*bimg.Options)
	if !ok {
		log.Error("context state bag does not contains the key ", skropImage)
		ctx.Serve(errorResponse())
		return errors.New("processing failed, initialization of options not successful")
	}

	if f.CanBeMerged(optionsFromStateBag, optionsFromRequest) {
		ctx.StateBag()[skropOptions] = f.Merge(optionsFromStateBag, optionsFromRequest)
		ctx.StateBag()[hasMergedFilters] = true
		log.Debug("Filter ", f, " merged in ", ctx.StateBag()[skropOptions])
		return nil
	}

	log.Debugf("Transform the image based on the options from request: %+v", optionsFromRequest)
	buf, err := transformImage(image, optionsFromRequest)
	if err != nil {
		log.Error("Failed to process image ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	newImage := bimg.NewImage(buf)
	if err != nil {
		log.Error("Failed to create new options ", err.Error())
		ctx.Serve(errorResponse())
		return err
	}

	ctx.StateBag()[skropImage] = newImage
	ctx.StateBag()[skropOptions] = &bimg.Options{}
	return nil
}

// FinalizeResponse is called at the end of the transformations on an image to empty the queue of
// operations to perform
func FinalizeResponse(ctx filters.FilterContext) {
	log.Debug("Finalize Response")
	//in case the response had an error from he backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return
	}

	if _, ok := ctx.StateBag()[skropInit]; !ok {
		return
	}

	image := ctx.StateBag()[skropImage].(*bimg.Image)
	opts := ctx.StateBag()[skropOptions].(*bimg.Options)

	buf := image.Image()

	var err error

	if ctx.StateBag()[hasMergedFilters] == true {
		buf, err = transformImage(image, opts)
		ctx.StateBag()[hasMergedFilters] = false
	}

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

	log.Debugf("successfully applied the following options on the image: %+v\n", opts)

	transformedImageBytes, err := image.Process(*defOpt)

	if err != nil {
		return nil, err
	}

	return transformedImageBytes, nil
}

func applyDefaults(o *bimg.Options) *bimg.Options {
	if (stripMetadata) {
		o.StripMetadata = true
	}
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
