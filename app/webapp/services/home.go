package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/webapp/models"
	pkgs "github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
)

var Logger = pkgs.DmLogger

func (x *WebappService) HomePageHandler(c echo.Context) error {
	response := models.HomePageResponse{
		Dashboard:           x.grpcClients.GetDashboard(c),
		HideBackListingLink: true,
	}
	globalContext, err := x.getHomeSectionGlobals(c)
	if err != nil {
		Logger.Err(err).Msg("error while getting home section types")
		return c.Render(http.StatusBadRequest, "403.html", map[string]any{
			"GlobalContext": &models.GlobalContexts{
				LoginUrl: routes.Login,
			},
		})
	}
	return c.Render(http.StatusOK, "home.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"Section":       "Dashboard",
		"SectionHeader": "Dashboard",
	})
}
