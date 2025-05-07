package aiface

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/core"
	gdrl "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/guardrails"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/prompts"
	aimodels "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/providers"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusdata/core/process"
)

type AIChatGateway struct {
	*processes.VapusInterfaceBase
	payload        *prompts.GenerativePrompterPayload
	resultMetaData []*models.Mapper
	request        *pb.ChatRequest
	modelNode      *models.AIModelNode
	streamServer   pb.AIStudio_ChatServer
	dbLog          *models.AIStudioLog
	usageLog       *models.AIStudioUsages
	finalInput     string
	*AIBaseInterface
	chatObject    *models.AIStudioChat
	chatThreadMap map[string]*AIChatThread
	isChat        bool
	internal      bool
}

type ChatGatewayOpts func(*AIChatGateway)

func WithChatGatewayAiBase(baseInterface *AIBaseInterface) ChatGatewayOpts {
	return func(agent *AIChatGateway) {
		agent.AIBaseInterface = baseInterface
	}
}

func WithChatGatewayRequest(request *pb.ChatRequest) ChatGatewayOpts {
	return func(agent *AIChatGateway) {
		agent.request = request
	}
}

func WithChatGatewayGrpcStream(stream pb.AIStudio_ChatServer) ChatGatewayOpts {
	return func(agent *AIChatGateway) {
		agent.streamServer = stream
	}
}

func WithChatEnabled(isChat bool) ChatGatewayOpts {
	return func(agent *AIChatGateway) {
		agent.isChat = isChat
	}
}

func WithInternal(internal bool) ChatGatewayOpts {
	return func(agent *AIChatGateway) {
		agent.internal = internal
	}
}

func WithCtxClaim(claim map[string]string) ChatGatewayOpts {
	return func(agent *AIChatGateway) {
		agent.CtxClaim = claim
	}
}

func NewAIChatGateway(ctx context.Context, logger zerolog.Logger, opts ...ChatGatewayOpts) (*AIChatGateway, error) {
	var err error
	agent := &AIChatGateway{
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			AgentId: dmutils.GetUUID(),
		},
		chatThreadMap: make(map[string]*AIChatThread),
	}
	for _, opt := range opts {
		opt(agent)
	}
	if !agent.internal {
		vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
		if !ok {
			logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
			return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
		}
		agent.CtxClaim = vapusPlatformClaim
	} else {
		if agent.CtxClaim == nil {
			return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
		}
	}
	if agent.request == nil {
		logger.Error().Err(err).Msg("error while getting model node")
		return nil, dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}
	if agent.isChat {
		if agent.request.ChatId == "" {
			chat := &models.AIStudioChat{
				ChatId:   dmutils.GetUUID(),
				Messages: make([]*models.AIStudioLog, 0),
			}
			chat.VapusID = chat.ChatId
			chat.PreSaveCreate(agent.CtxClaim)
			chat.Organization = agent.CtxClaim[encryption.ClaimOrganizationKey]
			err := agent.Dmstore.CreateAIStudioChat(ctx, chat, agent.CtxClaim)
			if err != nil {
				logger.Error().Err(err).Msg("error while creating chat object")
				return nil, dmerrors.DMError(apperr.ErrAIStudioChat400, err)
			}
			agent.chatObject = chat
		} else {
			chat, err := agent.Dmstore.GetAIStudioChat(ctx, agent.request.ChatId, agent.CtxClaim)
			if err != nil {
				logger.Error().Err(err).Msg("error while getting chat object")
				return nil, dmerrors.DMError(apperr.ErrAIStudioChat404, err)
			}
			agent.chatObject = chat
			for _, mess := range chat.Messages {
				agent.chatThreadMap[mess.VapusID] = &AIChatThread{
					dbLog: mess,
				}
			}
		}
	}
	agent.logger = dmlogger.GetSubDMLogger(logger, "AIInterface", agent.AgentId)
	return agent, nil
}

func (s *AIChatGateway) GetResult() *pb.ChatResponse {
	if s.payload == nil {
		return &pb.ChatResponse{
			Event: aicore.StreamEventEnd.String(),
		}
	}
	return s.payload.Response
}

