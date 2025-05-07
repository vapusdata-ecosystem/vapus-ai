package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusai/webapp/services"
)

func homeRouter(e *echo.Group) {
	e.GET("/", services.WebappServiceManager.HomePageHandler)
}

func CommonRouter(e *echo.Group) {
	e.GET(routes.OrganizationAuth, services.WebappServiceManager.AuthOrganizationHandler)
}
