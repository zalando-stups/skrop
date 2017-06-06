package filters

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/zalando/skipper/filters/filtertest"
)

func TestNewConvertImageType(t *testing.T) {
	name := NewConvertImageType().Name()
	assert.Equal(t, "convertImageType", name)
}

func TestConvertImageType_Name(t *testing.T) {
	c := &convertImageType{}
	if c.Name() != "convertImageType" {
		t.Error("Convert Image Type name incorrect")
	}
}

func TestConvertImageType_CreateOptions(t *testing.T) {
	c := &convertImageType{imageType: bimg.JPEG}
	options, _ := c.CreateOptions(nil)
	assert.Equal(t, bimg.JPEG, options.Type)
}

func TestConvertImageType_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewConvertImageType, []imagefiltertest.CreateTestItem{{
		"no args",
		nil,
		true,
	}, {
		"one valid arg",
		[]interface{}{"jpeg"},
		false,
	}, {
		"one invalid arg",
		[]interface{}{"jpeg1"},
		true,
	}, {
		"more than one args",
		[]interface{}{"jpeg", "webp"},
		true,
	}})
}

func TestConvertImageType_Response(t *testing.T) {

	fc1 := createFilterContext(t, "http://localhost:9090/images/shoe")
	c1:= convertImageType{imageType: bimg.ImageType(1)}

	c1.Response(fc1)
	rsp1 := fc1.FResponse

	assert.Equal(t, rsp1.Header.Get("Content-Type"), "image/jpeg")
	assert.Equal(t, rsp1.Header.Get("Content-Disposition"), "inline;filename=result.jpeg")
	rsp1.Body.Close()

	fc2 := createFilterContext(t, "http://localhost:9090/images/bag.png")
	fc2.Request().RequestURI = "/images/bag.png"
	c2 := convertImageType{imageType: bimg.ImageType(1)}

	c2.Response(fc2)
	rsp2 := fc2.FResponse

	assert.Equal(t, rsp2.Header.Get("Content-Type"), "image/jpeg")
	assert.Equal(t, rsp2.Header.Get("Content-Disposition"), "inline;filename=bag.jpeg")
	rsp2.Body.Close()
}

func createFilterContext(t *testing.T,  url string) *filtertest.Context {
	buffer, err := bimg.Read(imagefiltertest.PNGImageFile)
	assert.Nil(t, err, "Failed to read sample image")
	imageReader := ioutil.NopCloser(bytes.NewReader(buffer))
	response := &http.Response{Body: imageReader}
	response.Header = make(http.Header)
	response.Header.Add("Content-Length", "100")

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		t.Error(err)
	}

	return &filtertest.Context{FResponse: response, FRequest: req}
}