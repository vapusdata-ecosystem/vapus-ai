package aiface

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

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

type AIGateway struct {
	*processes.VapusInterfaceBase
	payload        *prompts.GenerativePrompterPayload
	resultMetaData []*models.Mapper
	*AIBaseInterface
	request        *pb.ChatRequest
	modelNode      *models.AIModelNode
	dbLog          *models.AIStudioLog
	usageLog       *models.AIStudioUsages
	finalInput     string
	stream         bool
	GatewayChannel chan *aicore.GwChatCompletionChunk
}

type GwOpts func(*AIGateway)

func WithGwAiBase(baseInterface *AIBaseInterface) GwOpts {
	return func(agent *AIGateway) {
		agent.AIBaseInterface = baseInterface
	}
}

func WithGwRequest(request *pb.ChatRequest) GwOpts {
	return func(agent *AIGateway) {
		agent.request = request
	}
}

func WithGwStream(stream bool) GwOpts {
	return func(agent *AIGateway) {
		agent.stream = stream
	}
}

func WithGwGatewayChannel(channel chan *aicore.GwChatCompletionChunk) GwOpts {
	return func(agent *AIGateway) {
		agent.GatewayChannel = channel
	}
}

func NewAIGateway(ctx context.Context, logger zerolog.Logger, opts ...GwOpts) (*AIGateway, error) {
	var err error
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	agent := &AIGateway{
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			AgentId:  dmutils.GetUUID(),
			CtxClaim: vapusPlatformClaim,
		},
	}
	for _, opt := range opts {
		opt(agent)
	}
	if agent.request == nil {
		logger.Error().Err(err).Msg("error while getting model node")
		return nil, dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}

	agent.logger = dmlogger.GetSubDMLogger(logger, "AIInterface", agent.AgentId)
	return agent, nil
}

func (s *AIGateway) GetResult() *pb.ChatResponse {
	if s.payload == nil {
		return &pb.ChatResponse{
			Event: aicore.StreamEventEnd.String(),
		}
	}
	return s.payload.Response
}

func (x *AIGateway) SendSSE(obj *aicore.GwChatCompletionChunk) {
	if x.GatewayChannel != nil {
		for {
			select {
			case x.GatewayChannel <- obj:
				return
			case <-time.After(50 * time.Millisecond):
				x.Logger.Info().Msg("waiting for channel to be ready")
			}
		}
	}
}

