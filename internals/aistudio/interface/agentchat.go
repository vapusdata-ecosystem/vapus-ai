package aiface

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	gdrl "github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	aimodels "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
)

type AIAgentChat struct {
	*processes.VapusInterfaceBase
	streamServer pb.AIStudio_BidiChatServer
	*AIBaseInterface
	chatObject    *models.AIStudioChat
	chatThreadMap map[string]*AIChatThread
}

type AIChatThread struct {
	payload    *prompts.GenerativePrompterPayload
	modelNode  *models.AIModelNode
	dbLog      *models.AIStudioLog
	usageLog   *models.AIStudioUsages
	finalInput string
	request    *pb.ChatRequest
	modelConn  aimodels.AIModelNodeInterface
}
type AgentChatOpts func(*AIAgentChat)

func WithAgentChatAiBase(baseInterface *AIBaseInterface) AgentChatOpts {
	return func(agent *AIAgentChat) {
		agent.AIBaseInterface = baseInterface
	}
}

func WithAgentChatGrpcStream(stream pb.AIStudio_BidiChatServer) AgentChatOpts {
	return func(agent *AIAgentChat) {
		agent.streamServer = stream
	}
}

func NewAIAgentChat(ctx context.Context, logger zerolog.Logger, opts ...AgentChatOpts) (*AIAgentChat, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	agent := &AIAgentChat{
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			AgentId:  dmutils.GetUUID(),
			CtxClaim: vapusPlatformClaim,
		},
	}
	for _, opt := range opts {
		opt(agent)
	}
	agent.logger = dmlogger.GetSubDMLogger(logger, "AIInterface", agent.AgentId)
	return agent, nil
}

func (s *AIAgentChat) Chat(ctx context.Context) error {
	idleTimer := time.NewTimer(60 * time.Minute)
	defer idleTimer.Stop()

	// This channel receives new client requests (ChatRequest).
	msgChan := make(chan *pb.ChatRequest)
	streamErr := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer close(msgChan)
		defer close(streamErr)
		defer wg.Done()
		for {
			req, err := s.streamServer.Recv()
			if err == io.EOF {
				s.logger.Info().Msg("Client closed connection")
				return
			}
			if err != nil {
				s.logger.Error().Err(err).Msg("error while receiving message from client")
				return
			}
			if s.chatObject == nil {
				if req.ChatId == "" {
					chat := &models.AIStudioChat{
						ChatId:   dmutils.GetUUID(),
						Messages: make([]*models.AIStudioLog, 0),
					}
					chat.PreSaveCreate(s.CtxClaim)
					chat.Organization = s.CtxClaim[encryption.ClaimOrganizationKey]
					err := s.Dmstore.CreateAIStudioChat(ctx, chat, s.CtxClaim)
					if err != nil {
						s.logger.Error().Err(err).Msg("error while creating chat object")
						streamErr <- dmerrors.DMError(apperr.ErrAIStudioChat400, err)
						return
					}
					s.chatObject = chat
				} else {
					chat, err := s.Dmstore.GetAIStudioChat(ctx, req.ChatId, s.CtxClaim)
					if err != nil {
						s.logger.Error().Err(err).Msg("error while getting chat object")
						streamErr <- dmerrors.DMError(apperr.ErrAIStudioChat404, err)
						return
					}
					s.chatObject = chat
					for _, mess := range chat.Messages {
						s.chatThreadMap[mess.VapusID] = &AIChatThread{
							dbLog: mess,
						}
					}
				}
			} else if s.chatObject != nil {
				if s.chatObject.ChatId != req.ChatId {
					s.logger.Error().Msg("error while processing action, chat id mismatch")
					streamErr <- dmerrors.DMError(apperr.ErrDifferentChatRequested, nil)
					return
				}
			}
			if !idleTimer.Stop() {
				<-idleTimer.C
			}
			idleTimer.Reset(60 * time.Minute)

			result, err := s.ModelPool.GetorSetNodeObject(req.ModelNodeId, nil, false)
			if err != nil || result == nil {
				s.logger.Error().Err(err).Msg("error while getting model node")
				streamErr <- dmerrors.DMError(apperr.ErrAIModelNode404, nil)
				return
			}
			thread := &AIChatThread{
				request:   req,
				modelNode: result,
				dbLog: &models.AIStudioLog{
					Input:            make([]*models.MessageLog, 0),
					Output:           make([]*models.MessageLog, 0),
					ParsedInput:      make([]*models.MessageLog, 0),
					ParsedOutput:     make([]*models.MessageLog, 0),
					AIModel:          req.Model,
					ModelNode:        result.VapusID,
					ToolCallSchema:   []*mpb.ToolCall{},
					ToolCallResponse: []*mpb.ToolCall{},
					ModelProvider:    result.ServiceProvider,
					VapusBase: models.VapusBase{
						VapusID: s.AgentId,
					},
				},
			}
			s.chatThreadMap[s.AgentId] = thread
			err = s.Act(ctx, thread)
			if err != nil {
				s.logger.Error().Err(err).Msg("error while processing action")
				streamErr <- err
				return
			}
			select {
			case <-s.streamServer.Context().Done():
				s.logger.Info().Msg("Client canceled or disconnected")
				streamErr <- s.streamServer.Context().Err()
				return
			default:
				// not canceled, continue to next Recv()
			}
		}
	}()
	wg.Wait()
	for err := range streamErr {
		if err != nil {
			s.logger.Error().Err(err).Msg("error while processing action")
			return err
		}
	}
	return nil
}

