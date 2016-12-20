package filters

import (
	"bytes"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"github.com/zalando/skipper/filters/filtertest"
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	sampleImageFile = "../images/lisbon-tram.jpg"
	widthTarget     = 400
	heightTarget    = 200
)

var optionsTarget = bimg.Options{
	Width:  widthTarget,
	Height: heightTarget,
}

func createSampleImageReader(t *testing.T) io.ReadCloser {
	buffer, err := bimg.Read(sampleImageFile)

	if err != nil {
		t.Error("Failed to read sample image")
	}

	return ioutil.NopCloser(bytes.NewReader(buffer))
}

func readResultImage(r io.Reader, t *testing.T) *bimg.Image {
	result, err := ioutil.ReadAll(r)

	if err != nil {
		t.Error("Error reading the result image")
	}

	return bimg.NewImage(result)
}

func assertCorrectImageSize(r io.Reader, t *testing.T) {
	resultImage := readResultImage(r, t)

	size, _ := resultImage.Size()

	if size.Width != widthTarget {
		t.Error("The width is not correct")
	}

	if size.Height != heightTarget {
		t.Error("The height is not correct ", size.Height)
	}
}

func TestTransformImage(t *testing.T) {
	buffer, _ := bimg.Read(sampleImageFile)
	image := bimg.NewImage(buffer)

	r, w := io.Pipe()

	go transformImage(w, image, &optionsTarget)

	assertCorrectImageSize(r, t)
}

func TestHandleResponse(t *testing.T) {
	imageReader := createSampleImageReader(t)
	response := &http.Response{Body: imageReader}
	response.Header = make(http.Header)
	response.Header.Add("Content-Length", "100")
	fc := &filtertest.Context{FResponse: response}
	imageFilter := imagefiltertest.FakeImageFilter(optionsTarget)

	handleResponse(fc, &imageFilter)

	assertCorrectImageSize(fc.Response().Body, t)
}

func TestParseEskipIntArgSuccess(t *testing.T) {
	result, _ := parseEskipIntArg(1.0)

	if result != 1 {
		t.Error("Result incorrect")
	}
}

func TestParseEskipIntArgFailure(t *testing.T) {
	_, err := parseEskipIntArg(1.2)

	if err == nil {
		t.Error("There should be an error")
	}
}
