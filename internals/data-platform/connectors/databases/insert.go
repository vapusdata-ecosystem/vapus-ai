package databases

import (
	"context"

	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func (svc *DataStoreClient) InsertBulkDataSet(ctx context.Context, request *datasvcpkgs.InsertDataRequest) (*datasvcpkgs.InsertDataResponse, error) {
	switch svc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_MYSQL.String():
		return svc.MysqlClient.InsertInBulk(ctx, request, svc.Logger)
	case types.StorageEngine_POSTGRES.String():
		return svc.PostgresClient.InsertInBulk(ctx, request, svc.Logger)
	default:
		return nil, dputils.ErrInvalidDataStorageEngine
	}
}

func (svc *DataStoreClient) InsertDataSet(ctx context.Context, request *datasvcpkgs.InsertDataRequest) (*datasvcpkgs.InsertDataResponse, error) {
	resp := &datasvcpkgs.InsertDataResponse{
		DataTable:       request.DataTable,
		RecordsInserted: 0,
	}
	var err error
	switch svc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_MYSQL.String():
		err = svc.MysqlClient.Insert(ctx, request, svc.Logger)
	case types.StorageEngine_POSTGRES.String():
		err = svc.PostgresClient.Insert(ctx, request, svc.Logger)
	default:
		return nil, dputils.ErrInvalidDataStorageEngine
	}
	if err != nil {
		resp.RecordsFailed = 1
		return resp, err
	} else {
		resp.RecordsInserted = 1
		return resp, nil
	}
}
