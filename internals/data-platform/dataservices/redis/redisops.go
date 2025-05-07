package redis

import (
	"context"
	"time"

	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

func (r *RedisStore) getPath(path string) string {
	if path == types.EMPTYSTR {
		return defaultJsonPath
	}
	return defaultJsonPathPtr + path
}

func (r *RedisStore) WrtiteData(ctx context.Context, key, path string, value interface{}) (interface{}, error) {

	if r.IsFullRedisStack {
		return r.writeJsonData(ctx, key, path, value)
	}
	return r.WriteKV(ctx, key, dmutils.AStructToAString(value))
}

func (r *RedisStore) ReadData(ctx context.Context, key, path string) (interface{}, error) {
	if r.IsFullRedisStack {
		return r.readJsonData(ctx, key, path)
	}
	return r.ReadKV(ctx, key)
}

func (r *RedisStore) DeleteData(ctx context.Context, key, path string) (interface{}, error) {
	if r.IsFullRedisStack {
		return r.deleteJsonData(ctx, key, path)
	}
	return r.DeleteKey(ctx, key)
}

func (r *RedisStore) KeyExists(ctx context.Context, key string) (bool, error) {
	return r.keyExists(ctx, key)
}

func (r *RedisStore) WriteKV(ctx context.Context, key, value string) (interface{}, error) {
	res, err := r.Client.Do(ctx, "SET", key, value).Result()
	if err != nil {
		return types.EMPTYSTR, dmerrors.DMError(ErrRedisWrite, err)
	}
	return res, nil
}

func (r *RedisStore) keyExists(ctx context.Context, key string) (bool, error) {
	res, err := r.Client.Do(ctx, "Exists", key).Bool()
	if !res {
		return false, dmerrors.DMError(ErrRedisKeyNotExists, err)
	}
	return true, nil
}

func (r *RedisStore) ReadKV(ctx context.Context, key string) (interface{}, error) {
	res, err := r.Client.Do(ctx, "GET", key).Result()
	if err != nil {
		return types.EMPTYSTR, dmerrors.DMError(ErrRedisRead, err)
	}
	return res, nil
}

func (r *RedisStore) DeleteKey(ctx context.Context, key string) (interface{}, error) {
	res, err := r.Client.Do(ctx, "DEL", key).Result()
	if err != nil {
		return types.EMPTYSTR, dmerrors.DMError(ErrRedisDelete, err)
	}
	return res, nil
}

func (r *RedisStore) writeJsonData(ctx context.Context, key, path string, value interface{}) (interface{}, error) {
	path = r.getPath(path)
	result := r.Client.JSONSet(ctx, key, path, value)
	if result.Err() != nil {
		return types.EMPTYSTR, dmerrors.DMError(ErrRedisWrite, result.Err())
	}
	return result.Val(), nil
}

func (r *RedisStore) readJsonData(ctx context.Context, key, path string) (interface{}, error) {
	path = r.getPath(path)
	result := r.Client.JSONGet(ctx, key, path)
	if result.Err() != nil {
		return types.EMPTYSTR, dmerrors.DMError(ErrRedisRead, result.Err())
	}
	return result.Result()
}

func (r *RedisStore) deleteJsonData(ctx context.Context, key, path string) (interface{}, error) {
	path = r.getPath(path)
	res, err := r.Client.Do(ctx, "JSON.DEL", key, path).Result()
	if err != nil {
		return types.EMPTYSTR, dmerrors.DMError(ErrRedisDelete, err)
	}
	return res, nil
}

func (r *RedisStore) BFExists(ctx context.Context, key, item string) (interface{}, error) {
	return r.Client.Do(ctx, "BF.EXISTS", key, item).Result()
}

func (r *RedisStore) BFAdd(ctx context.Context, key, item string) (interface{}, error) {
	return r.Client.Do(ctx, "BF.EXISTS", key, item).Result()
}

func (r *RedisStore) PushQueueElem(ctx context.Context, key string, value ...any) (any, error) {
	return r.Client.LPush(ctx, key, value...).Result()
}

func (r *RedisStore) PopQueue(ctx context.Context, timeout time.Duration, keys ...string) (any, error) {
	return r.Client.BLPop(ctx, timeout, keys...).Result()
}
