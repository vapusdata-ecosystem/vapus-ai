package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	status "google.golang.org/grpc/status"
)

// var AuthzMiddlewareMap = map[string]func(context.Context, string) (context.Context, error){
// 	dpb.OrganizationService_DataProductServer_FullMethodName: authnDataProductAccess,
// }

// Initiate authenticator function for DataMarketplace JWT Authenication
// This function will be used as a middleware to authenticate the request
func AuthnMiddleware(ctx context.Context) (context.Context, error) {
	methodName, _ := grpc.Method(ctx)
	if !needAuthn(methodName) {
		return ctx, nil
	}
	mds, has := metadata.FromIncomingContext(ctx)
	if !has {
		return nil, status.Error(codes.Unauthenticated, "Authentication bearer token not found in request metadata")
	}
	log.Println("mds------------------------------------", mds)

	logger = pkgs.GetSubDMLogger("Middleware", "Authn")
	logger.Info().Msgf("Authenticating request for method - %v", methodName)
	token, err := rpcauth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		logger.Err(err).Msg("error while obtaining token from request header")
		return nil, status.Error(codes.Unauthenticated, "Authentication bearer token not found in request metadata")
	}
	// if val, ok := AuthzMiddlewareMap[methodName]; ok {
	// 	return val(ctx, token)
	// } else {
	return authnPlatformAccess(ctx, token)
	// }
}

func HttpAuthnMiddleware(ctx context.Context, req *http.Request) metadata.MD {
	token := req.Header.Get("Authorization")
	if token == "" {
		return metadata.Pairs("error", "Authentication bearer token not found in request metadata")
	}
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)
	ctx, err := authnPlatformAccess(ctx, token)
	if err != nil {
		return metadata.Pairs("error", err.Error())
	}
	bbyte, err := json.Marshal(ctx.Value(encryption.JwtDPCtxClaimKey))
	if err != nil {
		return metadata.Pairs("error", err.Error())
	}
	return metadata.Pairs(encryption.JwtDPCtxClaimKey, string(bbyte))
	// return nil
	// }
}

func authnPlatformAccess(ctx context.Context, token string) (context.Context, error) {
	parsedClaims, err := pkgs.SvcPackageManager.VapusJwtAuth.ValidateAccessToken(token)
	if err != nil {
		logger.Err(err).Msg("error while validating access token from request header")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	logger.Info().Msgf("parsed organization claims - %v", parsedClaims)
	user, err := dmstores.DMStoreManager.GetUser(ctx, parsedClaims[encryption.ClaimUserIdKey], parsedClaims)
	if err != nil || user == nil {
		logger.Err(err).Msgf("error while validating access token against the claim organization, user - %v, organization - %v", parsedClaims[encryption.ClaimUserIdKey], parsedClaims[encryption.ClaimOrganizationKey])
		return nil, status.Error(codes.Unauthenticated, apperr.ErrOrganization404.Error())
	} else if !user.ValidateJwtClaim(parsedClaims) {
		logger.Err(err).Msgf("error while validating access token against the claim organization, user - %v, organization - %v, role - %v", parsedClaims[encryption.ClaimUserIdKey], parsedClaims[encryption.ClaimOrganizationKey], parsedClaims[encryption.ClaimOrganizationRolesKey])
		return nil, status.Error(codes.Unauthenticated, apperr.ErrUserOrganization403.Error())
	}

	parsedClaims[encryption.ClaimUserNameKey] = user.GetDisplayName()
	return encryption.SetCtxClaim(ctx, parsedClaims), nil
}

func needAuthn(funcName string) bool {
	AnonymousFuncs := []string{
		pb.UserManagementService_LoginHandler_FullMethodName,
		pb.UserManagementService_LoginCallback_FullMethodName,
		grpc_reflection_v1.ServerReflection_ServerReflectionInfo_FullMethodName,
		pb.VapusdataService_VapusdataServicesInfo_FullMethodName,
		pb.VapusdataService_PlatformPublicInfo_FullMethodName,
		pb.UserManagementService_RegisterUser_FullMethodName,
	}
	for _, f := range AnonymousFuncs {
		if f == funcName {
			return false
		}
	}
	return true
}
