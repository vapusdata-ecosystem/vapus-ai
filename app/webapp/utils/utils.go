package utils

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"google.golang.org/grpc/metadata"
)

func GetCookie(c echo.Context, name string) (string, error) {
	// Read the cookie for access token
	cookie, err := c.Cookie(name)
	if err != nil {
		return types.EMPTYSTR, err
	}
	return cookie.Value, nil
}

func GetHost(c echo.Context, withScheme bool) string {
	if withScheme {
		// Get the host with scheme from the request
		if c.Request().URL.Scheme == types.EMPTYSTR {
			return "http://" + c.Request().Host
		}
		return c.Request().URL.Scheme + "://" + c.Request().Host
	}
	// Get the host from the request
	return c.Request().Host
}

func SetTemplate(dir, name string) string {
	return "/" + dir + "/" + name + ".html"
}

func GetBearerCtx(ctx context.Context, token string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)
}