func (s *AIChatGateway) Act(ctx context.Context) error {
	var err error
	result, err := s.ModelPool.GetorSetNodeObject(s.request.ModelNodeId, nil, false)
	if err != nil || result == nil {
		s.logger.Error().Err(err).Msg("error while getting model node")
		return dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}
	s.dbLog = &models.AIStudioLog{
		Mode:             pb.AIInterfaceMode_CHAT_MODE.String(),
		Input:            make([]*models.MessageLog, 0),
		Output:           make([]*models.MessageLog, 0),
		ParsedInput:      make([]*models.MessageLog, 0),
		ParsedOutput:     make([]*models.MessageLog, 0),
		AIModel:          s.request.Model,
		ToolCallSchema:   []*mpb.ToolCall{},
		ToolCallResponse: []*mpb.ToolCall{},
		VapusBase: models.VapusBase{
			VapusID: dmutils.GetUUID(),
		},
	}
	s.modelNode = result
	s.dbLog.ModelNode = s.modelNode.VapusID
	s.dbLog.ModelProvider = s.modelNode.ServiceProvider
	modelConn, err := s.ModelPool.GetorSetConnection(s.modelNode, true, false)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while getting model node connection")
		return dmerrors.DMError(apperr.ErrAIModelNode404, err)
	}
	if s.modelNode != nil {
		if s.modelNode.GetScope() == mpb.ResourceScope_ORGANIZATION_SCOPE.String() {
			if slices.Contains(s.modelNode.ApprovedOrganizations, s.CtxClaim[encryption.ClaimOrganizationKey]) == false {
				s.logger.Error().Msg("error while processing action, model not available for Organization")
				return dmerrors.DMError(apperr.ErrAIModelNode403, nil)
			}
		}
	} else {
		s.logger.Error().Msg("error while processing action, model not found")
		return dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}
	if s.isChat {
		s.chatObject.MessageIds = append(s.chatObject.MessageIds, s.dbLog.VapusID)
	}
	if s.request != nil {
		if s.streamServer != nil {
			err = s.generateContentStream(ctx, modelConn)
		} else {
			err = s.generateContent(ctx, modelConn)
		}
	} else {
		s.logger.Error().Msg("error while processing action, invalid action")
		return dmerrors.DMError(apperr.ErrAIModelManagerAction404, nil)
	}

	if err != nil {
		s.dbLog.Error = err.Error()
		s.dbLog.ResponseStatus = "400"
		s.logger.Error().Err(err).Msg("error while processing action")
		return dmerrors.DMError(apperr.ErrAIModelManagerAction404, err)
	} else {
		s.dbLog.ResponseStatus = "200"
	}
	go func() {
		nCtx, ContextCancel := pbtools.NewInCancelCtxWithAuthToken(ctx)
		defer ContextCancel()
		if s.isChat {
			s.dbLog.ChatId = s.chatObject.ChatId
		}
		err := s.logRequest(nCtx, s.finalInput, s.payload.ParsedOutput, s.dbLog, s.usageLog, s.CtxClaim)
		if err != nil {
			s.logger.Error().Err(err).Msg("error while logging request")
		}
		// Structuring the chat and storing in the DB
		if s.isChat {
			s.chatObject.Messages = append(s.chatObject.Messages, s.dbLog)
			s.chatObject.CurrentLog = &models.Mapper{
				Key:   s.dbLog.VapusID,
				Value: s.finalInput,
			}

			err := s.Dmstore.PutAIStudioChat(nCtx, s.chatObject)
			if err != nil {
				s.logger.Error().Err(err).Msg("error while updating chat object")
			}
		}
	}()
	return nil
}

func (s *AIChatGateway) buildPayload(ctx context.Context, stream pb.AIStudio_ChatServer) error {
	var err error
	var payload *prompts.GenerativePrompterPayload
	if s.streamServer != nil {
		payload = prompts.NewPrompter(s.request, nil, stream, nil, s.logger)
	} else {
		payload = prompts.NewPrompter(s.request, nil, nil, nil, s.logger)
	}
	payload.Mode = pb.AIInterfaceMode_CHAT_MODE
	// if s.request.MaxCompletionTokens < 1 {
	// 	if stream == nil {
	// 		s.request.MaxCompletionTokens = prompts.DefaultMaxOPTokenLength
	// 	}
	// }
	payload.StudioLog = s.dbLog
	var wg sync.WaitGroup
	var errChan = make(chan error, 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if s.request.PromptId != "" {
			prompt, err := s.Dmstore.GetAIPrompt(ctx, s.request.PromptId, s.CtxClaim)
			if err != nil {
				s.logger.Error().Err(err).Msgf("error while getting prompt %v", s.request.PromptId)
				errChan <- err
			}
			payload.Prompt = prompt
		}
		payload.RenderPrompt()
	}()
	if s.isChat {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sessionMessages := []*models.AIStudioLog{}
			if len(s.chatObject.Messages) > 7 {
				// fmt.Println("chat Lenght is grater than 7")
				for _, mess := range s.chatObject.Messages[len(s.chatObject.Messages)-7:] {
					sessionMessages = append(sessionMessages, s.chatThreadMap[mess.VapusID].dbLog)
				}
			} else {
				// fmt.Println("chat Lenght is less than 7")
				for _, mess := range s.chatObject.Messages {
					sessionMessages = append(sessionMessages, s.chatThreadMap[mess.VapusID].dbLog)
				}
			}
			// Printing what is appending in the sessionMessage
			// fmt.Println(len(sessionMessages))
			// for index, val := range sessionMessages {
			// 	fmt.Println(val.Input[0].Content, "::::>>> ", index)
			// 	// if len(val.Output) > 0 {
			// 	// 	fmt.Println(val.Output[0].Content)
			// 	// }
			// }
			for _, session := range sessionMessages {
				for _, message := range session.Input {
					if message.Content == "" {
						continue
					}
					if strings.TrimSpace(message.Content) != "" {
						if len(message.Content) > 500 {
							message.Content = dmutils.StringSlicer(message.Content, 500)
						}
						payload.SessionContext = append(payload.SessionContext, &prompts.SessionMessage{
							Message: message.Content,
							Role:    message.Role,
						})
					}
				}
				for _, message := range session.Output {
					if message.Content == "" {
						continue
					}
					if strings.TrimSpace(message.Content) != "" {
						if len(message.Content) > 500 {
							message.Content = dmutils.StringSlicer(message.Content, 500)
						}
						payload.SessionContext = append(payload.SessionContext, &prompts.SessionMessage{
							Message: message.Content,
							Role:    message.Role,
						})
					}
				}
			}
		}()
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	s.finalInput = ""
	for i, mess := range s.request.Messages {
		mess.Content = strings.TrimSpace(mess.Content)
		mess.Content = strings.Trim(mess.Content, "\n")
		if mess.Content == "" {
			s.request.Messages = append(s.request.Messages[:i], s.request.Messages[i+1:]...)
			continue
		}
		if mess.Role != aicore.SYSTEM {
			s.finalInput += mess.Content
		}
	}
	s.payload = payload
	if len(s.request.Messages) < 1 || len(s.finalInput) < 1 {
		s.payload.GuardrailsFailed = true
		s.payload.ParsedOutput = "Please provide a valid input, input cannot be empty"
		return nil
	}
	err = s.guardrailChecks(ctx)
	if err != nil || payload.GuardrailsFailed {
		return nil
	}
	return nil
}

