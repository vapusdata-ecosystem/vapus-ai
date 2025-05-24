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

// type PangeaGuardrailModels string

// const (

// )
