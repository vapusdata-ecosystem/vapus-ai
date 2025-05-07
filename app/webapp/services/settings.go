package services

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"github.com/vapusdata-ecosystem/vapusai/webapp/models"
	pkgs "github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"

	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
	"github.com/vapusdata-ecosystem/vapusai/webapp/utils"
)

func (x *WebappService) SettingsProfile(c echo.Context) error {
	response := models.SettingsResponse{
		ActionRules:         []*models.ActionRule{},
		HideBackListingLink: true,
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsProfilePage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section types")
		return HandleGLobalContextError(c, err)
	}
	if globalContext.UserInfo != nil {
		response.ResourceId = globalContext.UserInfo.UserId
		obj := &dpb.UserManagerRequest{
			Spec:   globalContext.UserInfo,
			Action: dpb.UserManagerActions_PATCH_USER,
		}
		bytess, err := pbtools.ProtoYamlMarshal(obj)
		if err != nil {
			Logger.Err(err).Msg("error while marshaling user spec")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		response.ActionRules = append(response.ActionRules, &models.ActionRule{
			Action:     mpb.ResourceLcActions_UPDATE.String(),
			API:        fmt.Sprintf("%s/api/v1alpha1/users", pkgs.NetworkConfigManager.GatewayURL),
			Method:     http.MethodPost,
			YamlSpec:   string(bytess),
			ResourceId: globalContext.UserInfo.UserId,
			Title:      "Update " + globalContext.UserInfo.DisplayName,
		})
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-profile.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Your Profile",
	})
}

func (x *WebappService) SettingsVapusPlatform(c echo.Context) error {
	response := models.SettingsResponse{
		ActionParams:        &models.ResourceManagerParams{},
		ActionRules:         []*models.ActionRule{},
		HideBackListingLink: true,
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPlatformPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section types")
		return HandleGLobalContextError(c, err)
	}
	if globalContext.Account != nil {
		response.ResourceId = globalContext.Account.AccountId
		obj := &dpb.AccountManagerRequest{
			Spec:    globalContext.Account,
			Actions: dpb.AccountAgentActions_UPDATE_PROFILE,
		}
		bytess, err := pbtools.ProtoYamlMarshal(obj)
		if err != nil {
			Logger.Err(err).Msg("error while marshaling platform account spec")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		response.ActionRules = append(response.ActionRules, &models.ActionRule{
			Action:     mpb.ResourceLcActions_UPDATE.String(),
			API:        fmt.Sprintf("%s/api/v1alpha1/platform", pkgs.NetworkConfigManager.GatewayURL),
			Method:     http.MethodPost,
			YamlSpec:   string(bytess),
			ResourceId: globalContext.UserInfo.UserId,
			Title:      "Update " + globalContext.Account.Name,
		})
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-platform.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Platform Settings",
	})
}

