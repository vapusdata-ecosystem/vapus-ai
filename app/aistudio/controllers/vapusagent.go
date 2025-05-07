package dmcontrollers

import (
	"github.com/rs/zerolog"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type VapusAgentServer struct {
	pb.UnimplementedAgentServiceServer
	validator  *dmutils.DMValidator
	DMServices *services.AIStudioServices
	log        zerolog.Logger
}

var VapusAgentServerManager *VapusAgentServer

func NewVapusAgentServer() *VapusAgentServer {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIPrompts")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("NewAIPrompts: Error while loading validator")
	}

	l.Info().Msg("AIPrompts Controller initialized")
	return &VapusAgentServer{
		log:        pkgs.GetSubDMLogger("VapusAgentServer", "VapusAgentServer"),
		DMServices: services.AIStudioServiceManager,
		validator:  validator,
	}
}

func InitVapusAgentServerController() {
	if VapusAgentServerManager == nil {
		VapusAgentServerManager = NewVapusAgentServer()
	}
}

// func (x *VapusAgentServer) Create(ctx context.Context, request *pb.AgentManagerRequest) (*mpb.VapusCreateResponse, error) {
// 	log.Println("Create", request)
// 	agent, err := x.DMServices.NewVapusAgentManager(ctx, services.WithVapusAgentAction(mpb.ResourceLcActions_ADD.String()), services.WithVapusAgentRequest(&services.VapusAgentManagerRequest{
// 		ManageRequest: request,
// 	}))
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	response := agent.GetCreateResponse()
// 	response.DmResp = pbtools.HandleDMResponse(ctx, "VapusAgent create action executed successfully", "200")
// 	return response, nil
// }

// func (x *VapusAgentServer) Update(ctx context.Context, request *pb.AgentManagerRequest) (*pb.AgentResponse, error) {
// 	log.Println("Update", request)
// 	agent, err := x.DMServices.NewVapusAgentManager(ctx, services.WithVapusAgentAction(mpb.ResourceLcActions_UPDATE.String()), services.WithVapusAgentRequest(&services.VapusAgentManagerRequest{
// 		ManageRequest: request,
// 	}))
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	return agent.GetManagerResponse(), nil
// }

// func (x *VapusAgentServer) Archive(ctx context.Context, request *pb.AgentGetterRequest) (*pb.AgentResponse, error) {
// 	log.Println("Archive", request)
// 	agent, err := x.DMServices.NewVapusAgentManager(ctx, services.WithVapusAgentAction(mpb.ResourceLcActions_ARCHIVE.String()), services.WithVapusAgentRequest(&services.VapusAgentManagerRequest{
// 		GetterRequest: request,
// 	}))
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	return agent.GetManagerResponse(), nil
// }

// func (x *VapusAgentServer) Get(ctx context.Context, request *pb.AgentGetterRequest) (*pb.AgentResponse, error) {
// 	agent, err := x.DMServices.NewVapusAgentManager(ctx, services.WithVapusAgentAction(mpb.ResourceLcActions_GET.String()), services.WithVapusAgentRequest(&services.VapusAgentManagerRequest{
// 		GetterRequest: request,
// 	}))
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	return agent.GetManagerResponse(), nil
// }

// func (x *VapusAgentServer) List(ctx context.Context, request *pb.AgentGetterRequest) (*pb.AgentResponse, error) {
// 	agent, err := x.DMServices.NewVapusAgentManager(ctx, services.WithVapusAgentAction(mpb.ResourceLcActions_GET.String()), services.WithVapusAgentRequest(&services.VapusAgentManagerRequest{
// 		GetterRequest: request,
// 	}))
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	return agent.GetManagerResponse(), nil
// }

// func (x *VapusAgentServer) ManageState(ctx context.Context, request *pb.AgentStateRequest) (*pb.AgentResponse, error) {
// 	log.Println("ManageState", request)
// 	agent, err := x.DMServices.NewVapusAgentManager(ctx, services.WithVapusAgentAction(request.Action.String()), services.WithVapusAgentRequest(&services.VapusAgentManagerRequest{
// 		StateRequest: request,
// 	}))
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	return agent.GetManagerResponse(), nil
// }

// func (x *VapusAgentServer) Validate(request *pb.AgentSignalRequest, stream pb.AgentService_ValidateServer) error {
// 	log.Println("ManageState", request)
// 	ctx := stream.Context()
// 	agent, err := x.DMServices.NewVapusAgentManager(ctx, services.WithVapusAgentAction(mpb.ResourceLcActions_VALIDATE.String()), services.WithVapusAgentRequest(&services.VapusAgentManagerRequest{
// 		SignalRequest: request,
// 	}))
// 	if err != nil {
// 		return pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	return nil
// }
