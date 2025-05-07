package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/aistudio/services"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

type VapusAIGuardrails struct {
	pb.UnimplementedAIGuardrailsServer
	validator  *dmutils.DMValidator
	DMServices *services.AIStudioServices
	Logger     zerolog.Logger
}

var VapusAIGuardrailsManager *VapusAIGuardrails

func NewVapusAIGuardrails() *VapusAIGuardrails {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusAIGuardrails")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusAIGuardrails Controller initialized")
	return &VapusAIGuardrails{
		validator:  validator,
		Logger:     l,
		DMServices: services.AIStudioServiceManager,
	}
}

func InitVapusAIGuardrails() {
	if VapusAIGuardrailsManager == nil {
		VapusAIGuardrailsManager = NewVapusAIGuardrails()
	}
}

func (v *VapusAIGuardrails) Create(ctx context.Context, req *pb.GuardrailsManagerRequest) (*mpb.VapusCreateResponse, error) {

	agent, err := v.DMServices.NewAIGuardrailIntAgent(ctx, services.WithGuardrailManagerRequest(req), services.WithGuardrailManagerAction(mpb.ResourceLcActions_CREATE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while creating AI guardrails manager request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while processing AI guardrails creation request")
		return nil, err
	}
	response := agent.GetCreateResponse()
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIGuardrail create action executed successfully", "200")
	return response, nil
}

func (v *VapusAIGuardrails) Update(ctx context.Context, req *pb.GuardrailsManagerRequest) (*pb.GuardrailsResponse, error) {
	agent, err := v.DMServices.NewAIGuardrailIntAgent(ctx, services.WithGuardrailManagerRequest(req), services.WithGuardrailManagerAction(mpb.ResourceLcActions_UPDATE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while updating AI guardrails manager request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while processing AI guardrails update request")
		return nil, err
	}
	response := &pb.GuardrailsResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent update action executed successfully", "200")
	return response, nil
}

func (v *VapusAIGuardrails) Get(ctx context.Context, req *pb.GuardrailsGetterRequest) (*pb.GuardrailsResponse, error) {
	agent, err := v.DMServices.NewAIGuardrailIntAgent(ctx, services.WithGuardrailGetterRequest(req), services.WithGuardrailManagerAction(mpb.ResourceLcActions_GET))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while fetching AI guardrails")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while processing AI guardrails fetch request")
		return nil, err
	}
	response := &pb.GuardrailsResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent get action executed successfully", "200")
	return response, nil
}

func (v *VapusAIGuardrails) List(ctx context.Context, req *pb.GuardrailsGetterRequest) (*pb.GuardrailsResponse, error) {
	agent, err := v.DMServices.NewAIGuardrailIntAgent(ctx, services.WithGuardrailGetterRequest(req), services.WithGuardrailManagerAction(mpb.ResourceLcActions_LIST))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while listing AI guardrails")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while processing AI guardrails list request")
		return nil, err
	}
	response := &pb.GuardrailsResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent list action executed successfully", "200")
	return response, nil
}

func (v *VapusAIGuardrails) Archive(ctx context.Context, req *pb.GuardrailsGetterRequest) (*pb.GuardrailsResponse, error) {
	agent, err := v.DMServices.NewAIGuardrailIntAgent(ctx, services.WithGuardrailGetterRequest(req), services.WithGuardrailManagerAction(mpb.ResourceLcActions_ARCHIVE))
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while deleting AI guardrails")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("Error while processing AI guardrails delete request")
		return nil, err
	}
	response := &pb.GuardrailsResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "AIModelPromptConfigAgent delete action executed successfully", "200")
	return response, nil
}
