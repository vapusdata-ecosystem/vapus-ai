package dmcontrollers

import (
	"context"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	services "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type PluginController struct {
	dpb.UnimplementedPluginServiceServer
	DMServices *services.AIStudioServices
	logger     zerolog.Logger
}

var PluginControllerManager *PluginController

func NewPluginController() *PluginController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "PluginController")
	l.Info().Msg("Organization Controller initialized")
	return &PluginController{
		DMServices: services.AIStudioServiceManager,
		logger:     l,
	}
}

func InitPluginController() {
	if PluginControllerManager == nil {
		PluginControllerManager = NewPluginController()
	}
}

func (x *PluginController) Create(ctx context.Context, request *dpb.PluginManagerRequest) (*mpb.VapusCreateResponse, error) {
	agent, err := x.DMServices.NewPluginManagerAgent(ctx, services.WithPluginAgentManagerRequest(request), services.WithPluginAgentManagerAction(mpb.ResourceLcActions_ADD.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetCreateResponse()
	response.DmResp = pbtools.HandleDMResponse(ctx, "PluginManager action executed successfully", "201")
	return response, nil
}

func (x *PluginController) Update(ctx context.Context, request *dpb.PluginManagerRequest) (*dpb.PluginResponse, error) {
	agent, err := x.DMServices.NewPluginManagerAgent(ctx, services.WithPluginAgentManagerRequest(request), services.WithPluginAgentManagerAction(mpb.ResourceLcActions_UPDATE.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := &dpb.PluginResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "PluginManager action executed successfully", "201")
	return response, nil
}

func (x *PluginController) List(ctx context.Context, request *dpb.PluginGetterRequest) (*dpb.PluginResponse, error) {
	agent, err := x.DMServices.NewPluginManagerAgent(ctx, services.WithPluginAgentGetterRequest(request), services.WithPluginAgentManagerAction(mpb.ResourceLcActions_LIST.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := &dpb.PluginResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "PluginGetter action executed successfully", "201")
	return response, nil
}

func (x *PluginController) Get(ctx context.Context, request *dpb.PluginGetterRequest) (*dpb.PluginResponse, error) {
	agent, err := x.DMServices.NewPluginManagerAgent(ctx, services.WithPluginAgentGetterRequest(request), services.WithPluginAgentManagerAction(mpb.ResourceLcActions_GET.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := &dpb.PluginResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "PluginGetter action executed successfully", "201")
	return response, nil
}

func (x *PluginController) Archive(ctx context.Context, request *dpb.PluginGetterRequest) (*dpb.PluginResponse, error) {
	agent, err := x.DMServices.NewPluginManagerAgent(ctx, services.WithPluginAgentGetterRequest(request), services.WithPluginAgentManagerAction(mpb.ResourceLcActions_ARCHIVE.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := &dpb.PluginResponse{
		Output: agent.GetResult(),
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "PluginGetter action executed successfully", "201")
	return response, nil
}

func (x *PluginController) Action(ctx context.Context, request *dpb.PluginActionRequest) (*dpb.PluginActionResponse, error) {
	agent, err := x.DMServices.NewPluginActionsAgent(ctx, request)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	return &dpb.PluginActionResponse{
		DmResp: pbtools.HandleDMResponse(ctx, "Email sent successfully", "201"),
	}, nil
}
