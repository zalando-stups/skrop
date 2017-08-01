package filters

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
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

func HandleImageResponse(ctx filters.FilterContext, f ImageFilter) {
	image := ctx.StateBag()[SkropImage].(*bimg.Image)

	opt, err := f.CreateOptions(image)
	if err != nil {
		log.Error("Failed to create options ", err.Error())
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		})
		return
	}

	opts := ctx.StateBag()[SkropOptions].(*bimg.Options)

	if f.CanBeMerged(opts, opt) {
		log.Debug("Filter ", f, " merged in ", ctx.StateBag()[SkropOptions])
		ctx.StateBag()[SkropOptions] = f.Merge(opts, opt)
	}

	//transform the image
	buf, err := transformImage(image, opts)
	if err != nil {
		log.Error("Failed to process image ", err.Error())
		ctx.Serve(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		})
		return
	}

	// set opt in the stateBag
	ctx.StateBag()[SkropImage] = bimg.NewImage(buf)
	ctx.StateBag()[SkropOptions] = opt

}

func FinalizeResponse(ctx filters.FilterContext) {
	image := ctx.StateBag()[SkropImage].(*bimg.Image)
	opts := ctx.StateBag()[SkropOptions].(*bimg.Options)

	buf, err := transformImage(image, opts)
	if err != nil {
		log.Error("Failed to process image ", err.Error())
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
