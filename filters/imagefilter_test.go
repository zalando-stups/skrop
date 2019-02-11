package filters

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/zalando/skipper/filters/filtertest"
)

const (
	widthTarget  = 400
	heightTarget = 200
)

type FakeImageFilter bimg.Options

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
	fc.FStateBag[skropOptions] = &optionsTarget
	fc.FStateBag[skropInit] = true

	FinalizeResponse(fc)

	assertCorrectImageSize(fc.Response().Body, t)
}

func TestHandleResponse_InvalidImage(t *testing.T) {
	fc := createDefaultContext(t, "doesNotMatter.com")
	fc.FStateBag[skropOptions] = &optionsTarget
	fc.FStateBag[skropImage] = bimg.NewImage([]byte("invalid image"))
	fc.FStateBag[skropInit] = true

	FinalizeResponse(fc)

	assert.Equal(t, http.StatusInternalServerError, fc.FResponse.StatusCode)
}

func TestHandleImageResponse(t *testing.T) {
	fc := createDefaultContext(t, "doesNotMatter.com")
	imageFilter := FakeImageFilter(optionsTarget)

	err := HandleImageResponse(fc, &imageFilter)

	assert.Nil(t, err, "there should not be any error")
	assert.Equal(t, fc.FStateBag[skropOptions], &optionsTarget)
}

func TestHandleImageResponse_WithResponse304(t *testing.T) {
	fc := createDefaultContext(t, "doesNotMatter.com")
	imageFilter := FakeImageFilter(optionsTarget)

	fc.FResponse.StatusCode = 304

	err := HandleImageResponse(fc, &imageFilter)

	assert.NotNil(t, err, "should not able to process when the backend response is 304")
}

func TestInitResponse_ok(t *testing.T) {
	//given
	emptyBag := make(map[string]interface{})
	ctx := createContext(t, "GET", "zalando.de/image.jpg", imagefiltertest.PortraitImageFile, emptyBag)
	buffer, _ := bimg.Read(imagefiltertest.PortraitImageFile)
	original := bimg.NewImage(buffer)

	//when
	initResponse(ctx)

	//then
	image, ok := ctx.StateBag()[skropImage].(*bimg.Image)
	assert.True(t, ok)
	oriSiz, _ := original.Size()
	imgSiz, _ := image.Size()
	assert.Equal(t, oriSiz, imgSiz)
	_, ok = ctx.StateBag()[skropOptions].(*bimg.Options)
	assert.True(t, ok)
}

func TestInitResponse_ErrorReadingImg(t *testing.T) {
	//given
	emptyBag := make(map[string]interface{})
	ctx := createContext(t, "GET", "zalando.de/image.jpg", "nonExisting.jpg", emptyBag)

	//when
	initResponse(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Response().StatusCode)
}

func createDefaultContext(t *testing.T, url string) *filtertest.Context {
	buffer, _ := bimg.Read(imagefiltertest.PNGImageFile)
	bag := make(map[string]interface{})
	bag[skropImage] = bimg.NewImage(buffer)
	bag[skropOptions] = &bimg.Options{}
	bag[skropInit] = true
	bag[hasMergedFilters] = true
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

func (f *FakeImageFilter) CreateOptions(_ *ImageFilterContext) (*bimg.Options, error) {
	options := bimg.Options(*f)
	return &options, nil
}

func (f *FakeImageFilter) CanBeMerged(other *bimg.Options, self *bimg.Options) bool {
	return (other.Width == 0 && other.Height == 0) ||
		(other.Width == self.Width && other.Height == self.Height)
}

func (f *FakeImageFilter) Merge(other *bimg.Options, self *bimg.Options) *bimg.Options {
	other.Width = self.Width
	other.Height = self.Height
	other.Quality = self.Quality
	other.Background = self.Background
	return other
}
