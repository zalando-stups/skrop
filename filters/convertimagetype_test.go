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

func TestConvertImageType_Response_WithExtension(t *testing.T) {
	fc := createFilterContext(t, "http://localhost:9090/images/bag.png")
	fc.Request().RequestURI = "/images/bag.png"
	c := convertImageType{imageType: bimg.ImageType(1)}

	c.Response(fc)
	rsp := fc.FResponse

	assert.Equal(t, rsp.Header.Get("Content-Type"), "image/jpeg")
	assert.Equal(t, rsp.Header.Get("Content-Disposition"), "inline;filename=bag.jpeg")
	rsp.Body.Close()
}

func TestConvertImageType_Response_WithOutExtension(t *testing.T) {

	fc := createFilterContext(t, "http://localhost:9090/images/shoe")
	c := convertImageType{imageType: bimg.ImageType(1)}
	fc.Request().RequestURI = "/images/bag"

	c.Response(fc)
	rsp := fc.FResponse

	assert.Equal(t, rsp.Header.Get("Content-Type"), "image/jpeg")
	assert.Equal(t, rsp.Header.Get("Content-Disposition"), "inline;filename=bag.jpeg")
	rsp.Body.Close()
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