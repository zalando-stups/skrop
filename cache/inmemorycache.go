package cache

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type inMemoryCache struct {
	data map[string]*CacheContent
}

func NewInMemoryCache() Cache {
	d := make(map[string]*CacheContent)
	return &inMemoryCache{
		data: d,
	}
}

func (f *inMemoryCache) Write(cacheKey string, content *CacheContent) error {
	log.Infof("writing key:%q\nvalue:%+v", cacheKey, content)
	f.data[cacheKey] = content
	return nil
}

func (f *inMemoryCache) Read(cacheKey string) (*CacheContent, error) {
	val, ok := f.data[cacheKey]
	log.Infof("reading key:%q\nvalue:%+v", cacheKey, val)
	if !ok {
		return nil, errors.New("error")
	}

	return val, nil
}
