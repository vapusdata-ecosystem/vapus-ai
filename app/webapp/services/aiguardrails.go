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

func (x *WebappService) CreateGuardrails(c echo.Context) error {
	response := &models.AIStudioResponse{
		ActionRules: []*models.ActionRule{},
	}
	response.BackListingLink = fmt.Sprintf("%s/ui/ai/manage/guardrails", pkgs.NetworkConfigManager.ExternalURL)
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIGuardrailsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	globalContext.Manager = true
	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/guardrails", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AiGuardrailManagerRequest),
		Title:    "Enter your Agent Spec",
	}
	response.AIModelNodes = x.grpcClients.AIModelNodes(c)

	Enums := map[string][]string{}
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "Frequency", "ApiTokenType", "AIGuardrailScanMode", "GuardRailLevels", "ResourceScope", "ClassifiedTransformerActions":
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}
	return c.Render(http.StatusOK, "create-guardrails.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Create Guardrail",
		"Enums":         Enums,
	})
}

func (x *WebappService) ManageAIGuardrailsHandler(c echo.Context) error {
	response := &models.AIStudioResponse{
		AIGuardrails:        x.grpcClients.ListAIGuardrails(c),
		ActionRules:         []*models.ActionRule{},
		BackListingLink:     c.Request().URL.String(),
		HideBackListingLink: true,
	}
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIGuardrailsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	response.ActionRules = append(response.ActionRules, &models.ActionRule{
		Action:   mpb.ResourceLcActions_CREATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/guardrails", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPost,
		YamlSpec: GetProtoYamlString(appconfigs.AiGuardrailManagerRequest),
		Title:    "Enter your AI Guardrail Spec",
	})
	response.CreateActionParams = &models.ActionRule{
		Weblink: fmt.Sprintf("%s/ui/ai/manage/guardrails/create", pkgs.NetworkConfigManager.ExternalURL),
	}
	return c.Render(http.StatusOK, "manageai-guardrails.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "AI Guardrails",
	})
}

func (x *WebappService) UpdateGuardrails(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIGuardrailsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section globals")
		return HandleGLobalContextError(c, err)
	}

	guradrailId, err := GetUrlParams(c, types.GuardrailId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIGuardrail:        x.grpcClients.DescribeAIGuardrail(c, guradrailId),
		ActionRules:        []*models.ActionRule{},
		BackListingLink:    fmt.Sprintf("%s/ui/ai/manage/guardrails", pkgs.NetworkConfigManager.ExternalURL),
		CreateActionParams: &models.ActionRule{},
	}

	globalContext.Manager = true
	response.AIModelNodes = x.grpcClients.AIModelNodes(c)
	if response.AIGuardrail != nil {
		response.ResourceId = response.AIGuardrail.GuardrailId
		if response.AIGuardrail.ResourceBase.CreatedBy == globalContext.UserInfo.UserId && response.AIGuardrail.ResourceBase.Organization == globalContext.CurrentOrganization.OrganizationId {
			upObj := &pb.GuardrailsManagerRequest{
				Spec: response.AIGuardrail,
			}
			if upObj.Spec.Words == nil {
				upObj.Spec.Words = []*mpb.WordGuardRails{
					{
						Words: []string{},
					},
				}
			}
			if upObj.Spec.SensitiveDataset == nil {
				upObj.Spec.SensitiveDataset = []*mpb.SensitiveDataGuardrails{
					{},
				}
			}
			if upObj.Spec.Topics == nil {
				upObj.Spec.Topics = []*mpb.TopicGuardrails{
					{
						Samples: []string{},
					},
				}
			}
			bytess, err := pbtools.ProtoYamlMarshal(response.AIGuardrail)
			if err != nil {
				x.logger.Err(err).Msg("error while marshaling ai guardrails spec for publish")
			}
			response.YamlSpec = string(bytess)
		}
	} else {
		x.logger.Err(err).Msg("error while getting ai guardrail details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}

	response.CreateActionParams = &models.ActionRule{
		Action:   mpb.ResourceLcActions_UPDATE.String(),
		API:      fmt.Sprintf("%s/api/v1alpha1/aistudio/guardrails", pkgs.NetworkConfigManager.GatewayURL),
		Method:   http.MethodPut,
		YamlSpec: GetProtoYamlString(appconfigs.AinodeConfiguratorRequest),
		Title:    "Enter your Agent Spec",
	}

	Enums := map[string][]string{}
	KeyValueMap := make(map[string]map[string]int32)
	for key, value := range x.SvcPkgManager.ValidEnums {
		switch key {
		case "Frequency", "ApiTokenType", "AIGuardrailScanMode", "GuardRailLevels", "ResourceScope", "ClassifiedTransformerActions":
			KeyValueMap[key] = value
			Enums[key] = slices.Sorted(maps.Keys(value))
			Enums[key] = slices.DeleteFunc(Enums[key], func(s string) bool {
				return strings.Contains(s, "INVALID")
			})
		default:
			continue
		}
	}

	return c.Render(http.StatusOK, "update-guardrails.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Update Guardrail",
		"AIStudio":      routes.AIStudioHome,
		"ResourceBase":  response.AIGuardrail.ResourceBase,
		"Enums":         Enums,
		"KeyValueMap":   KeyValueMap,
	})
}

