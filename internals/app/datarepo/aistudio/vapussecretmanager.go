package aidmstore

import (
	"context"
	"fmt"

	"github.com/databricks/databricks-sql-go/logger"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (ds *AIStudioDMStore) CreateVapusSecret(ctx context.Context, obj *models.SecretStore, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.SecretStoreTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving vapus-secrets config in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.Resources_SECRETS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, ds.logger, ctxClaim)
	}()
	return nil
}

func (ds *AIStudioDMStore) ListVapusSecrets(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.SecretStore, error) {
	var result []*models.SecretStore
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.SecretStoreTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// err := ds.Db.Select(ctx, &datasvcpkgs.QueryOpts{
	// 	DataCollection: apppkgs.DataSourcesMetadataTable,
	// 	RawQuery:       query,
	// }, &result, logger)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching vapus-secrets by query from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *AIStudioDMStore) GetVapusSecret(ctx context.Context, iden string, ctxClaim map[string]string) (*models.SecretStore, error) {
	if iden == "" {
		return nil, apperr.ErrVapusSecret404
	}
	var result []*models.SecretStore
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.SecretStoreTable, apppkgs.GetByIdFilter("name", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// err := ds.Db.Select(ctx, &datasvcpkgs.QueryOpts{
	// 	DataCollection: apppkgs.DataSourcesMetadataTable,
	// 	RawQuery:       query,
	// }, &result, logger)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while fetching vapus-secrets by query from datastore")
		return nil, apperr.ErrVapusSecret404
	}
	return result[0], nil
}

func (ds *AIStudioDMStore) PutVapusSecret(ctx context.Context, obj *models.SecretStore, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.SecretStoreTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating vapus-secrets in datastore")
		return err
	}
	return nil
}

func (ds *AIStudioDMStore) DeleteVapusSecret(ctx context.Context, iden string, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewDelete().Model(&models.SecretStore{}).ModelTableExpr(apppkgs.SecretStoreTable).Where(apppkgs.VapusIdFilter(), iden).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while deleting vapus-secrets from datastore")
		return err
	}
	return nil
}
