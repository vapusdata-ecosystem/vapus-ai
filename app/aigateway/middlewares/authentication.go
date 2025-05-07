package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	rpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/vapusdata-ecosystem/vapusai/aigateway/pkgs"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	"google.golang.org/grpc/metadata"
)

func Authentication(c *fiber.Ctx) error {
	if c.Path() == "/metrics" {
		return c.Next()
	}
	tokenBytes := c.Request().Header.Peek("Authorization")
	if tokenBytes == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	token := string(tokenBytes)
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)
	ctx, err := authnPlatformAccess(c.Context(), token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.SetUserContext(ctx)
	return c.Next()
}

func authnPlatformAccess(ctx context.Context, token string) (context.Context, error) {
	pkgs.DmLogger.Info().Msg("authenticating request")
	parsedClaims, err := pkgs.SvcPackageManager.VapusJwtAuth.ValidateAccessToken(token)
	if err != nil {
		pkgs.DmLogger.Err(err).Msg("error while validating access token from request header")
		return nil, apperr.ErrUnAuthenticated
	}
	nctx := metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", "Bearer "+token))
	// defer nctx.Done()
	token, err = rpcauth.AuthFromMD(nctx, "bearer")
	if err != nil {
		token = ""
	}
	user, err := pkgs.VapusSvcInternalClientManager.GetUser(nctx, parsedClaims[encryption.ClaimUserIdKey], pkgs.DmLogger, 3)
	if err != nil {
		pkgs.DmLogger.Err(err).Msg("error while fetching user details from request header")
		return nil, apperr.ErrInvalidToken
	}
	if !user.ValidateJwtClaim(parsedClaims) {
		pkgs.DmLogger.Err(err).Msg("error while validating jwt claims")
		return nil, apperr.ErrInvalidToken
	}
	pkgs.DmLogger.Info().Msgf("parsed domain claims - %v", parsedClaims)
	parsedClaims[encryption.ClaimUserNameKey] = user.GetDisplayName()
	// uCtx := context.Background()
	return encryption.SetCtxClaim(nctx, parsedClaims), nil
}
