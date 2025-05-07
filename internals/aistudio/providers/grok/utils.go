package grok

import (
	"context"
	"fmt"
	"strings"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

func (o *Grok) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {
	var res struct {
		Models []map[string]interface{} `json:"models"`
	}

	params := make(map[string]string)

	err = o.httpClient.Get(ctx, "/language-models", params, &res, "application/json")
	if err != nil {
		o.OpenAI.Log.Err(err).Msg("error while getting models from openai")
		return nil, err
	}
	fmt.Println("Result is fetched Successfully")
	for _, model := range res.Models {

		modelID, _ := model["id"].(string)
		ownedBy, _ := model["organization"].(string)
		modelInputType, _ := model["input_modalities"].(string)
		displayName, _ := model["id"].(string)

		if strings.Contains(modelInputType, "image") {
			result = append(result, &models.AIModelBase{
				ModelId:   modelID,
				OwnedBy:   ownedBy,
				ModelType: mpb.AIModelType_LLM.String(),
				ModelName: displayName,
				ModelNature: []string{
					mpb.AIModelType_EMBEDDING.String(),
					mpb.AIModelType_LLM.String(),
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
