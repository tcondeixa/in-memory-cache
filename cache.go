package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Cache[T any] struct {
	cache  map[string]cacheEntry[T]
	lock   sync.RWMutex
	expire time.Duration
}

type cacheEntry[T any] struct {
	Value  T
	update time.Time `json:"-" yaml:"-"`
}

// Create a New Cache
func NewCache[T any](cacheExpiry time.Duration) *Cache[T] {

	return &Cache[T]{
		cache:  map[string]cacheEntry[T]{},
		lock:   sync.RWMutex{},
		expire: cacheExpiry,
	}
}

// Load from File
func (c *Cache[T]) Load(filePath string, fileFormat FileFormat) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	byteValue, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	switch fileFormat {
	case Json:
		err = json.Unmarshal(byteValue, &c.cache)
		if err != nil {
			return err
		}
	case Yaml:
		err = yaml.Unmarshal(byteValue, &c.cache)
		if err != nil {
			return err
		}
	default:
		return &UnknownFileFormatError{}
	}

	currentTime := time.Now()
	for _, v := range c.cache {
		v.update = currentTime
	}

	return nil
}

// Save to File
func (c *Cache[T]) Save(filePath string, fileFormat FileFormat) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	var bytesValue []byte
	switch fileFormat {
	case Json:
		bytesValue, err = json.MarshalIndent(c.cache, "", "  ")
		if err != nil {
			return err
		}
	case Yaml:
		bytesValue, err = yaml.Marshal(c.cache)
		if err != nil {
			return err
		}
	default:
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err := enc.Encode(c.cache)
		if err != nil {
			return err
		}

		bytesValue = buffer.Bytes()
	}

	err = os.WriteFile(absPath, bytesValue, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Get Value for a Key
func (c *Cache[T]) Get(key string) (*T, error) {
	c.lock.RLock()

	entry, ok := c.cache[key]
	if ok {
		if time.Since(entry.update) >= c.expire {
			delete(c.cache, key)
			c.lock.RUnlock()
			return nil, &OutdatedCacheEntryError{}
		}

		c.lock.RUnlock()

		return &entry.Value, nil

	}

	c.lock.RUnlock()
	return nil, &InvalidCacheKeyError{}

}

// Insert of Update (if exists) Value for Key
func (c *Cache[T]) Upsert(key string, value T) {
	currentTime := time.Now()
	c.lock.Lock()

	c.cache[key] = cacheEntry[T]{
		Value:  value,
		update: currentTime,
	}

	c.lock.Unlock()
}

// Insert Value in Cache, return error if already exists
func (c *Cache[T]) Insert(key string, value T) error {
	currentTime := time.Now()
	c.lock.Lock()

	if _, ok := c.cache[key]; ok {
		c.lock.Unlock()
		return &KeyAlreadyExistsError{}
	}

	c.cache[key] = cacheEntry[T]{
		Value:  value,
		update: currentTime,
	}

	c.lock.Unlock()

	return nil
}

// Bulk Upsert Values in Cache passing all objects at once
func (c *Cache[T]) BulkUpsert(values map[string]T) {

	c.lock.Lock()
	currentTime := time.Now()

	for key, value := range values {
		c.cache[key] = cacheEntry[T]{
			Value:  value,
			update: currentTime,
		}
	}

	c.lock.Unlock()
}

// Bulk Insert Values in Cache passing all objects at once
func (c *Cache[T]) BulkInsert(values map[string]T) (map[string]T, error) {

	existingObjects := map[string]T{}

	c.lock.Lock()
	currentTime := time.Now()

	for key, value := range values {
		if _, ok := c.cache[key]; ok {
			existingObjects[key] = value
			continue
		}

		c.cache[key] = cacheEntry[T]{
			Value:  value,
			update: currentTime,
		}
	}

	c.lock.Unlock()
	if len(existingObjects) != 0 {
		return existingObjects, &KeyAlreadyExistsError{}
	}

	return nil, nil
}

// Delete All Expired keys and values
func (c *Cache[T]) Expire() {
	c.lock.Lock()

	for key, value := range c.cache {
		if time.Since(value.update) >= c.expire {
			delete(c.cache, key)
		}
	}

	c.lock.Unlock()
}

// Flush Cache
func (c *Cache[T]) Flush() {
	c.lock.Lock()
	c.cache = map[string]cacheEntry[T]{}
	c.lock.Unlock()
}

// Delete Key and Object From Cache
func (c *Cache[T]) Delete(key string) {
	c.lock.Lock()
	delete(c.cache, key)
	c.lock.Unlock()
}

// Size of the Cache
func (c *Cache[T]) Size() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.cache)
}
