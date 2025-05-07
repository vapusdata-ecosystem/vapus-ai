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
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"github.com/vapusdata-ecosystem/vapusai/webapp/models"
	pkgs "github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
)

func (x *WebappService) ManageAIPromptsHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIPrompts:           x.grpcClients.AIModelPrompts(c),
		BackListingLink:     c.Request().URL.String(),
		HideBackListingLink: true,
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.AIPromptsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/prompts", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AiPromptManagerRequest),
		Title:    "Enter your AI Prompt Spec",
	})
	response.CreateActionParams = &models.ActionRule{
		Weblink: fmt.Sprintf("%s/ui/ai/manage/prompts/create", pkgs.NetworkConfigManager.ExternalURL),
	}
	return c.Render(http.StatusOK, "manageai-prompts.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "AI Model Prompts",
	})
}

func (x *WebappService) ManageAIPromptDetailHandler(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.AIPromptsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	promptId, err := GetUrlParams(c, types.PromptId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIPrompt:           x.grpcClients.AIModelPromptDetails(c, promptId),
		ActionRules:        []*models.ActionRule{},
		BackListingLink:    BuildBacklistingUrl(c.Request().URL.String()),
		CreateActionParams: &models.ActionRule{},
	}
	response.CreateActionParams = &models.ActionRule{
		AiStudioURL: fmt.Sprintf("%s/ui/ai-studio?promptId=%s", pkgs.NetworkConfigManager.ExternalURL, response.ResourceId),
	}
	if response.AIPrompt != nil {
		response.ResourceId = response.AIPrompt.PromptId
		// Checking is the user has editing rights or not
		if (response.AIPrompt.ResourceBase.CreatedBy == globalContext.UserInfo.UserId || slices.Contains(response.AIPrompt.ResourceBase.Editors, globalContext.UserInfo.UserId)) && response.AIPrompt.ResourceBase.Organization == globalContext.CurrentOrganization.OrganizationId {
			response.CreateActionParams = &models.ActionRule{
				AiStudioURL: fmt.Sprintf("%s/ui/ai-studio?promptId=%s", pkgs.NetworkConfigManager.ExternalURL, response.ResourceId),
				Weblink:     fmt.Sprintf("%s/ui/ai/manage/prompts/%s/update", pkgs.NetworkConfigManager.ExternalURL, response.AIPrompt.PromptId),
			}

			bytess, err := pbtools.ProtoYamlMarshal(response.AIPrompt)
			if err != nil {
				x.logger.Err(err).Msg("error while marshaling ai AIPrompt spec for publish")
			}
			response.YamlSpec = string(bytess)
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_ARCHIVE.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/aistudio/prompts/%s", pkgs.NetworkConfigManager.GatewayURL, response.AIPrompt.PromptId),
				Method:     http.MethodDelete,
				Title:      "Archive " + response.AIPrompt.Name,
				ResourceId: response.AIPrompt.PromptId,
			})
		}
	} else {
		x.logger.Err(err).Msg("error while getting ai prompt details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	return c.Render(http.StatusOK, "manageai-prompts-detail.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "AI Model Prompt Details",
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIPrompt.ResourceBase,
		"ResourceName":  response.AIPrompt.Name,
	})
}

func (x *WebappService) UpdatePrompts(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.AIPromptsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	promptId, err := GetUrlParams(c, types.PromptId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		BackListingLink:    fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.GatewayURL, routes.ManageAIGroup, routes.ManageAIPrompts),
		AIPrompt:           x.grpcClients.AIModelPromptDetails(c, promptId),
		ActionRules:        []*models.ActionRule{},
		CreateActionParams: &models.ActionRule{},
	}
	globalContext.Manager = true
	upObj := &pb.PromptManagerRequest{
		Spec: response.AIPrompt,
	}
	response.CreateActionParams = &models.ActionRule{
		Action:     mpb.ResourceLcActions_UPDATE.String(),
		API:        fmt.Sprintf("%s/api/v1alpha1/aistudio/prompts", pkgs.NetworkConfigManager.GatewayURL),
		Method:     http.MethodPut,
		YamlSpec:   GetProtoYamlString(upObj),
		Title:      "Enter your Prompt Spec",
		ResourceId: response.AIPrompt.PromptId,
	}
	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "AIResponseFormat":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "update-prompts.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Update Your Prompt",
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIPrompt.ResourceBase,
		"Enums":         Enums,
	})
}

func (x *WebappService) CreatePrompts(c echo.Context) error {
	response := &models.AIStudioResponse{
		BackListingLink: fmt.Sprintf("%s/ui%s%s", pkgs.NetworkConfigManager.ExternalURL, routes.ManageAIGroup, routes.ManageAIPrompts),
		ActionRules:     []*models.ActionRule{},
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.AIPromptsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}
	globalContext.Manager = true
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/prompts", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AiPromptManagerRequest),
		Title:    "Enter your Prompt Spec",
	}
	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "AIResponseFormat":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "create-prompts.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Build Your Prompt",
		"Enums":         Enums,
	})
}
