package redshift

import (
	"context"

	"database/sql"

	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
)

func (m *RedshiftStore) SelectInApp(ctx context.Context, query *string, result interface{}) error {
	// query the mysql database
	resp := m.Orm.Raw(*query)
	rows, err := resp.Rows()
	if err != nil {
		m.logger.Error().Ctx(ctx).Msgf("Error while executing query - %v", err)
		return err
	}

	m.logger.Info().Ctx(ctx).Msg("Query executed successfully")
	// defer rows.Close()
	err = m.Orm.ScanRows(rows, result)
	if err != nil {
		m.logger.Error().Ctx(ctx).Msgf("Error while scanning rows - %v", err)
		return err
	}
	return nil
}

// func (m *RedshiftStore) SelectWithFilter(Ctx context.Context, queryOpts *datasvcpkgs.QueryOpts, result interface{}) error {
// 	// query the mysql database
// 	query := m.DB.NewSelect().Table(queryOpts.DataCollection)
// 	for key, value := range queryOpts.QueryFilters {
// 		v := reflect.ValueOf(value)
// 		kind := v.Kind()
// 		switch kind {
// 		case reflect.Slice, reflect.Array:
// 			// Use IN clause for slice/array values
// 			query.Where(fmt.Sprintf("%s IN (?)", key), bun.In(value))
// 		// case reflect.Map:
// 		// 	// Use JSONB containment (@>) for map values
// 		// 	// Assuming the map is intended for a JSONB column named exactly as key
// 		// 	query.Where(fmt.Sprintf("%s @> ?", key), bun. (value))
// 		default:
// 			// Use equality for scalar values
// 			query.Where(fmt.Sprintf("%s = ?", key), value)
// 		}
// 	}
// 	// result := make([]map[string]any, 0)
// 	if queryOpts.Limit > 0 {
// 		query = query.Limit(int(queryOpts.Limit))
// 	}
// 	if queryOpts.OrderBy != nil {
// 		query = query.Order(queryOpts.OrderBy.Field + " " + queryOpts.OrderBy.Order.String())
// 	}
// 	if len(queryOpts.IncludeFields) > 0 {
// 		query = query.Column(queryOpts.IncludeFields...)
// 	}
// 	if err := query.Scan(Ctx, &result); err != nil {
// 		return err
// 	}
// 	// fmt.Println("Executed SQL:", tx.Statement.SQL.String())
// 	// fmt.Println("Variables:", tx.Statement.Vars)
// 	return nil
// }

func (m *RedshiftStore) Select(Ctx context.Context, query *string) (*sql.Rows, error) {
	// query the mysql database
	resp, err := m.Orm.Raw(*query).Rows()
	if err != nil {
		return nil, err
	}
	m.logger.Info().Ctx(Ctx).Msgf("Query executed successfully - %v", resp)
	return resp, nil
}

func (m *RedshiftStore) SelectWithFilter(Ctx context.Context, queryOpts *datasvcpkgs.QueryOpts) ([]map[string]any, error) {
	// query the mysql database
	query := m.Orm.Table(queryOpts.DataCollection)
	for key, value := range queryOpts.QueryFilters {
		query = query.Where(key, value)
	}
	result := make([]map[string]any, 0)
	if queryOpts.Limit > 0 {
		query = query.Limit(int(queryOpts.Limit))
	}
	if queryOpts.OrderBy != nil {
		query = query.Order(queryOpts.OrderBy.Field + " " + queryOpts.OrderBy.Order.String())
	}
	if len(queryOpts.IncludeFields) > 0 {
		for _, field := range queryOpts.IncludeFields {
			query = query.Select(field)
		}
	}
	if err := query.Find(&result).Error; err != nil {
		return result, err
	}
	// fmt.Println("Executed SQL:", tx.Statement.SQL.String())
	// fmt.Println("Variables:", tx.Statement.Vars)
	return result, nil
}
