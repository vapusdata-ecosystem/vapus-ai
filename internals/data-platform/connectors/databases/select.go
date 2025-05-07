package databases

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	dputils "github.com/vapusdata-ecosystem/vapusai/core/data-platform/utils"

	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func (dsc *DataStoreClient) prepareSqlQuery(qopts *datasvcpkgs.QueryOpts, logger zerolog.Logger) (*string, error) {
	query := "SELECT {field} FROM {table} where {condition}"
	query = strings.Replace(query, "{table}", qopts.DataCollection, -1)
	if qopts.QueryString != "" {
		query = strings.Replace(query, "{condition}", qopts.QueryString, -1)
	}
	if qopts.CountRecords {
		query = strings.Replace(query, "{field}", "count(*)", -1)
	} else {
		if len(qopts.IncludeFields) > 0 {
			query = strings.Replace(query, "{field}", strings.Join(qopts.IncludeFields, ","), -1)
		} else {
			query = strings.Replace(query, "{field}", "*", -1)
		}
	}
	logger.Info().Msgf("Query prepared by opts dataSVC - %v", query)
	return &query, nil
}

func (dsc *DataStoreClient) SelectWithFilter(ctx context.Context, qopts *datasvcpkgs.QueryOpts, resultObj interface{}, logger zerolog.Logger) ([]map[string]any, error) {
	var err error
	result := make([]map[string]any, 0)
	switch dsc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_REDIS.String():
		return result, nil
	case types.StorageEngine_MYSQL.String():
		result, err = dsc.MysqlClient.SelectWithFilter(ctx, qopts)
		if err != nil {
			return result, err
		}
		return result, nil
	case types.StorageEngine_POSTGRES.String():
		result, err = dsc.PostgresClient.SelectWithFilter(ctx, qopts)
		if err != nil {
			return result, err
		}
		return result, nil
	case types.StorageEngine_ELASTICSEARCH.String():
		return result, nil
	}
	return result, dputils.ErrInvalidDataStorageEngine
}

func (dsc *DataStoreClient) Select(ctx context.Context, qopts *datasvcpkgs.QueryOpts, dest interface{}, logger zerolog.Logger) error {
	// TO:DO - Add support for other data storage engines - https://github.com/Masterminds/squirrel
	switch dsc.DataStoreParams.DataSourceEngine {
	case types.StorageEngine_REDIS.String():
		return nil
	case types.StorageEngine_MYSQL.String():
		var query *string
		if qopts.RawQuery != "" {
			query = &qopts.RawQuery
		} else {
			q, err := dsc.prepareSqlQuery(qopts, logger)
			if err != nil {
				return err
			}
			query = q
		}
		rows, err := dsc.MysqlClient.Select(ctx, query)
		if err != nil {
			return err
		}
		return dputils.ScanSql(rows, dest, logger)
	case types.StorageEngine_POSTGRES.String():
		var query *string
		if qopts.RawQuery != "" {
			query = &qopts.RawQuery
		} else {
			q, err := dsc.prepareSqlQuery(qopts, logger)
			if err != nil {
				return err
			}
			query = q
		}
		rows, err := dsc.PostgresClient.Select(ctx, query)
		if err != nil {
			return err
		}
		return dputils.ScanSql(rows, dest, logger)
	case types.StorageEngine_ELASTICSEARCH.String():
		return nil
	}
	return dputils.ErrInvalidDataStorageEngine
}
