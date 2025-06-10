package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt		time.Time
	val				[]byte
}

type Cache struct {
	entries 	map[string]cacheEntry
	mu 			sync.Mutex
}


func NewCache(interval time.Duration) Cache {
	var cac =  Cache{
		entries: map[string]cacheEntry{},
	}
	cac.reapLoop(interval)
	return cac
}


func (cache *Cache) Add(key *string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.entries[*key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	};
}


func (cache *Cache) Get(key *string) (val []byte, ok bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.entries[*key]
	return entry.val, ok
}


func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for ;; {
			t := <-ticker.C;
			for k, v := range cache.entries {
				if (v.createdAt.Add(interval).Before(t)) {
					delete(cache.entries, k)
				}
			}
		}
	}()

}