func (s *AIGateway) Act(ctx context.Context) error {
	var err error
	result, err := s.ModelPool.GetorSetNodeObject(s.request.ModelNodeId, nil, false)
	if err != nil || result == nil {
		s.logger.Error().Err(err).Msg("error while getting model node")
		return dmerrors.DMError(apperr.ErrAIModelNode404, nil)
	}
	s.dbLog = &models.AIStudioLog{
		Mode:             pb.AIInterfaceMode_P2P.String(),
		Input:            make([]*models.MessageLog, 0),
		Output:           make([]*models.MessageLog, 0),
		ParsedInput:      make([]*models.MessageLog, 0),
		ParsedOutput:     make([]*models.MessageLog, 0),
		AIModel:          s.request.Model,
		ToolCallSchema:   []*mpb.ToolCall{},
		ToolCallResponse: []*mpb.ToolCall{},
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
	switch s.stream {
	case true:
		err = s.generateContentStream(ctx, modelConn)
	case false:
		err = s.generateContent(ctx, modelConn)
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
		s.logRequest(nCtx, s.finalInput, s.payload.ParsedOutput, s.dbLog, s.usageLog, s.CtxClaim)
	}()
	return nil
}

func (s *AIGateway) buildPayload(ctx context.Context) error {
	var err error
	var payload *prompts.GenerativePrompterPayload

	payload = prompts.NewPrompter(s.request, nil, nil, nil, s.logger)
	payload.Mode = pb.AIInterfaceMode_P2P
	if s.stream {
		payload.SSEChan = s.GatewayChannel
	}
	// if s.request.MaxCompletionTokens < 1 {
	// 	s.request.MaxCompletionTokens = prompts.DefaultMaxOPTokenLength
	// }
	payload.StudioLog = s.dbLog
	if s.request.PromptId != "" {
		prompt, err := s.Dmstore.GetAIPrompt(ctx, s.request.PromptId, s.CtxClaim)
		if err != nil {
			s.logger.Error().Err(err).Msgf("error while getting prompt %v", s.request.PromptId)
			return err
		}
		payload.Prompt = prompt
	}
	payload.RenderPrompt()
	s.finalInput = ""
	for i, mess := range s.request.Messages {
		log.Println("mess.Content++++++++++++++++++++++++++++++++++++++++", mess)
		mess.Content = strings.TrimSpace(mess.Content)
		mess.Content = strings.Trim(mess.Content, "\n")
		if mess.Content == "" && mess.StructuredContent == nil {
			slices.Delete(s.request.Messages, i, i+1)
			continue
		}
		if mess.Role != aicore.SYSTEM {
			if mess.Content != "" {
				s.finalInput += mess.Content
			} else if mess.StructuredContent != nil {
				for _, value := range mess.StructuredContent {
					switch value.Type {
					case aicore.AIResponseFormatText.String():
						s.finalInput += value.Text
					case aicore.AIResponseFormatImageUrl.String():
						s.finalInput += value.Text
						if value.ImageUrl.Url != "" {
							s.finalInput = value.ImageUrl.Url
						}
					case aicore.AIResponseFormatInputAudio.String():
						s.finalInput += value.Text
						if value.InputAudio.Data != "" {
							s.finalInput = value.InputAudio.Data
						}
					case aicore.AIResponseFormatInputFile.String():
						s.finalInput += value.Text
						// if value.File.Data != "" {
						// 	s.finalInput = value.InputAudio.Data
						// }
					default:
						continue
					}
				}
			}
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

func (s *AIGateway) guardrailChecks(ctx context.Context) error {
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
			s.logGuardrailRequest(ctx, s.finalInput, gdLog, scanResult.Usage, s.CtxClaim)
			return fmt.Errorf("guardrail failed for your input %v", scanResult)
		} else {
			gdLog.Output = append(gdLog.Output, "No guardrail failed")
			gdLog.Failed = false
			s.logGuardrailRequest(ctx, s.finalInput, gdLog, scanResult.Usage, s.CtxClaim)
		}
	}
	return nil
}

func (s *AIGateway) generateContent(ctx context.Context, modelConn aimodels.AIModelNodeInterface) error {
	var err error
	log.Println("=========>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	err = s.buildPayload(ctx)
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
	log.Println("========= +++++++++++++++++++++++++++++++++++++")
	err = modelConn.GenerateContent(ctx, s.payload)
	if err != nil {
		s.logger.Error().Err(err).Msgf("error while generating content from model %v", s.request.Model)
		return err
	}
	s.usageLog = s.payload.Usage
	return nil
}

func (s *AIGateway) generateContentStream(ctx context.Context, modelConn aimodels.AIModelNodeInterface) error {
	var err error
	log.Println("========= +++++++++++++++++++++++++++++++++++++ streraming")
	err = s.buildPayload(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while building payload")
		return err
	}
	if s.payload.GuardrailsFailed {
		s.SendSSE(&aicore.GwChatCompletionChunk{
			ID:      dmutils.GetUUID(),
			Created: dmutils.GetEpochTime(),
			Model:   s.request.Model,
			Choices: []aicore.GwChatChoice{
				{
					Delta: &aicore.GwChatDelta{
						Role:    aicore.VAPUSGUARD,
						Content: s.payload.ParsedOutput,
						Refusal: aicore.StreamGuardrailFailed.String(),
					},
					Index:        0,
					FinishReason: dmutils.Str2Ptr(aicore.StreamGuardrailFailed.String()),
				},
			},
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
