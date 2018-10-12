package filters

import (
	"github.com/zalando-stups/skrop/cache"
	"github.com/zalando-stups/skrop/filters/imagefiltertest"
	"github.com/zalando/skipper/filters"
	"github.com/zalando/skipper/filters/filtertest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type metricsHandler struct{}

var emptyMeta map[string]*string = make(map[string]*string)

func (f *metricsHandler) MeasureSince(key string, start time.Time) {}

func (f *metricsHandler) IncCounter(key string) {}

func TestLocalFileCache_NewLocalFileCache(t *testing.T) {
	cache := cache.NewFileSystemCache()
	cacheFilter := NewLocalFileCache(cache)

	assert.NotNil(t, cacheFilter)
}

func TestLocalFileCache_Name(t *testing.T) {
	cache := &localFileCache{}
	name := cache.Name()

	assert.Equal(t, "localFileCache", name)
}

func TestLocalFileCache_CreateFilter(t *testing.T) {
	cache := cache.NewFileSystemCache()
	imagefiltertest.TestCreate(t, func() filters.Spec { return NewLocalFileCache(cache) },
		[]imagefiltertest.CreateTestItem{{
			"no args",
			nil,
			true,
		}, {
			"one arg",
			[]interface{}{"/tmp"},
			false,
		}, {
			"two args",
			[]interface{}{"/tmp", "hello"},
			true,
		}})
}

func TestLocalFileCache_determineKey(t *testing.T) {
	resourcePath := "/LT/12/1A/01/39/51/LT121A013-951@25.jpg"
	requestPath := "http://www.example.org" + resourcePath
	queryParams := "?a=3"
	urlNoQuery := httptest.NewRequest("GET", requestPath, nil).URL
	urlQuery := httptest.NewRequest("GET", requestPath+queryParams, nil).URL
	cacheDir := "/images"

	extensionIndex := strings.LastIndex(resourcePath, ".")
	path := resourcePath[:extensionIndex]
	ext := resourcePath[extensionIndex:]

	keyQuery := determineKey(cacheDir, urlQuery)
	keyNoQuery := determineKey(cacheDir, urlNoQuery)

	assert.True(t, strings.HasPrefix(keyQuery, cacheDir+path))
	assert.True(t, strings.HasSuffix(keyQuery, ext))
	assert.False(t, resourcePath == keyQuery)
	assert.False(t, keyQuery == keyNoQuery)
}

func TestLocalFileCache_Request_NoCache(t *testing.T) {
	reqPath := "/LT/12/1A/01/39/51/LT121A013-951@25.jpg"

	cache := cache.NewFileSystemCache()
	cacheFilter := NewLocalFileCache(cache)
	f, _ := cacheFilter.CreateFilter([]interface{}{"/images"})
	req, _ := http.NewRequest("GET", "http://www.example.org"+reqPath+"?refresh=true", nil)

	ctx := &filtertest.Context{FRequest: req, FMetrics: &metricsHandler{}}
	f.Request(ctx)

	assert.False(t, ctx.Served())
}

func TestLocalFileCache_Request_NotInCache(t *testing.T) {
	reqPath := "/LT/12/1A/01/39/51/LT121A013-951@25.jpg"

	cache := cache.NewFileSystemCache()
	cacheFilter := NewLocalFileCache(cache)
	f, _ := cacheFilter.CreateFilter([]interface{}{"/images"})
	req, _ := http.NewRequest("GET", "http://www.example.org"+reqPath, nil)

	ctx := &filtertest.Context{FRequest: req, FMetrics: &metricsHandler{}}
	f.Request(ctx)

	assert.False(t, ctx.Served())
}

func TestLocalFileCache_Request_InCache(t *testing.T) {
	reqPath := "/lisbon-tram.jpg"

	cache := cache.NewFileSystemCache()
	cacheFilter := NewLocalFileCache(cache)

	f, _ := cacheFilter.CreateFilter([]interface{}{"../images"})
	req, _ := http.NewRequest("GET", "http://www.example.org"+reqPath, nil)

	ctx := &filtertest.Context{FRequest: req, FMetrics: &metricsHandler{}}
	f.Request(ctx)

	assert.True(t, ctx.Served())
	_, err := ioutil.ReadAll(ctx.Response().Body)
	assert.Nil(t, err)
}

func TestLocalFileCache_refresh(t *testing.T) {
	query := make(url.Values)
	assert.False(t, refreshCache(query))

	query = make(url.Values)
	query[refreshCacheQueryKey] = []string{""}
	assert.False(t, refreshCache(query))

	query = make(url.Values)
	query[refreshCacheQueryKey] = []string{"false"}
	assert.False(t, refreshCache(query))

	query = make(url.Values)
	query[refreshCacheQueryKey] = []string{"true"}
	assert.True(t, refreshCache(query))
}

func TestLocalFileCache_Response(t *testing.T) {
	reqPath := "/lisbon-tram.jpg"

	cache := cache.NewFileSystemCache()
	cacheFilter := NewLocalFileCache(cache)

	f, _ := cacheFilter.CreateFilter([]interface{}{"../tmpimages"})
	req, _ := http.NewRequest("GET", "http://www.example.org"+reqPath, nil)
	recorder := httptest.NewRecorder()

	ctx := &filtertest.Context{FRequest: req, FMetrics: &metricsHandler{}, FResponse: recorder.Result()}
	f.Response(ctx)

	assert.Equal(t, "200 OK", ctx.Response().Status)

	os.RemoveAll("../tmpimages")
}