func (x *WebappService) ManageAIGuardrailDetailsHandler(c echo.Context) error {
	globalContext, err := x.getAiStudioSectionGlobals(c, routes.ManageAIGuardrailsPage.String())
	if err != nil {
		x.logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	guradrailId, err := GetUrlParams(c, types.GuardrailId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	response := &models.AIStudioResponse{
		AIGuardrail:     x.grpcClients.DescribeAIGuardrail(c, guradrailId),
		ActionRules:     []*models.ActionRule{},
		BackListingLink: BuildBacklistingUrl(c.Request().URL.String()),
	}
	if response.AIGuardrail != nil {
		response.ResourceId = response.AIGuardrail.GuardrailId
		if (response.AIGuardrail.ResourceBase.CreatedBy == globalContext.UserInfo.UserId || slices.Contains(response.AIGuardrail.ResourceBase.Editors, globalContext.UserInfo.UserId)) && response.AIGuardrail.ResourceBase.Organization == globalContext.CurrentOrganization.OrganizationId {
			response.CreateActionParams = &models.ActionRule{
				Weblink: fmt.Sprintf("%s/ui/ai/manage/guardrails/%s/update", pkgs.NetworkConfigManager.ExternalURL, guradrailId),
			}
			upObj := &pb.GuardrailsManagerRequest{
				Spec: response.AIGuardrail,
			}
			if upObj.Spec.Words == nil {
				upObj.Spec.Words = []*mpb.WordGuardRails{
					{
						Words: []string{},
					},
				}
			}
			if upObj.Spec.SensitiveDataset == nil {
				upObj.Spec.SensitiveDataset = []*mpb.SensitiveDataGuardrails{
					{},
				}
			}
			if upObj.Spec.Topics == nil {
				upObj.Spec.Topics = []*mpb.TopicGuardrails{
					{
						Samples: []string{},
					},
				}
			}

			bytess, err := pbtools.ProtoYamlMarshal(response.AIGuardrail)
			if err != nil {
				x.logger.Err(err).Msg("error while marshaling ai guardrails spec for publish")
			}
			response.YamlSpec = string(bytess)
			response.ActionRules = append(response.ActionRules, &models.ActionRule{
				Action:     mpb.ResourceLcActions_ARCHIVE.String(),
				API:        fmt.Sprintf("%s/api/v1alpha1/aistudio/guardrails/%s", pkgs.NetworkConfigManager.GatewayURL, response.AIGuardrail.GuardrailId),
				Method:     http.MethodDelete,
				Title:      "Archive " + response.AIGuardrail.Name,
				ResourceId: response.AIGuardrail.GuardrailId,
			})
		}
	} else {
		x.logger.Err(err).Msg("error while getting ai guardrail details")
		return x.HandleError(c, err, http.StatusNotFound, globalContext)
	}
	return c.Render(http.StatusOK, "manageai-guardrail-detail.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "AI Guardrail Details",
		"ResourceBase":  response.AIGuardrail.ResourceBase,
		"ResourceName":  response.AIGuardrail.Name,
	})
}
