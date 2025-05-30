package aidmstore

import (
	"context"
	"fmt"
	"log"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (ds *AIStudioDMStore) ConfigureAIPrompts(ctx context.Context, obj *models.AIPrompt, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.AIModelPromptTable).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai model prompt to datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.Resources_AIPROMPTS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, ds.logger, ctxClaim)
	}()
	return nil
}

func (ds *AIStudioDMStore) ListAIPrompts(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.AIPrompt, error) {
	condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.AIModelPromptTable, condition)
	result := []*models.AIPrompt{}
	log.Println("query:", query)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while fetching ai prompt list from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *AIStudioDMStore) GetAIPrompt(ctx context.Context, iden string, ctxClaim map[string]string) (*models.AIPrompt, error) {
	if iden == "" {
		return nil, apperr.ErrAIPrompt404
	}
	result := []*models.AIPrompt{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.AIModelPromptTable, apppkgs.GetByIdFilter("", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		ds.logger.Err(err).Ctx(ctx).Msg("error while getting ai model prompt details from datastore")
		return nil, err
	}
	return result[0], err
}

func (ds *AIStudioDMStore) PutAIPrompts(ctx context.Context, obj *models.AIPrompt, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.AIModelPromptTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while updating ai model prompt to datastore")
		return err
	}
	return nil
}
