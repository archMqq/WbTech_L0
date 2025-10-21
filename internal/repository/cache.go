package repository

import (
	"L0/internal/database/models"
	"sync"
	"time"
)

type InMemoryCache struct {
	items             map[string]*CacheObject
	cleanInterval     time.Duration
	defaultExpiration time.Duration
	mutex             *sync.RWMutex
}

type CacheObject struct {
	value      *models.Order
	expiration int64
}

func (ch *InMemoryCache) StartCollector() {
	go ch.collect()
}

func (ch *InMemoryCache) collect() {
	for {
		<-time.After(ch.cleanInterval)

		if ch.items == nil {
			continue
		}

		if keys := ch.getExpired(); len(keys) > 0 {
			ch.clean(keys)
		}
	}
}

func (ch *InMemoryCache) getExpired() (keys []string) {
	ch.mutex.RLock()
	defer ch.mutex.RUnlock()

	for k, i := range ch.items {
		if time.Now().UnixNano() > i.expiration && i.expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (ch *InMemoryCache) clean(keys []string) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	for _, k := range keys {
		delete(ch.items, k)
	}
}

func (c *InMemoryCache) Set(order *models.Order, id string) {
	chObj := &CacheObject{
		value:      order,
		expiration: time.Now().Add(c.defaultExpiration).UnixNano(),
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[id] = chObj
}

func (c *InMemoryCache) Get(id string) *models.Order {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch, ok := c.items[id]
	if !ok {
		return nil
	}

	return ch.value
}
