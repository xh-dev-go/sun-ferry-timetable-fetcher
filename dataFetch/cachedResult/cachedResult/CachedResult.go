package cachedResult

import "net/http"

type CacheResult[T any] struct {
	Cached   bool
	Value    T
	Response *http.Response
	Error    error
}

func (CacheResult *CacheResult[T]) mapTo(mapping func(t T) U) *CacheResult[U] {
	if CacheResult.HasError() {
		return CacheResult
	}
	rs := mapping(CacheResult.Value)

}

func (CacheResult *CacheResult[T]) HasError() bool {
	if CacheResult.Error != nil {
		return true
	} else {
		return false
	}
}

func Cached[T](t T) CacheResult[T] {
	return CacheResult[T]{
		Value:  t,
		Cached: true,
	}
}

func NotCached[T](response *http.Response) CacheResult[T] {
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
