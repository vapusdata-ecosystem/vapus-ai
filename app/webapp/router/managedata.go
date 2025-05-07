package router

import (
	"github.com/labstack/echo/v4"
	routes "github.com/vapusdata-ecosystem/vapusai/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusai/webapp/services"
)

func domainRouter(e *echo.Group) {
	domainRoute := e.Group(routes.MyOrganizationGroup)
	domainRoute.GET(routes.DataSources, services.WebappServiceManager.OrganizationDataSourcesList)
	domainRoute.GET(routes.CreateDataSource, services.WebappServiceManager.CreateOrganizationDataSource)
	domainRoute.GET(routes.DataSourcesResource, services.WebappServiceManager.OrganizationDataSourcesDetail)
}
