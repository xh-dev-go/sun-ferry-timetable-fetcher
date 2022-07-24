package cachedResult

import (
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/cachedResult/cachedResult"
	"io"
	"net/http"
)

type Cache[T any] struct {
	Value T
	Key   string ""
}

func (cache *Cache[T]) NotInit() bool {
	return cache.Key == ""
}
func (cache *Cache[T]) HavingCache() bool {
	return !cache.NotInit()
}

type ETag = string
type SimpleHttpCache[T any] struct {
	Cache Cache[T]
	//HttpCast func() (T, ETag, int)
	//Cast     func(response http.Response) (T, error)
}

func (f *SimpleHttpCache[[]byte]) Intercept(r *http.Request, client *http.Client) cachedResult.CacheResult[[]byte] {
	return f.InterceptWithCast(r, client, func(response http.Response) ([]byte, error) {
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		} else {
			return bytes, nil
		}
	})
}
func (f *SimpleHttpCache[T]) InterceptWithCast(r *http.Request, client *http.Client, cast func(response http.Response) (T, error)) cachedResult.CacheResult[T] {
	if f.Cache.HavingCache() {
		r.Header.Set("If-None-Match", f.Cache.Key)
	}
	res, error := client.Do(r)
	if error != nil {
		return cachedResult.Error[T](error)
	}
	if res.StatusCode == 304 {
		return cachedResult.Cached(f.Cache.Value, f.Cache.Key)
	} else if res.StatusCode == 200 {
		t, err := cast(*res)
		if err != nil {
			return cachedResult.Error[T](err)
		}
		f.Cache.Key = res.Header.Get("ETag")
		f.Cache.Value = t
		return cachedResult.Cached(f.Cache.Value, f.Cache.Key)
	} else {
		return cachedResult.NotCached[T](res)
	}
}
