package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries  map[string]CacheEntry
	Mu       sync.Mutex
	Interval time.Duration
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		make(map[string]CacheEntry),
		sync.Mutex{},
		interval,
	}
	go cache.ReapLoop(interval)

	return cache
}

func (C *Cache) Add(key string, val []byte) {
	C.Mu.Lock()
	defer C.Mu.Unlock()
	C.Entries[key] = CacheEntry{
		time.Now(),
		val,
	}
}

func (C *Cache) Get(getKey string) ([]byte, bool) {
	C.Mu.Lock()
	defer C.Mu.Unlock()
	for entryKey, entry := range C.Entries {
		if entryKey == getKey {
			return entry.Val, true
		}
	}
	return nil, false
}

func (C *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		C.Mu.Lock()
		for key, entry := range C.Entries {
			if time.Since(entry.CreatedAt) >= interval {
				delete(C.Entries, key)
			}
		}
		C.Mu.Unlock()
	}
}