func (x *WebappService) OrganizationSettings(c echo.Context) error {
	var err error
	response := models.OrganizationSvcResponse{
		ActionParams:        &models.ResourceManagerParams{},
		ActionRules:         []*models.ActionRule{},
		HideBackListingLink: true,
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.OrganizationSettingsPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section types")
		return HandleGLobalContextError(c, err)
		// return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	if response.CurrentOrganization != nil {
		response.ResourceId = response.CurrentOrganization.OrganizationId
		if globalContext.IsOrganizationOwner {
			obj := &dpb.OrganizationManagerRequest{
				Spec: response.CurrentOrganization,
			}
			bytess, err := pbtools.ProtoYamlMarshal(obj)
			if err != nil {
				Logger.Err(err).Msg("error while marshaling domain spec")
				return x.HandleError(c, err, http.StatusNotFound, globalContext)
			}

			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_UPDATE.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/domains", pkgs.NetworkConfigManager.GatewayURL),
				Method:     http.MethodPut,
				YamlSpec:   string(bytess),
				ResourceId: response.CurrentOrganization.OrganizationId,
				Title:      "Update " + response.CurrentOrganization.Name,
			})

			addUserObj := &dpb.OrganizationAdduserRequest{
				OrganizationId: response.CurrentOrganization.OrganizationId,
				Users: []*mpb.OrganizationUserOps{{
					UserId:           "",
					Role:             []string{},
					ValidTill:        0,
					InviteIfNotFound: true,
				}},
			}
			uBytess, err := pbtools.ProtoYamlMarshal(addUserObj)
			if err != nil {
				Logger.Err(err).Msg("error while marshaling domain add user spec")
				return x.HandleError(c, err, http.StatusNotFound, globalContext)
			}
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     utils.ADD_USERS,
				API:        fmt.Sprintf("%s/api/v1alpha1/domains/%s/users", pkgs.NetworkConfigManager.GatewayURL, response.CurrentOrganization.OrganizationId),
				Method:     http.MethodPut,
				YamlSpec:   string(uBytess),
				ResourceId: response.CurrentOrganization.OrganizationId,
				Title:      "Add Users in " + response.CurrentOrganization.Name,
			})

			uBytess, err = pbtools.ProtoYamlMarshal(&dpb.OrganizationGetterRequest{
				OrganizationId: response.CurrentOrganization.OrganizationId,
			})
			if err != nil {
				Logger.Err(err).Msg("error while marshaling domain upgrade OS spec")
				return x.HandleError(c, err, http.StatusNotFound, globalContext)
			}
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     utils.UPGRADE,
				API:        fmt.Sprintf("%s/api/v1alpha1/domains/%s/upgrade-os", pkgs.NetworkConfigManager.GatewayURL, response.CurrentOrganization.OrganizationId),
				Method:     http.MethodPost,
				YamlSpec:   string(uBytess),
				ResourceId: response.CurrentOrganization.OrganizationId,
				Title:      "Upgrade OS for " + response.CurrentOrganization.Name,
			})
		}

	} else {
		Logger.Err(err).Msg("error while getting domain")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	return c.Render(http.StatusOK, "settings-domain.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Organization Settings",
		"ResourceName":  response.CurrentOrganization.Name,
		"ResourceBase":  response.CurrentOrganization.ResourceBase,
	})
}

func (x *WebappService) PlatformUsersList(c echo.Context) error {
	response := models.SettingsResponse{
		Users:           x.grpcClients.GetPlatformUsers(c),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: c.Request().URL.String(),
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsUsersPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section types")
		return HandleGLobalContextError(c, err)
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-users.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Platform Users",
	})
}

func (x *WebappService) OrganizationUsersList(c echo.Context) error {
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsUsersPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section types")
		return HandleGLobalContextError(c, err)
	}
	response := models.OrganizationSvcResponse{
		BackListingLink:     c.Request().URL.String(),
		ActionParams:        &models.ResourceManagerParams{},
		HideBackListingLink: true,
	}
	if globalContext.CurrentOrganization.OrganizationType == mpb.OrganizationType_SERVICE_ORGANIZATION {
		response.Users = x.grpcClients.GetPlatformUsers(c)
	} else {
		response.Users = x.grpcClients.GetMyOrganizationUsers(c)
	}

	response.CurrentOrganization = globalContext.CurrentOrganization
	response.CurrentOrganization = globalContext.CurrentOrganization
	if response.CurrentOrganization != nil {
		if globalContext.IsOrganizationOwner {
			addUserObj := &dpb.OrganizationAdduserRequest{
				OrganizationId: response.CurrentOrganization.OrganizationId,
				Users: []*mpb.OrganizationUserOps{{
					UserId:           "",
					Role:             []string{},
					ValidTill:        0,
					InviteIfNotFound: true,
				}},
			}
			uBytess, err := pbtools.ProtoYamlMarshal(addUserObj)
			if err != nil {
				Logger.Err(err).Msg("error while marshaling domain add user spec")
				return x.HandleError(c, err, http.StatusNotFound, globalContext)
			}
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     utils.ADD_USERS,
				API:        fmt.Sprintf("%s/api/v1alpha1/domains/%s/users", pkgs.NetworkConfigManager.GatewayURL, response.CurrentOrganization.OrganizationId),
				Method:     http.MethodPut,
				YamlSpec:   string(uBytess),
				ResourceId: response.CurrentOrganization.OrganizationId,
				Title:      "Add Users in " + response.CurrentOrganization.Name,
			})
		}
	} else {
		Logger.Err(err).Msg("error while getting domain")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	return c.Render(http.StatusOK, "settings-domain-users.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Organization Users",
	})
}

