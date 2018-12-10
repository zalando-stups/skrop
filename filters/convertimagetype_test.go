package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/h2non/bimg"
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
	fc := createDefaultContext(t, "http://localhost:9090/images/bag.png")
	fc.Request().RequestURI = "/images/bag.png"
	c := convertImageType{imageType: bimg.ImageType(1)}

	c.Response(fc)
	rsp := fc.FResponse

	assert.Equal(t, rsp.Header.Get("Content-Type"), "image/jpeg")
	assert.Equal(t, rsp.Header.Get("Content-Disposition"), "inline;filename=bag.jpeg")
	rsp.Body.Close()
}

func TestConvertImageType_Response_WithOutExtension(t *testing.T) {

	fc := createDefaultContext(t, "http://localhost:9090/images/shoe")
	c := convertImageType{imageType: bimg.ImageType(1)}
	fc.Request().RequestURI = "/images/bag"

	c.Response(fc)
	rsp := fc.FResponse

	assert.Equal(t, rsp.Header.Get("Content-Type"), "image/jpeg")
	assert.Equal(t, rsp.Header.Get("Content-Disposition"), "inline;filename=bag.jpeg")
	rsp.Body.Close()
}

func TestConvertImageType_Response_WithResponse304(t *testing.T) {
	fc := createDefaultContext(t, "http://localhost:9090/images/bag.png")
	fc.Request().RequestURI = "/images/bag.png"
	fc.FResponse.StatusCode = 304
	c := convertImageType{imageType: bimg.ImageType(1)}

	c.Response(fc)
	rsp := fc.FResponse
	defer rsp.Body.Close()

	assert.Equal(t, rsp.Header.Get("Content-Type"), "", "should not add header when the processing of image failed")
	assert.Equal(t, rsp.Header.Get("Content-Disposition"), "", "should not add header when the processing of image failed")
}
