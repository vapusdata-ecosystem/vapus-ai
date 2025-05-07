package snowflake

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusai/core/data-platform/dataservices/pkgs"
)

func (m *SnowflakeStore) Insert(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) error {
	// query the mysql database
	fields := make([]string, 0, len(param.DataSet))
	values := make([]string, 0, len(param.DataSet))

	for key, value := range param.DataSet {
		fields = append(fields, key)
		values = append(values, value.(string))
	}
	query := fmt.Sprintf("INSERT INTO (%v) (%v) VALUES (%v)", param.DataTable, fields, values)
	_, err := m.DB.Query(query, m.SchemaName)
	if err != nil {
		return err
	}
	return nil
}

// func (m *SnowflakeStore) InsertInBulk(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) (*datasvcpkgs.InsertDataResponse, error) {
// 	// resp is for response
// 	resp := &datasvcpkgs.InsertDataResponse{
// 		DataTable:       param.DataTable,
// 		RecordsInserted: 0,
// 	}
// 	// var wg sync.WaitGroup
// 	logger.Info().Msgf("Inserting ||||||||||||||>>>>>>>>>>>>>> %d records with sample %v with batch size %v", len(param.DataSets), param.DataSets[0], param.BatchSize)

// 	if len(param.DataSets) < int(param.BatchSize) {
// 		logger.Info().Msg("Inserting extracted row in single batch") //No gorm
// 		l := len(param.DataSets)                                     // no gorm
// 		// res,err
// 		res := m.Orm.Table(param.DataTable).CreateInBatches(param.DataSets, l)
// 		if res.Error != nil {
// 			resp.RecordsInserted += 0
// 		} else {
// 			if res.RowsAffected == int64(l) {
// 				// rows calculate kar rahe hai..
// 				resp.RecordsInserted += res.RowsAffected
// 			} else {
// 				resp.RecordsFailed += int64(l) - res.RowsAffected
// 				resp.RecordsInserted += res.RowsAffected
// 			}
// 		}
// 	} else {
// 		logger.Info().Msgf("Inserting extracted row in mulitple batch - %v - %v", len(param.DataSets), int(param.BatchSize))
// 		// Printed - Inserting extracted row in mulitple batch - 500 - 100
// 		var wg sync.WaitGroup
// 		for i := 0; i < len(param.DataSets); i += int(param.BatchSize) {
// 			end := i + int(param.BatchSize)
// 			if end > len(param.DataSets) {
// 				end = len(param.DataSets)
// 			}
// 			batch := param.DataSets[i:end]
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				res := m.Orm.Table(param.DataTable).CreateInBatches(batch, int(param.BatchSize))
// 				if res.Error != nil {
// 					resp.RecordsInserted += 0
// 					resp.RecordsFailed += int64(param.BatchSize)
// 				} else {
// 					resp.RecordsInserted += res.RowsAffected
// 					resp.RecordsFailed += int64(param.BatchSize) - res.RowsAffected
// 				}
// 			}()
// 		}
// 		wg.Wait()
// 	}
// 	return resp, nil
// }
