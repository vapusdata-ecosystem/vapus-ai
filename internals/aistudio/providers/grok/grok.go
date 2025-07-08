package grok

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/generic"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	httpCls "github.com/vapusdata-ecosystem/vapusai/core/pkgs/http"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type GrokInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error // No Embedding models for Grok till now
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error
	GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type Grok struct {
	logger     zerolog.Logger
	httpClient *httpCls.RestHttp
	*generic.OpenAI
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (GrokInterface, error) {
	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = defaultEndpoint // its working
	}
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}

	httpCl, _ := httpCls.New(logger,
		httpCls.WithAddress(defaultEndpoint),
		httpCls.WithBasePath(baseAPIPath),
		httpCls.WithBearerAuth(token),
	)
	client, err := generic.New(ctx, node, retries, logger)
	if err != nil {
		return nil, err
	}

	return &Grok{
		OpenAI:     client.(*generic.OpenAI),
		httpClient: httpCl,
		logger:     logger,
	}, nil
}

func (o *Grok) buildRequest(model string, payload *prompts.GenerativePrompterPayload, stream bool) ([]byte, error) {
	messages := make([]*Message, 0)

	if payload.Mode == pb.AIInterfaceMode_CHAT_MODE {
		for _, msg := range payload.SessionContext {
			if strings.ToLower(msg.Role) == aicore.USER {
				messages = append(messages, &Message{
					Role:    aicore.USER,
					Content: &msg.Message,
				})
			} else {
				messages = append(messages, &Message{
					Role:    aicore.ASSISTANT,
					Content: &msg.Message,
				})
			}
		}
	}
	// for P2P mode
	for _, c := range payload.Params.Messages {
		switch strings.ToLower(c.Role) {
		case aicore.USER:
			messages = append(messages, &Message{
				Role:    aicore.USER,
				Content: &c.Content,
			})
		case aicore.ASSISTANT:
			messages = append(messages, &Message{
				Role:    aicore.ASSISTANT,
				Content: &c.Content,
			})
		case aicore.SYSTEM:
			messages = append(messages, &Message{
				Role:    aicore.SYSTEM,
				Content: &c.Content,
			})
		default:
			return nil, aicore.ErrInvalidAIModelRole
		}
	}
	searchParameters := SearchParameters{}
	// Mode is Chat
	if payload.SessionContext[0].SearchParameters != nil {
		searchParameters = SearchParameters{
			Mode:             payload.SessionContext[0].SearchParameters.Mode,
			ReturnCitation:   payload.SessionContext[0].SearchParameters.ReturnCitation,
			FromDate:         payload.SessionContext[0].SearchParameters.FromDate,
			ToDate:           payload.SessionContext[0].SearchParameters.ToDate,
			MaxSearchResults: payload.SessionContext[0].SearchParameters.MaxSearchResult,
			Sources:          []*Source{},
		}
		if payload.SessionContext[0].SearchParameters.Sources != nil {
			tempSources := []*Source{}
			for _, val := range payload.SessionContext[0].SearchParameters.Sources {
				tempSource := &Source{
					ExcludeWebsite: []string{},
					Links:          []string{},
					XHandles:       []string{},
				}
				switch val.Type {
				case mpb.SearchParameterSources_WEB.String():
					tempSource.Type = mpb.SearchParameterSources_WEB.String()
					if val.Country != "" {
						tempSource.Country = val.Country
					}
					if val.ExcludeWebsite != nil {
						tempSource.ExcludeWebsite = val.ExcludeWebsite
					}
					if val.SafeSearch {
						tempSource.SafeSearch = val.SafeSearch
					}

				case mpb.SearchParameterSources_NEWS.String():
					tempSource.Type = mpb.SearchParameterSources_NEWS.String()
					if val.Country != "" {
						tempSource.Country = val.Country
					}
					if val.ExcludeWebsite != nil {
						tempSource.ExcludeWebsite = val.ExcludeWebsite
					}
					if val.SafeSearch {
						tempSource.SafeSearch = val.SafeSearch
					}

				case mpb.SearchParameterSources_X.String():
					tempSource.Type = mpb.SearchParameterSources_X.String()
					if val.XHandles != nil {
						tempSource.XHandles = val.XHandles
					}

				case mpb.SearchParameterSources_RES.String():
					tempSource.Type = mpb.SearchParameterSources_RES.String()
					if val.Links != nil {
						tempSource.Links = val.Links
					}
				}
				tempSources = append(tempSources, tempSource)
			}
		}
	}

	// if Mode is P2P
	if payload.Params.SearchParameters != nil {
		searchParameters = SearchParameters{
			Mode:             payload.Params.SearchParameters.Mode,
			ReturnCitation:   payload.Params.SearchParameters.ReturnCitation,
			FromDate:         payload.Params.SearchParameters.FromDate,
			ToDate:           payload.Params.SearchParameters.ToDate,
			MaxSearchResults: payload.Params.SearchParameters.MaxSearchResult,
			Sources:          []*Source{},
		}
		if payload.Params.SearchParameters.Sources != nil {
			tempSources := []*Source{}
			for _, val := range payload.Params.SearchParameters.Sources {
				tempSource := &Source{
					ExcludeWebsite: []string{},
					Links:          []string{},
					XHandles:       []string{},
				}
				switch val.Type {
				case mpb.SearchParameterSources_WEB.String():
					tempSource.Type = mpb.SearchParameterSources_WEB.String()
					if val.Country != "" {
						tempSource.Country = val.Country
					}
					if val.ExcludeWebsite != nil {
						tempSource.ExcludeWebsite = val.ExcludeWebsite
					}
					if val.SafeSearch {
						tempSource.SafeSearch = val.SafeSearch
					}

				case mpb.SearchParameterSources_NEWS.String():
					tempSource.Type = mpb.SearchParameterSources_NEWS.String()
					if val.Country != "" {
						tempSource.Country = val.Country
					}
					if val.ExcludeWebsite != nil {
						tempSource.ExcludeWebsite = val.ExcludeWebsite
					}
					if val.SafeSearch {
						tempSource.SafeSearch = val.SafeSearch
					}

				case mpb.SearchParameterSources_X.String():
					tempSource.Type = mpb.SearchParameterSources_X.String()
					if val.XHandles != nil {
						tempSource.XHandles = val.XHandles
					}

				case mpb.SearchParameterSources_RES.String():
					tempSource.Type = mpb.SearchParameterSources_RES.String()
					if val.Links != nil {
						tempSource.Links = val.Links
					}
				}
				tempSources = append(tempSources, tempSource)
			}
		}
	}
	tools := make([]Tool, 0)
	if len(payload.ToolCalls) > 0 {
		for _, tool := range payload.ToolCalls {
			if tool != nil && tool.FunctionSchema != nil {
				tools = append(tools, Tool{
					Type: ToolType(getToolType(tool.Type)),
					Function: Function{
						Name:        tool.FunctionSchema.Name,
						Description: tool.FunctionSchema.Description,
						Parameters:  tool.FunctionSchema.Parameters,
					},
				})
			}
		}
	}
	reqObj := &GenerativeRequest{
		Model:            model,
		Messages:         messages,
		SearchParameters: &searchParameters,
		Stream:           stream,
		Temperature:      float64(payload.Params.Temperature),
		ToolChoice:       "auto",
		MaxTokens:        int(payload.Params.MaxCompletionTokens),
		TopP: func() float64 {
			if payload.Params.TopP == 0.0 {
				return defaultTopP
			}
			return payload.Params.TopP
		}(),
	}
	if len(tools) > 0 {
		reqObj.Tools = tools
		reqObj.ToolChoice = "any"
	}
	reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		o.logger.Error().Err(err).Msg("Error marshalling request object")
		return nil, err
	}
	return reqBytes, err
}

