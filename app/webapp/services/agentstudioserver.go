package services

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/models"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/webapp/pkgs"
	routes "github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
)

func (x *WebappService) AgentStudioHandler(c echo.Context) error {
	globalContext, err := x.getStudioSectionGlobals(c, routes.AgentStudioPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting agent Studio section types")
		return c.Render(http.StatusBadRequest, "400.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}
	agentId := c.QueryParam(types.AgentId)

	response := &models.AgentStudioResponse{
		ActionParams: &models.ResourceManagerParams{
			API:     fmt.Sprintf("%s/api/v1alpha1/aistudio/agents/run", pkgs.NetworkConfigManager.GatewayURL),
			ChatAPI: fmt.Sprintf("%s/api/v1alpha1/aistudio/agents/run", pkgs.NetworkConfigManager.GatewayURL),
		},
		// AIModelNode: &mpb.AIModelNode{
		// 	ModelNodeId: aiModelNode,
		// },
		AIAgent: &mpb.VapusAgent{
			AgentId: agentId,
		},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIAgents = x.grpcClients.ListVapusAgents(c)
	}()
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	response.AIModelNodes = x.grpcClients.AIModelNodes(c)
	// }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIPrompts = x.grpcClients.AIModelPrompts(c)
	}()
	wg.Wait()

	return c.Render(http.StatusOK, "agent-studio.html", map[string]any{
		"GlobalContext": globalContext,
		"Response":      response,
		"SectionHeader": "Agent Studio",
		"StepsEnum":     mpb.AgentStepEnum_name,
	})
}