func (s *AIChatGateway) guardrailChecks(ctx context.Context) error {
	val, ok := s.GuardrailPool.ModelGuardRails.Load(s.modelNode.VapusID)
	if !ok {
		return nil
	}
	gdClients, valid := val.([]string)
	if !valid {
		return nil
	}
	if s.finalInput == "" {
		return nil
	}
	for _, guardId := range gdClients {
		val, ok := s.GuardrailPool.GuardrailClientMap.Load(guardId)
		if !ok {
			continue
		}
		guard, valid := val.(*gdrl.GuardRailClient)
		if !valid {
			continue
		}
		gdLog := &models.AIGuardrailsLog{
			Output: make([]string, 0),
			Input: []*models.MessageLog{
				{
					Role:    aicore.USER,
					Content: s.finalInput,
				},
			},
			StartedAt: dmutils.GetMilliEpochTime(),
		}
		scanResult := guard.Scan(ctx, s.finalInput, s.logger, s.CtxClaim)
		gdLog.PreSaveCreate(s.CtxClaim)
		gdLog.EndedAt = dmutils.GetMilliEpochTime()
		if len(scanResult.WordGuard) > 0 || len(scanResult.TopicGuard) > 0 || len(scanResult.ContentGuard) > 0 {
			s.payload.ParsedOutput = guard.Guardrail.FailureMessage
			s.payload.GuardrailsFailed = true
			gdLog.Output = append(gdLog.Output, scanResult.WordGuard...)
			gdLog.Output = append(gdLog.Output, scanResult.TopicGuard...)
			gdLog.Output = append(gdLog.Output, scanResult.ContentGuard...)
			gdLog.Failed = true
			gdLog.FailedMessage = s.payload.ParsedOutput
		} else {
			gdLog.Output = append(gdLog.Output, "No guardrail failed")
			gdLog.Failed = false
		}
		s.logGuardrailRequest(ctx, s.finalInput, gdLog, scanResult.Usage, s.CtxClaim)
		if gdLog.Failed {
			s.logger.Error().Msg("guardrail failed for your input")
			return fmt.Errorf("guardrail failed for your input")
		}
	}
	return nil
}

func (s *AIChatGateway) generateContentStream(ctx context.Context, modelConn aimodels.AIModelNodeInterface) error {
	var err error
	err = s.buildPayload(ctx, s.streamServer)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while building payload")
		return err
	}
	if s.payload.GuardrailsFailed {
		s.streamServer.Send(
			s.payload.BuildStreamResponseOP(aicore.StreamGuardrailFailed.String(), &prompts.PayloadgenericResponse{
				Data: s.payload.ParsedOutput,
				Role: aicore.VAPUSGUARD,
			}),
		)
		s.streamServer.SendMsg(&pb.ChatResponse{
			Event: aicore.StreamEventEnd.String(),
		})
		return nil
	}
	err = modelConn.GenerateContentStream(ctx, s.payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating content from model %v", s.request.Model)
		return err
	}
	s.usageLog = s.payload.Usage
	return nil
}

func (s *AIChatGateway) generateContent(ctx context.Context, modelConn aimodels.AIModelNodeInterface) error {
	var err error

	err = s.buildPayload(ctx, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while building payload")
		return err
	}
	if s.payload.GuardrailsFailed {
		s.payload.BuildResponseOP(aicore.StreamGuardrailFailed.String(), &prompts.PayloadgenericResponse{
			FinishReason: aicore.StreamGuardrailFailed.String(),
			Data:         s.payload.ParsedOutput,
			Role:         aicore.VAPUSGUARD,
		}, false)
		return nil
	}
	err = modelConn.GenerateContent(ctx, s.payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating content from model %v", s.request.Model)
		return err
	}
	s.usageLog = s.payload.Usage
	return nil
}
