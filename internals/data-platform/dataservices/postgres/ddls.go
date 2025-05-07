package postgres

import (
	"context"
	"log"
	"reflect"

	pkgs "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
)

func (m *PostgresStore) RunDDL(ctx context.Context, query *string) error {
	// query the mysql database
	if m.Orm != nil {
		_, err := m.Orm.Raw(*query).Rows()
		if err != nil {
			return err
		}
		return nil
	}
	m.logger.Debug().Msgf("Running DDL query: %s", *query)
	_, err := m.DB.NewRaw(*query).Exec(ctx)
	return err
}

func (m *PostgresStore) CreateDataTables(ctx context.Context, opts *pkgs.DataTablesOpts) error {
	if opts == nil {
		return pkgs.ErrInvalidDataTablesOpts
	}
	if opts.Query != "" {
		return m.DB.NewRaw(opts.Query).Scan(ctx)
	} else if opts.StructScheme != nil {
		return m.CreateTable(ctx, opts.StructScheme, opts.Name)
	} else if opts.MapsScheme != nil {
		return m.CreateTable(ctx, opts.MapsScheme, opts.Name)
	} else {
		return pkgs.ErrInvalidDataTablesOpts
	}
}

func (m *PostgresStore) CreateTable(ctx context.Context, modelStruct interface{}, tablename string) error {
	mTy := reflect.TypeOf(modelStruct)
	if mTy.Kind() != reflect.Ptr {
		return pkgs.ErrInvalidModelStruct
	}
	query := m.DB.NewCreateTable().Model(modelStruct).IfNotExists()
	if tablename != "" {
		query.ModelTableExpr(tablename)
	}
	log.Println(query.String())
	_, err := query.Exec(ctx)
	return err
}
