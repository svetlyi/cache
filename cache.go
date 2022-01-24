package cache

import "time"

type Cache interface {
	Get(key string) CacheItem
	Has(key string) bool
	Delete(key string) error
	Save(item CacheItem) error
	Close() error
	Serialize() ([]byte, error)
	Unserialize(data []byte) error
}

type CacheItem struct {
	Key       string
	Value     string
	hit       bool
	ExpiresAt time.Time
}

func (ci CacheItem) IsHit() bool {
	return ci.hit
}

type cacheItems struct {
	Items map[string]CacheItem
}
