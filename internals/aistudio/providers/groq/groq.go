package groq

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/generic"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	httpCls "github.com/vapusdata-ecosystem/vapusai/core/pkgs/http"
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

type GroqAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error
	GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type GroqAI struct {
	*generic.OpenAI
	httpClient *httpCls.RestHttp
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (GroqAIInterface, error) {
	// its Wokring
	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = "https://api.groq.com/openai/v1"
	}
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}

	httpCl, _ := httpCls.New(logger,
		httpCls.WithAddress("https://api.groq.com"),
		httpCls.WithBasePath("/openai/v1"),
		httpCls.WithBearerAuth(token),
	)

	client, err := generic.New(ctx, node, retries, logger)
	if err != nil {
		return nil, err
	}

	return &GroqAI{
		OpenAI:     client.(*generic.OpenAI),
		httpClient: httpCl,
	}, nil
}

func (o *GroqAI) GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error {
	o.OpenAI.Log.Info().Msgf("Generating transcription from openai with model %s", payload.Model)

	// timestampGranularities := []openai.AudioTranscriptionNewParamsTimestampGranularity{}
	// for _, val := range payload.TimestampGranularities {

	// 	temp, ok := OpenAITimestampGranularityMap[val]
	// 	if !ok {
	// 		timestampGranularities = append(timestampGranularities, openai.AudioTranscriptionNewParamsTimestampGranularitySegment)
	// 	} else {
	// 		timestampGranularities = append(timestampGranularities, temp)
	// 	}
	// }

	// resp, err := o.OpenAI.Client.Audio.Transcriptions.New(
	// 	ctx,
	// 	openai.AudioTranscriptionNewParams{
	// 		File:                   payload.File,
	// 		Model:                  payload.Model,
	// 		Language:               param.NewOpt(payload.Language),
	// 		Prompt:                 openai.F(payload.Prompt),
	// 		ResponseFormat:         openai.F(openai.AudioResponseFormat(payload.ResponseFormat)),
	// 		Temperature:            param.NewOpt(payload.Temperature),
	// 		TimestampGranularities: openai.F(timestampGranularities),
	// 	},
	// )

	// if err != nil {
	// 	o.OpenAI.Log.Err(err).Msg("error while generating audio transcriptions from openai")
	// 	return err
	// }
	// payload.ResponseText = resp.Text
	// payload.ResponseJson = resp.JSON.RawJSON()

	return nil

}

func (o *GroqAI) GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error {
	o.OpenAI.Log.Info().Msgf("Generating translation from openai with model %s", payload.Model)

	// resp, err := o.OpenAI.Client.Audio.Translations.New(
	// 	ctx,
	// 	openai.AudioTranslationNewParams{
	// 		File:           openai.F[io.Reader](payload.File),
	// 		Model:          openai.F(payload.Model),
	// 		Prompt:         openai.F(payload.Prompt),
	// 		ResponseFormat: openai.F(openai.AudioResponseFormat(payload.ResponseFormat)),
	// 		Temperature:    openai.F(payload.Temperature),
	// 	},
	// )

	// if err != nil {
	// 	o.OpenAI.Log.Err(err).Msg("error while generating audio transcriptions from openai")
	// 	return err
	// }
	// payload.ResponseText = resp.Text
	// payload.ResponseJson = resp.JSON.RawJSON()

	return nil

}

func (o *GroqAI) GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		// o.OpenAI.ModelNode.NetworkParams.Url = "https://api.groq.com/openai/v1"
		return o.OpenAI.GenerateContent(ctx, request)
	}
	// token := o.OpenAI.ModelNode.NetworkParams.Credentials.ApiToken // nil check is pending
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *GroqAI) GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		// o.OpenAI.ModelNode.NetworkParams.Url = "https://api.groq.com/openai/v1"
		return o.OpenAI.GenerateContentStream(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *GroqAI) GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateEmbeddings(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}
