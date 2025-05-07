package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusai/webapp/services"
)

func insightsRouters(e *echo.Group) {
	grp := e.Group(routes.InsightsGroup)
	grp.GET(routes.LLMInsights, services.WebappServiceManager.InsightsHandler)
}
