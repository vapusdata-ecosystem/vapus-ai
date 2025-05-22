package dmcontrollers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	grpccodes "google.golang.org/grpc/codes"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	dmsvc "github.com/vapusdata-ecosystem/vapusai/aistudio/services"
	utils "github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	dmauthn "github.com/vapusdata-ecosystem/vapusai/core/pkgs/authn"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	pbtools "github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type VapusDataUsers struct {
	pb.UnimplementedUserManagementServiceServer
	validator  *dmutils.DMValidator
	DMServices *dmsvc.AIStudioServices
	logger     zerolog.Logger
}

var VapusDataUsersManager *VapusDataUsers

func NewVapusDataUsers() *VapusDataUsers {
	l := pkgs.GetSubDMLogger(pkgs.CNTRLR, "VapusDataUsers")
	validator, err := dmutils.NewDMValidator()
	if err != nil {
		l.Panic().Err(err).Msg("Error while loading validator")
	}

	l.Info().Msg("VapusDataUsers Controller initialized")
	return &VapusDataUsers{
		validator:  validator,
		logger:     l,
		DMServices: dmsvc.AIStudioServiceManager,
	}
}

func InitVapusDataUsers() {
	if VapusDataUsersManager == nil {
		VapusDataUsersManager = NewVapusDataUsers()
	}
}

func (dmc *VapusDataUsers) AccessTokenInterface(ctx context.Context, request *pb.AccessTokenInterfaceRequest) (*pb.AccessTokenResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	return dmc.DMServices.AccessTokenAgentHandler(ctx, request)
}

func (dmc *VapusDataUsers) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.AccessTokenResponse, error) {
	dmc.logger.Info().Msg("Registering user.......")
	return dmc.DMServices.SignupHandler(ctx, request)
}

func (dmc *VapusDataUsers) UserManager(ctx context.Context, request *pb.UserManagerRequest) (*pb.UserResponse, error) {
	agent, err := dmc.DMServices.NewUserManagerAgent(ctx, request, nil)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act(ctx, "")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}
func (dmc *VapusDataUsers) UserGetter(ctx context.Context, request *pb.UserGetterRequest) (*pb.UserResponse, error) {
	agent, err := dmc.DMServices.NewUserManagerAgent(ctx, nil, request)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act(ctx, "")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (dmc *VapusDataUsers) RefreshTokenManager(ctx context.Context, request *pb.RefreshTokenManagerRequest) (*pb.RefreshTokenResponse, error) {
	response, err := dmc.DMServices.RefreshTokenAgentHandler(ctx, request, nil)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}
func (dmc *VapusDataUsers) RefreshTokenGetter(ctx context.Context, request *pb.RefreshTokenGetterRequest) (*pb.RefreshTokenResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	response, err := dmc.DMServices.RefreshTokenAgentHandler(ctx, nil, request)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (dmc *VapusDataUsers) AuthzManager(ctx context.Context, request *pb.AuthzManagerRequest) (*pb.AuthzResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	agent, err := dmc.DMServices.NewAuthzManagerAgent(ctx, request, nil)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act("")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (dmc *VapusDataUsers) AuthzGetter(ctx context.Context, request *pb.AuthzGetterRequest) (*pb.AuthzResponse, error) {
	dmc.logger.Info().Msg("Generating platform access token.......")
	agent, err := dmc.DMServices.NewAuthzManagerAgent(ctx, nil, request)
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	err = agent.Act("")
	if err != nil {
		return nil, pbtools.HandleGrpcError(err, grpccodes.Internal) //nolint:wrapcheck
	}
	response := agent.GetResult()
	agent.LogAgent()
	response.DmResp = pbtools.HandleDMResponse(ctx, "VdcDeployment action executed successfully", "200")
	return response, nil
}

func (dmc *VapusDataUsers) LoginHandler(ctx context.Context, request *mpb.EmptyRequest) (*pb.LoginHandlerResponse, error) {
	fmt.Println("I am in Login Handlerrrrrrrr")
	state, err := dmutils.GenerateRandomState()

	if err != nil {
		dmc.logger.Err(err).Msg(apperr.ErrAuthenticatorInitFailed.Error())
		return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrLoginGeneration, nil), grpccodes.Internal)
	}
	return &pb.LoginHandlerResponse{
		LoginUrl:    pkgs.AuthnManager.Authenticator.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "login")),
		CallbackUrl: pkgs.AuthnManager.RedirectURL,

		RedirectUri: pkgs.AuthnManager.Authenticator.Organization,
	}, nil

}

