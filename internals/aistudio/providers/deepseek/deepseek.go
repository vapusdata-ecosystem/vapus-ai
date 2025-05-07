package deepseek

import (
	"context"
	"log"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/generic"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	httpCls "github.com/vapusdata-ecosystem/vapusai/core/pkgs/http"
)

type DeepseekInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error // No Embedding models for deepseek till now
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error
	GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type Deepseek struct {
	httpClient *httpCls.RestHttp
	*generic.OpenAI
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (DeepseekInterface, error) {
	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = "https://api.deepseek.com" // its working
	}
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}

	httpCl, _ := httpCls.New(logger,
		httpCls.WithAddress("https://api.deepseek.com"),
		httpCls.WithBasePath(""),
		httpCls.WithBearerAuth(token),
	)
	client, err := generic.New(ctx, node, retries, logger)
	if err != nil {
		return nil, err
	}

	return &Deepseek{
		OpenAI:     client.(*generic.OpenAI),
		httpClient: httpCl,
	}, nil
}

func (o *Deepseek) GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	log.Println("Deepseek GenerateContent.............................")
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContent(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Deepseek) GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	log.Println("Deepseek GenerateContent.............................222222222")
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContentStream(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Deepseek) GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateEmbeddings(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Deepseek) GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranscription(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Deepseek) GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranslation(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}
