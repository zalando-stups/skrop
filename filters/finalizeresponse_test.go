package filters

import (
	"github.com/h2non/bimg"
	"github.com/stretchr/testify/assert"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"io/ioutil"
	"testing"
)

func TestNewFinalizeResponse(t *testing.T) {
	name := NewFinalizeResponse().Name()
	assert.Equal(t, "finalizeResponse", name)
}

func TestFinalizeResponse_Name(t *testing.T) {
	c := finalizeResponse{}
	assert.Equal(t, "finalizeResponse", c.Name())
}

func TestFinalizeResponse_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewFinalizeResponse, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  false,
	}, {
		Msg:  "one arg",
		Args: []interface{}{256.0},
		Err:  true,
	}})
}

func TestFinalizeResponse_Request(t *testing.T) {
	s := finalizeResponse{}
	ctx := createDefaultContext(t, "image.png")
	ctx1 := createDefaultContext(t, "image.png")

	s.Request(ctx1)

	assert.Equal(t, ctx, ctx1)
}

func TestFinalizeResponse_Response(t *testing.T) {
	//given
	s := finalizeResponse{}
	var bag map[string]interface{}
	bag = make(map[string]interface{})
	opt := &bimg.Options{Width: 100, Height: 200, Force: true}
	bag[skropOptions] = opt
	bag[skropImage] = imagefiltertest.PortraitImage()
	bag[skropInit] = true
	bag[hasMergedFilters] = true

	ctx := createContext(t, "GET", "url", imagefiltertest.PortraitImageFile, bag)

	//when
	s.Response(ctx)

	//then
	rsp := ctx.Response()
	buf, _ := ioutil.ReadAll(rsp.Body)
	size, _ := bimg.NewImage(buf).Size()
	assert.Equal(t, 100, size.Width)
	assert.Equal(t, 200, size.Height)
}
