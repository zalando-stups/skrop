package dataclient

import (
	"github.com/stretchr/testify/assert"
	"github.com/zalando/skipper/eskip"
	"testing"
)

func TestNewSkropDataClient(t *testing.T) {
	s := NewSkropDataClient("path/to/eskip.file")
	assert.NotNil(t, s)
}

func TestSkropDataClient_LoadAll_NonShunt(t *testing.T) {
	//given
	s := NewSkropDataClient("eskip_test.eskip")

	//when
	routes, _ := s.LoadAll()

	//then
	for _, route := range routes {
		if route.BackendType != eskip.ShuntBackend {
			assert.Equal(t, "finalizeResponse", route.Filters[0].Name)
			assert.Equal(t, "setupResponse", route.Filters[len(route.Filters)-1].Name)
		}
	}
}

func TestSkropDataClient_LoadAll_Shunt(t *testing.T) {
	//given
	s := NewSkropDataClient("eskip_test.eskip")

	//when
	routes, _ := s.LoadAll()

	//then
	for _, route := range routes {
		if route.BackendType == eskip.ShuntBackend {
			assert.Equal(t, 1, len(route.Filters))
		}
	}
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
