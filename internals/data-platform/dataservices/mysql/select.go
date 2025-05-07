package mysql

import (
	"context"
	"database/sql"

	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
)

func (m *MysqlStore) Select(Ctx context.Context, query *string) (*sql.Rows, error) {
	// query the mysql database
	resp, err := m.DB.Raw(*query).Rows()
	if err != nil {
		return nil, err
	}
	m.logger.Info().Ctx(Ctx).Msgf("Query executed successfully - %v", resp)
	return resp, nil
}

func (m *MysqlStore) SelectWithFilter(Ctx context.Context, queryOpts *datasvcpkgs.QueryOpts) ([]map[string]any, error) {
	// query the mysql database
	result := make([]map[string]any, 0)
	query := m.DB.Table(queryOpts.DataCollection)
	for key, value := range queryOpts.QueryFilters {
		query = query.Where(key, value)
	}
	if queryOpts.Limit > 0 {
		query = query.Limit(int(queryOpts.Limit))
	}
	if len(queryOpts.IncludeFields) > 0 {
		for _, field := range queryOpts.IncludeFields {
			query = query.Select(field)
		}
	}
	if err := query.Find(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}