func (o *Grok) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		o.logger.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	req, err := o.buildRequest(payload.Params.Model, payload, false)
	if err != nil {
		o.logger.Err(err).Msg("Error building request object for Mistral")
		return err
	}
	if req == nil {
		o.logger.Error().Msg("Error building request object for Mistral")
		return aicore.ErrInvalidAIModelRequest
	}
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	resp := &GenerativeResponse{}
	err = o.httpClient.Post(ctx, generatePath, req, resp, jsonContentType)
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	if err != nil {
		o.logger.Error().Err(err).Msg("Error generating content from Mistral completion")
		return err
	}
	// o.buildResponse(resp, payload)

	// // Note := Building Response is not done yet, we need to make it.
	// // Since I don't have the api key, I can't hit the request.

	// // Response will be changed because we will be getting citations also...
	// // Website previous mention response already created

	payload.LogUsage(&prompts.UsageMetrics{
		TotalTokens:  int64(resp.Usage.TotalTokens),
		InputTokens:  int64(resp.Usage.PromptTokens),
		OutputTokens: int64(resp.Usage.CompletionTokens),
	})
	return nil
}

func (o *Grok) GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContentStream(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Grok) GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateEmbeddings(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Grok) GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranscription(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Grok) GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error {
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateTranslation(ctx, payload)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}
