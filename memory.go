package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type memoryCache struct {
	sync.Mutex
	cache cacheItems
}

func NewMemoryCache() Cache {
	cacheItems := cacheItems{
		Items: make(map[string]CacheItem),
	}
	return &memoryCache{
		cache: cacheItems,
	}
}

func (mc *memoryCache) Get(key string) CacheItem {
	mc.Lock()
	defer mc.Unlock()

	mc.deleteExpired()

	return mc.cache.Items[key]
}

func (mc *memoryCache) Has(key string) bool {
	mc.Lock()
	defer mc.Unlock()

	mc.deleteExpired()

	_, exists := mc.cache.Items[key]

	return exists
}

func (mc *memoryCache) Delete(key string) error {
	mc.Lock()
	defer mc.Unlock()

	delete(mc.cache.Items, key)

	return nil
}

func (mc *memoryCache) Save(item CacheItem) error {
	mc.Lock()
	defer mc.Unlock()

	mc.cache.Items[item.Key] = CacheItem{
		hit:       true,
		Key:       item.Key,
		Value:     item.Value,
		ExpiresAt: item.ExpiresAt,
	}

	return nil
}

func (mc *memoryCache) Close() error {
	return nil
}

func (mc *memoryCache) Serialize() ([]byte, error) {
	jsonBytes, err := json.Marshal(mc.cache)
	if err != nil {
		return nil, fmt.Errorf("could not encode cache items: %v", err)
	}

	return jsonBytes, nil
}

func (mc *memoryCache) Unserialize(data []byte) error {
	err := json.Unmarshal(data, &mc.cache)
	if err != nil {
		return fmt.Errorf("error decoding file: %v", err)
	}
	for _, item := range mc.cache.Items {
		mc.cache.Items[item.Key] = CacheItem{
			hit:       true,
			Key:       item.Key,
			Value:     item.Value,
			ExpiresAt: item.ExpiresAt,
		}
	}

	return nil
}

// deleteExpired should be run under lock
func (mc *memoryCache) deleteExpired() {
	now := time.Now()

	for key, item := range mc.cache.Items {
		if item.ExpiresAt.Before(now) && !item.ExpiresAt.IsZero() {
			delete(mc.cache.Items, key)
		}
	}
}
