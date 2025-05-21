package generic

import (
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
)

func BuildRequest(payload *prompts.GenerativePrompterPayload) *GuardRailRequest {
	resp := &GuardRailRequest{}
	for _, msg := range payload.Params.Messages {
		temp := Messages{}
		switch msg.Role {
		case aicore.USER:
			temp.Role = "user"
			temp.Content = msg.Content
		case aicore.ASSISTANT:
			temp.Role = "assistant"
			temp.Content = msg.Content
		case aicore.SYSTEM:
			temp.Role = "system"
			temp.Content = msg.Content
		}
		resp.Messages = append(resp.Messages, &temp)
	}

	for _, msg := range payload.SessionContext {
		temp := Messages{}
		switch msg.Role {
		case aicore.USER:
			temp.Role = "user"
			temp.Content = msg.Message
		case aicore.ASSISTANT:
			temp.Role = "assistant"
			temp.Content = msg.Message
		case aicore.SYSTEM:
			temp.Role = "system"
			temp.Content = msg.Message
		}
		resp.Messages = append(resp.Messages, &temp)
	}
	return resp
}