func (x *WebappService) UserDetails(c echo.Context) error {
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsUsersPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section types")
		return HandleGLobalContextError(c, err)
	}
	userId, err := GetUrlParams(c, types.UserId)
	user, _, err := x.grpcClients.GetUserInfo(c, userId)
	if err != nil || user == nil {
		Logger.Err(err).Msg("error while getting user info")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	response := models.SettingsResponse{
		User:            user,
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: BuildBacklistingUrl(c.Request().URL.String()),
	}

	cRoles := GetUserCurrentOrganizationRole(globalContext.UserInfo, globalContext.CurrentOrganization.OrganizationId)
	if slices.Contains(cRoles, mpb.UserRoles_ORG_OWNER.String()) {
		var updateSpec = ""
		bbtes, err := pbtools.ProtoYamlMarshal(&dpb.UserManagerRequest{
			Spec:   user,
			Action: dpb.UserManagerActions_PATCH_USER,
		})
		if err != nil {
			Logger.Err(err).Msg("error while marshaling user spec")
			return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		}
		updateSpec = string(bbtes)
		response.ActionParams.API = fmt.Sprintf("%s/api/v1alpha1/users", pkgs.NetworkConfigManager.GatewayURL)
		response.ActionParams.ActionMap = map[string]string{
			utils.UPDATE: updateSpec,
		}
	}
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "settings-user-details.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "User Details",
		"ResourceName":  "User - " + user.UserId,
	})
}

func (x *WebappService) ManagePluginsHandler(c echo.Context) error {
	response := &models.SettingsResponse{
		Plugins:             x.grpcClients.ListPlugins(c),
		ActionParams:        &models.ResourceManagerParams{},
		BackListingLink:     c.Request().URL.String(),
		HideBackListingLink: true,
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPluginsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting plugin section types")
		return HandleGLobalContextError(c, err)
	}

	// Specs := make(map[string]string)
	// for _, plugin := range response.Plugins {
	// 	plugin.NetworkParams.Credentials = &mpb.GenericCredentialObj{
	// 		AwsCreds:   &mpb.AWSCreds{},
	// 		GcpCreds:   &mpb.GCPCreds{},
	// 		AzureCreds: &mpb.AzureCreds{},
	// 	}
	// 	obj := &dpb.PluginManagerRequest{
	// 		Spec:   plugin,
	// 		Action: dpb.PluginAgentAction_PATCH_PLUGIN,
	// 	}
	// 	bytess, err := pbtools.ProtoYamlMarshal(obj)
	// 	if err != nil {
	// 		x.logger.Err(err).Msg("error while marshaling plugin spec")
	// 		Specs[plugin.PluginId] = ""
	// 	} else {
	// 		Specs[plugin.PluginId] = string(bytess)
	// 	}
	// }
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/plugins", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.PluginManagerRequest),
		Title:    "Enter your Plugin Spec",
	})
	response.CreateActionParams = &models.ActionRule{
		Weblink: fmt.Sprintf("%s/ui/settings/plugins/create", pkgs.NetworkConfigManager.ExternalURL),
	}
	return c.Render(http.StatusOK, "settings-plugins.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Plugins",
	})
}

