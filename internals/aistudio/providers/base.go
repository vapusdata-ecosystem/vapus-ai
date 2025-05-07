package ai

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	anthropic "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/anthropic"
	bedrock "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/bedrock"
	deepseek "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/deepseek"
	generic "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/generic"
	google "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/google"
	grok "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/grok"
	groq "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/groq"
	mistral "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/mistral"
	openai "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/openai"
	perplexity "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/perplexity"
	together "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/togetherai"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
)

const (
	defaultRetries = 3
)

var AILogger zerolog.Logger

type AIModelNodeInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type AIModelNodeClient struct {
	node   *models.AIModelNode
	logger zerolog.Logger
}

type AiModelOpts func(*AIModelNodeClient)

func WithAIModelNode(c *models.AIModelNode) AiModelOpts {
	return func(opts *AIModelNodeClient) {
		opts.node = c
	}
}

func WithLogger(logger zerolog.Logger) AiModelOpts {
	return func(opts *AIModelNodeClient) {
		opts.logger = logger
	}
}

func NewAIModelNode(opts ...AiModelOpts) (AIModelNodeInterface, error) {
	configurator := &AIModelNodeClient{}
	for _, opt := range opts {
		opt(configurator)
	}
	AILogger = dmlogger.GetSubDMLogger(configurator.logger, "ailogger", "base")
	switch configurator.node.ServiceProvider {
	case mpb.ServiceProvider_OPENAI.String():
		return openai.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_ANTHROPIC.String():
		return anthropic.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_MISTRAL.String():
		return mistral.New(configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_GEMINI.String():
		return google.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_TOGETHER.String():
		return together.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_GROQ.String():
		return groq.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_GENERIC.String():
		return generic.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_BEDROCK.String():
		return bedrock.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_DEEPSEEK.String():
		return deepseek.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_PERPLEXITY.String():
		return perplexity.New(configurator.node, defaultRetries, AILogger)
	case mpb.ServiceProvider_GROK.String():
		return grok.New(context.TODO(), configurator.node, defaultRetries, AILogger)
	default:
		configurator.logger.Error().Msgf("Unknown service provider: %s", configurator.node.ServiceProvider)
		return nil, aicore.ErrUnknownServiceProvider
	}
}
