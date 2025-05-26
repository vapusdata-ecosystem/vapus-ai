package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/guardrails/generic"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	httpCls "github.com/vapusdata-ecosystem/vapusai/core/pkgs/http"
)

type MistralModeration struct {
	client *httpCls.RestHttp
	log    zerolog.Logger
}

func NewMistralGuardrail(ctx context.Context, token string, logger zerolog.Logger) (*MistralModeration, error) {
	httpCl, err := httpCls.New(logger,
		httpCls.WithAddress(defaultEndpoint),
		httpCls.WithBasePath(baseAPIPath),
		httpCls.WithBearerAuth(token),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating http client for Mistral")
		return nil, err
	}

	return &MistralModeration{
		client: httpCl,
		log:    logger,
	}, nil
}

func (m MistralModeration) ValidateChatInput(ctx context.Context, payload *prompts.GenerativePrompterPayload) (*MistralGuardrailResponse, error) {
	type guardInput struct {
		Input []*generic.Messages `json:"input"`
		Model string              `json:"model"`
	}

	inp := generic.BuildRequest(payload)
	input := &guardInput{
		Input: inp.Messages,
		Model: "mistral-moderation-latest",
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		m.log.Error().Err(err).Msg("Error marshalling request object")
		return nil, err
	}

	fmt.Println("Input: ", reflect.ValueOf(input))

	resp := &MistralGuardrailResponse{}
	err = m.client.Post(ctx, generatePath, inputBytes, resp, jsonContentType)
	if err != nil {
		m.log.Error().Err(err).Msg("Error while getting the response")
		return nil, err
	}
	// respBytes, err := json.Marshal(resp)
	// if err != nil {
	// 	m.log.Error().Err(err).Msg("Error marshalling request object")
	// 	return err
	// }

	return resp, nil
}
