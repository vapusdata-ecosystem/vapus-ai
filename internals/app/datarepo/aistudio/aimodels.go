package aidmstore

import (
	"context"
	"fmt"
	"log"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	llm_observability "github.com/vapusdata-ecosystem/vapusdata/core/observability/llm"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
)

func (ds *AIStudioDMStore) ConfigureGetAIModelNode(ctx context.Context, obj *models.AIModelNode, ctxClaim map[string]string) error {
	_, err := ds.Db.PostgresClient.DB.NewInsert().Model(obj).ModelTableExpr(apppkgs.AIModelsNodesTable).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai model node metadata in datastore")
		return err
	}
	go func() {
		mCtx := context.TODO()
		_ = apppkgs.AddResourceArn(mCtx, ds.Db, &models.VapusResourceArn{
			ResourceId:   obj.VapusID,
			ResourceName: mpb.Resources_AIMODELS.String(),
			VapusBase: models.VapusBase{
				Editors: obj.Editors,
			},
		}, ds.logger, ctxClaim)
	}()
	return nil
}

func (ds *AIStudioDMStore) SaveAIInterfaceLog(ctx context.Context, obj *models.AIStudioLog, usageLogs *models.AIStudioUsages, ctxClaim map[string]string) error {
	obj.PreSaveCreate(ctxClaim)
	_, err := ds.Db.PostgresClient.DB.NewInsert().
		Model(obj).
		ModelTableExpr(apppkgs.AIStudioLogsTable).
		Returning("id").
		Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai model studio log in datastore")
		return err
	}

	if usageLogs != nil {
		usageLogs.PreSaveCreate(ctxClaim)
		usageLogs.AIStudioLogId = obj.VapusID
		_, err = ds.Db.PostgresClient.DB.NewInsert().Model(usageLogs).ModelTableExpr(apppkgs.AIStudioUsagesTable).Exec(ctx)
		if err != nil {
			ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai studio usage log in datastore")
			return err
		}
	}
	return nil
}

func (ds *AIStudioDMStore) SaveAIGuardrailLog(ctx context.Context, obj *models.AIGuardrailsLog, usageLogs *models.AIStudioUsages, ctxClaim map[string]string) error {
	obj.PreSaveCreate(ctxClaim)
	_, err := ds.Db.PostgresClient.DB.NewInsert().
		Model(obj).
		ModelTableExpr(apppkgs.AIGuardrailsLogsTable).
		Returning("id").
		Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai guardrail log in datastore")
		return err
	}

	if usageLogs != nil {
		usageLogs.PreSaveCreate(ctxClaim)
		usageLogs.GuardrailLogId = obj.VapusID
		_, err = ds.Db.PostgresClient.DB.NewInsert().Model(usageLogs).ModelTableExpr(apppkgs.AIStudioUsagesTable).Exec(ctx)
		if err != nil {
			ds.logger.Err(err).Ctx(ctx).Msg("error while saving ai guardrail log in datastore")
			return err
		}
	}
	return nil
}

func (ds *AIStudioDMStore) ListAIInterfaceLogByUser(ctx context.Context, limit int, ctxClaim map[string]string) ([]*models.AIStudioLog, error) {
	resp := []*models.AIStudioLog{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE created_by='%s' AND Organization='%s' AND owner_account='%s' ORDER BY id DESC LIMIT %v",
		apppkgs.AIStudioLogsTable,
		ctxClaim[encryption.ClaimUserIdKey],
		ctxClaim[encryption.ClaimOrganizationKey],
		ctxClaim[encryption.ClaimAccountKey],
		limit)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &resp)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while fetching ai model studio logs from datastore")
		return nil, err
	}
	return resp, nil
}

func (ds *AIStudioDMStore) PutAIModelNode(ctx context.Context, obj *models.AIModelNode) error {
	_, err := ds.Db.PostgresClient.DB.NewUpdate().Model(obj).ModelTableExpr(apppkgs.AIModelsNodesTable).Where(apppkgs.VapusIdFilter(), obj.VapusID).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while updating ai model node metadata in datastore")
		return err
	}
	return nil
}

