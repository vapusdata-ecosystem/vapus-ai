package middlewares

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	pkgs "github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
)

// Middleware to check for cookie existence
func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ignoredRoutes := map[string]bool{
			pkgs.WebAppConfigManager.URIs.Callback: true,
			pkgs.WebAppConfigManager.URIs.Login:    true,
			routes.LoginRedirect:                   true,
		}

		if _, ok := ignoredRoutes[c.Path()]; ok {
			return next(c)
		}

		// Try to get the cookie
		_, err := c.Cookie(types.ACCESS_TOKEN)
		if err != nil {
			if err == http.ErrNoCookie {
				// Cookie not found, return a 401 Unauthorized response
				log.Println(err, "Access token not found")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Access token not found",
				})
			}
			// Some other error occurred
			return err
		}

		// If cookie exists, continue with the next handler
		return next(c)
	}
}
