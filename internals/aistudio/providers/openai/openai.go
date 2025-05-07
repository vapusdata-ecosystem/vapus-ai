package openaistd

import (
	"context"

	"github.com/rs/zerolog"
	openai "github.com/sashabaranov/go-openai"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/providers/generic"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

var (
	OpenAILLMMap = map[string]string{
		"gpt-4o":      openai.GPT4o,
		"gpt-4":       openai.GPT4,
		"gpt-4o-mini": openai.GPT4oMini,
	}
	OpenAIEmbeddingMap = map[string]openai.EmbeddingModel{
		"text-embedding-3-large":  openai.LargeEmbedding3,
		"text-embedding-3-small":  openai.SmallEmbedding3,
		"ttext-embedding-3-small": openai.AdaEmbeddingV2,
	}
)

var ResponseFormatMap = map[string]openai.ChatCompletionResponseFormatType{
	mpb.AIResponseFormat_TEXT.String():        openai.ChatCompletionResponseFormatTypeText,
	mpb.AIResponseFormat_JSON_SCHEMA.String(): openai.ChatCompletionResponseFormatTypeJSONObject,
}

type OpenAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type OpenAI struct {
	*generic.OpenAI
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (OpenAIInterface, error) {
	// if node.NetworkParams.Url == "" {
	// node.NetworkParams.Url = "https://api.openai.com/api/v1"
	// }
	client, err := generic.New(ctx, node, retries, logger)
	if err != nil {
		return nil, err
	}

	return &OpenAI{
		OpenAI: client.(*generic.OpenAI),
	}, nil
}

func (o *OpenAI) GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil {
		return o.OpenAI.GenerateContent(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *OpenAI) GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil {
		if request.Params.Stream {
			request.Params.StreamOptions = &pb.StreamOptions{
				IncludeUsage: true,
			}
		}
		return o.OpenAI.GenerateContentStream(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *OpenAI) GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil {
		return o.OpenAI.GenerateEmbeddings(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *OpenAI) CrawlModels(ctx context.Context) ([]*models.AIModelBase, error) {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil {
		return o.OpenAI.CrawlModels(ctx)
	}
	return nil, apperr.ErrInvalidOrMissingPodelAPIKey
}
