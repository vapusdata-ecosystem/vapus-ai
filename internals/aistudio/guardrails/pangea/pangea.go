package pangea

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ai_guard"
	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails/generic"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
)

type PanegaOpts struct {
	Token  string
	Domain string
}

type PanegaStore struct {
	client ai_guard.Client
	logger zerolog.Logger
}

// In Panega
//
//	// ai_guard:= Post-generation (output evaluation)
//	// prompt_guard:= Pre-generation (input prompt modification)
//  // Recipe := https://pangea.cloud/docs/ai-guard/recipes
//	// https://console.pangea.cloud/service/ai-guard/recipes?referrer=docs-ai-guard-recipes&_gl=1*io6mwm*_gcl_au*MTE1MzYxMDMzOC4xNzQ3MjExNzY5*_ga*MTMxMzM4NzY1Ny4xNzQ3MjExNzcy*_ga_TWJLY45X2R*czE3NDczODkyNDQkbzQkZzEkdDE3NDczOTE3MDEkajAkbDAkaDA.*_ga_2DXDPDGJ1L*czE3NDczODkyNDQkbzQkZzEkdDE3NDczOTE3MDEkajAkbDAkaDA.*_ga_BKCWGP6DYY*czE3NDczODkyNDQkbzQkZzEkdDE3NDczOTE3MDEkajAkbDAkaDA.*_ga_FBJYCW7BT0*czE3NDczODkyNDQkbzQkZzEkdDE3NDczOTE3MDEkajAkbDAkaDA.

// ["pangea_prompt_guard", "pangea_llm_response_guard", "pangea_ingestion_guard", "pangea_agent_pre_plan_guard", "pangea_agent_pre_tool_guard", "pangea_agent_post_tool_guard"]

func NewPangeaGuardrail(ctx context.Context, cred *PanegaOpts, logger zerolog.Logger) (*PanegaStore, error) {
	if cred == nil || cred.Token == "" || cred.Domain == "" {
		return nil, errors.New("invalid or nil pangea credential")
	}
	// Client configuration.
	aiGuardClient := ai_guard.New(&pangea.Config{
		Token:  cred.Token,
		Domain: cred.Domain,
	})

	return &PanegaStore{
		client: aiGuardClient,
	}, nil
}

// https://github.com/pangeacyber/pangea-go/blob/main/pangea-sdk/service/ai_guard/integration_test.go

func (p PanegaStore) ValidateChatInput(ctx context.Context, payload *prompts.GenerativePrompterPayload) (*PangeaResponse, error) {
	inp := generic.BuildRequest(payload)
	input := &ai_guard.TextGuardRequest{
		Messages: inp.Messages,
		Recipe:   "pangea_llm_response_guard",
	}

	fmt.Println("Input: ", reflect.ValueOf(input))

	out, err := p.client.GuardText(ctx, input)
	if err != nil {
		p.logger.Err(err).Msg("failed to pangea guard messages")
		return nil, err
	}
	result, err := p.buildResponse(ctx, out)
	if err != nil {
		p.logger.Err(err).Msg("failed to convert the pangea guard message")
		return nil, err
	}

	return result, nil
}

func (p PanegaStore) buildResponse(ctx context.Context, resp *pangea.PangeaResponse[ai_guard.TextGuardResult]) (*PangeaResponse, error) {
	if resp == nil || resp.Result == nil {
		return nil, errors.New("invalid or nil response")
	}
	p.logger.Log().Msg("I Building response of panega")

	var promptMessages []*PromptMessages
	if raw, ok := resp.Result.PromptMessages.([]interface{}); ok {
		for _, item := range raw {
			if m, ok := item.(map[string]interface{}); ok {
				promptMessages = append(promptMessages, &PromptMessages{
					Content: m["content"].(string),
					Role:    m["role"].(string),
				})
			}
		}
	}

	d := resp.Result.Detectors
	detectors := &Detectors{
		PromptInjection:      convertDetector(d.PromptInjection, analyze),
		Gibberish:            convertDetector(d.Gibberish, classify),
		Sentiment:            convertDetector(d.Sentiment, classify),
		SelfHarm:             convertDetector(d.SelfHarm, classify),
		PiiEntity:            convertPiiEntities(d.PiiEntity),
		MaliciousEntity:      convertMaliciousEntities(d.MaliciousEntity),
		CustomEntity:         convertRedactEntities(d.CustomEntity),
		SecretsDetection:     convertSecretsEntities(d.SecretsDetection),
		Competitors:          convertSingleEntity(d.Competitors),
		ProfanityAndToxicity: convertDetector(d.ProfanityAndToxicity, classify),
		LanguageDetection:    convertDetector(d.LanguageDetection, lang),
		TopicDetection:       convertDetector(d.TopicDetection, topic),
		CodeDetection:        convertDetector(d.CodeDetection, code),
	}

	result := &Result{
		Recipe:         resp.Result.Recipe,
		Blocked:        resp.Result.Blocked,
		PromptMessages: promptMessages,
		Detectors:      detectors,
	}

	return &PangeaResponse{
		RequestID:    *resp.RequestID,
		RequestTime:  *resp.RequestTime,
		ResponseTime: *resp.ResponseTime,
		Status:       *resp.Status,
		Summary:      *resp.Summary,
		Result:       result,
	}, nil
}
