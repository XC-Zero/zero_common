package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"sync"
)

const (
	Max128MB = 2 << 9 << 10 << 7
	Max256MB = 2 << 9 << 10 << 8
)

var cache *fastcache.Cache
var once sync.Once

func getCache() *fastcache.Cache {
	once.Do(func() {
		cache = fastcache.New(Max128MB)
	})
	return cache
}

func GetCache(key string) []byte {

	return getCache().Get(nil, []byte(key))
}

func SetCache(key string, val []byte) {
	getCache().Set([]byte(key), val)
}

func DelCache(key string) {
	getCache().Del([]byte(key))
}
