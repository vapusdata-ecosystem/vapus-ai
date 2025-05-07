package aidmstore

import (
	"context"
	"fmt"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (ds *AIStudioDMStore) ConfigureAIGuardrails(ctx context.Context, obj *models.AIGuardrails, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.VapusGuardrailsTable).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai guardrail to datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.Resources_AIAGENTS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, ds.logger, ctxClaim)
	}()
	return nil
}

func (ds *AIStudioDMStore) SaveGuardrailThread(ctx context.Context, obj *models.AIGuardrails, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.VapusGuardrailsTable).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai guardrail thread to datastore")
		return err
	}
	return nil
}

func (ds *AIStudioDMStore) ListAIGuardrails(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.AIGuardrails, error) {
	if ctxClaim == nil {
		condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	}
	if condition == "" {
		condition = "deleted_at IS NULL AND status = 'ACTIVE'"
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.VapusGuardrailsTable, condition)
	result := []*models.AIGuardrails{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while fetching ai guardrail list from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *AIStudioDMStore) GetAIGuardrail(ctx context.Context, iden string, ctxClaim map[string]string) (*models.AIGuardrails, error) {
	if iden == "" {
		return nil, apperr.ErrAIGuardrail404
	}
	result := []*models.AIGuardrails{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.VapusGuardrailsTable, apppkgs.GetByIdFilter("", iden, ctxClaim))
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		ds.logger.Err(err).Ctx(ctx).Msg("error while getting ai guardrail details from datastore")
		return nil, err
	}
	return result[0], err
}

func (ds *AIStudioDMStore) PutAIGuardrails(ctx context.Context, obj *models.AIGuardrails, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.VapusGuardrailsTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while updating ai guardrail to datastore")
		return err
	}
	return nil
}
