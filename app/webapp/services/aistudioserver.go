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

func (x *WebappService) AIStudioHandler(c echo.Context) error {
	globalContext, err := x.getStudioSectionGlobals(c, routes.AIStudioPage.String())
	if err != nil {
		Logger.Err(err).Msg("error while getting aistudio section types")
		return HandleGLobalContextError(c, err)
	}
	promptId := c.QueryParam(types.PromptId)
	aiModelNode := c.QueryParam("aiModelNode")
	chatId := c.QueryParam(types.AIStudioChatId)
	createNewChat := c.QueryParam("createNewChat")
	if createNewChat == "true" {
		resp, err := x.grpcClients.CreateAIGatewayChat(c)
		if err != nil || len(resp) == 0 {
			Logger.Err(err).Msg("error while creating new chat")
			return c.Render(http.StatusBadRequest, "400.html", map[string]any{
				"GlobalContext": globalContext,
			})
		}
		chatId = resp[0].ChatId
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s?%s=%s", globalContext.CurrentUrl, types.AIStudioChatId, chatId))
	}
	response := &models.AIStudioResponse{
		HideBackListingLink: true,
		AIPrompt: &mpb.AIPrompt{
			PromptId: promptId,
		},
		AIModelNode: &mpb.AIModelNode{
			Name: aiModelNode,
		},
		ActionParams: &models.ResourceManagerParams{
			API:     fmt.Sprintf("%s/gateway/v1/chat/completions", pkgs.NetworkConfigManager.GatewayURL),
			ChatAPI: fmt.Sprintf("%s/api/v1alpha1/aistudio/chat", pkgs.NetworkConfigManager.GatewayURL),
		},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIModelNodes = x.grpcClients.AIModelNodes(c)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIPrompts = x.grpcClients.AIModelPrompts(c)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		response.AIStudioChats = x.grpcClients.ListAIGatewayChats(c)
	}()
	if chatId != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := x.grpcClients.GetAIGatewayChat(c, chatId)
			if err != nil || len(res) < 1 {
				Logger.Err(err).Msg("error while getting aistudio chat")
				return
			}
			response.AIStudioChat = res[0]
		}()
	}
	wg.Wait()
	if chatId != "" && response.AIStudioChat == nil {
		Logger.Err(err).Msg("error while getting aistudio chat")
		return c.Render(http.StatusNotFound, "404.html", map[string]any{
			"GlobalContext": globalContext,
		})
	}

	return c.Render(http.StatusOK, "aistudio.html", map[string]any{
		"GlobalContext":   globalContext,
		"Response":        response,
		"SectionHeader":   "AI Studio Interface",
		"StartNewChatAPI": fmt.Sprintf("%s/api/v1alpha1/aistudio/chats", pkgs.NetworkConfigManager.GatewayURL),
		"AIPromptAPI":     fmt.Sprintf("%s/api/v1alpha1/aistudio/prompts", pkgs.NetworkConfigManager.GatewayURL),
	})
}
