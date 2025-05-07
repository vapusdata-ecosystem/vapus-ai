package vertex

import (
	"context"

	"cloud.google.com/go/vertexai/genai"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
)

// var defaultModel = "gemini-1.5-flash-001"

type VertexAIInterface interface {
	// GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	// GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	// GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	// CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type VertexAI struct {
	// bedrockService *bedrockService.Client
	// *googlegenai.GoogleGenAI
	client    *genai.Client
	log       zerolog.Logger
	modelNode *models.AIModelNode
	// maxRetries int
	// params     map[string]any
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (VertexAIInterface, error) {
	location := node.NetworkParams.Credentials.GcpCreds.Region
	projectID := node.NetworkParams.Credentials.GcpCreds.ProjectId

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating google vertex ai client")
	}

	return &VertexAI{
		client:    client,
		modelNode: node,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Vertex AI"),
	}, nil
}
