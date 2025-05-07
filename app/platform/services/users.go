package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	apppdrepo "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	authn "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	utils "github.com/vapusdata-ecosystem/vapusdata/platform/utils"
	grpccodes "google.golang.org/grpc/codes"
)

func (dms *DMServices) AccessTokenAgentHandler(ctx context.Context, request *pb.AccessTokenInterfaceRequest) (*pb.AccessTokenResponse, error) {
	validTill := time.Now().Add(utils.DEFAULT_PLATFORM_AT_VALIDITY)
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		dms.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	switch request.GetUtility() {
	case pb.AccessTokenAgentUtility_ORGANIZATION_LOGIN:
		accessToken, scope, err := dms.GeneratePlatformAccessToken(ctx, request.GetIdToken(), "", request.GetOrganization(), vapusPlatformClaim)
		if err != nil {
			log.Println("error while generating token", err, err == apperr.ErrUser404)
			if errors.Is(err, apperr.ErrUser404) {
				return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.NotFound)
			}
			dms.logger.Err(err).Msg("error while generating token")
			return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.Unauthenticated)
		}
		return &pb.AccessTokenResponse{
			Token: &pb.AccessToken{
				AccessToken: accessToken,
				ValidTill:   validTill.Unix(),
			},
			TokenScope: scope,
		}, nil
	case pb.AccessTokenAgentUtility_REFRESH_TOKEN_LOGIN:
		accessToken, scope, err := dms.GeneratePlatformAccessToken(ctx, "", request.GetRefreshToken(), request.GetOrganization(), vapusPlatformClaim)
		if err != nil {
			dms.logger.Err(err).Msg("error while generating token")
			return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.Unauthenticated)
		}
		return &pb.AccessTokenResponse{
			Token: &pb.AccessToken{
				AccessToken: accessToken,
				ValidTill:   validTill.Unix(),
			},
			TokenScope: scope,
		}, nil
	default:
		return nil, dmerrors.DMError(apperr.ErrInvalidAccessTokenAgentUtility, nil)
	}
}

func (dms *DMServices) SignupHandler(ctx context.Context, request *pb.RegisterUserRequest) (*pb.AccessTokenResponse, error) {
	validTill := time.Now().Add(utils.DEFAULT_PLATFORM_AT_VALIDITY)
	vapusPlatformClaim := map[string]string{
		encryption.ClaimAccountKey: dmstores.AccountPool.VapusID,
	}
	claims, err := authn.ValidateOIDCAuth(request.GetIdToken(), dms.logger)
	if err != nil {
		dms.logger.Err(err).Msg("invalid token, validation failed")
		return nil, dmerrors.DMError(err, nil)
	}
	lu, err := claimToLocal(claims)
	if err != nil {
		dms.logger.Err(err).Msg("invalid token, relevant claims are missing")
		return nil, dmerrors.DMError(err, nil)
	}
	userObj, err := dms.DMStore.CreateUser(ctx, lu, nil, vapusPlatformClaim)
	if err != nil {
		dms.logger.Err(err).Msg("error while creating user")
		return nil, dmerrors.DMError(err, nil)
	}
	if request.GetOrganization() == "" {
		dm := strings.Split(userObj.Email, "@")
		request.Organization = strings.ToTitle(dm[0]) + "." + strings.ToTitle(dm[1]) + " organization"
	}
	vapusPlatformClaim[encryption.ClaimUserIdKey] = userObj.UserId
	organization, err := organizationConfigureTool(ctx, &models.Organization{
		Name: request.GetOrganization(),
	}, dms.DMStore, dms.logger, vapusPlatformClaim)
	if err != nil {
		dms.logger.Err(err).Msg("error while configuring organization")
		return nil, dmerrors.DMError(err, nil)
	}
	userObj.Roles = []*models.UserOrganizationRole{
		{
			OrganizationId: organization.VapusID,
			RoleArns:       []string{mpb.UserRoles_ORG_OWNER.String()},
			InvitedOn:      dmutils.GetEpochTime(),
		},
	}
	dms.DMStore.PutUser(ctx, userObj, vapusPlatformClaim)
	accessToken, scope, err := dms.GeneratePlatformAccessToken(ctx, request.GetIdToken(), "", organization.VapusID, vapusPlatformClaim)
	if err != nil {
		return nil, pbtools.HandleGrpcError(dmerrors.DMError(apperr.ErrUnAuthenticated, err), grpccodes.Unauthenticated)
	}
	return &pb.AccessTokenResponse{
		Token: &pb.AccessToken{
			AccessToken: accessToken,
			ValidTill:   validTill.Unix(),
		},
		TokenScope: scope,
	}, nil
}

