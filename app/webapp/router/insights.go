package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusdata/webapp/services"
)

func insightsRouters(e *echo.Group) {
	grp := e.Group(routes.InsightsGroup)
	grp.GET(routes.LLMInsights, services.WebappServiceManager.InsightsHandler)
}
