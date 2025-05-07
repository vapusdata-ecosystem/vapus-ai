package aidmstore

import (
	"context"

	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func (ds *AIStudioDMStore) GetAIModelNodeNetworkParams(ctx context.Context, aiModelNode *models.AIModelNode) (*models.AIModelNodeNetworkParams, error) {
	secrets, err := apppkgs.ReadCredentialFromStore(ctx, aiModelNode.NetworkParams.SecretName, ds.VapusStore, ds.logger)
	if err != nil {
		ds.logger.Err(err).Msg("error while reading credentials from store")
		return nil, err
	}

	return &models.AIModelNodeNetworkParams{
		Url:                 aiModelNode.NetworkParams.GetUrl(),
		Credentials:         secrets,
		ApiVersion:          aiModelNode.NetworkParams.GetApiVersion(),
		LocalPath:           aiModelNode.NetworkParams.GetLocalPath(),
		SecretName:          aiModelNode.NetworkParams.SecretName,
		IsAlreadyInSecretBs: aiModelNode.NetworkParams.IsAlreadyInSecretBs,
	}, nil
}

func (ds *AIStudioDMStore) CacheFilter(ctx context.Context, action, key string, value ...string) (interface{}, error) {
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

func (n *AIStudioDMStore) GetDataSourceCreds(ctx context.Context, dataSource *models.DataSource) (*models.DataSourceCredsParams, error) {
	return apppkgs.GetDataSourceCreds(ctx, dataSource, n.VapusStore, n.logger)
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
