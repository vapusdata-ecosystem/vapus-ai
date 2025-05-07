package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/models"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
)

func (x *WebappService) SettingsResources(c echo.Context) error {
	response := models.SettingsResponse{
		SpecMap:            appconfigs.SpecMap,
		ResourceActionsMap: appconfigs.ResourceActionsMap,
	}
	globalContext, err := x.getDeveloperSectionGlobals(c, routes.DevelopersResourcesPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section types")
		return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-resources.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Resources & Specs",
	})
}

func (x *WebappService) SettingsEnums(c echo.Context) error {
	// enums := maps.Keys(appconfigs.EnumSpecs)
	response := models.SettingsResponse{
		Enums: x.SvcPkgManager.ValidEnums,
	}
	globalContext, err := x.getDeveloperSectionGlobals(c, routes.DevelopersEnumsPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section types")
		return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-enums.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Enums",
	})
}
