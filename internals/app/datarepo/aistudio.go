package datarepo

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func ListAIModelNodes(ctx context.Context, dmstroes *apppkgs.VapusStore, logger zerolog.Logger, condition string, ctxClaim map[string]string) ([]*models.AIModelNode, error) {
	if condition == "" {
		condition = "deleted_at IS NULL AND status = 'ACTIVE'"
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.AIModelsNodesTable, condition)
	result := []*models.AIModelNode{}
	err := dmstroes.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// _, err = ds.Db.PostgresClient.DB.NewSelect().Model(&result).ModelTableExpr(vapussvc.AIModelsNodesTable).Where(condition).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching ai model nodes from datastore")
		return nil, err
	}
	return result, nil
}

func ListAIGuardrails(ctx context.Context, dmstroes *apppkgs.VapusStore, logger zerolog.Logger, condition string, ctxClaim map[string]string) ([]*models.AIGuardrails, error) {
	if ctxClaim == nil {
		condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	}
	if condition == "" {
		condition = "deleted_at IS NULL AND status = 'ACTIVE'"
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.VapusGuardrailsTable, condition)
	result := []*models.AIGuardrails{}
	err := dmstroes.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while fetching ai guardrail list from datastore")
		return nil, err
	}
	return result, nil
}
