package aidmstore

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func SaveAIToolCallLog(ctx context.Context, dmstore *apppkgs.VapusStore, obj *models.AIToolCallLog, logger zerolog.Logger, ctxClaim map[string]string) error {
	obj.PreSaveCreate(ctxClaim)
	_, err := dmstore.Db.PostgresClient.DB.NewInsert().
		Model(obj).
		ModelTableExpr(apppkgs.AIToolCallLogTable).Returning("id").Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving fabric chat tool log in datastore")
		return err
	}
	upQ := fmt.Sprintf("UPDATE %s SET input = to_tsvector('english', plain_input) WHERE id = %d", apppkgs.AIToolCallLogTable, obj.ID)
	_, err = dmstore.Db.PostgresClient.DB.NewRaw(upQ).Exec(ctx)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while saving data source symantics and FTS vectors in datastore")
		return err
	}
	return nil
}

func GetAIToolCallLog(ctx context.Context, dmstore *apppkgs.VapusStore, input string, actionAnalyzer, paramAnalyzer bool, logger zerolog.Logger, ctxClaim map[string]string) (*models.AIToolCallLog, error) {
	type localVal struct {
		rank          float64 `bun:"rank"`
		output_schema string  `bun:"output_schema"`
	}
	result := []map[string]any{}
	fCon := fmt.Sprintf("param_analyzer=%t AND action_analyzer=%t", paramAnalyzer, actionAnalyzer)
	query := fmt.Sprintf(`
	SELECT rank,output_schema
	FROM (
		SELECT output_schema,ts_rank(input, plainto_tsquery('english', '%v')) AS rank
		FROM %s
		WHERE %s
	) subquery
	WHERE rank > 1.0
	ORDER BY rank DESC
	LIMIT 1;
	`,
		input, apppkgs.AIToolCallLogTable, fCon)
	err := dmstore.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		logger.Err(err).Ctx(ctx).Msg("error while getting fabric chat tool log from datastore")
		return nil, err
	}
	return &models.AIToolCallLog{
		OutputSchema: result[0]["output_schema"].(string),
	}, err
}
