package redshift

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
	builder "github.com/vapusdata-ecosystem/vapusai/core/tools/sqlops"
)

func (m *RedshiftStore) Insert(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) error {
	// fmt.Sprintf("%s.%s", param.TableSchema, param.DataTable)
	var columns []any
	var value []interface{}
	var values [][]interface{}
	for key, val := range param.DataSet {
		columns = append(columns, key)
		value = append(value, val)
	}
	values = append(values, value)

	// Calling InsertQueryBuilder
	query, args := builder.InsertQueryBuilder(&builder.SqlBuilder{
		Table:   fmt.Sprintf("%s.%s", param.TableSchema, param.DataTable),
		Columns: columns,
		Values:  values,
	}, "POSTGRES", logger)

	_, err := m.Conn.Query(query, args...)
	if err != nil {
		logger.Err(err).Msgf("Error while inserting values")
	}
	// var temp map[string]any

	return nil
}

func (m *RedshiftStore) InsertInBulk(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) (*datasvcpkgs.InsertDataResponse, error) {
	resp := &datasvcpkgs.InsertDataResponse{
		DataTable:       param.DataTable,
		RecordsInserted: 0,
	}
	// var wg sync.WaitGroup
	var columns []any
	var values [][]interface{}

	for key, _ := range param.DataSets[0] {
		columns = append(columns, key)
	}

	for _, vals := range param.DataSets {
		var value []interface{}
		for _, col := range columns {
			val := vals[col.(string)]
			value = append(value, val)
		}
		values = append(values, value)
	}

	// Wrong Logic: Map is not order
	// for i, vals := range param.DataSets {
	// 	var value []interface{}
	// 	for key, val := range vals {
	// 		if i == 0 {
	// 			columns = append(columns, key)
	// 		}
	// 		value = append(value, val)
	// 	}
	// 	values = append(values, value)
	// }

	logger.Info().Msgf("Inserting ||||||||||||||>>>>>>>>>>>>>> %d records with sample %v with batch size %v", len(param.DataSets), param.DataSets[0], param.BatchSize)
	if len(param.DataSets) < int(param.BatchSize) {
		logger.Info().Msg("Inserting extracted row in single batch")
		l := len(param.DataSets)

		// Insert Query Builder
		query, args := builder.InsertQueryBuilder(&builder.SqlBuilder{
			Table:   fmt.Sprintf("%s.%s", param.TableSchema, param.DataTable),
			Columns: columns,
			Values:  values,
		}, "POSTGRES", logger)
		res, err := m.Conn.Exec(query, args...)
		if err != nil {
			logger.Err(err).Msgf("Error while inserting values")
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			resp.RecordsInserted += 0
		} else {
			if rowsAffected == int64(l) {
				resp.RecordsInserted += rowsAffected
			} else {
				resp.RecordsFailed += int64(l) - rowsAffected
				resp.RecordsInserted += rowsAffected
			}
		}
	} else {
		logger.Info().Msgf("Inserting extracted row in mulitple batch - %v - %v", len(param.DataSets), int(param.BatchSize))
		// Printed - Inserting extracted row in mulitple batch - 500 - 100
		var wg sync.WaitGroup
		for i := 0; i < len(param.DataSets); i += int(param.BatchSize) {
			end := i + int(param.BatchSize)
			if end > len(param.DataSets) {
				end = len(param.DataSets)
			}
			//batch := param.DataSets[i:end]
			batch := values[i:end]
			batchLen := len(batch)
			wg.Add(1)
			go func() {
				defer wg.Done()
				query, args := builder.InsertQueryBuilder(&builder.SqlBuilder{
					Table:   fmt.Sprintf("%s.%s", param.TableSchema, param.DataTable),
					Columns: columns,
					Values:  batch,
				}, "POSTGRES", logger)
				res, err := m.Conn.Exec(query, args...)
				if err != nil {
					logger.Err(err).Msgf("Error while inserting values")
				}
				rowsAffected, err := res.RowsAffected()
				// res := m.Orm.Table(param.DataTable).CreateInBatches(batch, int(param.BatchSize))
				if err != nil {
					resp.RecordsInserted += 0
					resp.RecordsFailed += int64(param.BatchSize)
				} else {
					resp.RecordsInserted += rowsAffected
					if rowsAffected != int64(batchLen) {
						resp.RecordsFailed += int64(param.BatchSize) - rowsAffected
					}
				}
			}()
		}
		wg.Wait()
	}
	return resp, nil
}
