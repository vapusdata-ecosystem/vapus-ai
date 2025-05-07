package redis

type RedisOpts struct {
	// RedisConfig is the configuration for the Redis
	URL, Username, Password string
	Db                      int
	IsFullRedisStack        bool
	Port                    int
	MaxConnections          int32
}
