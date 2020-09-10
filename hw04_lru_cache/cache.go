package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*Item
	mu       sync.Mutex
}

// NewCache creates new Cache instance.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Item, capacity),
	}
}

// Get returns {value, true} whether element with key stored in cache and {nil, false} pair otherwise.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)

		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

// Set stores value in cache with key, returns whether element was already in cache.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value.(*cacheItem).value = value

		return true
	}

	if c.capacity == c.queue.Len() {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		delete(c.items, lastItem.Value.(*cacheItem).key)
	}

	c.queue.PushFront(&cacheItem{key, value})
	c.items[key] = c.queue.Front()

	return false
}

// Clear removes all elements from cache.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*Item, c.capacity)
}
