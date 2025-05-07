package dmcontrollers

import (
	"context"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	pbtools "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	services "github.com/vapusdata-ecosystem/vapusdata/platform/services"
	grpccodes "google.golang.org/grpc/codes"
)

type SecretsController struct {
	dpb.UnimplementedSecretServiceServer
	DMServices *services.DMServices
	logger     zerolog.Logger
}

var SecretsControllerManager *SecretsController

func NewSecretsController() *SecretsController {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "SecretsController")
	l.Info().Msg("SecretsController initialized")
	return &SecretsController{
		DMServices: services.DMServicesManager,
		logger:     l,
	}
}

func InitSecretsController() {
	if SecretsControllerManager == nil {
		SecretsControllerManager = NewSecretsController()
	}
}

func (x *SecretsController) Create(ctx context.Context, request *dpb.SecretManagerRequest) (*mpb.VapusCreateResponse, error) {
	smAgent, err := x.DMServices.NewSecretsManagerAgent(ctx, services.WithSecreManagerAction(mpb.ResourceLcActions_CREATE.String()), services.WithSecretManagerRequest(request))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = smAgent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := smAgent.GetCreateResponse()
	response.DmResp = pbtools.HandleDMResponse(ctx, "Secrets Create action executed successfully", "201")
	return response, nil
}

func (x *SecretsController) Update(ctx context.Context, request *dpb.SecretManagerRequest) (*dpb.VapusSecretsResponse, error) {
	smAgent, err := x.DMServices.NewSecretsManagerAgent(ctx, services.WithSecreManagerAction(mpb.ResourceLcActions_UPDATE.String()), services.WithSecretManagerRequest(request))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = smAgent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	resp := smAgent.GetResult()
	return resp, nil
}

func (x *SecretsController) List(ctx context.Context, request *dpb.SecretGetterRequest) (*dpb.VapusSecretsResponse, error) {
	smAgent, err := x.DMServices.NewSecretsManagerAgent(ctx, services.WithSecreManagerAction(mpb.ResourceLcActions_LIST.String()), services.WithSecretGetterRequest(request))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = smAgent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	resp := smAgent.GetResult()
	return resp, nil
}

func (x *SecretsController) Get(ctx context.Context, request *dpb.SecretGetterRequest) (*dpb.VapusSecretsResponse, error) {
	smAgent, err := x.DMServices.NewSecretsManagerAgent(ctx, services.WithSecreManagerAction(mpb.ResourceLcActions_GET.String()), services.WithSecretGetterRequest(request))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = smAgent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	resp := smAgent.GetResult()
	return resp, nil
}

func (x *SecretsController) Archive(ctx context.Context, request *dpb.SecretGetterRequest) (*dpb.VapusSecretsResponse, error) {
	smAgent, err := x.DMServices.NewSecretsManagerAgent(ctx, services.WithSecreManagerAction(mpb.ResourceLcActions_ARCHIVE.String()), services.WithSecretGetterRequest(request))
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = smAgent.Act(ctx)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	resp := smAgent.GetResult()
	return resp, nil
}
