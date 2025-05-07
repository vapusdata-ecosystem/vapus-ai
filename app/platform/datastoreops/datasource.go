package dmstores

import (
	"context"
	"encoding/json"
	"fmt"

	// mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

// CreateDataSource creates a new data source in the data organization store
func (ds *DMStore) CreateDataSource(ctx context.Context, dataSource *models.DataSource, ctxClaim map[string]string) error {
	dataSource.SetAccountId(ctxClaim[encryption.ClaimAccountKey])
	if ds.Cacher != nil {
		_, err := ds.Cacher.RedisClient.WrtiteData(ctx, dataSource.VapusID, types.EMPTYSTR, dataSource)
		if err != nil {
			logger.Err(err).Ctx(ctx).Msg(apperr.ErrCreateDataSource.Error())
			return err
		}
	}

	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(dataSource).ModelTableExpr(apppkgs.DataSourcesTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving data source in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   dataSource.VapusID,
			ResourceName: mpb.Resources_DATASOURCES.String(),
			VapusBase: models.VapusBase{
				Editors: dataSource.Editors,
			},
		}, logger, ctxClaim)
	}()
	return nil
}

func (ds *DMStore) PutDataSource(ctx context.Context, obj *models.DataSource, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.DataSourcesTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating data source in datastore")
		return err
	}
	return nil
}

func (ds *DMStore) ListDataSources(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.DataSource, error) {
	result := []*models.DataSource{}
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.DataSourcesTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// err := ds.Db.Select(ctx, &datasvcpkgs.QueryOpts{
	// 	DataCollection: apppkgs.DataSourcesTable,
	// 	RawQuery:       query,
	// }, &result, logger)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting data sources from datastore")
		return nil, err
	}
	return result, err
}

func (ds *DMStore) CountDataSources(ctx context.Context, condition string, ctxClaim map[string]string) (int64, error) {
	var result int64
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", apppkgs.DataSourcesTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting data sources count from datastore")
		return 0, err
	}
	return result, err
}

func (ds *DMStore) GetDataSource(ctx context.Context, iden string, ctxClaim map[string]string) (*models.DataSource, error) {
	if iden == "" {
		return nil, apperr.ErrDataSource404
	}

	result := []*models.DataSource{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.DataSourcesTable, apppkgs.GetByIdFilter("", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// err := ds.Db.Select(ctx, &datasvcpkgs.QueryOpts{
	// 	DataCollection: apppkgs.DataSourcesTable,
	// 	RawQuery:       query,
	// }, &result, logger)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting data sources from datastore")
		return nil, apperr.ErrDataSource404
	}
	return result[0], err
}

func (ds *DMStore) GetDataCredsFromSecret(ctx context.Context, secretName string) (*models.DataSourceCreds, error) {
	logger.Debug().Ctx(ctx).Msgf("Getting secret for %v", secretName)
	secretStr, err := ds.SecretStore.ReadSecret(ctx, secretName)
	gCred := &models.GenericCredentialModel{}
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while reading the secrets from secret store")
		return nil, err
	}
	err = json.Unmarshal([]byte(dmutils.AnyToStr(secretStr)), gCred)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while unmarshaling the secrets from secret store")
		return nil, err
	}
	return &models.DataSourceCreds{
		SecretName:          secretName,
		Credentials:         gCred,
		IsAlreadyInSecretBs: true,
		Name:                secretName,
	}, nil
}

func (ds *DMStore) ListDataSourceSyncLogs(ctx context.Context, dsId string, ctxClaim map[string]string) ([]*models.NabrunnerLog, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE resource_id='%s'", apppkgs.NabrunnerLogTable, dsId)
	result := []*models.NabrunnerLog{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching nab runner schedules")
		return nil, err
	}
	return result, nil
}

func (ds *DMStore) GetDataSourceSyncSchedules(ctx context.Context, dsId string, ctxClaim map[string]string) ([]*models.NabrunnerLog, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted_at is NULL AND is_recurring is true AND resource_id='%s'", apppkgs.NabrunnerLogTable, dsId)
	result := []*models.NabrunnerLog{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching nab runner schedules")
		return nil, err
	}
	return result, nil
}
