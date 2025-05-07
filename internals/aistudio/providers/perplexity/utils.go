package perplexity

import mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"

func getHardcodedModels() ModelList {
	return ModelList{
		Data: []*Models{
			{
				ID:                      "sonar",
				OwnedBy:                 "perplexity",
				Name:                    "sonar",
				Type:                    mpb.AIModelType_LLM.String(),
				Description:             "General-purpose model for chat completions",
				MaxContextLength:        128000,
				DefaultModelTemperature: 0.7,
				Capabilities: &ModelCapabilities{
					CompletionChat:  true,
					CompletionFim:   false,
					FunctionCalling: false,
					FineTuning:      false,
					Vision:          false,
				},
			},
			{
				ID:                      "sonar-deep-research",
				OwnedBy:                 "perplexity",
				Name:                    "sonar-deep-research",
				Type:                    mpb.AIModelType_LLM.String(),
				Description:             "Designed for deep research queries with long context support",
				MaxContextLength:        128000,
				DefaultModelTemperature: 0.7,
				Capabilities: &ModelCapabilities{
					CompletionChat:  true,
					CompletionFim:   false,
					FunctionCalling: false,
					FineTuning:      false,
					Vision:          false,
				},
			},
			{
				ID:                      "sonar-reasoning-pro",
				OwnedBy:                 "perplexity",
				Name:                    "sonar-reasoning-pro",
				Type:                    mpb.AIModelType_LLM.String(),
				Description:             "Advanced reasoning capabilities with extended context",
				MaxContextLength:        128000,
				DefaultModelTemperature: 0.7,
				Capabilities: &ModelCapabilities{
					CompletionChat:  true,
					CompletionFim:   false,
					FunctionCalling: false,
					FineTuning:      false,
					Vision:          false,
				},
			},
			{
				ID:                      "sonar-reasoning",
				OwnedBy:                 "perplexity",
				Name:                    "sonar-reasoning",
				Type:                    mpb.AIModelType_LLM.String(),
				Description:             "Reasoning focused model with long context support",
				MaxContextLength:        128000,
				DefaultModelTemperature: 0.7,
				Capabilities: &ModelCapabilities{
					CompletionChat:  true,
					CompletionFim:   false,
					FunctionCalling: false,
					FineTuning:      false,
					Vision:          false,
				},
			},
			{
				ID:                      "sonar-pro",
				OwnedBy:                 "perplexity",
				Name:                    "sonar-pro",
				Type:                    mpb.AIModelType_LLM.String(),
				Description:             "Professional model with very large context window",
				MaxContextLength:        200000,
				DefaultModelTemperature: 0.7,
				Capabilities: &ModelCapabilities{
					CompletionChat:  true,
					CompletionFim:   false,
					FunctionCalling: false,
					FineTuning:      false,
					Vision:          false,
				},
			},
			{
				ID:                      "r1-1776",
				OwnedBy:                 "perplexity",
				Name:                    "r1-1776",
				Type:                    mpb.AIModelType_LLM.String(),
				Description:             "Specialized model for research-based queries",
				MaxContextLength:        128000,
				DefaultModelTemperature: 0.7,
				Capabilities: &ModelCapabilities{
					CompletionChat:  true,
					CompletionFim:   false,
					FunctionCalling: false,
					FineTuning:      false,
					Vision:          false,
				},
			},
		},
	}
}

func getToolType(t string) string {
	val, ok := ToolTypeMap[t]
	if !ok {
		return ToolTypeFunction.String()
	}
	return val.String()
}
