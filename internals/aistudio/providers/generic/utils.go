package generic

import (
	"context"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

func (o *OpenAI) CrawlModels(ctx context.Context) (result []*models.AIModelBase, err error) {
	openaiModels, err := o.Client.Models.List(ctx)
	if err != nil {
		o.Log.Err(err).Msg("error while getting models from openai")
		return nil, err
	}

	for {
		for _, model := range openaiModels.Data {
			if strings.Contains(model.ID, "embed") {
				result = append(result, &models.AIModelBase{
					ModelId:   model.ID,
					OwnedBy:   model.OwnedBy,
					ModelType: mpb.AIModelType_EMBEDDING.String(),
					ModelName: model.ID,
					ModelNature: []string{
						mpb.AIModelType_EMBEDDING.String(),
					},
				})
			} else {
				result = append(result, &models.AIModelBase{
					ModelId:   model.ID,
					OwnedBy:   model.OwnedBy,
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

func ConvertToCompletionUserPart[T string | []openai.ChatCompletionContentPartUnionParam](obj *pb.ChatMessageObject) (response T) {
	if obj == nil {
		return response
	}
	if obj.Content != "" {
		response = any(obj.Content).(T)
		return response
	}
	if obj.StructuredContent != nil {
		res := []openai.ChatCompletionContentPartUnionParam{}
		for _, value := range obj.StructuredContent {
			switch value.Type {
			case aicore.AIResponseFormatText.String():
				res = append(res, openai.TextContentPart(value.Text))
			case aicore.AIResponseFormatImageUrl.String():
				res = append(res, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL:    value.ImageUrl.Url,
					Detail: value.ImageUrl.Detail,
				}))
			case aicore.AIResponseFormatInputAudio.String():
				res = append(res, openai.InputAudioContentPart(openai.ChatCompletionInputAudioDataParam{
					Data:   value.InputAudio.Data,
					Format: strings.ToLower(value.InputAudio.Format.String()),
				}))
			case aicore.AIResponseFormatInputFile.String():
				res = append(res, openai.FileContentPart(openai.ChatCompletionFileContentPartParam{
					FileData: param.NewOpt(value.File.FileData),
					FileID:   param.NewOpt(value.File.FileId),
					Filename: param.NewOpt(value.File.Filename),
				}))
			default:
				res = append(res, openai.TextContentPart(value.Text))
			}
		}
		response = any(res).(T)
		return response
	}

	return response
}

func ConvertToCompletionAssistantPart[T string | []openai.ChatCompletionAssistantMessagePartUnion](obj *pb.ChatMessageObject) (response T) {
	if obj == nil {
		return response
	}
	if obj.Content != "" {
		response = any(obj.Content).(T)
		return response
	}
	if obj.StructuredContent != nil {
		res := []openai.ChatCompletionAssistantMessagePartUnion{}
		for _, value := range obj.StructuredContent {
			switch value.Type {
			case aicore.AIResponseFormatText.String():
				res = append(res, openai.ChatCompletionAssistantMessagePartUnion{
					OfText: &openai.ChatCompletionContentPartTextParam{
						Text: value.Text,
					},
				})
			case aicore.AIResponseFormatRefusal.String():
				res = append(res, openai.ChatCompletionAssistantMessagePartUnion{
					OfRefusal: &openai.ChatCompletionContentPartRefusalParam{
						Refusal: value.Text,
					},
				})
			default:
				res = append(res, openai.ChatCompletionAssistantMessagePartUnion{
					OfText: &openai.ChatCompletionContentPartTextParam{
						Text: value.Text,
					},
				})
			}
		}
		response = any(res).(T)
		return response
	}

	return response
}

func ConvertToCompletionSystemPart[T string | []openai.ChatCompletionContentPartTextParam](obj *pb.ChatMessageObject) (response T) {
	if obj == nil {
		return response
	}
	if obj.Content != "" {
		response = any(obj.Content).(T)
		return response
	}
	if obj.StructuredContent != nil {
		res := []openai.ChatCompletionContentPartTextParam{}
		for _, value := range obj.StructuredContent {
			switch value.Type {
			case aicore.AIResponseFormatText.String():
				res = append(res, openai.ChatCompletionContentPartTextParam{
					Text: value.Text,
				})
			default:
				res = append(res, openai.ChatCompletionContentPartTextParam{
					Text: value.Text,
				})
			}
		}
		response = any(res).(T)
		return response
	}

	return response
}
