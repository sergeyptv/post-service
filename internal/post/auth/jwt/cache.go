package jwt

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrPublicKeyTtlExpired = errors.New("public key ttl is expired")
	ErrPublicKeyNotSet     = errors.New("public key is not set")
)

type JwtCache interface {
	Set(publicKey string, ttl time.Duration)
	Get() (string, error)
	Stop()
}

type inMemoryCache struct {
	mu        sync.RWMutex
	publicKey string
	ttl       time.Time
	syncChan  chan struct{}
}

func NewInMemoryCache(keyLivingTime time.Duration) *inMemoryCache {
	cache := inMemoryCache{
		syncChan: make(chan struct{}),
	}

	go cache.cleanMemory(keyLivingTime)

	return &cache
}

func (c *inMemoryCache) Set(publicKey string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.publicKey = publicKey
	c.ttl = time.Now().Add(ttl)
}

func (c *inMemoryCache) Get() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.publicKey == "" {
		return "", ErrPublicKeyNotSet
	}

	if c.ttl.Before(time.Now()) {
		return "", ErrPublicKeyTtlExpired
	}

	return c.publicKey, nil
}

func (c *inMemoryCache) Stop() {
	close(c.syncChan)
}

func (c *inMemoryCache) cleanMemory(keyLivingTime time.Duration) {
	ticker := time.NewTicker(keyLivingTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()

			if c.ttl.Before(time.Now()) {
				c.publicKey = ""
			}

			c.mu.Unlock()

		case <-c.syncChan:
			return
		}
	}
}