func (s *AIAgentChat) Act(ctx context.Context, thread *AIChatThread) error {
	var err error
	modelConn, err := s.ModelPool.GetorSetConnection(thread.modelNode, true, false)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while getting model node connection")
		return dmerrors.DMError(apperr.ErrAIModelNode404, err)
	}
	if thread.modelNode != nil {
		if thread.modelNode.GetScope() == mpb.ResourceScope_ORGANIZATION_SCOPE.String() {
			if slices.Contains(thread.modelNode.ApprovedOrganizations, s.CtxClaim[encryption.ClaimOrganizationKey]) == false {
				s.logger.Error().Msg("error while processing action, model not available for Organization")
				return dmerrors.DMError(apperr.ErrAIModelNode403, nil)
			}
		}
	} else {
		s.logger.Error().Msg("error while processing action, model not found")
		return dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}
	thread.modelConn = modelConn
	s.chatObject.MessageIds = append(s.chatObject.MessageIds, thread.dbLog.VapusID)
	err = s.generateContentStream(ctx, thread)
	if err != nil {
		thread.dbLog.Error = err.Error()
		thread.dbLog.ResponseStatus = "400"
		s.logger.Error().Err(err).Msg("error while processing action")
		return dmerrors.DMError(apperr.ErrAIModelManagerAction404, err)
	} else {
		thread.dbLog.ResponseStatus = "200"
	}
	go func() {
		nCtx, ContextCancel := pbtools.NewInCancelCtxWithAuthToken(ctx)
		defer ContextCancel()
		s.logRequest(nCtx, thread.finalInput, thread.payload.ParsedOutput, thread.dbLog, thread.usageLog, s.CtxClaim)
		s.chatObject.Messages = append(s.chatObject.Messages, thread.dbLog)
		s.chatObject.CurrentLog = &models.Mapper{
			Key:   thread.dbLog.VapusID,
			Value: thread.finalInput,
		}
		err := s.Dmstore.PutAIStudioChat(nCtx, s.chatObject)
		if err != nil {
			s.logger.Error().Err(err).Msg("error while updating chat object")
		}
	}()
	return nil
}

