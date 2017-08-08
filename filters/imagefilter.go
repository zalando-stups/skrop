package filters

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/zalando/skipper/filters"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
	"github.com/zalando-incubator/skrop/messages"
)

const (
	North        = "north"
	South        = "south"
	East         = "east"
	West         = "west"
	Center       = "center"
	Quality      = 100
	doNotEnlarge = "DO_NOT_ENLARGE"
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

func HandleImageResponse(ctx filters.FilterContext, f ImageFilter) {
	//in case the response had an error fromt he backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return
	}

	image, ok := ctx.StateBag()[SkropImage].(*bimg.Image)
	if !ok {
		log.Error("context state bag does not contains the key ", SkropImage)
		ctx.Serve(errorResponse())
		return
	}

	opt, err := f.CreateOptions(image)
	if err != nil {
		log.Error("Failed to create options ", err.Error())
		ctx.Serve(errorResponse())
		return
	}

	opts, ok := ctx.StateBag()[SkropOptions].(*bimg.Options)
	if !ok {
		log.Error("context state bag does not contains the key ", SkropImage)
		ctx.Serve(errorResponse())
		return
	}

	if f.CanBeMerged(opts, opt) {
		ctx.StateBag()[SkropOptions] = f.Merge(opts, opt)
		log.Debug("Filter ", f, " merged in ", ctx.StateBag()[SkropOptions])
		return
	}

	//transform the image
	buf, err := transformImage(image, opts)
	if err != nil {
		log.Error("Failed to process image ", err.Error())
		ctx.Serve(errorResponse())
		return
	}

	// set opt in the stateBag
	newImage := bimg.NewImage(buf)
	newOption, err := f.CreateOptions(newImage)
	if err != nil {
		log.Error("Failed to create new options ", err.Error())
		ctx.Serve(errorResponse())
		return
	}

	ctx.StateBag()[SkropImage] = newImage
	ctx.StateBag()[SkropOptions] = newOption

}

func FinalizeResponse(ctx filters.FilterContext) {
	//in case the response had an error fromt he backend or from a previous filter
	if ctx.Response().StatusCode > 300 {
		return
	}

	image := ctx.StateBag()[SkropImage].(*bimg.Image)
	opts := ctx.StateBag()[SkropOptions].(*bimg.Options)

	buf, err := transformImage(image, opts)
	if err != nil {
		log.Error("failed to process image ", err.Error())
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
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
	return o
}
