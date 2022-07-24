package cachedResult

import (
	"net/http"
)

type CacheKey = string
type Cache[T any] struct {
	value T
	key   CacheKey ""
}

func (cache *Cache[T]) Value() T {
	return cache.value
}
func (cache *Cache[T]) Key() CacheKey {
	return cache.key
}

func (cache *Cache[T]) Match(key string) bool {
	return cache.HavingCache() && cache.Key() == key
}

func (cache *Cache[T]) Update(key CacheKey, value T) *Cache[T] {
	cache.key = key
	cache.value = value
	return cache
}
func (cache *Cache[T]) NotInit() bool {
	return cache.key == ""
}
func (cache *Cache[T]) HavingCache() bool {
	return !cache.NotInit()
}

func (cache *Cache[T]) WithError(response *http.Response, err error) CacheResult[T] {
	return CacheResult[T]{
		cache:    cache,
		response: response,
		error:    err,
	}
}

func (cache *Cache[T]) UpdateCache(key string, value T, response *http.Response) CacheResult[T] {
	cache.Update(key, value)
	return CacheResult[T]{
		cache:    cache,
		response: response,
		error:    nil,
	}
}

func (cache *Cache[T]) NoUpdate(response *http.Response) CacheResult[T] {
	return CacheResult[T]{
		cache:    cache,
		response: response,
		error:    nil,
	}
}

func (cache *Cache[T]) HttpCaching(
	r *http.Request,
	client *http.Client,
	cast func(response http.Response) (T, error)) CacheResult[T] {

	return HttpCache(cache, r, client, cast)
}

type CacheResult[T any] struct {
	cache    *Cache[T]
	response *http.Response
	error    error
}

func (cr *CacheResult[T]) Cache() *Cache[T] {
	return cr.cache
}
func (cr *CacheResult[T]) Response() *http.Response {
	return cr.response
}
func (cr *CacheResult[T]) Error() error {
	return cr.error
}

func (cr *CacheResult[T]) HasError() bool {
	return cr.error != nil
}
func (cr *CacheResult[T]) IsResultCached() bool {
	return cr.response.StatusCode == 304 || cr.response.StatusCode == 200
}

func HttpCache[T any](
	cache *Cache[T],
	r *http.Request,
	client *http.Client,
	cast func(response http.Response) (T, error)) CacheResult[T] {
	if cache.HavingCache() {
		r.Header.Set("If-None-Match", cache.Key())
	}
	res, err := client.Do(r)
	if err != nil {
		return cache.WithError(res, err)
	}
	if res.StatusCode == 304 {
		return cache.NoUpdate(res)
	} else if res.StatusCode == 200 {
		t, err := cast(*res)
		if err != nil {
			return cache.WithError(res, err)
		}
		return cache.UpdateCache(res.Header.Get("ETag"), t, res)
	} else {
		return cache.NoUpdate(res)
	}
}
