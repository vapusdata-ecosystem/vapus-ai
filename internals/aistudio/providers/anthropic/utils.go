package anthropic

import (
	"context"
	"encoding/base64"
	"log"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
)

func (o *AnthropicAI) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {

	openaiModels, err := o.Client.Models.List(ctx, anthropic.ModelListParams{})
	if err != nil {
		o.Log.Err(err).Msg("error while getting models from openai")
		return nil, err
	}

	for {
		for _, model := range openaiModels.Data {
			if strings.Contains(model.ID, "embed") {
				result = append(result, &models.AIModelBase{
					ModelId:   model.ID,
					OwnedBy:   "anthropic",
					ModelType: mpb.AIModelType_EMBEDDING.String(),
					ModelName: model.ID,
					ModelNature: []string{
						mpb.AIModelType_EMBEDDING.String(),
					},
				})
			} else {
				result = append(result, &models.AIModelBase{
					ModelId:   model.ID,
					OwnedBy:   "anthropic",
					ModelType: mpb.AIModelType_LLM.String(),
					ModelName: model.ID,
					ModelNature: []string{
						mpb.AIModelType_LLM.String(),
					},
				})
			}
		}

		openaiModels, err := openaiModels.GetNextPage()

		if err != nil {
			o.Log.Err(err).Msg("error while paginating models")
			return nil, err
		}

		if openaiModels == nil {
			break
		}
	}
	return result, nil
}

func BuildContentBlock(obj *pb.ChatMessageObject) []anthropic.ContentBlockParamUnion {
	response := make([]anthropic.ContentBlockParamUnion, 0)
	if obj.Content != "" {
		response = append(response, anthropic.NewTextBlock(obj.Content))
		return response
	}

	if obj.StructuredContent != nil {
		for _, value := range obj.StructuredContent {
			switch value.Type {
			case aicore.AIResponseFormatText.String():
				response = append(response, anthropic.NewTextBlock(value.Text))
			case aicore.AIResponseFormatImageUrl.String():
				isBase64Encoded := dmutils.IsBase64Encoded(value.ImageUrl.Data)
				mimeType, ok := filetools.FileMimeMap[value.ImageUrl.Format]
				if !ok {
					log.Printf("Failed to get format for %s: ", value.ImageUrl.Format)
					continue
				}
				if isBase64Encoded {
					response = append(response, anthropic.NewImageBlockBase64(mimeType[0], value.ImageUrl.Data))
				} else {
					bytes, err := base64.StdEncoding.DecodeString(value.ImageUrl.Data)
					if err != nil {
						log.Println("Error decoding base64 string:", err)
						return nil
					}
					response = append(response, anthropic.NewImageBlockBase64(mimeType[0], string(bytes)))
				}
			case aicore.AIResponseFormatInputFile.String():
				log.Println("Input file format not supported")
			default:
				response = append(response, anthropic.NewTextBlock(value.Text))
			}
		}
		return response
	}
	return response
}
