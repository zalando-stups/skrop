package filters

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
	"gopkg.in/h2non/bimg.v1"
	"net/http"
	"testing"
)

func TestNewSetupResponse(t *testing.T) {
	name := NewSetupResponse().Name()
	assert.Equal(t, "setupResponse", name)
}

func TestSetupResponse_Name(t *testing.T) {
	c := setupResponse{}
	assert.Equal(t, "setupResponse", c.Name())
}

func TestSetupResponse_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewSetupResponse, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  false,
	}, {
		Msg:  "one arg",
		Args: []interface{}{256.0},
		Err:  true,
	}})
}

func TestSetupResponse_Request(t *testing.T) {
	s := setupResponse{}
	ctx := createDefaultContext(t, "image.png")
	ctx1 := createDefaultContext(t, "image.png")

	s.Request(ctx1)

	assert.Equal(t, ctx, ctx1)
}

func TestSetupResponse_Response(t *testing.T) {
	//given
	s := setupResponse{}
	emptyBag := make(map[string]interface{})
	ctx := createContext(t, "GET", "zalando.de/image.jpg", imagefiltertest.PortraitImageFile, emptyBag)
	buffer, _ := bimg.Read(imagefiltertest.PortraitImageFile)
	original := bimg.NewImage(buffer)

	//when
	s.Response(ctx)

	//then
	image, ok := ctx.StateBag()[SkropImage].(*bimg.Image)
	assert.True(t, ok)
	oriSiz, _ := original.Size()
	imgSiz, _ := image.Size()
	assert.Equal(t, oriSiz, imgSiz)
	_, ok = ctx.StateBag()[SkropOptions].(*bimg.Options)
	assert.True(t, ok)
}

func TestSetupResponse_Response_ErrorReadingImg(t *testing.T) {
	//given
	s := setupResponse{}
	emptyBag := make(map[string]interface{})
	ctx := createContext(t, "GET", "zalando.de/image.jpg", "nonExisting.jpg", emptyBag)

	//when
	s.Response(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Response().StatusCode)
}
