package filters

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zalando-stups/skrop/cache"
	"github.com/zalando-stups/skrop/parse"
	"github.com/zalando/skipper/filters"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	localFileCacheName     = "localFileCache"
	refreshCacheQueryKey   = "refresh"
	cacheName              = "cache"
	servedByCacheKey       = "served-by-cache"
	pathSeparator          = "/"
	cacheControlValue      = "public,max-age=31536000,immutable"
	servedByCacheVal       = "true"
	lastModifiedDateFormat = "Mon, 2 Jan 2006 15:04:05 MST"
)

type localFileCache struct {
	cacheDir string
	cache    cache.Cache
}

func NewLocalFileCache(cache cache.Cache) filters.Spec {
	return &localFileCache{cache: cache}
}

func (c *localFileCache) Name() string {
	return localFileCacheName
}

func (c *localFileCache) CreateFilter(args []interface{}) (filters.Filter, error) {
	var err error

	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of args: %d, expected 1", len(args))
	}

	c.cacheDir, err = parse.EskipStringArg(args[0])
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *localFileCache) Request(ctx filters.FilterContext) {

	start := time.Now()

	req := determineRquest(ctx)

	// request with no-cache skip cache
	if refreshCache(req.URL.Query()) {
		return
	}

	key := determineKey(c.cacheDir, req.URL)

	log.Debug("Cache key: ", key)

	cached, err := c.cache.Read(key)

	if err != nil {
		log.Debug("Cache was not hit for ", key)
		// the image was not found in the cache or there was an error
		cache.ReportCacheTime(start, cacheName, "read", false)
		return
	}
	cache.ReportCacheTime(start, cacheName, "read", true)

	// serve the image from the cache
	headers := constructHeaders(req.URL.Path, cached)

	log.Debug("Cache was hit for ", key)

	ctx.Serve(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(cached.Content)),
		Header:     headers,
	})
}

func determineRquest(ctx filters.FilterContext) *http.Request {
	req := ctx.OriginalRequest()
	if req == nil {
		req = ctx.Request()
	}
	return req
}

func (c *localFileCache) Response(ctx filters.FilterContext) {

	log.Debug("Response called")

	rsp := ctx.Response()

	// cached images are served by the Request method.
	// Once executed, the Response method is called anyway
	if rsp.Header.Get(servedByCacheKey) == servedByCacheVal {
		log.Debug("Request was served by cache")
		rsp.Header.Del(servedByCacheKey)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		log.Warnf("backend returned with %d status code", rsp.StatusCode)
		return
	}

	img, err := ioutil.ReadAll(rsp.Body)

	log.Debug("Content Length: ", rsp.ContentLength)

	if err != nil {
		log.Warn("failed while writing image to cache with error: ", err.Error())
		return
	}

	content := &cache.CacheContent{
		Content:            img,
		ContentType:        ctx.Response().Header.Get(cache.HContentTypeKey),
		ContentDisposition: ctx.Response().Header.Get(cache.HContentDispositionKey),
	}

	key := determineKey(c.cacheDir, determineRquest(ctx).URL)

	go cacheImage(c.cache, key, content)

	defer rsp.Body.Close()
	rsp.Body = ioutil.NopCloser(bytes.NewReader(img))
}

func cacheImage(handler cache.Cache, key string, content *cache.CacheContent) {
	start := time.Now()
	err := handler.Write(key, content)

	if err != nil {
		cache.ReportCacheTime(start, cacheName, "write", false)
		return
	}
	cache.ReportCacheTime(start, cacheName, "write", true)
}

func refreshCache(queryValues url.Values) bool {
	refresh, ok := queryValues[refreshCacheQueryKey]

	return ok && len(refresh) >= 1 && refresh[0] == "true"
}

/*
The cache key is generally equal to the path:
e.g. GU/15/2O/00/0C/11/GU152O000-C11@3.1.jpg
In case there are query parameters specified a hash of the parameters is also added:
e.g. GU/15/2O/00/0C/11/GU152O000-C11@3.1.123456.jpg
*/
func determineKey(cacheDir string, url *url.URL) string {

	// this checks that the queryparams have been specified.
	par := url.Query()
	if len(par) == 0 {
		return cacheDir + url.Path
	}

	h := fnv.New32a()
	h.Write([]byte(url.RawQuery))
	queryHash := "." + strconv.Itoa(int(h.Sum32()))

	extensionIndex := strings.LastIndex(url.Path, ".")
	path := url.Path[:extensionIndex]
	extension := url.Path[extensionIndex:]

	return cacheDir + path + queryHash + extension
}

func constructHeaders(path string, content *cache.CacheContent) http.Header {

	headers := make(http.Header)

	// Cache Control
	headers.Add(cache.HCacheControlKey, cacheControlValue)
	headers.Add(servedByCacheKey, servedByCacheVal)

	if content == nil {
		return headers
	}

	headers.Add(cache.HContentTypeKey, content.ContentType)

	// Content Disposition
	if content.ContentDisposition == "" {
		contentDisp := "inline;filename=" + path[strings.LastIndex(path, pathSeparator):]
		headers.Add(cache.HContentDispositionKey, contentDisp)
	} else {
		headers.Add(cache.HContentDispositionKey, content.ContentDisposition)
	}

	// adding last-modified
	if !content.LastModified.IsZero() {
		// format: Tue, 15 Nov 1994 12:45:26 GMT as defined: https://tools.ietf.org/html/rfc2616
		stringDate := content.LastModified.Format(lastModifiedDateFormat)
		headers.Add(cache.HLastModified, stringDate)
	}

	return headers
}
