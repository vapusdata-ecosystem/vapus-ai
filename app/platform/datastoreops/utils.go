package dmstores

import (
	"context"

	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

func (ds *DMStore) CacheFilter(ctx context.Context, action, key string, value ...string) (interface{}, error) {
	switch action {
	case types.LIST:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case types.ADD:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.ADD", key, value[0]).Result()
	case types.EXISTS:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	case types.COUNT:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.CARD", key).Result()
	case types.MADD:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.MADD", key, value).Result()
	case types.DEL:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.DEL", key, value[0]).Result()
	default:
		return ds.Cacher.RedisClient.Client.Do(ctx, "CF.EXISTS", key, value[0]).Result()
	}
}

func (n *DMStore) GetDataSourceCreds(ctx context.Context, dataSource *models.DataSource) (*models.DataSourceCredsParams, error) {
	return apppkgs.GetDataSourceCreds(ctx, dataSource, n.VapusStore, logger)
}

func GetAccountAIAttributes(ctx context.Context, dmStore *DMStore, ctxClaim map[string]string) *models.AccountAIAttributes {
	poolItem := GetAccountFromPool(ctx, dmStore, ctxClaim)
	return poolItem.AIAttributes
}

const (
	PgCosignSimilarity = "COSINE"
	PgEuclidean        = "EUCLIDEAN"
	PgInnerProduct     = "INNER_PRODUCT"
)

var PGVectorAlgoMap = map[string]string{
	"COSINE":        "<=>",
	"EUCLIDEAN":     "<->",
	"INNER_PRODUCT": "<#>",
}
