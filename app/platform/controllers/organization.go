package dmcontrollers

// Package controllers provides the implementation of organization controllers.
// These controllers handle the business logic for the organization package.

import (
	"context"

	"github.com/rs/zerolog"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	services "github.com/vapusdata-ecosystem/vapusdata/platform/services"
	utils "github.com/vapusdata-ecosystem/vapusdata/platform/utils"
	grpccodes "google.golang.org/grpc/codes"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
)

type OrganizationController struct {
	pb.UnimplementedOrganizationServiceServer
	DMServices *services.DMServices
	Logger     zerolog.Logger
}

var OrganizationControllerManager *OrganizationController

func NewOrganizationController() *OrganizationController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "OrganizationController")
	l.Info().Msg("Organization Controller initialized")
	return &OrganizationController{
		DMServices: services.DMServicesManager,
		Logger:     l,
	}
}

func InitOrganizationController() {
	if OrganizationControllerManager == nil {
		OrganizationControllerManager = NewOrganizationController()
	}
}

func (dmc *OrganizationController) Dashboard(ctx context.Context, request *mpb.EmptyRequest) (*pb.OrganizationDashboardResponse, error) {
	resp, err := dmc.DMServices.DashboardSvc(ctx)
	if err != nil {
		dmc.Logger.Err(err).Ctx(ctx).Msg("Error while getting dashboard info")
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal)
	}
	return &pb.OrganizationDashboardResponse{
		Output: resp,
		DmResp: pbtools.HandleDMResponse(ctx, utils.ACCOUNT_CREATED, "200"),
	}, nil
}

func (nc *OrganizationController) Create(ctx context.Context, request *pb.OrganizationManagerRequest) (*mpb.VapusCreateResponse, error) {
	agent, err := nc.DMServices.NewOrganizationAgent(ctx, services.WithDmAgentManagerRequest(request), services.WithDmAgentManagerAction(mpb.ResourceLcActions_ADD.String()))
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
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (nc *OrganizationController) Update(ctx context.Context, request *pb.OrganizationManagerRequest) (*pb.OrganizationResponse, error) {
	agent, err := nc.DMServices.NewOrganizationAgent(ctx, services.WithDmAgentManagerRequest(request), services.WithDmAgentManagerAction(mpb.ResourceLcActions_UPDATE.String()))
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
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (nc *OrganizationController) List(ctx context.Context, request *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	agent, err := nc.DMServices.NewOrganizationAgent(ctx, services.WithDmAgentGetterRequest(request), services.WithDmAgentManagerAction(mpb.ResourceLcActions_LIST.String()))
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
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (nc *OrganizationController) Get(ctx context.Context, request *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	agent, err := nc.DMServices.NewOrganizationAgent(ctx, services.WithDmAgentGetterRequest(request), services.WithDmAgentManagerAction(mpb.ResourceLcActions_GET.String()))
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
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (nc *OrganizationController) Archive(ctx context.Context, request *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	agent, err := nc.DMServices.NewOrganizationAgent(ctx, services.WithDmAgentGetterRequest(request), services.WithDmAgentManagerAction(mpb.ResourceLcActions_ARCHIVE.String()))
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
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (nc *OrganizationController) UpgradeOS(ctx context.Context, request *pb.OrganizationGetterRequest) (*pb.OrganizationResponse, error) {
	agent, err := nc.DMServices.NewOrganizationAgent(ctx, services.WithDmAgentGetterRequest(request), services.WithDmAgentManagerAction(mpb.ResourceLcActions_UPGRADE.String()))
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
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (nc *OrganizationController) AddUsers(ctx context.Context, request *pb.OrganizationAdduserRequest) (*pb.OrganizationResponse, error) {
	agent, err := nc.DMServices.NewOrganizationAgent(ctx, services.WithDmAgentUserManagerRequest(request), services.WithDmAgentManagerAction(mpb.ResourceLcActions_ADD_USERS.String()))
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
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}
