package router

import (
	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
	services "github.com/vapusdata-ecosystem/vapusdata/webapp/services"
)

func manageAIRoutes(e *echo.Group) {
	aistudio := e.Group(routes.ManageAIGroup)
	aistudio.GET(routes.ManageAIModelNodes, services.WebappServiceManager.ManageAIModelNodesHandler)
	aistudio.GET(routes.ManageAIManageModelNodeResource, services.WebappServiceManager.ManageAIModelNodesDetailHandler)
	aistudio.GET(routes.ManageAIPrompts, services.WebappServiceManager.ManageAIPromptsHandler)
	aistudio.GET(routes.CreateAIPrompts, services.WebappServiceManager.CreatePrompts)
	aistudio.GET(routes.ManageAIPromptResource, services.WebappServiceManager.ManageAIPromptDetailHandler)
	aistudio.GET(routes.UpdateAIPrompts, services.WebappServiceManager.UpdatePrompts)
	aistudio.GET(routes.ManageAIAgents, services.WebappServiceManager.ManageVapusAgentsHandler)
	aistudio.GET(routes.ManageAIAgentsResource, services.WebappServiceManager.ManageVapusAgentDetailHandler)
	aistudio.GET(routes.CreateAIAgents, services.WebappServiceManager.CreateVapusAgent)
	aistudio.GET(routes.UpdateAIAgents, services.WebappServiceManager.UpdateVapusAgent)
	aistudio.GET(routes.ManageAIGuardrails, services.WebappServiceManager.ManageAIGuardrailsHandler)
	aistudio.GET(routes.CreateAIGuardrails, services.WebappServiceManager.CreateGuardrails)
	aistudio.GET(routes.ManageAIGuardrailResource, services.WebappServiceManager.ManageAIGuardrailDetailsHandler)
	aistudio.GET(routes.CreateAIModelNodes, services.WebappServiceManager.CreateModels)
	aistudio.GET(routes.UpdateAIModelNodes, services.WebappServiceManager.UpdateModels)
	aistudio.GET(routes.UpdateAIGuardrails, services.WebappServiceManager.UpdateGuardrails)
}
