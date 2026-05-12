package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/sergeyptv/post_service/platform/cache"
	"sync"
	"time"
)

var (
	ErrPublicKeyTtlExpired = errors.New("public key ttl is expired")
	ErrPublicKeyNotSet     = errors.New("public key is not set")
)

type JwtCache interface {
	Set(publicKey *rsa.PublicKey)
	Get() (*rsa.PublicKey, error)
}

type inMemoryCache struct {
	mu        sync.RWMutex
	publicKey *rsa.PublicKey
	ttl       time.Time
	keyTtl    time.Duration
}

func NewInMemoryCache(c cache.Config) *inMemoryCache {
	memCache := inMemoryCache{
		keyTtl: c.Ttl,
	}

	return &memCache
}

func (c *inMemoryCache) Set(publicKey *rsa.PublicKey) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.publicKey = publicKey
	c.ttl = time.Now().Add(c.keyTtl)
}

func (c *inMemoryCache) Get() (*rsa.PublicKey, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.publicKey == nil {
		return nil, ErrPublicKeyNotSet
	}

	if c.ttl.Before(time.Now()) {
		return c.publicKey, ErrPublicKeyTtlExpired
	}

	return c.publicKey, nil
}
