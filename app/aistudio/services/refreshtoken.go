package services

import (
	"context"
	"log"
	"time"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	utils "github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	grpccodes "google.golang.org/grpc/codes"
)

func (dms *AIStudioServices) RefreshTokenAgentHandler(ctx context.Context, managerRequest *pb.RefreshTokenManagerRequest, getterRequest *pb.RefreshTokenGetterRequest) (*pb.RefreshTokenResponse, error) {
	validTill := time.Now().Add(utils.DEFAULT_PLATFORM_RT_VALIDITY)
	log.Println("validTill: ", validTill)
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	if managerRequest != nil {
		switch managerRequest.GetUtility() {
		case pb.RefreshTokenAgentUtility_GENERATE_REFRESH_TOKEN:
			_, _, err := dms.GeneratePlatformAccessToken(ctx, "", "", "request.GetOrganization()", vapusPlatformClaim)
			if err != nil {
				return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.Unauthenticated)
			}
			return nil, nil
		case pb.RefreshTokenAgentUtility_REVOKE_REFRESH_TOKEN:
			_, _, err := dms.GeneratePlatformAccessToken(ctx, "", "", "request.GetDomai n()", vapusPlatformClaim)
			if err != nil {
				return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.Unauthenticated)
			}
			return nil, nil
		default:
			return nil, dmerrors.DMError(apperr.ErrInvalidAccessTokenAgentUtility, nil)
		}
	}
	return nil, nil
}
