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

func (x *WebappService) ManageVapusAgentsHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIAgents:            x.grpcClients.ListVapusAgents(c),
		ActionRules:         []*models.ActionRule{},
		BackListingLink:     c.Request().URL.String(),
		HideBackListingLink: true,
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIAgentsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/agents", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AIAgentsManagerRequest),
		Title:    "Enter your AI Agent Spec",
	})
	response.CreateActionParams = &models.ActionRule{
		Weblink: fmt.Sprintf("%s/ui/ai/manage/agents/create", pkgs.NetworkConfigManager.ExternalURL),
	}
	return c.Render(http.StatusOK, "manageai-agents.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "AI Agents",
	})
}

func (x *WebappService) ManageVapusAgentDetailHandler(c echo.Context) error {
	var vBytes []byte
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIAgentsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	agentId, err := GetUrlParams(c, types.AgentId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIAgent:         x.grpcClients.VapusAgentDetail(c, agentId),
		ActionRules:     []*models.ActionRule{},
		BackListingLink: BuildBacklistingUrl(c.Request().URL.String()),
	}
	if response.AIAgent != nil {
		response.ResourceId = response.AIAgent.AgentId
		if response.AIAgent.ResourceBase.CreatedBy == globalContext.UserInfo.UserId && response.AIAgent.ResourceBase.Organization == globalContext.CurrentOrganization.OrganizationId {
			response.CreateActionParams = &models.ActionRule{
				Weblink: fmt.Sprintf("%s/ui/ai/manage/agents/%s/update", pkgs.NetworkConfigManager.ExternalURL, response.AIAgent.AgentId),
			}
			// upObj := &pb.AgentManagerRequest{
			// 	Spec: response.AIAgent,
			// }
			// bytess, err := pbtools.ProtoYamlMarshal(upObj)
			// if err != nil {
			// 	x.logger.Err(err).Msg("error while marshaling ai agent spec for publish")
			// 	return x.HandleError(c, err, http.StatusBadRequest, globalContext)
			// }
			// response.ActionRules = append(response.ActionRules, &models.ActionRule{
			// 	Action:     mpb.ResourceLcActions_UPDATE.String(),
			// 	API:        fmt.Sprintf("%s/api/v1alpha1/agents", pkgs.NetworkConfigManager.GatewayURL),
			// 	Method:     http.MethodPut,
			// 	YamlSpec:   string(bytess),
			// 	Title:      "Update " + response.AIAgent.Name,
			// 	ResourceId: response.AIAgent.AgentId,
			// })
			bytess, err := pbtools.ProtoYamlMarshal(response.AIAgent)
			if err != nil {
				x.logger.Err(err).Msg("error while marshaling ai agent spec for publish")
			}
			response.YamlSpec = string(bytess)
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_ARCHIVE.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/agents/%s", pkgs.NetworkConfigManager.GatewayURL, response.AIAgent.AgentId),
				Method:     http.MethodDelete,
				Title:      "Archive " + response.AIAgent.Name,
				ResourceId: response.AIAgent.AgentId,
			})
			vObj := &pb.AgentSignalRequest{
				VapusAgentId: response.AIAgent.AgentId,
			}
			pObj := &pb.AgentStateRequest{
				VapusAgentId: response.AIAgent.AgentId,
				Action:       mpb.ResourceLcActions_VALIDATE,
			}
			vBytes, err = pbtools.ProtoJsonMarshal(vObj)
			if err != nil {
				x.logger.Err(err).Msg("error while marshaling ai agent spec for publish")
				return x.HandleError(c, err, http.StatusBadRequest, globalContext)
			}

			if response.AIAgent.ResourceBase.Status == mpb.CommonStatus_READY.String() {
				pObj.Action = mpb.ResourceLcActions_UNPUBLISH
				bytess, err = pbtools.ProtoJsonMarshal(pObj)
				if err != nil {
					x.logger.Err(err).Msg("error while marshaling ai agent spec for publish")
					return x.HandleError(c, err, http.StatusBadRequest, globalContext)
				}
				response.ActionRules = append(response.ActionRules, &models.ActionRule{
					Action:     mpb.ResourceLcActions_UNPUBLISH.String(),
					API:        fmt.Sprintf("%s/api/v1alpha1/agents/%s/state", pkgs.NetworkConfigManager.GatewayURL, response.AIAgent.AgentId),
					Method:     http.MethodPost,
					Title:      "UnPublish " + response.AIAgent.Name,
					YamlSpec:   string(bytess),
					ResourceId: response.AIAgent.AgentId,
				})
			} else if response.AIAgent.ResourceBase.Status == mpb.CommonStatus_VALIDATED.String() {
				// } else {
				pObj.Action = mpb.ResourceLcActions_PUBLISH
				bytess, err = pbtools.ProtoJsonMarshal(pObj)
				if err != nil {
					x.logger.Err(err).Msg("error while marshaling ai agent spec for publish")
					return x.HandleError(c, err, http.StatusBadRequest, globalContext)
				}
				response.ActionRules = append(response.ActionRules, &models.ActionRule{
					Action:     mpb.ResourceLcActions_PUBLISH.String(),
					API:        fmt.Sprintf("%s/api/v1alpha1/agents/%s/state", pkgs.NetworkConfigManager.GatewayURL, response.AIAgent.AgentId),
					Method:     http.MethodPost,
					Title:      "Publish " + response.AIAgent.Name,
					YamlSpec:   string(bytess),
					ResourceId: response.AIAgent.AgentId,
				})
			}
		}
	} else {
		x.logger.Err(err).Msg("error while getting ai agents details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	return c.Render(http.StatusOK, "manageai-agents-detail.html", map[string]any{
		"GlobalContext":  globalContext,
		"Response":       response,
		"SectionHeader":  "AI Agent Details",
		"AgentStudio":    routes.AgentStudioHome,
		"ResourceBase":   response.AIAgent.ResourceBase,
		"ResourceName":   response.AIAgent.Name,
		"DownloadUrl":    fmt.Sprintf("%s/api/v1alpha1/fabric/download", pkgs.NetworkConfigManager.GatewayURL),
		"FileUploadAPI":  fmt.Sprintf("%s/api/v1alpha1/utility/upload", pkgs.NetworkConfigManager.GatewayURL),
		"ValidationSpec": string(vBytes),
		"ValidationAPI":  fmt.Sprintf("%s/api/v1alpha1/agents/%s/validate", pkgs.NetworkConfigManager.GatewayURL, response.AIAgent.AgentId),
	})
}

