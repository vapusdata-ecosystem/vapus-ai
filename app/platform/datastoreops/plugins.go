package dmstores

import (
	"context"
	"fmt"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func (ds *DMStore) ConfigurePlugin(ctx context.Context, obj *models.Plugin, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.PluginsTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving plugin config in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.Resources_PLUGINS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, logger, ctxClaim)
	}()
	return nil
}

func (ds *DMStore) ListPlugins(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.Plugin, error) {
	var result []*models.Plugin
	if ctxClaim != nil {
		condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.PluginsTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching plugin by query from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *DMStore) CountPlugins(ctx context.Context, condition string, ctxClaim map[string]string) (int64, error) {
	var result int64
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", apppkgs.PluginsTable, condition)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while getting plugins count from datastore")
		return 0, err
	}
	return result, err
}

func (ds *DMStore) GetPlugin(ctx context.Context, iden string, ctxClaim map[string]string) (*models.Plugin, error) {
	if iden == "" {
		return nil, apperr.ErrPlugin404
	}
	var result []*models.Plugin
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.PluginsTable, apppkgs.GetByIdFilter("", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while fetching plugin by query from datastore")
		return nil, apperr.ErrPlugin404
	}
	return result[0], nil
}

func (ds *DMStore) PutPlugin(ctx context.Context, obj *models.Plugin, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.PluginsTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while updating plugin in datastore")
		return err
	}
	return nil
}

// func (ds *DMStore) SetPluginCreds(ctx context.Context, secretName string, creds *models.GenericCredentialModel, ctxClaim map[string]string) error {
// 	result, err := dmutils.StructToMap(creds)
// 	if err != nil {
// 		logger.Err(err).Ctx(ctx).Msgf("error while converting struct to map")
// 		return err
// 	}

// 	err = ds.SecretStore.WriteSecret(ctx, result, secretName)
// 	if err != nil {
// 		logger.Err(err).Ctx(ctx).Msgf("error while writing secret %v", secretName)
// 		return err
// 	}
// 	return nil
// }

// func (ds *DMStore) GetPluginCreds(ctx context.Context, secretName string, creds *models.GenericCredentialModel, ctxClaim map[string]string) error {
// 	result, err := dmutils.StructToMap(creds)
// 	if err != nil {
// 		logger.Err(err).Ctx(ctx).Msgf("error while converting struct to map")
// 		return err
// 	}

// 	err = ds.SecretStore.WriteSecret(ctx, result, secretName)
// 	if err != nil {
// 		logger.Err(err).Ctx(ctx).Msgf("error while writing secret %v", secretName)
// 		return err
// 	}
// 	return nil
// }
