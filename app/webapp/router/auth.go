package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/pkgs"
	routes "github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusdata/webapp/services"
)

func authnRouter(e *echo.Echo) {
	services.InitAuthnService(pkgs.WebAppConfigManager.Path)
	e.GET(routes.Login, services.AuthnServiceManager.LoginHandler)
	e.GET(routes.LoginCallBack, services.AuthnServiceManager.CallbackHandler)
	e.GET(routes.LoginRedirect, services.AuthnServiceManager.LoginRedirect)
	e.GET(routes.Logout, services.AuthnServiceManager.Logout)
}
