package cache

import (
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration *time.Time
}

type LocalCache struct {
	stop chan struct{}

	wg sync.WaitGroup
	mu sync.RWMutex

	data map[string]*cacheItem
}

func NewLocalCache(invalidationInterval time.Duration) *LocalCache {
	c := &LocalCache{
		stop: make(chan struct{}),
		data: make(map[string]*cacheItem),
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		t := time.NewTicker(invalidationInterval)
		defer t.Stop()

		for {
			select {
			case <-c.stop:
				return
			case <-t.C:
				c.mu.Lock()
				for k, v := range c.data {
					if v.expiration != nil && v.expiration.Before(time.Now()) {
						delete(c.data, k)
					}
				}
				c.mu.Unlock()
			}
		}
	}()

	return c
}

func (c *LocalCache) Stop() {
	close(c.stop)
	c.wg.Wait()
}

func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	return item.value, true
}

func (c *LocalCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Now().Add(expiration)

	c.data[key] = &cacheItem{
		value:      value,
		expiration: &expiresAt,
	}
}

func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