func (dms *VapusDataServices) GeneratePlatformAccessToken(ctx context.Context, idToken, refreshToken, organization string, ctxClaim map[string]string) (string, mpb.AccessTokenScope, error) {
	var userObj *models.Users
	var organizationObj *models.Organization
	var err error
	var lu *apppdrepo.LocalUserM
	if idToken != "" {
		claims, err := authn.ValidateOIDCAuth(idToken, dms.Logger)
		if err != nil {
			dms.Logger.Err(err).Msg("invalid token, validation failed")
			return "", 0, dmerrors.DMError(err, nil)
		}
		lu, err = claimToLocal(claims)
		if err != nil {
			dms.Logger.Err(err).Msg("invalid token, relevant claims are missing")
			return "", 0, dmerrors.DMError(err, nil)
		}
	} else if refreshToken != "" {
		refreshToken = encryption.GenerateSHA256(refreshToken, "")
		rtObj, err := dms.DMStore.GetPlatformRTinfo(ctx, refreshToken, ctxClaim)
		if err != nil {
			dms.Logger.Err(err).Msg("invalid token, refresh token doesn't exists")
			return "", 0, dmerrors.DMError(apperr.ErrRefreshToken404, err)
		}
		if rtObj.ValidTill < dmutils.GetEpochTime() {
			dms.Logger.Err(err).Msg("invalid token, refresh token expired")
			return "", 0, dmerrors.DMError(apperr.ErrRefreshTokenExpired, nil)
		}
		if rtObj.Status != mpb.CommonStatus_ACTIVE.String() {
			dms.Logger.Err(err).Msg("invalid token, refresh token is not active")
			return "", 0, dmerrors.DMError(apperr.ErrRefreshTokenInactive, nil)
		}
		lu = &apppdrepo.LocalUserM{
			Email: rtObj.UserId,
		}
		organization = rtObj.Organization
	} else {
		lu = &apppdrepo.LocalUserM{
			Email: ctxClaim[encryption.ClaimUserIdKey],
		}
	}
	validTill := time.Now().Add(utils.DEFAULT_PLATFORM_AT_VALIDITY)
	if organization == "" {
		userObj, err = dms.DMStore.GetOrUpdateUser(ctx, lu, true, true, ctxClaim)
		if err != nil {
			dms.Logger.Err(err).Msg("invalid token, user doesn't exists")
			return idToken, 0, apperr.ErrUser404
		}
		organizationObj, err = dms.DMStore.GetOrganization(ctx, userObj.Roles[0].OrganizationId, nil)
		if err != nil {
			dms.Logger.Err(err).Msg("invalid token, organization doesn't exists")
			return "", 0, dmerrors.DMError(apperr.ErrOrganization404, err)
		}
		dms.Logger.Info().Msgf("User obtained for generating token with default organization- %v ", userObj)
	} else {
		lu.Organization = organization
		organizationObj, err = dms.DMStore.GetOrganization(ctx, lu.Organization, nil)
		if err != nil {
			dms.Logger.Err(err).Msg("invalid token, organization doesn't exists")
			return "", 0, dmerrors.DMError(apperr.ErrOrganization404, err)
		}
		userObj, err = dms.DMStore.GetOrUpdateUser(ctx, lu, false, false, ctxClaim)
		if err != nil {
			dms.Logger.Err(err).Msg("invalid token, user doesn't exists")
			return "", 0, dmerrors.DMError(err, nil)
		}

		if !userObj.IsValidUserByOrganization(organizationObj.VapusID) {
			dms.Logger.Err(err).Msg("invalid request, user is not present in this organization")
			return "", 0, dmerrors.DMError(err, nil)
		}
		dms.Logger.Info().Msgf("User obtained - %v ", userObj)
	}
	dms.Logger.Info().Msgf("Organization obtained - %v ", organizationObj)
	dms.Logger.Info().Msgf("User obtained here - %v ", userObj)
	tokenId, token, err := generateOrganizationAccessToken(userObj, organizationObj.VapusID, validTill)
	if err != nil {
		return "", 0, dmerrors.DMError(err, nil)
	}

	newCtx := context.TODO()
	go dms.DMStore.LogPlatformJwtinfo(newCtx, &models.JwtLog{
		JwtId:        tokenId,
		UserId:       userObj.UserId,
		Organization: organizationObj.VapusID,
		Scope:        mpb.AccessTokenScope_ORG_TOKEN.String(),
	}, ctxClaim)
	return token, mpb.AccessTokenScope_ORG_TOKEN, nil
}
