package cachedResult

import (
	"github.com/xh-dev-go/sun-ferry-timetable-fetcher/dataFetch/cachedResult/cachedResult"
	"net/http"
)

type ETag = string
type SingleHttpCache[T any] struct {
	Value T
	ETag  string ""
	Func  func() (T, ETag, int)
	Cast  func(response http.Response) (T, error)
}

func (f *SingleHttpCache[T]) Intercept(r *http.Request, client *http.Client) cachedResult.CacheResult[T] {
	if f.ETag != "" {
		r.Header.Set("If-None-Match", f.ETag)
	}
	res, error := client.Do(r)
	if error != nil {
		return cachedResult.Error[T](error)
	}
	if res.StatusCode == 304 {
		return cachedResult.Cached(f.Value)
	} else if res.StatusCode == 200 {
		t, err := f.Cast(res)
		if err != nil {
			return cachedResult.Error[T](err)
		}
		f.ETag = res.Header.Get("ETag")
		f.Value = t
		return cachedResult.Cached(f.Value)
	} else {
		return cachedResult.NotCached[T](res)
	}
}
