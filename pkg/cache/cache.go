package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/go-redis/redis/v8"
	"github.com/aydinnyunus/wallet-tracker/pkg/logger"
)

// Cache interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(host string, port int, password string, db int, prefix string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return &RedisCache{
		client: client,
		prefix: prefix,
	}, nil
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	key = c.prefixKey(key)
	
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}
	
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("failed to unmarshal cached data: %w", err)
	}
	
	logger.Debugf("Cache hit for key: %s", key)
	return nil
}

// Set stores a value in cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	key = c.prefixKey(key)
	
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}
	
	logger.Debugf("Cached key: %s with TTL: %v", key, ttl)
	return nil
}

// Delete removes a value from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	key = c.prefixKey(key)
	
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	
	logger.Debugf("Deleted key from cache: %s", key)
	return nil
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	key = c.prefixKey(key)
	
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	
	return exists > 0, nil
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// prefixKey adds prefix to cache key
func (c *RedisCache) prefixKey(key string) string {
	if c.prefix != "" {
		return fmt.Sprintf("%s:%s", c.prefix, key)
	}
	return key
}

// MemoryCache implements in-memory cache (for testing or when Redis is not available)
type MemoryCache struct {
	data map[string]*memoryItem
}

type memoryItem struct {
	value     []byte
	expiresAt time.Time
}

// NewMemoryCache creates a new in-memory cache instance
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]*memoryItem),
	}
}

// Get retrieves a value from memory cache
func (m *MemoryCache) Get(ctx context.Context, key string, value interface{}) error {
	item, exists := m.data[key]
	if !exists {
		return ErrCacheMiss
	}
	
	if time.Now().After(item.expiresAt) {
		delete(m.data, key)
		return ErrCacheMiss
	}
	
	return json.Unmarshal(item.value, value)
}

// Set stores a value in memory cache
func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	m.data[key] = &memoryItem{
		value:     data,
		expiresAt: time.Now().Add(ttl),
	}
	
	return nil
}

// Delete removes a value from memory cache
func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// Exists checks if a key exists in memory cache
func (m *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.data[key]
	return exists, nil
}

// Error types
var (
	ErrCacheMiss = fmt.Errorf("cache miss")
)

// CacheKeyBuilder helps build consistent cache keys
type CacheKeyBuilder struct {
	namespace string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder(namespace string) *CacheKeyBuilder {
	return &CacheKeyBuilder{
		namespace: namespace,
	}
}

// WalletKey builds a cache key for wallet data
func (b *CacheKeyBuilder) WalletKey(network, address string) string {
	return fmt.Sprintf("%s:wallet:%s:%s", b.namespace, network, address)
}

// TransactionKey builds a cache key for transaction data
func (b *CacheKeyBuilder) TransactionKey(hash string) string {
	return fmt.Sprintf("%s:tx:%s", b.namespace, hash)
}

// ExchangeKey builds a cache key for exchange data
func (b *CacheKeyBuilder) ExchangeKey(exchange string) string {
	return fmt.Sprintf("%s:exchange:%s", b.namespace, exchange)
}
