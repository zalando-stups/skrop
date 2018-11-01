package dataclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSkropDataClient(t *testing.T) {
	s := NewSkropDataClient("path/to/eskip.file")
	assert.NotNil(t, s)
}

func TestSkropDataClient_LoadAll_FileErr(t *testing.T) {
	//given
	s := NewSkropDataClient("nonexistent.eskip")

	//when
	routes, err := s.LoadAll()

	//then
	assert.Nil(t, routes)
	assert.NotNil(t, err)
}

func TestSkropDataClient_LoadUpdate(t *testing.T) {
	//given
	s := NewSkropDataClient("eskip_test.eskip")

	//when
	a, b, err := s.LoadUpdate()

	//then
	assert.Nil(t, a)
	assert.Nil(t, b)
	assert.Nil(t, err)
}
