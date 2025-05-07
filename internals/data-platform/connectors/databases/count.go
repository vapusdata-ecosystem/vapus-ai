package databases

import (
	"context"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func (dsc *DataStoreClient) CountRows(ctx context.Context, qopts *datasvcpkgs.QueryOpts, logger zerolog.Logger) (int64, error) {
	switch dsc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_POSTGRES.String():
		return dsc.PostgresClient.Count(ctx, qopts)
	case types.StorageEngine_ELASTICSEARCH.String():
		return dsc.ElasticSearchClient.Count(ctx, qopts)
	default:
		logger.Error().Msg("Invalid data storage engine")
		return 0, dputils.ErrInvalidDataStorageEngine
	}
}
