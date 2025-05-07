package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/aistudio/pkgs"
	services "github.com/vapusdata-ecosystem/vapusdata/aistudio/services"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

type AIModels struct {
	pb.UnimplementedAIModelsServer
	validator  *dmutils.DMValidator
	DMServices *services.AIStudioServices
	Logger     zerolog.Logger
}

var AIModelsManager *AIModels

func NewAIModels() *AIModels {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIModels")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("NewAIModels: Error while loading validator")
	}

	l.Info().Msg("NewAIModels: AIModels Controller initialized")
	return &AIModels{
		validator:  validator,
		Logger:     l,
		DMServices: services.AIStudioServiceManager,
	}
}

func InitAIModels() {
	if AIModelsManager == nil {
		AIModelsManager = NewAIModels()
	}
}

func (v *AIModels) Create(ctx context.Context, req *pb.AIModelNodeManagerRequest) (*mpb.VapusCreateResponse, error) {
	agent, err := v.DMServices.NewAIModelIntAgent(ctx, services.WithModelManagerRequest(req), services.WithModelManagerAction(mpb.ResourceLcActions_CREATE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Create: error while creating new AIStudio Node Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Create: error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetCreateResponse()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Create: AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIModels) Update(ctx context.Context, req *pb.AIModelNodeManagerRequest) (*pb.AIModelNodeResponse, error) {
	agent, err := v.DMServices.NewAIModelIntAgent(ctx, services.WithModelManagerRequest(req), services.WithModelManagerAction(mpb.ResourceLcActions_UPDATE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Update: error while creating new AIStudio Node Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Update: error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Update: AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIModels) Get(ctx context.Context, req *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	agent, err := v.DMServices.NewAIModelIntAgent(ctx, services.WithModelGetterRequest(req), services.WithModelManagerAction(mpb.ResourceLcActions_GET))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Get: error while creating new AIStudio Node Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Get: error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Get: AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIModels) List(ctx context.Context, req *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	agent, err := v.DMServices.NewAIModelIntAgent(ctx, services.WithModelGetterRequest(req), services.WithModelManagerAction(mpb.ResourceLcActions_LIST))
	if err != nil {
		v.Logger.Error().Err(err).Msg("List: error while creating new AIStudio Node Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("List: error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = pbtools.HandleDMResponse(ctx, "List: AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIModels) Archive(ctx context.Context, req *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	agent, err := v.DMServices.NewAIModelIntAgent(ctx, services.WithModelGetterRequest(req), services.WithModelManagerAction(mpb.ResourceLcActions_ARCHIVE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Archive: error while creating new AIStudio Node Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Archive: error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Archive: AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}

func (v *AIModels) Sync(ctx context.Context, req *pb.AIModelNodeGetterRequest) (*pb.AIModelNodeResponse, error) {
	agent, err := v.DMServices.NewAIModelIntAgent(ctx, services.WithModelGetterRequest(req), services.WithModelManagerAction(mpb.ResourceLcActions_SYNC))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Sync: error while creating new AIStudio Node Agent request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Sync: error while processing AIStudio Node Agent request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Sync: AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}

// Insights
func (v *AIModels) ListInsights(ctx context.Context, req *pb.AIModelNodeInsightsRequest) (*pb.AIModelNodeInsightsResponse, error) {
	agent, err := v.DMServices.NewAIModelIntAgent(ctx, services.WithModelNodeInsightsRequest(req), services.WithModelManagerAction(mpb.ResourceLcActions_INSIGHTS))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Insights: error while fetching Model usage Insights request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Insights: error while processing Model usage Insights request")
		return nil, err
	}
	response := agent.GetInsightResult()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Insights: AIModelNodeConfigAgent action executed successfully", "200")
	return response, nil
}
