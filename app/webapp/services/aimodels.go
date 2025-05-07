package services

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/models"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/webapp/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
)

func (x *WebappService) ManageAIModelNodesHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIModelNodes:        x.grpcClients.AIModelNodes(c),
		ActionRules:         []*models.ActionRule{},
		BackListingLink:     c.Request().URL.String(),
		HideBackListingLink: true,
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIModelNodesPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/models", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AinodeConfiguratorRequest),
		Title:    "Enter your AI Model Node Spec",
	})
	response.CreateActionParams = &models.ActionRule{
		Weblink: fmt.Sprintf("%s/ui/ai/manage/model-nodes/create", pkgs.NetworkConfigManager.ExternalURL),
	}
	return c.Render(http.StatusOK, "manageai-modelnodes.html", map[string]any{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "AI Models",
		"CreateTemplate": GetProtoYamlString(appconfigs.AinodeConfiguratorRequest),
	})
}

func (x *WebappService) ManageAIModelNodesDetailHandler(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIModelNodesPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	nodeId, err := GetUrlParams(c, types.AIModelNodeId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIModelNode:        x.grpcClients.AIModelNodesDetails(c, nodeId),
		ActionRules:        []*models.ActionRule{},
		BackListingLink:    BuildBacklistingUrl(c.Request().URL.String()),
		CreateActionParams: &models.ActionRule{},
	}
	response.CreateActionParams = &models.ActionRule{
		AiStudioURL: fmt.Sprintf("%s/ui/ai-studio?aiModelNode=%s", pkgs.NetworkConfigManager.ExternalURL, response.ResourceId),
	}
	if response.AIModelNode != nil {
		response.ResourceId = response.AIModelNode.ModelNodeId
		if (response.AIModelNode.ResourceBase.CreatedBy == globalContext.UserInfo.UserId || slices.Contains(response.AIModelNode.ResourceBase.Editors, globalContext.UserInfo.UserId)) && response.AIModelNode.ResourceBase.Organization == globalContext.CurrentOrganization.OrganizationId {

			bbytes, err := pbtools.ProtoYamlMarshal(response.AIModelNode)
			if err != nil {
				x.logger.Err(err).Msg("error while marshaling ai model node spec for download")
				return x.HandleError(c, err, http.StatusBadRequest, globalContext)
			}
			response.YamlSpec = string(bbytes)
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_ARCHIVE.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/aistudio/models/%s", pkgs.NetworkConfigManager.GatewayURL, response.AIModelNode.ModelNodeId),
				Method:     http.MethodDelete,
				ResourceId: response.AIModelNode.ModelNodeId,
				Title:      "Archive " + response.AIModelNode.Name,
			})
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_SYNC.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/aistudio/models/%s/sync", pkgs.NetworkConfigManager.GatewayURL, response.AIModelNode.ModelNodeId),
				Method:     http.MethodPost,
				ResourceId: response.AIModelNode.ModelNodeId,
				Title:      "Sync " + response.AIModelNode.Name,
			})
			response.CreateActionParams = &models.ActionRule{
				AiStudioURL: fmt.Sprintf("%s/ui/ai-studio?aiModelNode=%s", pkgs.NetworkConfigManager.ExternalURL, response.ResourceId),
				Weblink:     fmt.Sprintf("%s/ui/ai/manage/model-nodes/%s/update", pkgs.NetworkConfigManager.ExternalURL, nodeId),
			}
		}
	} else {
		x.logger.Err(err).Msg("error while getting ai model node details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	return c.Render(http.StatusOK, "manageai-modelnodes-detail.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "AI Model Node Details",
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIModelNode.ResourceBase,
		"ResourceName":  response.AIModelNode.Name,
	})
}

func (x *WebappService) CreateModels(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIGuardrails:    x.grpcClients.ListAIGuardrails(c),
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.ExternalURL, routes.ManageAIGroup, routes.ManageAIModelNodes),
		ActionRules:     []*models.ActionRule{},
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIModelInterfacePage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	// mps := x.grpcClients.Datamarketplacelist(c)
	// if len(mps) == 0 {
	// 	return c.Render(http.StatusNoContent, "404.html", map[string]any{
	// 		"GlobalContext": globalContext,
	// 	})
	// }
	globalContext.Manager = true
	// response.AIModelNodes = x.grpcClients.AIModelNodes(c)
	// response.DataProducts = x.grpcClients.ListAccessibleDataproducts(c, mps[0].MarketplaceId)
	// response.DataSources = x.grpcClients.ListDataSources(c)
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/models", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AinodeConfiguratorRequest),
		Title:    "Enter your Agent Spec",
	}
	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "Frequency", "ApiTokenType", "SvcProvider", "AIModelNodeHosting", "ResourceScope":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})

		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "create-model.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Build Your Models",
		"Enums":         Enums,
	})
}

func (x *WebappService) UpdateModels(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIModelInterfacePage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	nodeId, err := GetUrlParams(c, types.AIModelNodeId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIGuardrails:       x.grpcClients.ListAIGuardrails(c),
		AIModelNode:        x.grpcClients.AIModelNodesDetails(c, nodeId),
		ActionRules:        []*models.ActionRule{},
		BackListingLink:    fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.ExternalURL, routes.ManageAIGroup, routes.ManageAIModelNodes),
		CreateActionParams: &models.ActionRule{},
	}

	globalContext.Manager = true
	upObj := &pb.AIModelNodeManagerRequest{
		Spec: response.AIModelNode,
	}
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_UPDATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/models", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPut,
		YamlSpec: GetProtoYamlString(upObj),
		Title:    "Update your model Spec",
	}

	Enums := map[string][]string{}
	KeyValueMap := make(map[string]map[string]int32)
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "Frequency", "ApiTokenType", "SvcProvider", "AIModelNodeHosting", "ResourceScope":
			KeyValueMap[key] = value
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}

	return c.Render(http.StatusOK, "update-model.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Update Model",
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIModelNode.ResourceBase,
		"Enums":         Enums,
		"KeyValueMap":   KeyValueMap,
	})
}
