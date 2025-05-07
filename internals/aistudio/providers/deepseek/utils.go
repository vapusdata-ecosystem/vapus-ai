package deepseek

import (
	"context"
	"strings"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (o *Deepseek) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {
	var res struct {
		Data []map[string]interface{} `json:"data"`
	}

	params := make(map[string]string)

	err = o.httpClient.Get(ctx, "/models", params, &res, "application/json")
	if err != nil {
		o.OpenAI.Log.Err(err).Msg("error while getting models from openai")
		return nil, err
	}
	for _, model := range res.Data {

		modelID, _ := model["id"].(string)
		ownedBy, _ := model["organization"].(string)
		modelType, _ := model["type"].(string)
		displayName, _ := model["display_name"].(string)

		if strings.Contains(modelID, "embed") {
			result = append(result, &models.AIModelBase{
				ModelId:   modelID,
				OwnedBy:   ownedBy,
				ModelType: modelType,
				ModelName: displayName,
				ModelNature: []string{
					mpb.AIModelType_EMBEDDING.String(),
				},
			})
		} else {
			result = append(result, &models.AIModelBase{
				ModelId:   modelID,
				OwnedBy:   ownedBy,
				ModelType: mpb.AIModelType_LLM.String(),
				ModelName: modelID,
				ModelNature: []string{
					mpb.AIModelType_LLM.String(),
				},
			})
		}
	}
	return result, nil
}
