package snowflake

import (
	"context"
	"fmt"

	pkgs "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
)

// RunDDL executes a DDL query on the Snowflake database
func (m *SnowflakeStore) RunDDL(ctx context.Context, query *string) error {
	if query == nil || *query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	// Execute the query using the database connection
	_, err := m.DB.Query(*query, *m.SchemaName)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to execute DDL query")
		return err
	}

	m.logger.Info().Msg("DDL query executed successfully")
	return nil
}

// CreateDataTables creates data tables based on the provided options
func (m *SnowflakeStore) CreateDataTables(ctx context.Context, opts *pkgs.DataTablesOpts) error {
	if opts == nil {
		return pkgs.ErrInvalidDataTablesOpts
	}

	if opts.Query != "" {
		return m.RunDDL(ctx, &opts.Query)
	} else if opts.StructScheme != nil {
		return m.CreateTable(ctx, opts.StructScheme, opts.Name)
	} else if opts.MapsScheme != nil {
		return m.CreateTable(ctx, opts.MapsScheme, opts.Name)
	}

	return pkgs.ErrInvalidDataTablesOpts
}

// CreateTable creates a table in Snowflake based on the provided schema and table name
func (m *SnowflakeStore) CreateTable(ctx context.Context, modelStruct interface{}, tablename string) error {
	if tablename == "" {
		return fmt.Errorf("tablename cannot be empty")
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", tablename)

	_, err := m.DB.Query(query, *m.SchemaName)
	if err != nil {
		m.logger.Error().Err(err).Msgf("Failed to create table %s", tablename)
		return err
	}

	m.logger.Info().Msgf("Table %s created successfully", tablename)
	return nil
}
