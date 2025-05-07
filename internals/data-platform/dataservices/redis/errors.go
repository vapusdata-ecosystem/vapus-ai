package redis

import "errors"

var (
	// Error constants for Redis operations
	ErrRedisConnection   = errors.New("error while connecting to Redis")
	ErrRedisPing         = errors.New("error while pinging Redis")
	ErrRedisWrite        = errors.New("error while writing to Redis")
	ErrRedisDelete       = errors.New("error while deleting from Redis")
	ErrRedisRead         = errors.New("error while reading from Redis")
	ErrRedisKeyNotExists = errors.New("error while checking if key exists in Redis")
)
