package grok

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/providers/generic"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	httpCls "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/http"
)

type GrokInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error // No Embedding models for Grok till now
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error
	GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type Grok struct {
	httpClient *httpCls.RestHttp
	*generic.OpenAI
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (GrokInterface, error) {
	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = "https://api.x.ai" // its working
	}
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}

	httpCl, _ := httpCls.New(logger,
		httpCls.WithAddress("https://api.x.ai"),
		httpCls.WithBasePath("/v1"),
		httpCls.WithBearerAuth(token),
	)
	client, err := generic.New(ctx, node, retries, logger)
	if err != nil {
		return nil, err
	}

	return &Grok{
		OpenAI:     client.(*generic.OpenAI),
		httpClient: httpCl,
	}, nil
}

func (o *Grok) GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContent(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Grok) GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContentStream(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Grok) GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateEmbeddings(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Grok) GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranscription(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Grok) GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranslation(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}
