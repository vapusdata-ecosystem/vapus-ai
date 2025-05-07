package bedrock

import (
	"context"
	"os"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (o *Bedrock) ListInferenceProfiles(ctx context.Context) (result []*models.AIModelBase, err error) {
	modelsIter, _ := o.bedrockService.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})

	for _, model := range modelsIter.ModelSummaries {
		obj := &models.AIModelBase{
			ModelId:   *model.ModelId,
			OwnedBy:   *model.ProviderName,
			ModelName: *model.ModelId,
		}
		result = append(result, obj)
	}
	return result, nil
}

func (o *Bedrock) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {
	// Unable to fetch the Active models List even tried different methods
	modelsIter, err := o.bedrockService.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})
	if err != nil {
		logger := zerolog.New(os.Stdout)
		logger.Error().Err(err).Msg("Error while listing the bedrock models")
		return nil, err
	}

	for _, model := range modelsIter.ModelSummaries {
		obj := &models.AIModelBase{
			ModelId:     *model.ModelId,
			ModelName:   *model.ModelName,
			ModelArn:    *model.ModelArn,
			OwnedBy:     *model.ProviderName,
			ModelNature: []string{},
		}
		if slices.Contains(model.InputModalities, "EMBEDDING") {
			obj.ModelType = mpb.AIModelType_EMBEDDING.String()
			obj.ModelNature = append(obj.ModelNature, mpb.AIModelType_EMBEDDING.String())
		} else if slices.Contains(model.InputModalities, "IMAGE") {
			obj.ModelType = mpb.AIModelType_VISION.String()
			obj.ModelNature = append(obj.ModelNature, mpb.AIModelType_VISION.String())
		} else {
			obj.ModelType = mpb.AIModelType_LLM.String()
			obj.ModelNature = append(obj.ModelNature, mpb.AIModelType_LLM.String())
		}
		result = append(result, obj)
	}

	return result, nil
}
