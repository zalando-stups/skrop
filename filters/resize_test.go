package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando-incubator/skrop/filters/imagefiltertest"
)

func TestNewResize(t *testing.T) {
	name := NewResize().Name()
	assert.Equal(t, "resize", name)
}

func TestResize_Name(t *testing.T) {
	r := resize{}
	assert.Equal(t, "resize", r.Name())
}

func TestResize_CreateOptions(t *testing.T) {
	r := resize{width: 800, height: 600}
	options, _ := r.CreateOptions(nil)

	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 600, options.Height)
}

func TestResize_CreateFilter(t *testing.T) {
	imagefiltertest.TestCreate(t, NewResize, []imagefiltertest.CreateTestItem{{
		Msg:  "no args",
		Args: nil,
		Err:  true,
	}, {
		Msg:  "two args",
		Args: []interface{}{800.0, 600.0},
		Err:  false,
	}, {
		Msg:  "more than 2 args",
		Args: []interface{}{800.0, 600.0, "whaaat?"},
		Err:  true,
	}, {
		Msg:  "less than 2 args",
		Args: []interface{}{800.0},
		Err:  true,
	}})
}
