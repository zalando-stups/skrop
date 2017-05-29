package filters

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"gopkg.in/h2non/bimg.v1"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
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