package cachedResult

import "net/http"

type ETag = string
type CacheResult[T any] struct {
	ETag     string
	Cached   bool
	Value    T
	Response *http.Response
	Error    error
}

func (CacheResult *CacheResult[T]) HasError() bool {
	if CacheResult.Error != nil {
		return true
	} else {
		return false
	}
}

func Cached[T any](t T, etag string) CacheResult[T] {
	return CacheResult[T]{
		ETag:   etag,
		Value:  t,
		Cached: true,
	}
}

func NotCached[T any](response *http.Response) CacheResult[T] {
	return CacheResult[T]{
		Response: response,
		Cached:   true,
	}
}

func Error[T any](err error) CacheResult[T] {
	return CacheResult[T]{
		Error:  err,
		Cached: true,
	}
}