func (ds *AIStudioDMStore) ListAIModelNodes(ctx context.Context, condition string, ctxClaim map[string]string) ([]*models.AIModelNode, error) {
	if ctxClaim == nil {
		condition = apppkgs.GetAccountFilter(ctxClaim, condition)
	}
	if condition == "" {
		condition = "deleted_at IS NULL AND status = 'ACTIVE'"
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.AIModelsNodesTable, condition)
	result := []*models.AIModelNode{}
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	// _, err = ds.Db.PostgresClient.DB.NewSelect().Model(&result).ModelTableExpr(vapussvc.AIModelsNodesTable).Where(condition).Exec(ctx)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while fetching ai model nodes from datastore")
		return nil, err
	}
	return result, nil
}

func (ds *AIStudioDMStore) GetAIModelNode(ctx context.Context, iden string, ctxClaim map[string]string) (*models.AIModelNode, error) {
	if iden == "" {
		return nil, apperr.ErrAIModelNode404
	}
	result := []*models.AIModelNode{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", apppkgs.AIModelsNodesTable, apppkgs.GetByIdFilter("", iden, ctxClaim))
	// err := ds.Db.Select(ctx, &datasvcpkgs.QueryOpts{
	// 	DataCollection: vapussvc.DataSourcesTable,
	// 	RawQuery:       query,
	// }, &result, ds.logger)
	err := ds.Db.PostgresClient.SelectInApp(ctx, &query, &result)
	if err != nil || len(result) == 0 {
		ds.logger.Err(err).Ctx(ctx).Msg("error while getting ai models from datastore")
		return nil, err
	}
	return result[0], err
}

// Insights/Obeservability

func (ds *AIStudioDMStore) GetAIModelNodeId(ctx context.Context, iden string, ctxClaim map[string]string) ([]*string, error) {
	if iden == "" {
		return nil, apperr.ErrAIModelNode404
	}

	query := fmt.Sprintf(`SELECT model_node FROM ai_studio_logs
			WHERE model_node IN 
			(SELECT vapus_id FROM ai_models 
			WHERE created_by = '%s' AND Organization = '%s' 
			AND owner_account = '%s' GROUP BY vapus_id) 
			GROUP BY model_node;`, ctxClaim[encryption.ClaimUserIdKey],
		ctxClaim[encryption.ClaimOrganizationKey],
		ctxClaim[encryption.ClaimAccountKey])

	rows, err := ds.Db.PostgresClient.Conn.Query(query)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error while getting ai model node IDs")
		return nil, err
	}
	defer rows.Close()

	var modelNodes []*string
	if rows != nil {
		for rows.Next() {
			var val string
			if err := rows.Scan(&val); err != nil {
				log.Fatalf("Error scanning row: %v", err)
			}
			modelNodes = append(modelNodes, &val)
		}
	}
	return modelNodes, nil

}