func (dmc *VapusDataUsers) LoginCallback(ctx context.Context, request *pb.LoginCallBackRequest) (*pb.AccessTokenResponse, error) {
	validTill := time.Now().Add(utils.DEFAULT_PLATFORM_AT_VALIDITY)
	dmc.logger.Info().Msg("Exchanging code for token.......")
	dmc.logger.Info().Msgf("Callback URL: %v ----- %v", request.GetHost(), request.GetCode())
	token, err := pkgs.AuthnManager.Authenticator.Exchange(ctx, request.GetCode(), oauth2.SetAuthURLParam("redirect_uri", request.GetHost()))
	if err != nil {
		dmc.logger.Err(err).Msg(dmauthn.ErrTokenExchangeFailed.Error())
		return nil, pbtools.HandleGrpcError(dmerrors.DMError(dmauthn.ErrTokenExchangeFailed, nil), grpccodes.Unauthenticated)

	}
	_, err = pkgs.AuthnManager.Authenticator.VerifyIDToken(ctx, token)

	if err != nil {
		dmc.logger.Err(err).Msg(dmauthn.ErrIDTokenVerificationFailed.Error())
		return nil, pbtools.HandleGrpcError(dmerrors.DMError(dmauthn.ErrIDTokenVerificationFailed, nil), grpccodes.Unauthenticated)

	}
	accessToken, scope, err := dmc.DMServices.GeneratePlatformAccessToken(ctx, token.Extra("id_token").(string), "", "", nil)
	if err != nil {
		log.Println("error while generating token", err, err == apperr.ErrUser404)
		if errors.Is(err, apperr.ErrUser404) {
			log.Println("user not found ", token.Extra("id_token").(string))
			return dmc.DMServices.SignupHandler(ctx, &pb.RegisterUserRequest{
				IdToken: token.Extra("id_token").(string),
			})
		}
		dmc.logger.Err(err).Msg("error while generating token")
		return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.Unauthenticated)
	}
	return &pb.AccessTokenResponse{
		Token: &pb.AccessToken{
			AccessToken: accessToken,
			ValidTill:   validTill.Unix(),
			ValidFrom:   dmutils.GetEpochTime(),
			IdToken:     token.Extra("id_token").(string),
		},
		TokenScope: scope,
		DmResp:     pbtools.HandleDMResponse(ctx, utils.ACCESS_TOKEN_CREATED, "201"),
	}, nil
}

// func (dmc *VapusDataUsers) LoginRefreshToken(ctx context.Context, request *pb.LoginCallBackRequest) (*pb.AccessTokenResponse, error) {
// 	validTill := time.Now().Add(utils.DEFAULT_PLATFORM_AT_VALIDITY)
// 	accessToken, scope, err := dmc.DMServices.GeneratePlatformAccessToken(ctx, "", "", "", nil)
// 	if err != nil {
// 		return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.Unauthenticated)
// 	}
// 	return &pb.AccessTokenResponse{
// 		Token: &pb.AccessToken{
// 			AccessToken: accessToken,
// 			ValidTill:   validTill.Unix(),
// 			ValidFrom:   dmutils.GetEpochTime(),
// 		},
// 		TokenScope: scope,
// 		DmResp:     pbtools.HandleDMResponse(ctx, utils.ACCESS_TOKEN_CREATED, "201"),
// 	}, nil
// }
