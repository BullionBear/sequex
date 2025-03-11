package redisxcache

import (
	"context"
	"errors"

	"github.com/BullionBear/sequex/pkg/cache/xcache"
	"github.com/redis/go-redis/v9"
)

var _ xcache.XCache = (*RedisXCache)(nil)

// RedisXCache implements XCache using Redis Streams
type RedisXCache struct {
	client *redis.Client
	ctx    context.Context
	topic  string
}

// NewRedisXCache creates a new RedisXCache instance with an injected Redis client
func NewRedisXCache(client *redis.Client, topic string) *RedisXCache {
	if client == nil {
		panic("Redis client cannot be nil")
	}

	return &RedisXCache{
		client: client,
		ctx:    context.Background(),
		topic:  topic,
	}
}

// Set adds a new entry to the Redis Stream
func (r *RedisXCache) Set(key int64, data interface{}) error {
	_, err := r.client.XAdd(r.ctx, &redis.XAddArgs{
		Stream: r.topic,
		Values: map[string]interface{}{
			"key":  key,
			"data": data,
		},
		ID: "*", // Auto-generate ID
	}).Result()
	return err
}

// GetLatest retrieves the latest 'size' elements from the Redis Stream
func (r *RedisXCache) GetLatest(size int64) (interface{}, error) {
	entries, err := r.client.XRevRangeN(r.ctx, r.topic, "+", "-", size).Result()
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, errors.New("no data found")
	}

	// Convert Redis stream entries to a list of maps
	result := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		dataMap := make(map[string]interface{})
		for k, v := range entry.Values {
			dataMap[k] = v
		}
		result[i] = dataMap
	}
	return result, nil
}

// Size returns the number of elements in the stream
func (r *RedisXCache) Size() uint64 {
	size, err := r.client.XLen(r.ctx, r.topic).Result()
	if err != nil {
		return 0
	}
	return uint64(size)
}

// Clear removes all entries from the Redis Stream
func (r *RedisXCache) Clear() error {
	_, err := r.client.Del(r.ctx, r.topic).Result()
	return err
}

// RemoveOldest removes the oldest 'size' elements from the stream
func (r *RedisXCache) RemoveOldest(size int64) error {
	if size <= 0 {
		return errors.New("invalid size: must be greater than zero")
	}

	// Get the first 'size' entries
	entries, err := r.client.XRangeN(r.ctx, r.topic, "-", "+", size).Result()
	if err != nil || len(entries) == 0 {
		return errors.New("no entries to remove")
	}

	// Get the last ID in this range (we will trim everything before it)
	oldestID := entries[len(entries)-1].ID

	// Trim the stream to remove everything older than the last retrieved entry
	_, err = r.client.XTrimMinID(r.ctx, r.topic, oldestID).Result()
	return err
}

// Close closes the Redis client connection
func (r *RedisXCache) Close() error {
	return nil
}