func (s *AIAgentChat) buildPayload(ctx context.Context, thread *AIChatThread) error {
	var err error
	var payload *prompts.GenerativePrompterPayload

	payload = prompts.NewPrompter(thread.request, nil, nil, s.streamServer, s.logger)

	// if thread.request.MaxCompletionTokens < 1 {
	// 	thread.request.MaxCompletionTokens = prompts.DefaultMaxOPTokenLength
	// }
	payload.StudioLog = thread.dbLog
	var wg sync.WaitGroup
	var errChan = make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		if thread.request.PromptId != "" {
			prompt, err := s.Dmstore.GetAIPrompt(ctx, thread.request.PromptId, s.CtxClaim)
			if err != nil {
				s.logger.Error().Err(err).Msgf("error while getting prompt %v", thread.request.PromptId)
				errChan <- err
			}
			payload.Prompt = prompt
		}
		payload.RenderPrompt()
	}()
	go func() {
		defer wg.Done()
		sessionMessages := []*models.AIStudioLog{}
		if len(s.chatObject.Messages) > 5 {
			for _, mess := range s.chatObject.Messages[:5] {
				sessionMessages = append(sessionMessages, s.chatThreadMap[mess.VapusID].dbLog)
			}
		} else {
			for _, mess := range s.chatObject.Messages {
				sessionMessages = append(sessionMessages, s.chatThreadMap[mess.VapusID].dbLog)
			}
		}

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
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	thread.finalInput = ""
	for i, mess := range thread.request.Messages {
		mess.Content = strings.TrimSpace(mess.Content)
		mess.Content = strings.Trim(mess.Content, "\n")
		if mess.Content == "" {
			thread.request.Messages = append(thread.request.Messages[:i], thread.request.Messages[i+1:]...)
			continue
		}
		if mess.Role != aicore.SYSTEM {
			thread.finalInput += mess.Content
		}
	}
	thread.payload = payload
	if len(thread.request.Messages) < 1 || len(thread.finalInput) < 1 {
		thread.payload.GuardrailsFailed = true
		thread.payload.ParsedOutput = "Please provide a valid input, input cannot be empty"
		return nil
	}
	err = s.guardrailChecks(ctx, thread)
	if err != nil || payload.GuardrailsFailed {
		return nil
	}
	return nil
}

func (s *AIAgentChat) guardrailChecks(ctx context.Context, thread *AIChatThread) error {
	// TODO: store guardrails log and usage
	val, ok := s.GuardrailPool.ModelGuardRails.Load(thread.modelNode.VapusID)
	if !ok {
		return nil
	}
	gdClients, valid := val.([]string)
	if !valid {
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
					Content: thread.finalInput,
				},
			},
			StartedAt: dmutils.GetMilliEpochTime(),
		}
		scanResult := guard.Scan(ctx, thread.finalInput, s.logger, s.CtxClaim)
		gdLog.PreSaveCreate(s.CtxClaim)
		gdLog.EndedAt = dmutils.GetMilliEpochTime()
		if len(scanResult.WordGuard) > 0 || len(scanResult.TopicGuard) > 0 || len(scanResult.ContentGuard) > 0 {
			thread.payload.ParsedOutput = guard.Guardrail.FailureMessage
			thread.payload.GuardrailsFailed = true
			gdLog.Output = append(gdLog.Output, scanResult.WordGuard...)
			gdLog.Output = append(gdLog.Output, scanResult.TopicGuard...)
			gdLog.Output = append(gdLog.Output, scanResult.ContentGuard...)
			gdLog.Failed = true
			gdLog.FailedMessage = thread.payload.ParsedOutput
			return fmt.Errorf("guardrail failed for your input %v", scanResult)
		} else {
			gdLog.Output = append(gdLog.Output, "Guardrails passed")
			gdLog.Failed = false
		}
		s.logGuardrailRequest(ctx, thread.finalInput, gdLog, scanResult.Usage, s.CtxClaim)
		if gdLog.Failed {
			s.logger.Error().Msg("guardrail failed for your input")
			return fmt.Errorf("guardrail failed for your input")
		}
	}
	return nil
}

func (s *AIAgentChat) generateContentStream(ctx context.Context, thread *AIChatThread) error {
	var err error
	err = s.buildPayload(ctx, thread)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while building payload")
		return err
	}
	if thread.payload.GuardrailsFailed {
		s.streamServer.Send(
			thread.payload.BuildStreamResponseOP(aicore.StreamGuardrailFailed.String(), &prompts.PayloadgenericResponse{
				Data: thread.payload.ParsedOutput,
				Role: aicore.VAPUSGUARD,
			}),
		)
		s.streamServer.SendMsg(&pb.ChatResponse{
			Event: aicore.StreamEventEnd.String(),
		})
		return nil
	}
	err = thread.modelConn.GenerateContentStream(ctx, thread.payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating content from model %v", thread.request.Model)
		return err
	}
	thread.usageLog = thread.payload.Usage
	return nil
}
