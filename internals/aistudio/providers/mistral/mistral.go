package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"

	// pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers/generic"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	httpCls "github.com/vapusdata-ecosystem/vapusai/core/pkgs/http"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	// dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type MistralInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
	FIM(ctx context.Context, payload *prompts.GenerativePrompterPayload, model string) error
}

type Mistral struct {
	client     *httpCls.RestHttp
	log        zerolog.Logger
	modelNode  *models.AIModelNode
	maxRetries int
	params     map[string]any
	*generic.OpenAI
}

func New(node *models.AIModelNode, retries int, logger zerolog.Logger) (MistralInterface, error) {
	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = "https://api.mistral.ai/v1" // Already defined
	}

	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	httpCl, err := httpCls.New(logger,
		httpCls.WithAddress(defaultEndpoint),
		httpCls.WithBasePath(baseAPIPath),
		httpCls.WithBearerAuth(token),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating http client for Mistral")
		return nil, err
	}

	client, err := generic.New(context.Background(), node, retries, logger)
	if err != nil {
		return nil, err
	}

	return &Mistral{
		OpenAI:    client.(*generic.OpenAI),
		client:    httpCl,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Mistral"),
		modelNode: node,
	}, nil
}

func (o *Mistral) GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	fmt.Println("----------I am using Mistralllllll-------")
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContent(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Mistral) GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error {
	fmt.Println("----------I am using Mistralllllll-------")
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateContentStream(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

func (o *Mistral) GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error {
	fmt.Println("----------I am using Mistralllllll-------")
	if o.OpenAI.ModelNode != nil && o.OpenAI.ModelNode.NetworkParams != nil && o.OpenAI.ModelNode.NetworkParams.Credentials != nil {
		return o.OpenAI.GenerateEmbeddings(ctx, request)
	}
	return apperr.ErrInvalidOrMissingPodelAPIKey
}

// func (x *Mistral) buildRequest(model string, payload *prompts.GenerativePrompterPayload, stream bool) ([]byte, error) {
// 	messages := make([]*Message, 0)

// 	if payload.Mode == pb.AIInterfaceMode_CHAT_MODE {
// 		for _, msg := range payload.SessionContext {
// 			if strings.ToLower(msg.Role) == aicore.USER {
// 				messages = append(messages, &Message{
// 					Role:    aicore.USER,
// 					Content: &msg.Message,
// 				})
// 			} else {
// 				messages = append(messages, &Message{
// 					Role:    aicore.ASSISTANT,
// 					Content: &msg.Message,
// 				})
// 			}
// 		}
// 	}
// 	for _, c := range payload.Params.Messages {
// 		switch strings.ToLower(c.Role) {
// 		case aicore.USER:
// 			messages = append(messages, &Message{
// 				Role:    aicore.USER,
// 				Content: &c.Content,
// 			})
// 		case aicore.ASSISTANT:
// 			messages = append(messages, &Message{
// 				Role:    aicore.ASSISTANT,
// 				Content: &c.Content,
// 			})
// 		case aicore.SYSTEM:
// 			messages = append(messages, &Message{
// 				Role:    aicore.SYSTEM,
// 				Content: &c.Content,
// 			})
// 		default:
// 			return nil, aicore.ErrInvalidAIModelRole
// 		}
// 	}

// 	tools := make([]Tool, 0)
// 	if len(payload.ToolCalls) > 0 {
// 		for _, tool := range payload.ToolCalls {
// 			if tool != nil && tool.FunctionSchema != nil {
// 				tools = append(tools, Tool{
// 					Type: ToolType(getToolType(tool.Type)),
// 					Function: Function{
// 						Name:        tool.FunctionSchema.Name,
// 						Description: tool.FunctionSchema.Description,
// 						Parameters:  tool.FunctionSchema.Parameters,
// 					},
// 				})
// 			}
// 		}
// 	}
// 	reqObj := &GenerativeRequest{
// 		Model:       model,
// 		Messages:    messages,
// 		Stream:      stream,
// 		Temperature: float64(payload.Params.Temperature),
// 		ToolChoice:  "auto",
// 		MaxTokens:   int(payload.Params.MaxCompletionTokens),
// 		TopP: func() float64 {
// 			if payload.Params.TopP == 0.0 {
// 				return defaultTopP
// 			}
// 			return payload.Params.TopP
// 		}(),
// 	}
// 	if len(tools) > 0 {
// 		reqObj.Tools = tools
// 		reqObj.ToolChoice = "any"
// 	}
// 	reqBytes, err := json.Marshal(reqObj)
// 	if err != nil {
// 		x.log.Error().Err(err).Msg("Error marshalling request object")
// 		return nil, err
// 	}
// 	return reqBytes, err
// }

func (x *Mistral) buildFimRequest(model string, payload *prompts.GenerativePrompterPayload) []byte {
	ip := ""
	for _, c := range payload.Params.Messages {
		if c.Role == aicore.USER {
			ip = ip + "\n" + c.Content
		}
	}
	req := &FIMRequests{
		Model:       model,
		Prompt:      ip,
		Suffix:      payload.Suffix,
		Temperature: float64(payload.Params.Temperature),
		MaxTokens:   int(payload.Params.MaxCompletionTokens),
	}
	bbytes, err := json.Marshal(req)
	if err != nil {
		x.log.Error().Err(err).Msg("Error marshalling request object for Mistral")
		return nil
	}
	return bbytes
}

// func (x *Mistral) buildResponse(resp *GenerativeResponse, payload *prompts.GenerativePrompterPayload) error {
// 	if len(resp.Choices) == 0 {
// 		x.log.Warn().Msg("No choices found in response")
// 		return aicore.ErrNoResponseFromAIModel
// 	}
// 	for _, choice := range resp.Choices {
// 		if choice.FinishReason == "tool_calls" {
// 			for _, tool := range choice.Message.ToolCalls {
// 				payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
// 					Id:   tool.Id,
// 					Type: tool.Type.String(),
// 					FunctionSchema: &mpb.FunctionCall{
// 						Name:       tool.Function.Name,
// 						Parameters: tool.Function.Arguments,
// 					},
// 				})
// 			}
// 			payload.ParseToolCallResponse()
// 			continue
// 		} else {
// 			payload.ParseOutput(&prompts.PayloadgenericResponse{
// 				FinishReason: string(resp.Choices[0].FinishReason),
// 				Data:         *resp.Choices[0].Message.Content,
// 				Role:         resp.Choices[0].Message.Role,
// 			})
// 		}
// 	}
// 	return nil
// }

// func (x *Mistral) buildStreamResponse(resp *http.Response, payload *prompts.GenerativePrompterPayload) (string, map[string]string) {
// 	// go func() {
// 	defer resp.Body.Close()
// 	reader := bufio.NewReader(resp.Body)
// 	content := ""
// 	toolCallsParams := map[string]string{}
// 	for {
// 		line, err := reader.ReadBytes('\n')
// 		if err != nil {
// 			if err == io.EOF {
// 				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
// 					Data:    "",
// 					IsEnd:   true,
// 					IsError: false,
// 					Created: dmutils.GetEpochTime(),
// 					Id:      dmutils.GetUUID(),
// 				})
// 				break
// 			} else {
// 				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
// 					Data:    "",
// 					IsEnd:   true,
// 					IsError: false,
// 					Created: dmutils.GetEpochTime(),
// 					Id:      dmutils.GetUUID(),
// 				})
// 				break
// 			}
// 		}

// 		if bytes.Equal(line, []byte("\n")) {
// 			continue
// 		}

// 		// Check if the line starts with "data: ".
// 		if bytes.HasPrefix(line, []byte("data: ")) {
// 			// Trim the prefix and any leading or trailing whitespace.
// 			jsonLine := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data: ")))
// 			// Check for the special "[DONE]" message.
// 			if bytes.Equal(jsonLine, []byte("[DONE]")) {
// 				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
// 					Data:    "",
// 					IsEnd:   true,
// 					IsError: false,
// 					Created: dmutils.GetEpochTime(),
// 					Id:      dmutils.GetUUID(),
// 				})
// 				break
// 			}

// 			streamResponse := &GenerativeResponseStream{}
// 			if err := json.Unmarshal(jsonLine, streamResponse); err != nil {
// 				continue
// 			}
// 			if len(streamResponse.Choices) == 0 {
// 				continue
// 			}
// 			lcContent := ""
// 			for _, c := range streamResponse.Choices {
// 				if len(c.Delta.Content) > 0 {
// 					lcContent = lcContent + c.Delta.Content
// 				}
// 				if len(c.Delta.ToolCalls) > 0 {
// 					for _, t := range c.Delta.ToolCalls {
// 						val, ok := toolCallsParams[t.Function.Name]
// 						if ok {
// 							toolCallsParams[t.Function.Name] = val + t.Function.Arguments
// 						} else {
// 							toolCallsParams[t.Function.Name] = t.Function.Arguments
// 						}
// 					}
// 				}
// 			}
// 			payload.LogUsage(&prompts.UsageMetrics{
// 				TotalTokens:  int64(streamResponse.Usage.TotalTokens),
// 				InputTokens:  int64(streamResponse.Usage.PromptTokens),
// 				OutputTokens: int64(streamResponse.Usage.CompletionTokens),
// 			}, nil)
// 			content = content + lcContent
// 			err = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
// 				Data:    lcContent,
// 				IsEnd:   false,
// 				IsError: false,
// 				Created: dmutils.GetEpochTime(),
// 				Id:      dmutils.GetUUID(),
// 			})
// 			if err != nil {
// 				x.log.Err(err).Msg("error while sending stream response")
// 				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
// 					Data:    "",
// 					IsEnd:   false,
// 					IsError: true,
// 					Created: dmutils.GetEpochTime(),
// 					Id:      dmutils.GetUUID(),
// 				})
// 				continue
// 			}
// 		}
// 	}
// 	// return
// 	// }()
// 	return content, toolCallsParams
// }

// func (x *Mistral) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload) error {
// 	if payload.EmbeddingModel == "" {
// 		x.log.Warn().Msg("Model name is empty, using default model")
// 		payload.EmbeddingModel = defaultEmbedModel
// 	}
// 	req := &EmbeddingRequest{
// 		Model: defaultEmbedModel,
// 		Input: payload.Input,
// 	}
// 	bbytes, err := json.Marshal(req)
// 	if err != nil {
// 		x.log.Error().Err(err).Msg("Error marshalling request object for Mistral")
// 	}
// 	resp := &EmbeddingResponse{}
// 	err = x.client.Post(ctx, embeddingsPath, bbytes, resp, jsonContentType)
// 	if err != nil {
// 		x.log.Error().Err(err).Msg("Error generating content from Mistral completion")
// 		return err
// 	}
// 	vectors := []float64{}
// 	for _, data := range resp.Data {
// 		vectors = append(vectors, data.Embedding...)
// 	}
// 	payload.Embeddings = &models.VectorEmbeddings{
// 		Vectors64: vectors,
// 	}
// 	return nil
// }

// func (x *Mistral) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
// 	if payload.Params.Model == "" {
// 		x.log.Warn().Msg("Model name is empty, using default model")
// 		payload.Params.Model = defaultModel
// 	}
// 	req, err := x.buildRequest(payload.Params.Model, payload, false)
// 	if err != nil {
// 		x.log.Err(err).Msg("Error building request object for Mistral")
// 		return err
// 	}
// 	if req == nil {
// 		x.log.Error().Msg("Error building request object for Mistral")
// 		return aicore.ErrInvalidAIModelRequest
// 	}
// 	if payload.StudioLog != nil {
// 		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
// 	}
// 	resp := &GenerativeResponse{}
// 	err = x.client.Post(ctx, generatePath, req, resp, jsonContentType)
// 	if payload.StudioLog != nil {
// 		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
// 	}
// 	if err != nil {
// 		x.log.Error().Err(err).Msg("Error generating content from Mistral completion")
// 		return err
// 	}
// 	x.buildResponse(resp, payload)
// 	payload.LogUsage(&prompts.UsageMetrics{
// 		TotalTokens:  int64(resp.Usage.TotalTokens),
// 		InputTokens:  int64(resp.Usage.PromptTokens),
// 		OutputTokens: int64(resp.Usage.CompletionTokens),
// 	}, nil)
// 	return nil
// }

// func (x *Mistral) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
// 	if payload.Params.Model == "" {
// 		x.log.Warn().Msg("Model name is empty, using default model")
// 		payload.Params.Model = defaultModel
// 	}
// 	req, err := x.buildRequest(payload.Params.Model, payload, true)
// 	if err != nil {
// 		x.log.Err(err).Msg("Error building request object for Mistral")
// 		return err
// 	}
// 	if req == nil {
// 		x.log.Error().Msg("Error building request object for mistral")
// 		return aicore.ErrInvalidAIModelRequest
// 	}
// 	if payload.StudioLog != nil {
// 		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
// 	}
// 	response, err := x.client.StreamPost(ctx, generatePath, req, jsonContentType)
// 	if err != nil {
// 		x.log.Error().Err(err).Msg("Error generating content from mistral completion")
// 		return err
// 	}
// 	if payload.StudioLog != nil {
// 		payload.StudioLog.TTFBAt = dmutils.GetMilliEpochTime()
// 	}
// 	if response == nil || response.Body == nil {
// 		x.log.Error().Msg("Error getting response from mistral")
// 		return aicore.ErrInvalidAIModelRequest
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		x.log.Error().Msg("Error getting response from mistral")
// 		return aicore.ErrInvalidAIModelRequest
// 	}
// 	content, toolCallsParams := x.buildStreamResponse(response, payload)
// 	if payload.StudioLog != nil {
// 		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
// 		payload.StudioLog.Output = []*models.MessageLog{
// 			{
// 				Role:    aicore.ASSISTANT,
// 				Content: content,
// 			},
// 		}
// 		for k, v := range toolCallsParams {
// 			payload.StudioLog.ToolCallResponse = append(payload.StudioLog.ToolCallResponse, &mpb.ToolCall{
// 				Type: aicore.FUNCTION,
// 				FunctionSchema: &mpb.FunctionCall{
// 					Name:       k,
// 					Parameters: v,
// 				},
// 			})
// 		}
// 	}
// 	return nil
// }

func (x *Mistral) CrawlModels(ctx context.Context) ([]*models.AIModelBase, error) {
	resp := &ModelList{}
	result := []*models.AIModelBase{}
	err := x.client.Get(ctx, modelsPath, nil, resp, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error crawling models from Mistral")
		return nil, err
	}
	for _, model := range resp.Data {
		result = append(result, &models.AIModelBase{
			ModelId:   model.ID,
			OwnedBy:   model.OwnedBy,
			ModelName: model.Name,
			ModelType: func() string {
				if strings.Contains(model.Type, "embed") {
					return mpb.AIModelType_EMBEDDING.String()
				}
				return mpb.AIModelType_LLM.String()

			}(),
			ModelNature: func() []string {
				if strings.Contains(model.Type, "embed") {
					return []string{
						mpb.AIModelType_EMBEDDING.String(),
					}
				}
				return []string{
					mpb.AIModelType_LLM.String(),
				}

			}(),
		})
	}
	return result, nil
}

func (x *Mistral) FIM(ctx context.Context, payload *prompts.GenerativePrompterPayload, model string) error {
	var err error
	if model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		model = defaultModel
	}
	req := x.buildFimRequest(model, payload)
	if req == nil {
		x.log.Error().Msg("Error marshalling request object for Mistral")
	}
	resp := &FIMResponse{}
	err = x.client.Post(ctx, fimPath, req, resp, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating FIM content from Mistral completion")
		return err
	}
	payload.ParseOutput(&prompts.PayloadgenericResponse{
		FinishReason: string(resp.Choices[0].FinishReason),
		Data:         resp.Object,
		Role:         resp.Choices[0].Message.Role,
	})
	return nil
}
