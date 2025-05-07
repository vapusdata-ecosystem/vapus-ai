package snowflake

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
)

// Select query chal raha hai...  with or without filter
func (m *SnowflakeStore) Select(Ctx context.Context, query *string) (*sql.Rows, error) {
	// query the mysql database
	resp, err := m.DB.Query(*query, *m.SchemaName)
	if err != nil {
		return nil, err
	}
	m.logger.Info().Ctx(Ctx).Msgf("Query executed successfully - %v", resp)
	return resp, nil
}

func (m *SnowflakeStore) SelectWithFilter(ctx context.Context, queryOpts *datasvcpkgs.QueryOpts) ([]map[string]any, error) {

	result := make([]map[string]any, 0)

	// Build the SELECT clause
	fields := "*"
	if len(queryOpts.IncludeFields) > 0 {
		fields = strings.Join(queryOpts.IncludeFields, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", fields, queryOpts.DataCollection)

	var args []interface{}
	// QueryFilter konsa filter lgana hai and kon sa args hai..
	if len(queryOpts.QueryFilters) > 0 {
		filterClauses := []string{}
		for key, value := range queryOpts.QueryFilters {
			filterClauses = append(filterClauses, fmt.Sprintf("%s = ?", key))
			args = append(args, value)
		}
		query += " WHERE " + strings.Join(filterClauses, " AND ")
	}

	if queryOpts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", queryOpts.Limit)
	}

	log.Printf("Executing query: %s with args: %v", query, args)

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, col := range columns {
			rowMap[col] = values[i]
		}
		result = append(result, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
