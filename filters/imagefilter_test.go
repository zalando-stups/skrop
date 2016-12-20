package filters

import (
	"bytes"
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
	"testing"
)

const sampleImageFile = "../images/lisbon-tram.jpg"

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

func readSampleImage(t *testing.T) io.ReadCloser {
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

func TestTransformImage(t *testing.T) {
	imageReader := readSampleImage(t)

	var options = &bimg.Options{
		Width: 400}

	r, w := io.Pipe()

	go transformImage(w, imageReader, options)

	resultImage := readResultImage(r, t)

	size, _ := resultImage.Size()

	if size.Width != 400 {
		t.Error("The width is not correct")
	}
}
