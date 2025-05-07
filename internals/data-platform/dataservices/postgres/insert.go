package postgres

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	datasvcpkgs "github.com/vapusdata-ecosystem/vapusdata/core/data-platform/dataservices/pkgs"
)

// func (m *PostgresStore) InsertInBulk(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) (*datasvcpkgs.InsertDataResponse, error) {
// 	resp := &datasvcpkgs.InsertDataResponse{
// 		DataTable:       param.DataTable,
// 		RecordsInserted: 0,
// 	}
// 	var batchErr = make(chan error, 1)
// 	var wg sync.WaitGroup
// 	logger.Info().Msgf("Inserting ||||||||||||||>>>>>>>>>>>>>> %d records with sample %v", len(param.DataSets), param.DataSets)
// 	for i := 0; i < len(param.DataSets); i += int(param.BatchSize) {
// 		end := i + int(param.BatchSize)
// 		if end > len(param.DataSets) {
// 			end = len(param.DataSets)
// 		}

// 		batch := (param.DataSets)[i:end]

// 		var vv []*map[string]any
// 		for _, v := range batch {
// 			vv = append(vv, &v)
// 		}

// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			res, err := m.DB.NewInsert().Model(vv).ModelTableExpr(param.DataTable).Exec(ctx)
// 			if err != nil {
// 				batchErr <- err
// 			} else {
// 				cv, err := res.RowsAffected()
// 				if err != nil {
// 					batchErr <- err
// 				} else {
// 					resp.RecordsInserted += cv
// 					if cv != int64(end) {
// 						resp.RecordsFailed += int64(end) - cv
// 					}
// 				}
// 			}
// 		}()
// 	}
// 	go func() {
// 		wg.Wait()
// 		close(batchErr)
// 	}()
// 	err := <-batchErr
// 	if err != nil {
// 		return resp, err
// 	}
// 	return resp, nil
// }

// func (m *PostgresStore) Insert(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) error {
// 	// query the mysql database
// 	var err error
// 	if param.DataSet != nil {
// 		_, err = m.DB.NewInsert().Model(param.DataSet).Table(param.DataTable).Exec(ctx)
// 	} else if param.StructDataSet != nil {
// 		_, err = m.DB.NewInsert().Model(param.StructDataSet).Table(param.DataTable).Exec(ctx)
// 	}
// 	return err
// }

func (m *PostgresStore) InsertInBulk(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) (*datasvcpkgs.InsertDataResponse, error) {
	resp := &datasvcpkgs.InsertDataResponse{
		DataTable:       param.DataTable,
		RecordsInserted: 0,
	}
	// var wg sync.WaitGroup
	logger.Info().Msgf("Inserting ||||||||||||||>>>>>>>>>>>>>> %d records with sample %v with batch size %v", len(param.DataSets), param.DataSets[0], param.BatchSize)
	if len(param.DataSets) < int(param.BatchSize) {
		logger.Info().Msg("Inserting extracted row in single batch")
		l := len(param.DataSets)
		res := m.Orm.Table(param.DataTable).CreateInBatches(param.DataSets, l)
		if res.Error != nil {
			resp.RecordsInserted += 0
		} else {
			if res.RowsAffected == int64(l) {
				resp.RecordsInserted += res.RowsAffected
			} else {
				resp.RecordsFailed += int64(l) - res.RowsAffected
				resp.RecordsInserted += res.RowsAffected
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
			batch := param.DataSets[i:end]
			wg.Add(1)
			go func() {
				defer wg.Done()
				res := m.Orm.Table(param.DataTable).CreateInBatches(batch, int(param.BatchSize))
				if res.Error != nil {
					resp.RecordsInserted += 0
					resp.RecordsFailed += int64(param.BatchSize)
				} else {
					resp.RecordsInserted += res.RowsAffected
					resp.RecordsFailed += int64(param.BatchSize) - res.RowsAffected
				}
			}()
		}
		wg.Wait()
	}
	return resp, nil
}

func (m *PostgresStore) Insert(ctx context.Context, param *datasvcpkgs.InsertDataRequest, logger zerolog.Logger) error {
	// query the mysql database
	return m.Orm.Table(param.DataTable).Create(param.DataSet).Error
}
