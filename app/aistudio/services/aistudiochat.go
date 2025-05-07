package services

import (
	"context"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type AIStudioChatAgentRequest struct {
	ManageRequest *pb.ManageAIChatRequest
	GetterRequest *pb.GetAIChatRequest
	Response      *pb.AIChatResponse
}

type AIStudioChatAgent struct {
	*AIStudioChatAgentRequest
	AgentError error
	*aidmstore.AIStudioDMStore
	*processes.VapusInterfaceBase
}

func (s *AIStudioServices) NewAIStudioChatManager(ctx context.Context, request *AIStudioChatAgentRequest) (*AIStudioChatAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	agent := &AIStudioChatAgent{
		AIStudioChatAgentRequest: request,
		AIStudioDMStore:          s.DMStore,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			InitAt:   dmutils.GetEpochTime(),
		},
	}
	agent.Response = &pb.AIChatResponse{}
	agent.SetAgentId()
	if agent.ManageRequest != nil {
		agent.Action = agent.ManageRequest.GetAction().String()
	}
	if agent.ManageRequest != nil {
		agent.Action = agent.ManageRequest.GetAction().String()
	} else {
		agent.Action = ""
	}
	agent.Logger = pkgs.GetSubDMLogger(types.NABHIKAGENT.String(), agent.AgentId)
	return agent, nil
}

func (a *AIStudioChatAgent) GetResponse() *pb.AIChatResponse {
	a.FinishAt = dmutils.GetEpochTime()
	a.FinalLog()
	return a.Response
}

func (a *AIStudioChatAgent) Act(ctx context.Context) error {
	switch a.Action {
	case pb.AIChatAction_CREATE.String():
		return a.createAIStudioChat(ctx)
	case pb.AIChatAction_ARCHIVE.String():
		return a.archiveAIStudioChat(ctx)
	default:
		return a.ListAIStudioChats(ctx)
	}
}

func (a *AIStudioChatAgent) ListAIStudioChats(ctx context.Context) error {
	if a.GetterRequest == nil {
		a.Logger.Error().Msg("error while getting request for listing AI studio chat")
		return dmerrors.DMError(apperr.ErrInvalidAIStudioChatRequested, nil)
	}
	if a.GetterRequest.GetChatId() != "" {
		return a.GetAIStudioChat(ctx)
	}
	chat, err := a.AIStudioDMStore.ListAIStudioChats(ctx, 0, a.CtxClaim)
	if err != nil || chat == nil {
		a.Logger.Error().Msg("error while getting AI studio chat history from datastore")
		return dmerrors.DMError(apperr.ErrInvalidAIStudioChatRequested, err)
	}
	a.Response.Output = utils.AISCLToPbChat(chat)
	return nil
}

func (a *AIStudioChatAgent) GetAIStudioChat(ctx context.Context) error {
	chat, err := a.AIStudioDMStore.GetAIStudioChat(ctx, a.GetterRequest.GetChatId(), a.CtxClaim)
	if err != nil || chat == nil {
		a.Logger.Error().Msg("error while getting AI studio chat from datastore")
		return dmerrors.DMError(apperr.ErrInvalidAIStudioChatRequested, err)
	}
	a.Response.Output = []*pb.AIStudioChat{utils.AISCOToPbChat(chat)}
	return nil
}

func (a *AIStudioChatAgent) createAIStudioChat(ctx context.Context) error {
	chat := &models.AIStudioChat{
		ChatId:   dmutils.GetUUID(),
		Messages: make([]*models.AIStudioLog, 0),
	}

	chat.PreSaveCreate(a.CtxClaim)
	chat.VapusID = chat.ChatId
	chat.Organization = a.CtxClaim[encryption.ClaimOrganizationKey]
	err := a.AIStudioDMStore.CreateAIStudioChat(ctx, chat, a.CtxClaim)
	if err != nil {
		a.Logger.Error().Msg("error while saving thread in datastore")
		return dmerrors.DMError(apperr.ErrCreatingAIStudioChat, err)
	}
	a.Response.Output = utils.AISCLToPbChat([]*models.AIStudioChat{chat})
	return nil
}

func (a *AIStudioChatAgent) archiveAIStudioChat(ctx context.Context) error {
	if a.ManageRequest == nil {
		a.Logger.Error().Msg("error while getting request for archiving AI studio chat")
		return dmerrors.DMError(apperr.ErrInvalidAIStudioChatRequested, nil)
	}
	chat, err := a.AIStudioDMStore.GetAIStudioChat(ctx, a.ManageRequest.GetChatId(), a.CtxClaim)
	if err != nil || chat == nil {
		a.Logger.Error().Msg("error while getting AI studio chat from datastore")
		return dmerrors.DMError(apperr.ErrInvalidAIStudioChatRequested, err)
	}

	chat.PreSaveDelete(a.CtxClaim)
	err = a.AIStudioDMStore.PutAIStudioChat(ctx, chat)
	if err != nil {
		a.Logger.Error().Msg("error while saving thread in datastore")
		return dmerrors.DMError(apperr.ErrArchiveFabricChat, err)
	}
	a.Response.Output = nil
	return nil
}
