package services

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/webapp/models"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
)

func (x *WebappService) InsightsHandler(c echo.Context) error {
	aiModelNodeId := c.QueryParams()["aiModelNodeId"]
	model := c.QueryParams()["model"]
	if len(aiModelNodeId) <= 0 {
		aiModelNodeId = []string{}
	} else if len(model) <= 0 {
		model = []string{}
	}
	globalContext, _ := x.getSectionGlobals(c, routes.InsightsNav.String(), routes.LLMInsightsPage.String())
	response := &models.AIStudioResponse{
		AIModelNodeInsights: x.grpcClients.InsightsList(c, aiModelNodeId, model),
		ActionRules:         []*models.ActionRule{},
		HideBackListingLink: true,
	}

	return c.Render(http.StatusOK, "llm-insights.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Model Nodes Insights",
	})
}
