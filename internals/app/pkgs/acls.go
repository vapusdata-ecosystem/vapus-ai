package apppkgs

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/data-platform/connectors/databases"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
)

func AddResourceArn(ctx context.Context,
	dbCl *databases.DataStoreClient,
	obj *models.VapusResourceArn,
	logger zerolog.Logger,
	ctxClaim map[string]string) error {

	obj.PreSaveCreate(ctxClaim)
	obj.ResourceARN = fmt.Sprintf(types.RESOURCE_ARN, obj.ResourceName, obj.ResourceId)
	_, err := dbCl.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(VapusResourceArnTable).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving resource arn to datastore")
		return err
	}
	return nil
}

func LoadAllResourceArns(ctx context.Context,
	dbCl *databases.DataStoreClient,
	logger zerolog.Logger,
	ctxClaim map[string]string) ([]*models.VapusResourceArn, error) {
	var objs []*models.VapusResourceArn
	_, err := dbCl.PostgresClient.DB.NewSelect().
		Model(objs).
		ModelTableExpr(VapusResourceArnTable).
		Where("deleted_at IS NULL").
		Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while  loading all resource arns from datastore")
		return nil, err
	}
	return objs, nil
}
