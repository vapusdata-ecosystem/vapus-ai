package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	services "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type AIPrompts struct {
	pb.UnimplementedAIPromptsServer
	validator  *dmutils.DMValidator
	DMServices *services.AIStudioServices
	Logger     zerolog.Logger
}

var AIPromptsManager *AIPrompts

func NewAIPrompts() *AIPrompts {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIPrompts")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("NewAIPrompts: Error while loading validator")
	}

	l.Info().Msg("AIPrompts Controller initialized")
	return &AIPrompts{
		validator:  validator,
		Logger:     l,
		DMServices: services.AIStudioServiceManager,
	}
}

func InitAIPrompts() {
	if AIPromptsManager == nil {
		AIPromptsManager = NewAIPrompts()
	}
}

func (v *AIPrompts) Create(ctx context.Context, req *pb.PromptManagerRequest) (*mpb.VapusCreateResponse, error) {
	agent, err := v.DMServices.NewAIPromptIntAgent(ctx, services.WithPromptManagerRequest(req), services.WithPromptManagerAction(mpb.ResourceLcActions_CREATE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Create: error while creating new AIStudio Prompt Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Create: error while processing AIStudio Prompt Agent request")
		return nil, err
	}
	response := agent.GetCreateResponse()
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIPrompts) Update(ctx context.Context, req *pb.PromptManagerRequest) (*pb.PromptResponse, error) {
	agent, err := v.DMServices.NewAIPromptIntAgent(ctx, services.WithPromptManagerRequest(req), services.WithPromptManagerAction(mpb.ResourceLcActions_UPDATE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Update: error while creating new AIStudio Prompt Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Update: error while processing AIStudio Prompt Agent request")
		return nil, err
	}
	response := &pb.PromptResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIPrompts) Get(ctx context.Context, req *pb.PromptGetterRequest) (*pb.PromptResponse, error) {
	agent, err := v.DMServices.NewAIPromptIntAgent(ctx, services.WithPromptGetterRequest(req), services.WithPromptManagerAction(mpb.ResourceLcActions_GET))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Get: error while creating new AIStudio Prompt Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Get: error while processing AIStudio Prompt Agent request")
		return nil, err
	}
	response := &pb.PromptResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIPrompts) List(ctx context.Context, req *pb.PromptGetterRequest) (*pb.PromptResponse, error) {
	agent, err := v.DMServices.NewAIPromptIntAgent(ctx, services.WithPromptGetterRequest(req), services.WithPromptManagerAction(mpb.ResourceLcActions_LIST))
	if err != nil {
		v.Logger.Error().Err(err).Msg("List: error while creating new AIStudio Prompt Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("List: error while processing AIStudio Prompt Agent request")
		return nil, err
	}
	response := &pb.PromptResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIPrompts) Archive(ctx context.Context, req *pb.PromptGetterRequest) (*pb.PromptResponse, error) {
	agent, err := v.DMServices.NewAIPromptIntAgent(ctx, services.WithPromptGetterRequest(req), services.WithPromptManagerAction(mpb.ResourceLcActions_ARCHIVE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Archive: error while creating new AIStudio Prompt Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Archive: error while processing AIStudio Prompt Agent request")
		return nil, err
	}
	response := &pb.PromptResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent action executed successfully", "200")
	return response, nil
}
