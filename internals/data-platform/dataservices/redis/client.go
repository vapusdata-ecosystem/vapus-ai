package redis

import (
	"context"
	"log"
	"strconv"
	"strings"

	redis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
)

type RedisInterface struct{}

type RedisStore struct {
	Client           *redis.Client
	IsFullRedisStack bool
	logger           zerolog.Logger
}

var (
	defaultJsonPath    = "$"
	defaultJsonPathPtr = "$."
)

// Function to connect with redis server and generate a client
func NewRedisStore(ctx context.Context, conf *RedisOpts, l zerolog.Logger) (*RedisStore, error) {
	addr := conf.URL
	log.Println("Connecting to Redis Server at: ", addr)
	if !strings.Contains(addr, ":"+strconv.Itoa(conf.Port)) {
		addr = addr + ":" + strconv.Itoa(conf.Port)
	}
	log.Println("Connecting to Redis Server at: ", addr)
	log.Println("Port: ", conf.Port)
	cl := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Username: conf.Username,
			Password: conf.Password,
			DB:       0,
		})

	// Check if the client is nil
	if cl == nil {
		return nil, dmerrors.DMError(ErrRedisConnection, nil)
	}

	// Ping the redis server to check the connection
	_, err := cl.Ping(ctx).Result()
	if err != nil {
		return nil, dmerrors.DMError(ErrRedisPing, err)
	}

	return &RedisStore{
		Client:           cl,
		IsFullRedisStack: conf.IsFullRedisStack,
		logger:           l,
	}, nil
}

func (r *RedisStore) Close() {
	r.Client.Close()
}
