package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusai/webapp/services"
)

func developersRouters(e *echo.Group) {
	developersGroup := e.Group(routes.DevelopersGroup)
	developersGroup.GET(routes.DevelopersResources, services.WebappServiceManager.SettingsResources)
	developersGroup.GET(routes.DevelopersEnums, services.WebappServiceManager.SettingsEnums)
}