func (x *WebappService) ManagePluginDetailHandler(c echo.Context) error {
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPluginsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting plugin section types")
		return HandleGLobalContextError(c, err)
	}
	PluginId, err := GetUrlParams(c, types.PluginId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.SettingsResponse{
		Plugin:             x.grpcClients.GetPlugin(c, PluginId),
		ActionRules:        []*models.ActionRule{},
		BackListingLink:    BuildBacklistingUrl(c.Request().URL.String()),
		CreateActionParams: &models.ActionRule{},
	}
	if response.Plugin != nil {
		response.ResourceId = response.Plugin.PluginId
		// obj := &dpb.PluginManagerRequest{
		// 	Spec: response.Plugin,
		// }
		response.CreateActionParams = &models.ActionRule{
			Weblink: fmt.Sprintf("%s/ui/settings/plugins/%s/update", pkgs.NetworkConfigManager.ExternalURL, PluginId),
		}
		// bytess, err := pbtools.ProtoYamlMarshal(obj)
		// if err != nil {
		// 	Logger.Err(err).Msg("error while marshaling data source spec")
		// 	return x.HandleError(c, err, http.StatusBadRequest, globalContext)
		// }
		// response.ActionRules = append(response.ActionRules, &models.ActionRule{
		// 	Action:     mpb.ResourceLcActions_UPDATE.String(),
		// 	API:        fmt.Sprintf("%s/api/v1alpha1/plugins", pkgs.NetworkConfigManager.GatewayURL),
		// 	Method:     http.MethodPut,
		// 	YamlSpec:   string(bytess),
		// 	ResourceId: response.Plugin.PluginId,
		// 	Title:      "Update " + response.Plugin.Name,
		// })
		response.ActionRules = append(response.ActionRules, &models.ActionRule{
			Action:     mpb.ResourceLcActions_ARCHIVE.String(),
			API:        fmt.Sprintf("%s/api/v1alpha1/plugins/%s", pkgs.NetworkConfigManager.GatewayURL, response.Plugin.PluginId),
			Method:     http.MethodDelete,
			ResourceId: response.Plugin.PluginId,
			Title:      "Archive " + response.Plugin.Name,
		})
	} else {
		x.logger.Err(err).Msg("error while getting plugin details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	return c.Render(http.StatusOK, "settings-plugin-detail.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Plugin Details",
		"ResourceName":  response.Plugin.Name,
		"ResourceBase":  response.Plugin.ResourceBase,
	})
}

// Create Plugins
func (x *WebappService) CreatePlugins(c echo.Context) error {
	response := &models.SettingsResponse{
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.ExternalURL, routes.SettingsGroup, routes.SettingsPlugins),
		ActionRules:     []*models.ActionRule{},
		PluginTypeMap:   x.grpcClients.ResourceGetter(c),
	}

	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPluginsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	globalContext.Manager = true
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/plugins", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.PluginManagerRequest),
		Title:    "Create Plugins",
	}
	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "IntegrationPluginTypes", "ResourceScope", "ApiTokenType":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "settings-create-plugins.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Build Your Plugins",
		"Enums":         Enums,
	})
}

// Plugin Update
func (x *WebappService) UpdatePlugins(c echo.Context) error {
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPluginsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting plugin section types")
		return HandleGLobalContextError(c, err)
	}
	PluginId, err := GetUrlParams(c, types.PluginId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.SettingsResponse{
		Plugin:          x.grpcClients.GetPlugin(c, PluginId),
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.ExternalURL, routes.SettingsGroup, routes.SettingsPluginResource),
		ActionRules:     []*models.ActionRule{},
		PluginTypeMap:   x.grpcClients.ResourceGetter(c),
	}

	globalContext.Manager = true
	upObj := &dpb.PluginManagerRequest{
		Spec:   response.Plugin,
		Action: dpb.PluginAgentAction_PATCH_PLUGIN,
	}
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_UPDATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/plugins", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPut,
		YamlSpec: GetProtoYamlString(upObj),
		Title:    "Update Plugins",
	}
	Enums := map[string][]string{}
	KeyValueMap := make(map[string]map[string]int32)
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "IntegrationPluginTypes", "ResourceScope", "ApiTokenType":
			KeyValueMap[key] = value
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "settings-update-plugins.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Build Your Plugins",
		"Enums":         Enums,
		"KeyValueMap":   KeyValueMap,
	})
}

func (x *WebappService) PlatformOrganizationsList(c echo.Context) error {
	response := models.ExploreResponse{
		Organizations:       x.grpcClients.GetOrganizations(c),
		ActionParams:        &models.ResourceManagerParams{},
		HideBackListingLink: true,
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsPlatformOrganizationsPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting explore page global context")
		return HandleGLobalContextError(c, err)
	}
	for _, domain := range globalContext.UserInfo.Roles {
		response.YourOrganizations = append(response.YourOrganizations, domain.OrganizationId)
	}
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/domains", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.OrganizationManagerRequest),
		Title:    "",
	})
	return c.Render(http.StatusOK, "settings-platform-domains.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Organizations",
	})
}

