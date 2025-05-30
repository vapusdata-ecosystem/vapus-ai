package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	services "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	grpccodes "google.golang.org/grpc/codes"
)

type UtilityController struct {
	dpb.UnimplementedUtilityServiceServer
	DMServices *services.AIStudioServices
	logger     zerolog.Logger
}

var UtilityControllerManager *UtilityController

func NewUtilityController() *UtilityController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "UtilityController")
	l.Info().Msg("UtilityController initialized")
	return &UtilityController{
		DMServices: services.AIStudioServiceManager,
		logger:     l,
	}
}

func InitUtilityController() {
	if UtilityControllerManager == nil {
		UtilityControllerManager = NewUtilityController()
	}
}

func (x *UtilityController) Upload(ctx context.Context, request *dpb.UploadRequest) (*dpb.UploadResponse, error) {
	utAgent, err := x.DMServices.NewUtilityAgent(ctx, request, nil, nil)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = utAgent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	resp := utAgent.GetUploadedResult()
	return resp, nil
}

func (x *UtilityController) UploadStream(stream dpb.UtilityService_UploadStreamServer) error {
	ctx := stream.Context()
	utAgent, err := x.DMServices.NewUtilityAgent(stream.Context(), nil, stream, nil)
	if err != nil {
		return pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = utAgent.Act(ctx)
	if err != nil {
		return pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	return nil
}

func (x *UtilityController) Download(ctx context.Context, request *dpb.DownloadRequest) (*dpb.DownloadResponse, error) {
	utAgent, err := x.DMServices.NewUtilityAgent(ctx, nil, nil, request)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = utAgent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	resp := utAgent.GetDownloadResult()
	return resp, nil
}
