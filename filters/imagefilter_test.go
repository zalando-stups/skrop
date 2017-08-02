package filters

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"github.com/zalando/skipper/filters/filtertest"
	"gopkg.in/h2non/bimg.v1"
)

const (
	widthTarget  = 400
	heightTarget = 200
)

var optionsTarget = bimg.Options{
	Width:  widthTarget,
	Height: heightTarget,
}

func createSampleImageReader(t *testing.T) io.ReadCloser {
	buffer, err := bimg.Read(imagefiltertest.LandscapeImageFile)
	assert.Nil(t, err, "Failed to read sample image")
	return ioutil.NopCloser(bytes.NewReader(buffer))
}

func readResultImage(r io.Reader, t *testing.T) *bimg.Image {
	result, err := ioutil.ReadAll(r)
	assert.Nil(t, err, "Error reading the result image")
	return bimg.NewImage(result)
}

func assertCorrectImageSize(r io.Reader, t *testing.T) {
	resultImage := readResultImage(r, t)
	size, _ := resultImage.Size()

	assert.Equal(t, widthTarget, size.Width, "The width is not correct")
	assert.Equal(t, heightTarget, size.Height, "The height is not correct")

}

func TestFinalizeResponse(t *testing.T) {
	fc := createDefaultContext(t, "doesNotMatter.com")
	//fc.FResponse = response
	fc.FStateBag[SkropOptions] = &optionsTarget

	FinalizeResponse(fc)

	assertCorrectImageSize(fc.Response().Body, t)
}

func TestHandleResponseWithInvalidImage(t *testing.T) {
	fc := createDefaultContext(t, "doesNotMatter.com")
	fc.FStateBag[SkropOptions] = &optionsTarget
	fc.FStateBag[SkropImage] = bimg.NewImage([]byte("invalid image"))

	FinalizeResponse(fc)

	assert.Equal(t, http.StatusInternalServerError, fc.FResponse.StatusCode)
}

func TestHandleImageResponse(t *testing.T) {
	imageReader := createSampleImageReader(t)
	response := &http.Response{Body: imageReader}
	response.Header = make(http.Header)
	response.Header.Add("Content-Length", "100")
	fc := createDefaultContext(t, "doesNotMatter.com")
	imageFilter := imagefiltertest.FakeImageFilter(optionsTarget)

	HandleImageResponse(fc, &imageFilter)

	assert.Equal(t, fc.FStateBag[SkropOptions], &optionsTarget)
}

func createDefaultContext(t *testing.T, url string) *filtertest.Context {
	buffer, _ := bimg.Read(imagefiltertest.PNGImageFile)
	bag := make(map[string]interface{})
	bag[SkropImage] = bimg.NewImage(buffer)
	bag[SkropOptions] = &bimg.Options{}
	return createContext(t, "GET", url, imagefiltertest.PNGImageFile, bag)
}

func createContext(t *testing.T, method string, url string, image string, stateBag map[string]interface{}) *filtertest.Context {
	buffer, _ := bimg.Read(image)

	imageReader := ioutil.NopCloser(bytes.NewReader(buffer))
	response := &http.Response{Body: imageReader}
	response.Header = make(http.Header)
	response.Header.Add("Content-Length", "100")

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Error(err)
	}

	return &filtertest.Context{FResponse: response, FRequest: req, FStateBag: stateBag}
}