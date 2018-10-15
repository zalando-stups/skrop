package cache

import (
	"fmt"
	"github.com/zalando/skipper/metrics"
	"time"
)

const (
	HCacheControlKey       = "Cache-Control"
	HContentDispositionKey = "Content-Disposition"
	HContentTypeKey        = "Content-Type"
	HCacheKey              = "X-Cache"
	HLastModified          = "Last-Modified"

	successStr = "success"
	failureStr = "failure"
)

type CacheContent struct {
	Content            []byte
	ContentType        string
	ContentDisposition string
	LastModified       time.Time
}

type Cache interface {
	Read(cacheKey string) (*CacheContent, error)
	Write(cacheKey string, content *CacheContent) error
}

func ReportCacheTime(start time.Time, cacheName string, op string, success bool) {
	rsp := successStr
	if !success {
		rsp = failureStr
	}
	metrics.Default.MeasureSince(fmt.Sprintf("cache.%s.%s.%s", cacheName, op, rsp), start)
}
