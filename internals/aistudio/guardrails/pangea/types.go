package pangea

type Recipe string

// Recipe is nothing but Guardrails Models of Pangea

const (
	PanegaPromptGuard        Recipe = "pangea_prompt_guard"
	PanegaLlmResponseGuard   Recipe = "pangea_llm_response_guard"
	PanegaIngestionGuard     Recipe = "pangea_ingestion_guard"
	PanegaAgentPrePlanGuard  Recipe = "pangea_agent_pre_plan_guard"
	PanegaAgentPreToolGuard  Recipe = "pangea_agent_pre_tool_guard"
	PanegaAgentPostToolGuard Recipe = "pangea_agent_post_tool_guard"
)

func (x Recipe) String() string {
	return string(x)
}
