package util

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisMutex struct {
	client     *redis.Client
	key        string
	timeout    time.Duration
	timewait   time.Duration
	maxRetries int
	locked     bool
}

func NewRedisMutex(cli *redis.Client, key string, timeout int, timewait int, retry int) *RedisMutex {
	mu := new(RedisMutex)

	mu.client = cli
	mu.key = "mutex:" + key
	mu.timeout = time.Duration(timeout) * time.Second
	mu.timewait = time.Duration(timewait) * time.Millisecond
	mu.maxRetries = retry
	mu.locked = false

	return mu
}

func (m *RedisMutex) Lock(ctx context.Context) bool {
	retryTime := 0
	for !m.locked {
		m.locked, _ = m.client.SetNX(ctx, m.key, 1, m.timeout).Result()
		if m.locked {
			break
		}
		retryTime++
		if retryTime > m.maxRetries {
			ttl, _ := m.client.TTL(ctx, m.key).Result()
			if ttl.String() == "-1s" {
				m.client.Del(ctx, m.key)
			}
			break
		}
		time.Sleep(m.timewait)
	}

	return m.locked
}

func (m *RedisMutex) Unlock(ctx context.Context) {
	if !m.locked {
		return
	}
	m.client.Del(ctx, m.key)
}
