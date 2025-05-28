package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type VapusGuardrailPlugins struct {
	pb.UnimplementedGuardrailPluginsServer
	validator  *dmutils.DMValidator
	DMServices *services.AIStudioServices
	Logger     zerolog.Logger
}

var VapusGuardrailPluginsManager *VapusGuardrailPlugins

func NewVapusGuardrailPlugins() *VapusGuardrailPlugins {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusGuardrailPlugins")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusGuardrailPlugins Controller initialized")
	return &VapusGuardrailPlugins{
		validator:  validator,
		Logger:     l,
		DMServices: services.AIStudioServiceManager,
	}
}

func InitVapusGuardrailPlugins() {
	if VapusGuardrailPluginsManager == nil {
		VapusGuardrailPluginsManager = NewVapusGuardrailPlugins()
	}
}

func (v *AIModels) ListBedrock(ctx context.Context, req *pb.GuardrailsTypeGetterRequest) (*pb.GuardrailsTypeResponse, error) {
	agent, err := v.DMServices.NewGuardrailPluginsIntAgent(ctx, services.WithGuardrailPluginsManagerRequest(req), services.WithGuardrailPluginsManagerAction(mpb.ResourceLcActions_LIST))
	if err != nil {
		v.Logger.Error().Err(err).Msg("List: error while creating new guardrail plugin request")
		return nil, err
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		v.Logger.Error().Err(err).Msg("List: error while processing guardrail plugin request")
		return nil, err
	}
	response := agent.GetResult()
	response.DmResp = pbtools.HandleDMResponse(ctx, "List: guardrail plugin action executed successfully", "200")
	return response, nil
}
