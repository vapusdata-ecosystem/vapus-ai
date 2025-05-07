package services

import (
	"context"
	"slices"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	aimodels "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
)

type EmbeddingAgent struct {
	*processes.VapusInterfaceBase
	embeddings *mpb.Embeddings
	request    *pb.EmbeddingsInterface
	logger     zerolog.Logger
	modelNode  *models.AIModelNode
	*AIStudioServices
}

func (s *AIStudioServices) NewEmbeddingAgent(ctx context.Context,
	request *pb.EmbeddingsInterface) (*EmbeddingAgent, error) {
	var err error
	modelsNodeId := ""
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	modelsNodeId = request.GetModelNodeId()
	result, err := pkgs.AIModelNodeConnectionPoolManager.GetorSetNodeObject(modelsNodeId, nil, false)
	if err != nil || result == nil {
		s.Logger.Error().Err(err).Msg("error while getting model node")
		return nil, dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}

	agent := &EmbeddingAgent{
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			AgentId: dmutils.GetUUID(),
			// Ctx:      ctx,
			CtxClaim: vapusPlatformClaim,
		},
		modelNode: result,
		request:   request,
	}
	agent.logger = pkgs.GetSubDMLogger("NewEmbeddingAgent", agent.AgentId)
	return agent, nil
}

func (s *EmbeddingAgent) GetEmbeddings() *mpb.Embeddings {
	return s.embeddings
}

func (s *EmbeddingAgent) Act(ctx context.Context) error {
	var err error
	modelConn, err := pkgs.AIModelNodeConnectionPoolManager.GetorSetConnection(s.modelNode, true, false)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while getting model node connection")
		return dmerrors.DMError(apperr.ErrAIModelNode404, err)
	}
	if s.modelNode != nil {
		if s.modelNode.GetScope() == mpb.ResourceScope_ORGANIZATION_SCOPE.String() {
			if slices.Contains(s.modelNode.ApprovedOrganizations, s.CtxClaim[encryption.ClaimOrganizationKey]) == false {
				s.logger.Error().Msg("error while processing action, model not available for domain")
				return dmerrors.DMError(apperr.ErrAIModelNode403, nil)
			}
		}
	} else {
		s.logger.Error().Msg("error while processing action, model not found")
		return dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}
	if s.request != nil {
		err = s.generateEmbeddings(ctx, modelConn)
	} else {
		s.logger.Error().Msg("error while processing action, invalid action")
		return dmerrors.DMError(apperr.ErrAIModelManagerAction404, nil)
	}
	return nil
}

func (s *EmbeddingAgent) generateEmbeddings(ctx context.Context, modelConn aimodels.AIModelNodeInterface) error {
	payload := &prompts.AIEmbeddingPayload{
		Dimensions:     int(s.request.GetDimension()),
		EmbeddingModel: s.request.GetAiModel(),
		Input:          s.request.GetInputText(),
	}
	err := modelConn.GenerateEmbeddings(ctx, payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating embeddings from model %v", s.request.GetAiModel())
		return err
	}
	if payload.Embeddings.Vectors32 != nil {
		s.embeddings = &mpb.Embeddings{
			Embeddings32: payload.Embeddings.Vectors32,
			Type:         mpb.EmbeddingType_FLOAT_32,
			Dimension:    int64(payload.Dimensions),
		}
	} else {
		s.embeddings = &mpb.Embeddings{
			Embeddings64: payload.Embeddings.Vectors64,
			Type:         mpb.EmbeddingType_FLOAT_64,
			Dimension:    int64(payload.Dimensions),
		}
	}
	return nil
}
