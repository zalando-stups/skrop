package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemory(t *testing.T) {
	cache := NewInMemoryCache()
	assert.NotNil(t, cache)
}

func TestInMemory_Write(t *testing.T) {
	cache := NewInMemoryCache()
	cont := &CacheContent{
		Content: []byte("fakeImage"),
	}
	err := cache.Write("testKey", cont)

	assert.Nil(t, err)
}

func TestInMemory_Read(t *testing.T) {
	cache := NewInMemoryCache()
	cont := &CacheContent{
		Content: []byte("fakeImage"),
	}
	err := cache.Write("testKey", cont)

	content, err := cache.Read("testKey")
	assert.Nil(t, err)
	assert.Equal(t, []byte("fakeImage"), content.Content)
}
