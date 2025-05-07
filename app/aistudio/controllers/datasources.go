package dmcontrollers

// Package controllers provides the implementation of organization controllers.
// These controllers handle the business logic for the organization package.

import (
	"context"

	"github.com/rs/zerolog"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	services "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	grpccodes "google.golang.org/grpc/codes"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
)

type DataSourcesController struct {
	dpb.UnimplementedDatasourceServiceServer
	DMServices *services.AIStudioServices
	Logger     zerolog.Logger
}

var DataSourcesControllerManager *DataSourcesController

func NewDataSourcesController() *DataSourcesController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "DataSourcesController")
	l.Info().Msg("Organization Controller initialized")
	return &DataSourcesController{
		DMServices: services.AIStudioServiceManager,
		Logger:     l,
	}
}

func InitDataSourcesController() {
	if DataSourcesControllerManager == nil {
		DataSourcesControllerManager = NewDataSourcesController()
	}
}

func (nc *DataSourcesController) Create(ctx context.Context, request *dpb.DataSourceManagerRequest) (*mpb.VapusCreateResponse, error) {
	agent, err := nc.DMServices.NewDataSourceAgent(ctx, services.WithDsAgentManagerRequest(request), services.WithDsAgentManagerAction(mpb.ResourceLcActions_ADD.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx, "")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetCreateResponse()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Data Source creation is initiated successfully", "201")
	return response, nil
}

func (nc *DataSourcesController) Get(ctx context.Context, request *dpb.DataSourceGetterRequest) (*dpb.DataSourceResponse, error) {
	agent, err := nc.DMServices.NewDataSourceAgent(ctx, services.WithDsAgentGetterRequest(request), services.WithDsAgentManagerAction(mpb.ResourceLcActions_GET.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx, "")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "NewDataSourceAgent action executed successfully", "200")
	return response, nil
}

func (nc *DataSourcesController) Update(ctx context.Context, request *dpb.DataSourceManagerRequest) (*dpb.DataSourceResponse, error) {
	agent, err := nc.DMServices.NewDataSourceAgent(ctx, services.WithDsAgentManagerRequest(request), services.WithDsAgentManagerAction(mpb.ResourceLcActions_UPDATE.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx, "")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "NewDataSourceAgent action executed successfully", "200")
	return response, nil
}

func (nc *DataSourcesController) List(ctx context.Context, request *dpb.DataSourceGetterRequest) (*dpb.DataSourceResponse, error) {
	agent, err := nc.DMServices.NewDataSourceAgent(ctx, services.WithDsAgentGetterRequest(request), services.WithDsAgentManagerAction(mpb.ResourceLcActions_LIST.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx, "")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "NewDataSourceAgent action executed successfully", "200")
	return response, nil
}

func (nc *DataSourcesController) Archive(ctx context.Context, request *dpb.DataSourceGetterRequest) (*dpb.DataSourceResponse, error) {
	agent, err := nc.DMServices.NewDataSourceAgent(ctx, services.WithDsAgentGetterRequest(request), services.WithDsAgentManagerAction(mpb.ResourceLcActions_ARCHIVE.String()))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	defer func() {
		dmutils.CleanPointers(agent)
	}()
	err = agent.Act(ctx, "")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Data Source creation is archived successfully", "200")
	return response, nil
}
