package services

import (
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/models"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/webapp/pkgs"
	routes "github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
)

func (x *WebappService) OrganizationDataSourcesList(c echo.Context) error {
	response := models.OrganizationSvcResponse{
		DataSources:         x.grpcClients.ListDataSources(c),
		BackListingLink:     c.Request().URL.String(),
		ActionRules:         []*models.ActionRule{},
		HideBackListingLink: true,
	}
	globalContext, err := x.getDatamanagerSectionGlobals(c, routes.OrganizationDataSourcesPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting domain section types")
		return HandleGLobalContextError(c, err)
	}
	globalContext.Manager = true
	response.CurrentOrganization = globalContext.CurrentOrganization
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/datasources", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.DatasourceManagerRequest),
		Title:    "Enter your Data Source spec in the yaml editor below",
	})
	response.CreateActionParams = &models.ActionRule{
		Weblink: fmt.Sprintf("%s/ui/data/manage/data-sources/create", pkgs.NetworkConfigManager.ExternalURL),
	}
	return c.Render(http.StatusOK, "datasources.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Data Sources",
	})
}

func (x *WebappService) OrganizationDataSourcesDetail(c echo.Context) error {
	var err error
	globalContext, err := x.getDatamanagerSectionGlobals(c, routes.OrganizationDataSourcesPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting home section types")
		return HandleGLobalContextError(c, err)
	}
	dataSourceId, err := GetUrlParams(c, types.DataSourceId)
	if err != nil {
		Logger.Err(err).Msg("error while getting data source id")
		return x.HandleError(c, err, http.StatusBadRequest, globalContext)
	}
	response := models.OrganizationSvcResponse{
		DataSource:      x.grpcClients.DescribeDataSources(c, dataSourceId),
		ActionRules:     []*models.ActionRule{},
		BackListingLink: BuildBacklistingUrl(c.Request().URL.String()),
	}

	if response.DataSource != nil {
		response.ResourceId = response.DataSource.DataSourceId
		if response.DataSource.Status != mpb.CommonStatus_DELETED.String() {
			obj := &dpb.DataSourceManagerRequest{
				Spec: response.DataSource,
			}
			bytess, err := pbtools.ProtoYamlMarshal(obj)
			if err != nil {
				Logger.Err(err).Msg("error while marshaling data source spec")
				return x.HandleError(c, err, http.StatusBadRequest, globalContext)
			}
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_UPDATE.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/datasources", pkgs.NetworkConfigManager.GatewayURL),
				Method:     http.MethodPut,
				YamlSpec:   string(bytess),
				ResourceId: response.DataSource.DataSourceId,
				Title:      "Update " + response.DataSource.Name,
			})

			bytess, err = pbtools.ProtoYamlMarshal(response.DataSource)
			if err != nil {
				Logger.Err(err).Msg("error while marshaling data source spec")
				return x.HandleError(c, err, http.StatusBadRequest, globalContext)
			}
			response.YamlSpec = string(bytess)
			bytess, err = pbtools.ProtoYamlMarshal(&dpb.DataSourceGetterRequest{
				DataSourceId: response.DataSource.DataSourceId,
			})
			if err != nil {
				Logger.Err(err).Msg("error while marshaling data source spec for sync")
				return x.HandleError(c, err, http.StatusBadRequest, globalContext)
			}
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_ARCHIVE.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/datasources/%s", pkgs.NetworkConfigManager.GatewayURL, response.DataSource.DataSourceId),
				Method:     http.MethodDelete,
				ResourceId: response.DataSource.DataSourceId,
				Title:      "Archive " + response.DataSource.Name,
			})
		}

	} else {
		Logger.Err(err).Msg("error while getting data source")
		return x.HandleError(c, err, http.StatusBadRequest, globalContext)
	}
	globalContext.Manager = true
	response.CurrentOrganization = globalContext.CurrentOrganization
	return c.Render(http.StatusOK, "datasource-details.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Data Source Details",
		"ResourceBase":  response.DataSource.ResourceBase,
		"ResourceName":  response.DataSource.Name,
	})
}

func (x *WebappService) CreateOrganizationDataSource(c echo.Context) error {
	response := models.OrganizationSvcResponse{
		DataSources:     x.grpcClients.ListDataSources(c),
		BackListingLink: BuildBacklistingUrl(c.Request().URL.String()),
		ActionRules:     []*models.ActionRule{},
	}
	globalContext, err := x.getDatamanagerSectionGlobals(c, routes.OrganizationDataSourcesPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting domain section types")
		return HandleGLobalContextError(c, err)
	}
	globalContext.Manager = true
	response.CurrentOrganization = globalContext.CurrentOrganization
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/datasources", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.DatasourceManagerRequest),
		Title:    "Enter your Data Source spec in the yaml editor below",
	}
	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "DataSourceType", "Frequency", "StorageEngine", "StorageService", "SvcProvider", "ApiTokenType", "DataSourceAccessScope":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = Enums[key]
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "create-datasource.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Data Sources",
		"Enums":         Enums,
	})
}
