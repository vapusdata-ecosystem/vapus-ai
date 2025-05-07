package appcl

import (
	"context"
	"log"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (x *VapusSvcInternalClients) Chat(ctx context.Context, req *aipb.ChatRequest, logger zerolog.Logger, retryCount int) (*aipb.ChatResponse, error) {
	if x == nil {
		return nil, ErrAIStudioConnNotInitialized
	}
	if x.AIStudioConn == nil {
		err := x.PingTestAndReconnect(ctx, x.AIStudioDns, logger)
		if err != nil {
			return nil, ErrAIStudioConnNotInitialized
		}
	}
	resp, err := x.AIStudioConn.Completions(pbtools.SwapNewContextWithAuthToken(ctx), req)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while Generating content.")
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if status.Code(err).String() == codes.Unavailable.String() {
				log.Println("Retry count", retryCount)
				if retryCount > 3 {
					return nil, ErrGeneratingContent
				}
				retryCount++
				logger.Err(err).Ctx(ctx).Msgf("error while calling AIStudio to GenerateInterface, retrying with new connection. Count = %v", retryCount)
				return x.Chat(ctx, req, logger, retryCount)
			}
		}
	}
	if err != nil {
		return nil, ErrGeneratingContent
	}
	return resp, nil
}

func (x *VapusSvcInternalClients) GenerateEmbeddings(ctx context.Context, req *aipb.EmbeddingsInterface, logger zerolog.Logger, retryCount int) (*aipb.EmbeddingsResponse, error) {
	if x == nil {
		return nil, ErrAIStudioConnNotInitialized
	}
	if x.AIStudioConn == nil {
		err := x.PingTestAndReconnect(ctx, x.AIStudioDns, logger)
		if err != nil {
			return nil, ErrAIStudioConnNotInitialized
		}
	}
	resp, err := x.AIStudioConn.GenerateEmbeddings(pbtools.SwapNewContextWithAuthToken(ctx), req)
	if err != nil {
		logger.Err(err).Ctx(ctx).Msg("error while Generating embedding vectors.")
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if status.Code(err).String() == codes.Unavailable.String() {
				log.Println("Retry count", retryCount)
				if retryCount > 3 {
					return nil, err
				}
				retryCount++
				logger.Err(err).Ctx(ctx).Msgf("error while calling AIStudio to GenerateEmbeddings, retrying with new connection. Count = %v", retryCount)
				return x.GenerateEmbeddings(ctx, req, logger, retryCount)
			}
		}
	}
	return resp, nil
}

// func (x *VapusSvcInternalClients) GenerateContentAndVectors(ctx context.Context, content []byte, dimension int, aiModelParams *models.AccountAIAttributes, logger zerolog.Logger) (string, []float32, error) {
// 	resp, err := x.GenerateInterface(ctx, &aipb.ChatRequest{
// 		Actions: []aipb.AIModelNodeAction{aipb.AIModelNodeAction_GENERATE_CONTENT},
// 		Spec: &aipb.GeneratePromptParams{
// 			AiModel:     aiModelParams.GenerativeModel,
// 			ModelNodeId: aiModelParams.GenerativeModelNode,
// 			Temperature: 0.9,
// 			InputText:   string(content),
// 		},
// 	}, logger, ClientRetryStart)
// 	if err != nil || len(resp.GetOutput()) == 0 {
// 		return "", nil, err
// 	}
// 	vectors, err := x.GenerateEmbeddingsClient(ctx, resp.GetOutput()[0].GetContent(), dimension, aiModelParams, logger)
// 	if err != nil {
// 		logger.Err(err).Ctx(ctx).Msg("error while generating content.")
// 		return "", nil, err
// 	}
// 	return resp.GetOutput()[0].GetContent(), vectors, err
// }

func (x *VapusSvcInternalClients) GenerateEmbeddingsClient(ctx context.Context, content string, dimension int, aiModelParams *models.AccountAIAttributes, logger zerolog.Logger) ([]float32, error) {
	result, err := x.GenerateEmbeddings(ctx, &aipb.EmbeddingsInterface{
		AiModel:     aiModelParams.EmbeddingModel,
		ModelNodeId: aiModelParams.EmbeddingModelNode,
		InputText:   content,
		Dimension:   int64(dimension),
	}, logger, types.ClientRetryStart)
	if err != nil || result == nil || result.GetOutput() == nil {
		logger.Err(err).Ctx(ctx).Msg("error while generating embeddings for schema description")
		return []float32{}, err
	}
	if result.GetOutput().Type == mpb.EmbeddingType_FLOAT_64 {
		float32Slice := make([]float32, len(result.GetOutput().Embeddings32)) // Create a slice with the same length
		for i, v := range result.GetOutput().Embeddings64 {
			float32Slice[i] = float32(v) // Explicitly convert each element
		}
		return float32Slice, nil
	}
	return result.GetOutput().Embeddings32, nil
}