func (x *WebappService) CreateVapusAgent(c echo.Context) error {
	response := &models.AIStudioResponse{
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.GatewayURL, routes.ManageAIGroup, routes.ManageAIAgents),
		ActionRules:     []*models.ActionRule{},
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIAgentsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}

	globalContext.Manager = true
	response.AIModelNodes = x.grpcClients.AIModelNodes(c)
	response.DataSources = x.grpcClients.ListDataSources(c)
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/agents", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AIAgentsManagerRequest),
		Title:    "Enter your Agent Spec",
	}
	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "Frequency", "ApiTokenType":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "create-agents.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Build Your Agent",
		"Enums":         Enums,
		"DetailUrl":     routes.AgentStudioHome,
	})
}

func (x *WebappService) UpdateVapusAgent(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIAgentsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	agentId, err := GetUrlParams(c, types.AgentId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}

	response := &models.AIStudioResponse{
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.GatewayURL, routes.ManageAIGroup, routes.ManageAIAgents),
		AIAgent:         x.grpcClients.VapusAgentDetail(c, agentId),
		ActionRules:     []*models.ActionRule{},
		// CreateActionParams: &models.ActionRule{},
	}

	globalContext.Manager = true
	response.AIModelNodes = x.grpcClients.AIModelNodes(c)
	response.DataSources = x.grpcClients.ListDataSources(c)

	globalContext.Manager = true
	upObj := &pb.AgentManagerRequest{
		Spec: response.AIAgent,
	}
	response.CreateActionParams = &models.ActionRule{
		Action:     mpb.ResourceLcActions_UPDATE.String(),
		API:        fmt.Sprintf("%s/api/v1alpha1/agents/%s", pkgs.NetworkConfigManager.GatewayURL, response.AIAgent.AgentId),
		Method:     http.MethodPut,
		YamlSpec:   GetProtoYamlString(upObj),
		Title:      "Update your Agent Spec",
		ResourceId: response.AIAgent.AgentId,
	}
	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "Frequency", "ApiTokenType":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "update-agents.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Update Your Agent",
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIAgent.ResourceBase,
		"Enums":         Enums,
	})
}
