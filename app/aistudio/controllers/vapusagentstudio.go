package dmcontrollers

import (
	"github.com/rs/zerolog"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/aistudio/services"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
)

type VapusAgentStudio struct {
	pb.UnimplementedAgentStudioServer
	validator  *dmutils.DMValidator
	DMServices *services.AIStudioServices
	log        zerolog.Logger
}

var VapusAgentStudioManager *VapusAgentStudio

func NewVapusAgentStudio() *VapusAgentStudio {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "AIPrompts")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("NewAIPrompts: Error while loading validator")
	}

	l.Info().Msg("AIPrompts Controller initialized")
	return &VapusAgentStudio{
		log:        pkgs.GetSubDMLogger("VapusAgentStudio", "VapusAgentStudio"),
		DMServices: services.AIStudioServiceManager,
		validator:  validator,
	}
}

func InitVapusAgentStudioController() {
	if VapusAgentStudioManager == nil {
		VapusAgentStudioManager = NewVapusAgentStudio()
	}
}

// func (x *VapusAgentStudio) DownloadFiles(ctx context.Context, request *pb.AgentDownloadFileRequest) (*pb.AgentDownloadFileResponse, error) {
// 	agent, err := x.DMServices.NewVapusAgentStudio(ctx, services.WithVapusAgentStudioAction(mpb.ResourceLcActions_DOWNLOAD.String()), services.WithVapusAgentSignal(&services.VapusAgentStudioRequest{
// 		DownloadRequest: request,
// 	}))
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.InvalidArgument) //nolint:wr
// 	}
// 	err = agent.Act(ctx)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wr
// 	}
// 	return agent.GetDownloadFileResponse(), nil
// }

// func (x *VapusAgentStudio) Signal(request *pb.AgentSignalRequest, stream pb.AgentStudio_SignalServer) error {
// 	ctx := stream.Context()
// 	agent, err := x.DMServices.NewVapusAgentStudio(ctx, services.WithVapusAgentSignal(&services.VapusAgentStudioRequest{
// 		SignalRequest: request,
// 		Stream:        stream,
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
