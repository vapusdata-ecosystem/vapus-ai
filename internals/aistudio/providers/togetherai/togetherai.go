package togetherai

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/rs/zerolog"
	aicore "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/providers/generic"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	httpCls "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/http"
)

var (
	OpenAIEmbeddingMap = map[string]openai.EmbeddingModel{
		"text-embedding-3-large":  openai.EmbeddingModelTextEmbedding3Large,
		"text-embedding-3-small":  openai.EmbeddingModelTextEmbedding3Small,
		"ttext-embedding-3-small": openai.EmbeddingModelTextEmbeddingAda002,
	}

	// OpenAITimestampGranularityMap = map[string]openai.AudioTranscriptionNewParamsTimestampGranularity{
	// 	"word":    openai.AudioTranscriptionNewParamsTimestampGranularityWord,
	// 	"segment": openai.AudioTranscriptionNewParamsTimestampGranularitySegment,
	// }
)

// var ResponseFormatMap = map[string]openai.ChatCompletionNewParamsResponseFormatType{
// 	mpb.AIResponseFormat_TEXT.String():        openai.ChatCompletionNewParamsResponseFormatTypeText,
// 	mpb.AIResponseFormat_JSON_SCHEMA.String(): openai.ChatCompletionNewParamsResponseFormatTypeJSONObject,
// }

type TogetherAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error
	GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type TogetherAI struct {
	*generic.OpenAI
	httpClient *httpCls.RestHttp
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (TogetherAIInterface, error) {

	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = "https://api.together.xyz/v1" // wokring
	}

	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}

	httpCl, _ := httpCls.New(logger,
		httpCls.WithAddress("https://api.together.xyz"),
		httpCls.WithBasePath("/v1"),
		httpCls.WithBearerAuth(token),
	)
	client, err := generic.New(ctx, node, retries, logger)
	if err != nil {
		return nil, err
	}

	return &TogetherAI{
		OpenAI:     client.(*generic.OpenAI),
		httpClient: httpCl,
	}, nil
}

func (o *TogetherAI) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload) error {

	o.OpenAI.Log.Info().Msgf("Generating embeddings from openai with model %s", payload.EmbeddingModel)
	reqObj := openai.EmbeddingNewParams{
		Model:      payload.EmbeddingModel,
		User:       param.NewOpt(string(aicore.USER)),
		Dimensions: param.NewOpt(int64(payload.Dimensions)),
	}
	if len(payload.InputArray) > 0 {
		reqObj.Input = openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: payload.InputArray,
		}
	} else {
		reqObj.Input = openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(payload.Input),
		}
	}
	resp, err := o.OpenAI.Client.Embeddings.New(
		ctx,
		reqObj,
	)
	if err != nil {
		o.OpenAI.Log.Err(err).Msg("error while generating embeddings from openai")
		return err
	}
	float32Slice := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		float32Slice[i] = float32(v)
	}
	payload.Embeddings = &models.VectorEmbeddings{
		Vectors32: float32Slice,
	}
	return nil
}

func (o *TogetherAI) GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContent(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *TogetherAI) GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContentStream(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *TogetherAI) GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranscription(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *TogetherAI) GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranslation(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}
