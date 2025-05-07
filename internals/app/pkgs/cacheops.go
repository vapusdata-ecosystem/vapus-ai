package apppkgs

import (
	"context"

	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

func CacheFilter(ctx context.Context, stores *VapusStore, action, key string, value ...string) (interface{}, error) {
	switch action {
	case types.LIST:
		return stores.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case types.ADD:
		return stores.Cacher.RedisClient.Client.Do(ctx, "CF.ADD", key, value[0]).Result()
	case types.EXISTS:
		return stores.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case types.COUNT:
		return stores.Cacher.RedisClient.Client.Do(ctx, "CF.CARD", key).Result()
	case types.MADD:
		return stores.Cacher.RedisClient.Client.Do(ctx, "CF.MADD", key, value).Result()
	case types.DEL:
		return stores.Cacher.RedisClient.Client.Do(ctx, "CF.DEL", key, value[0]).Result()
	default:
		return stores.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	}
}
