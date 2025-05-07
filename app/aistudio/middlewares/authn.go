package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/aistudio/pkgs"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	status "google.golang.org/grpc/status"
)

var AuthzMiddlewareMap = map[string]func(context.Context, string) (context.Context, error){}

// Initiate authenticator function for DataMesh JWT Authenication
// This function will be used as a middleware to authenticate the request
func AuthnMiddleware(ctx context.Context) (context.Context, error) {
	methodName, _ := grpc.Method(ctx)
	if !needAuthn(methodName) {
		return ctx, nil
	}
	logger = pkgs.GetSubDMLogger("Middleware", "Authn")
	logger.Info().Msgf("Authenticating request for method - %v", methodName)
	token, err := rpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		logger.Err(err).Msg("error while obtaining token from request header")
		return nil, status.Error(codes.Unauthenticated, "Authentication bearer token not found in request metadata")
	}
	if val, ok := AuthzMiddlewareMap[methodName]; ok {
		return val(ctx, token)
	} else {
		return authnPlatformAccess(ctx, token)
	}
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
	user, err := pkgs.VapusSvcInternalClientManager.GetUser(ctx, parsedClaims[encryption.ClaimUserIdKey], logger, 3)
	if err != nil {
		logger.Err(err).Msg("error while fetching user details from request header")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if !user.ValidateJwtClaim(parsedClaims) {
		logger.Err(err).Msg("error while validating jwt claims")
		return nil, status.Error(codes.Unauthenticated, "Invalid JWT claims")
	}
	logger.Info().Msgf("parsed domain claims - %v", parsedClaims)
	parsedClaims[encryption.ClaimUserNameKey] = user.GetDisplayName()
	return encryption.SetCtxClaim(ctx, parsedClaims), nil
}

func needAuthn(funcName string) bool {
	AnonymousFuncs := []string{
		grpc_reflection_v1.ServerReflection_ServerReflectionInfo_FullMethodName,
	}
	for _, f := range AnonymousFuncs {
		if f == funcName {
			return false
		}
	}
	return true
}
