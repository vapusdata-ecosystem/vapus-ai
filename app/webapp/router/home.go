package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusdata/webapp/services"
)

func homeRouter(e *echo.Group) {
	e.GET("/", services.WebappServiceManager.HomePageHandler)
}

func CommonRouter(e *echo.Group) {
	e.GET(routes.OrganizationAuth, services.WebappServiceManager.AuthOrganizationHandler)
}
