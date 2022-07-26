package cachedResult

import (
	"github.com/eko/gocache/v3/cache"
	"github.com/eko/gocache/v3/store"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

func CacheC[T any]() *cache.Cache[Cache[T]] {
	cc := gocache.New(5*time.Minute, 10*time.Minute)
	s := store.NewGoCache(cc)
	return cache.New[Cache[T]](s)
}
