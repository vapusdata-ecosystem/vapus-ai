package sqlops

import (
	"fmt"
	"log"

	"github.com/doug-martin/goqu/v9"
	Builder "github.com/huandu/go-sqlbuilder"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type SqlBuilder struct {
	Table      string
	Schema     string
	DataBase   string
	Filters    map[string]any
	Limit      int64
	OrderField string
	OrderBy    int64
	Columns    []any
	Values     [][]interface{}
	Where      map[string]map[string]any
}

// TODO: https://github.com/huandu/go-sqlbuilder Use this instead of goqu
func BuildSql(opts *SqlBuilder, logger zerolog.Logger) (string, error) {
	// Use "default" as a neutral dialect if there's no specific SQL dialect
	dialect := goqu.Dialect("default")

	// Build a generic SQL query
	queryOpt := dialect.From(opts.Table)
	if len(opts.Columns) < 1 {
		queryOpt = queryOpt.Select("*")
	} else {
		log.Println("opts.Columns ==============", opts.Columns)
		queryOpt = queryOpt.Select(opts.Columns...)
	}
	whereParam := goqu.Ex{}
	log.Println("opts.Filters ==============", opts.Filters)
	for k, v := range opts.Filters {
		whereParam[k] = v
	}
	log.Println("whereParam ==============", whereParam)
	queryOpt = queryOpt.Where(whereParam)

	if opts.Limit > 0 {
		queryOpt = queryOpt.Limit(uint(opts.Limit))
	}

	if opts.OrderField != "" {
		if opts.OrderBy == 0 {
			queryOpt = queryOpt.Order(goqu.I(opts.OrderField).Asc())
		} else {
			queryOpt = queryOpt.Order(goqu.I(opts.OrderField).Desc())
		}
	}

	sql, _, err := queryOpt.ToSQL()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to generate SQL query")
		return "", err
	}

	// Print the generated SQL query
	return sql, nil
}

func InsertQueryBuilder(opts *SqlBuilder, DataSourceEngine string, logger zerolog.Logger) (string, []interface{}) {
	var insertBuilder *Builder.InsertBuilder
	// We will create a new insert builder
	switch DataSourceEngine {
	case types.StorageEngine_MYSQL.String():
		insertBuilder = Builder.MySQL.NewInsertBuilder()
	case types.StorageEngine_POSTGRES.String():
		insertBuilder = Builder.PostgreSQL.NewInsertBuilder()
	case types.StorageEngine_CLICKHOUSE.String():
		insertBuilder = Builder.ClickHouse.NewInsertBuilder()
	case types.StorageEngine_ORACLE.String():
		insertBuilder = Builder.ClickHouse.NewInsertBuilder()
	default:
		insertBuilder = Builder.NewInsertBuilder()
	}
	// Table Creation
	insertBuilder.InsertInto(opts.Table)

	// Column insertion
	var columns []string
	for _, val := range opts.Columns {
		str := fmt.Sprintf("%v", val)
		columns = append(columns, str)
	}
	insertBuilder.Cols(columns...)

	// Value insertion
	for _, val := range opts.Values {
		insertBuilder.Values(val...)
	}
	sql, args := insertBuilder.Build()

	fmt.Println(sql)
	fmt.Println(args)

	return sql, args
}

func UpdateQueryBuilder(opts *SqlBuilder, logger zerolog.Logger) (string, []interface{}) {
	updateBuilder := Builder.Update(opts.Table)

	// map[string]map[string]any
	// col: {operator: val}
	// "id": {"=": 123}
	for col, cond := range opts.Where {
		for op, val := range cond {
			switch op {
			case "=":
				updateBuilder.Where(
					updateBuilder.Equal(col, val))
			case "!=":
				updateBuilder.Where(
					updateBuilder.NotEqual(col, val))
			case ">":
				updateBuilder.Where(
					updateBuilder.GreaterThan(col, val))
			case "<":
				updateBuilder.Where(
					updateBuilder.LessThan(col, val))
			case ">=":
				updateBuilder.Where(
					updateBuilder.GreaterEqualThan(col, val))
			case "<=":
				updateBuilder.Where(
					updateBuilder.LessEqualThan(col, val))
			case "LIKE":
				updateBuilder.Where(
					updateBuilder.Like(col, val))
			case "NOT LIKE":
				updateBuilder.Where(
					updateBuilder.NotLike(col, val))
			}
		}
	}

	sql, args := updateBuilder.Build()

	fmt.Println(sql)
	fmt.Println(args)

	return sql, args
}

func DeleteQueryBuilder(opts *SqlBuilder, logger zerolog.Logger) (string, []interface{}) {
	deleteBuilder := Builder.DeleteFrom(opts.Table)

	for col, cond := range opts.Where {
		for op, val := range cond {
			switch op {
			case "=":
				deleteBuilder.Where(deleteBuilder.Equal(col, val))
			case "!=":
				deleteBuilder.Where(deleteBuilder.NotEqual(col, val))
			case "<":
				deleteBuilder.Where(deleteBuilder.LessThan(col, val))
			case ">":
				deleteBuilder.Where(deleteBuilder.GreaterThan(col, val))
			case "<=":
				deleteBuilder.Where(deleteBuilder.LessEqualThan(col, val))
			case ">=":
				deleteBuilder.Where(deleteBuilder.GreaterEqualThan(col, val))
			case "LIKE":
				deleteBuilder.Where(deleteBuilder.Like(col, val))
			case "NOT LIKE":
				deleteBuilder.Where(deleteBuilder.NotLike(col, val))
			}
		}
	}
	sql, args := deleteBuilder.Build()

	fmt.Println(sql)
	fmt.Println(args)

	return sql, args

}
