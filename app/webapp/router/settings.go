package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusdata/webapp/services"
)

func settingsRouter(e *echo.Group) {
	settingsGroup := e.Group(routes.SettingsGroup)
	settingsGroup.GET("", services.WebappServiceManager.SettingsProfile)
	settingsGroup.GET(routes.SettingsPlatform, services.WebappServiceManager.SettingsVapusPlatform)
	settingsGroup.GET(routes.SettingsUsers, services.WebappServiceManager.OrganizationUsersList)
	settingsGroup.GET(routes.SettingsUserResource, services.WebappServiceManager.UserDetails)

	settingsGroup.GET(routes.SecretStoreList, services.WebappServiceManager.SecretServiceList)
	settingsGroup.GET(routes.CreateSecretStore, services.WebappServiceManager.CreateSecretService)
	settingsGroup.GET(routes.SecretStoreDetails, services.WebappServiceManager.ManageSecretServiceDetails)

	settingsGroup.GET(routes.SettingsPlugins, services.WebappServiceManager.ManagePluginsHandler)
	settingsGroup.GET(routes.SettingsPluginResource, services.WebappServiceManager.ManagePluginDetailHandler)
	settingsGroup.GET(routes.SettingsOrganizations, services.WebappServiceManager.OrganizationSettings)
	settingsGroup.GET(routes.SettingsPlatformOrganizations, services.WebappServiceManager.PlatformOrganizationsList)
	settingsGroup.GET(routes.SettingsPluginsCreate, services.WebappServiceManager.CreatePlugins)
	settingsGroup.GET(routes.SettingsPluginsUpdate, services.WebappServiceManager.UpdatePlugins)
}
