package main

import (
	"sync"

	"github.com/transip/gotransip/v6/authenticator"
)

type threadSafeTokenCache struct {
	mu    sync.Mutex
	cache authenticator.TokenCache
}

func newThreadSafeTokenCache(cache authenticator.TokenCache) *threadSafeTokenCache {
	return &threadSafeTokenCache{cache: cache}
}

func (t *threadSafeTokenCache) Get() (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.cache.Get()
}

func (t *threadSafeTokenCache) Store(token string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.cache.Store(token)
}