func (x *WebappService) SecretServiceList(c echo.Context) error {
	response := models.SecretServiceResponse{
		SecretStores:        x.grpcClients.GetSecretServiceList(c),
		ActionParams:        &models.ResourceManagerParams{},
		BackListingLink:     c.Request().URL.String(),
		HideBackListingLink: true,
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsSecretStorePage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section globals")
		return HandleGLobalContextError(c, err)
	}
	response.CreateActionParams = &models.ActionRule{
		Weblink: fmt.Sprintf("%s/ui/settings/secretstores/create", pkgs.NetworkConfigManager.ExternalURL),
	}
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action: mpb.ResourceLcActions_CREATE.String(),
		API:    fmt.Sprintf("%s/api/v1alpha1/secrets", pkgs.NetworkConfigManager.GatewayURL),
		Method: http.MethodPost,
		Title:  "",
	})

	return c.Render(http.StatusOK, "settings-secret-store.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Secret Service",
	})
}

func (x WebappService) CreateSecretService(c echo.Context) error {
	response := models.SecretServiceResponse{
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.ExternalURL, routes.SettingsGroup, routes.SecretStoreList),
		ActionParams:    &models.ResourceManagerParams{},
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsSecretStorePage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section globals")
		return HandleGLobalContextError(c, err)
	}
	globalContext.Manager = true
	response.CreateActionParams = &models.ActionRule{
		Action: mpb.ResourceLcActions_CREATE.String(),
		API:    fmt.Sprintf("%s/api/v1alpha1/secrets", pkgs.NetworkConfigManager.GatewayURL),
		Method: http.MethodPost,
		Title:  "",
	}

	Enums := map[string][]string{}
	KeyValueMap := make(map[string]map[string]int32)
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "VapusSecretType", "ApiTokenType", "DataSourceAccessScope":
			KeyValueMap[key] = value
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "settings-create-secret-store.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Create Secret Service",
		"Enums":         Enums,
	})
}

func (x *WebappService) ManageSecretServiceDetails(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.SettingsSecretStorePage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	SecretStoreName, err := GetUrlParams(c, types.SecretStoreName)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.SecretServiceResponse{
		SecretStore:     x.grpcClients.SecretServiceDetails(c, SecretStoreName),
		ActionParams:    &models.ResourceManagerParams{},
		BackListingLink: BuildBacklistingUrl(c.Request().URL.String()),
		ActionRules:     []*models.ActionRule{},
	}
	if response.SecretStore != nil {
		response.CreateActionParams = &models.ActionRule{
			Weblink: fmt.Sprintf("%s/ui/settings/secretstores/%s/update", pkgs.NetworkConfigManager.ExternalURL, response.SecretStore.Name),
		}

		bytess, err := pbtools.ProtoYamlMarshal(response.SecretStore)
		if err != nil {
			x.logger.Err(err).Msg("error while marshaling ai AIPrompt spec for publish")
		}
		response.YamlSpec = string(bytess)
		response.ActionRules = append(response.ActionRules, &models.ActionRule{
			Action: mpb.ResourceLcActions_ARCHIVE.String(),
			API:    fmt.Sprintf("%s/api/v1alpha1/secrets/%s", pkgs.NetworkConfigManager.GatewayURL, response.SecretStore.Name),
			Method: http.MethodDelete,
			Title:  "Archive " + response.SecretStore.Name,
		})
	} else {
		x.logger.Err(err).Msg("error while getting Secret Service Details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	return c.Render(http.StatusOK, "settings-secret-details.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Secret Service Details",
		"SecretService": routes.SecretServiceHome,
		"ResourceBase":  response.SecretStore.ResourceBase,
		"ResourceName":  response.SecretStore.Name,
	})
}

func (x WebappService) UpdateSecretService(c echo.Context) error {
	response := models.SecretServiceResponse{
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.GatewayURL, routes.SettingsGroup, routes.SecretStoreList),
		ActionParams:    &models.ResourceManagerParams{},
	}
	globalContext, err := x.getSettingsSectionGlobals(c, routes.SettingsSecretStorePage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting settings section globals")
		return HandleGLobalContextError(c, err)
	}
	globalContext.Manager = true
	response.CreateActionParams = &models.ActionRule{
		Action: mpb.ResourceLcActions_UPDATE.String(),
		API:    fmt.Sprintf("%s/api/v1alpha1/secrets", pkgs.NetworkConfigManager.GatewayURL),
		Method: http.MethodPut,
		Title:  "",
	}

	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "VapusSecretType":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "settings-update-secret-store.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Update Secret Service",
		"Enums":         Enums,
	})
}
