package types

import mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"

type GuardrailsProvider string

const (
	Bedrock GuardrailsProvider = "bedrock"
	Mistral GuardrailsProvider = "mistral"
	Pangea  GuardrailsProvider = "pangea"
	Vapus   GuardrailsProvider = "vapus"
	// Nemo    GuardrailsProvider = "nemo"
)

func (x GuardrailsProvider) String() string {
	return string(x)
}

// Grok -> SearchParameter Source Constraints
var GrokSearchParameterSourceType = map[mpb.SearchParameterSources][]string{
	mpb.SearchParameterSources_WEB: {
		"country", "excluded_websites", "safe_search",
	},
	mpb.SearchParameterSources_X: {
		"x_handles",
	},
	mpb.SearchParameterSources_NEWS: {
		"country", "excluded_websites", "safe_search",
	},
	mpb.SearchParameterSources_RES: {
		"links",
	},
}

type PangeaGuardrailModels string

const (
	PangeaPromptGuard        PangeaGuardrailModels = "pangea_prompt_guard"
	PangeaLlmResponseGuard   PangeaGuardrailModels = "pangea_llm_response_guard"
	PangeaIngestionGuard     PangeaGuardrailModels = "pangea_ingestion_guard"
	PangeaAgentPrePlanGuard  PangeaGuardrailModels = "pangea_agent_pre_plan_guard"
	PangeaAgentPreToolGuard  PangeaGuardrailModels = "pangea_agent_pre_tool_guard"
	PangeaAgentPostToolGuard PangeaGuardrailModels = "pangea_agent_post_tool_guard"
)

func (a PangeaGuardrailModels) String() string {
	return string(a)
}

var PanegaGuardrailList = []PangeaGuardrailModels{
	PangeaPromptGuard, PangeaLlmResponseGuard, PangeaIngestionGuard, PangeaAgentPrePlanGuard, PangeaAgentPreToolGuard, PangeaAgentPostToolGuard,
}

// ["pangea_prompt_guard", "pangea_llm_response_guard", "pangea_ingestion_guard", "pangea_agent_pre_plan_guard", "pangea_agent_pre_tool_guard", "pangea_agent_post_tool_guard"]

type MistralGuardrailModels string

const (
	MistralModerationLatest MistralGuardrailModels = "mistral-moderation-latest"
)

func (a MistralGuardrailModels) String() string {
	return string(a)
}

var MistralGuardrailList = []MistralGuardrailModels{
	MistralModerationLatest,
}

// "Mistral": {"mistral-moderation-latest"},
