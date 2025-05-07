package middlewares

import (
	"github.com/labstack/echo/v4"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/webapp/pkgs"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log the request
		pkgs.DmLogger.Info().Fields(map[string]interface{}{
			"method": c.Request().Method,
			"uri":    c.Request().URL.Path,
			"query":  c.Request().URL.RawQuery,
		}).Msg("Request")

		// call the next middleware/handler
		err := next(c)
		if err != nil {
			pkgs.DmLogger.Error().Fields(map[string]interface{}{
				"error": err.Error(),
			}).Msg("Response")

			return err
		}

		return nil
	}
}
