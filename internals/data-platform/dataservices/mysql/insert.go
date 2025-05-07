package mysql

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
)

func (m *MysqlStore) prepareMysqlInsertStmt(Ctx context.Context, table string, dataset *map[string]any, bulk bool) *string {
	// insert data set into mysql
	fields := ""
	values := ""
	for key, value := range *dataset {
		fields += key + ","
		valType := reflect.TypeOf(value)
		if valType.Kind() == reflect.Ptr {
			valType = valType.Elem()
		}
		switch valType.Kind() {
		case reflect.String:
			values += "'" + value.(string) + "',"
		case reflect.Int:
			values += strconv.Itoa(value.(int)) + ","
		case reflect.Int8:
			values += strconv.Itoa(int(value.(int8))) + ","
		case reflect.Int16:
			values += strconv.Itoa(int(value.(int16))) + ","
		case reflect.Int32:
			values += strconv.Itoa(int(value.(int32))) + ","
		case reflect.Int64:
			values += strconv.Itoa(int(value.(int64))) + ","
		case reflect.Float32:
			values += strconv.FormatFloat(value.(float64), 'f', -1, 64) + ","
		case reflect.Float64:
			values += strconv.FormatFloat(value.(float64), 'f', -1, 64) + ","
		case reflect.Bool:
			values += strconv.FormatBool(value.(bool)) + ","
		default:
			continue // TO:DO: handle other types and default types as well
		}
	}
	fields = strings.TrimSuffix(fields, ",")
	values = strings.TrimSuffix(values, ",")
	stmt := fmt.Sprintf("INSERT INTO `"+table+"` (%v) VALUES (%v)", fields, values)
	return &stmt
}

func (m *MysqlStore) InsertInBulk(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) (*datasvcpkgs.InsertDataResponse, error) {
	resp := &datasvcpkgs.InsertDataResponse{
		DataTable:       param.DataTable,
		RecordsInserted: 0,
	}
	// var batchErr = make(chan error, 1)
	// var wg sync.WaitGroup
	// for i := 0; i < len(param.DataSet); i += int(param.BatchSize) {
	// 	end := i + int(param.BatchSize)
	// 	if end > len(param.DataSet) {
	// 		end = len(param.DataSet)
	// 	}

	// 	batch := param.DataSet[i:end]
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		res := m.DB.Table(param.DataTable).CreateInBatches(batch, int(param.BatchSize))
	// 		if res.Error != nil {
	// 			batchErr <- res.Error
	// 		} else {
	// 			if res.RowsAffected == int64(end) {
	// 				resp.RecordsInserted += res.RowsAffected
	// 			} else {
	// 				resp.RecordsFailed += int64(end) - res.RowsAffected
	// 				resp.RecordsInserted += res.RowsAffected
	// 			}
	// 		}
	// 	}()
	// }
	// go func() {
	// 	wg.Wait()
	// 	close(batchErr)
	// }()
	// err := <-batchErr
	// if err != nil {
	// 	return resp, err
	// }
	var wg sync.WaitGroup
	for i := 0; i < len(param.DataSets); i += int(param.BatchSize) {
		end := i + int(param.BatchSize)
		if end > len(param.DataSets) {
			end = len(param.DataSets)
		}

		batch := (param.DataSets)[i:end]
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := m.DB.Table(param.DataTable).CreateInBatches(batch, int(param.BatchSize))
			if res.Error != nil {
				resp.RecordsInserted += 0
				resp.RecordsFailed += int64(param.BatchSize)
			} else {
				resp.RecordsInserted += res.RowsAffected
				resp.RecordsFailed += int64(param.BatchSize) - res.RowsAffected
			}
			logger.Info().Msgf("Number of rows inserted successfully - %v", res.RowsAffected)
		}()
	}
	go func() {
		wg.Wait()
	}()
	return resp, nil
}

func (m *MysqlStore) Insert(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) error {
	// query the mysql database
	return m.DB.Table(param.DataTable).Create(param.DataSet).Error
}
