package filters

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"github.com/zalando/skipper/filters/filtertest"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
	"testing"
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

func TestConvertImageType_CanBeMerged_True(t *testing.T) {
	s := convertImageType{}
	opt := &bimg.Options{}
	self := &bimg.Options{Type: 1}

	assert.True(t, s.CanBeMerged(opt, self))
}

func TestConvertImageType_CanBeMerged_False(t *testing.T) {
	s := convertImageType{}
	opt := &bimg.Options{Type: 1}
	self := &bimg.Options{Type: 2}

	assert.False(t, s.CanBeMerged(opt, self))
}

func TestConvertImageType_Merge(t *testing.T) {
	s := convertImageType{}
	self := &bimg.Options{Type: 3}

	opt := s.Merge(&bimg.Options{}, self)

	assert.Equal(t, self.Background, opt.Background)
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

func createFilterContext(t *testing.T, url string) *filtertest.Context {
	buffer, err := bimg.Read(imagefiltertest.PNGImageFile)
	assert.Nil(t, err, "Failed to read sample image")
	imageReader := ioutil.NopCloser(bytes.NewReader(buffer))
	response := &http.Response{Body: imageReader}
	response.Header = make(http.Header)
	response.Header.Add("Content-Length", "100")

	bag := make(map[string]interface{})
	bag[SkropImage] = bimg.NewImage(buffer)
	bag[SkropOptions] = &bimg.Options{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		t.Error(err)
	}

	return &filtertest.Context{FResponse: response, FRequest: req, FStateBag: bag}
}