func (ds *AIStudioDMStore) GetAIModelUsageInsights(ctx context.Context, iden string, modelNode *string, modelName []string, ctxClaim map[string]string) (*mpb.ModelNodeObservability, error) {
	type LocalSt struct {
		ModelNodeID        string  `sql:"column:model_node"`
		ModelName          string  `sql:"column:ai_model"`
		ModelProvider      string  `sql:"column:model_provider"`
		Requests           int32   `sql:"column:requests"`
		InputTokens        int32   `sql:"column:input_tokens"`
		OutputTokens       int32   `sql:"column:output_tokens"`
		InputCachedTokens  int32   `sql:"column:input_cached_tokens"`
		OutputCachedTokens int32   `sql:"column:output_cached_tokens"`
		InputAudioTokens   int32   `sql:"column:input_audio_tokens"`
		OutputAudioTokens  int32   `sql:"column:output_audio_tokens"`
		AverageTokens      float32 `sql:"column:avg_tokens"`
	}

	if iden == "" {
		return nil, apperr.ErrAIModelNode404
	}
	query := fmt.Sprintf(`SELECT
    l.model_node,
    l.ai_model,
    l.model_provider,
    COUNT(l.ai_model) AS requests,
    SUM(u.input_tokens) AS input_tokens,
    SUM(u.output_tokens) AS output_tokens,
    SUM(u.input_cached_tokens) AS input_cached_tokens,
    SUM(u.output_cached_tokens) AS output_cached_tokens,
    SUM(u.input_audio_tokens) AS input_audio_tokens,
    SUM(u.output_audio_tokens) AS output_audio_tokens,
    SUM(u.total_tokens) / COUNT(l.ai_model) AS avg_tokens
	FROM public.ai_studio_logs l
	JOIN public.ai_studio_usages u
	ON l.vapus_id = u.ai_studio_log_id
	WHERE l.model_node = '%v'
	GROUP BY l.model_node, l.ai_model, l.model_provider;`, *modelNode)

	rows, err := ds.Db.PostgresClient.Conn.QueryContext(ctx, query)
	if err != nil {
		ds.logger.Err(err).Ctx(ctx).Msg("error executing query")
		return nil, err
	}
	defer rows.Close()

	var vals []LocalSt
	if rows != nil {
		defer rows.Close()

		for rows.Next() {
			var val LocalSt
			err := rows.Scan(
				&val.ModelNodeID,
				&val.ModelName,
				&val.ModelProvider,
				&val.Requests,
				&val.InputTokens,
				&val.OutputTokens,
				&val.InputCachedTokens,
				&val.OutputCachedTokens,
				&val.InputAudioTokens,
				&val.OutputAudioTokens,
				&val.AverageTokens,
			)
			if err != nil {
				ds.logger.Err(err).Ctx(ctx).Msg("error scanning row")
				return nil, err
			}
			vals = append(vals, val)
		}

		if len(vals) == 0 {
			ds.logger.Info().Ctx(ctx).Msg("No records found for model node")
			return nil, nil
		}
	}

	var modelObservabilities []*mpb.ModelObservability
	if len(modelName) <= 0 {
		for _, val := range vals {
			modelObs := &mpb.ModelObservability{
				ModelNodeId:            val.ModelNodeID,
				ModelName:              val.ModelName,
				ModelProvider:          val.ModelProvider,
				Request:                val.Requests,
				InputTokens:            val.InputTokens,
				OutputTokens:           val.OutputTokens,
				InputCachedTokens:      val.InputCachedTokens,
				OutputCachedTokens:     val.OutputCachedTokens,
				InputAudioTokens:       val.InputAudioTokens,
				OutputAudioTokens:      val.OutputAudioTokens,
				AverageTokenPerRequest: val.AverageTokens,
			}

			modelObservabilities = append(modelObservabilities, modelObs)
		}
	} else {
		set := make(map[string]bool)
		for _, temp := range modelName {
			set[temp] = true
		}
		for _, val := range vals {
			if set[val.ModelName] {
				modelObs := &mpb.ModelObservability{
					ModelNodeId:            val.ModelNodeID,
					ModelName:              val.ModelName,
					ModelProvider:          val.ModelProvider,
					Request:                val.Requests,
					InputTokens:            val.InputTokens,
					OutputTokens:           val.OutputTokens,
					InputCachedTokens:      val.InputCachedTokens,
					OutputCachedTokens:     val.OutputCachedTokens,
					InputAudioTokens:       val.InputAudioTokens,
					OutputAudioTokens:      val.OutputAudioTokens,
					AverageTokenPerRequest: val.AverageTokens,
				}
				modelObservabilities = append(modelObservabilities, modelObs)
			}
		}
	}
	// Compute cost-related fields
	if modelObservabilities != nil {
		modelNodeObservability := &mpb.ModelNodeObservability{}
		for _, modelObservability := range modelObservabilities {
			priceStore, err := llm_observability.NewLLMPricingStore(ctx, ds.VapusStore, modelObservability, ds.logger)
			if err != nil {
				ds.logger.Err(err).Ctx(ctx).Msg("error fetching pricing info")
				return nil, err
			}

			amount := priceStore.Calculate(ctx, modelObservability, ds.logger)

			modelObservability.Cost = float32(amount.Charges)
			if modelObservability.Request > 0 {
				modelObservability.AverageCostPerRequest = modelObservability.Cost / float32(modelObservability.Request)
			}

			modelNodeObservability.TotalCost += modelObservability.Cost
			modelNodeObservability.TotalRequests += int64(modelObservability.Request)
		}

		// Assign values to modelNodeObservability
		modelNodeObservability.ModelNodeId = modelObservabilities[0].ModelNodeId
		modelNodeObservability.ModelProvider = modelObservabilities[0].ModelProvider
		modelNodeObservability.ModelObservability = modelObservabilities

		return modelNodeObservability, nil
	}
	return nil, nil

}
