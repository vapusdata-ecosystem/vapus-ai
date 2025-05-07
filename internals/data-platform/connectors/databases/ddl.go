package databases

import (
	"context"

	pkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func (svc *DataStoreClient) RunDDLs(ctx context.Context, request *string) error {
	switch svc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_MYSQL.String():
		return svc.MysqlClient.RunDDL(ctx, request)
	case types.StorageEngine_POSTGRES.String():
		return svc.PostgresClient.RunDDL(ctx, request)
	default:
		return dputils.ErrInvalidDataStorageEngine
	}
}

func (svc *DataStoreClient) CreateDataTables(ctx context.Context, opts *pkgs.DataTablesOpts) error {
	switch svc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_POSTGRES.String():
		return svc.PostgresClient.CreateDataTables(ctx, opts)
	case types.StorageEngine_ELASTICSEARCH.String():
		return svc.ElasticSearchClient.CreateIndexWithMapping(ctx, opts)
	default:
		return dputils.ErrInvalidDataStorageEngine
	}
}
