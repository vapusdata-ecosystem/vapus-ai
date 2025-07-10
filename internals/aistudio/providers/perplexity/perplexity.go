package perplexity

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	httpCls "github.com/vapusdata-ecosystem/vapusai/core/pkgs/http"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"

	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
)

type PerplexityInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type Perplexity struct {
	client    *httpCls.RestHttp
	log       zerolog.Logger
	modelNode *models.AIModelNode
}

func New(node *models.AIModelNode, retries int, logger zerolog.Logger) (PerplexityInterface, error) {
	// if node.NetworkParams.Url == "" {
	// 	node.NetworkParams.Url = "https://api.perplexity.ai"	// Already defined
	// }

	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	httpCl, err := httpCls.New(logger,
		httpCls.WithAddress(defaultEndpoint),
		httpCls.WithBearerAuth(token),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating http client for Perplexity")
		return nil, err
	}

	return &Perplexity{
		client:    httpCl,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Perplexity"),
		modelNode: node,
	}, nil
}

func (x *Perplexity) GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error {
	return nil
}

func (x *Perplexity) buildRequest(model string, payload *prompts.GenerativePrompterPayload, stream bool) ([]byte, error) {
	messages := make([]*Message, 0)
	fmt.Println("I am in built request")
	if payload.Mode == pb.AIInterfaceMode_CHAT_MODE {
		for _, msg := range payload.SessionContext {
			if strings.ToLower(msg.Role) == aicore.USER {
				message := &Message{
					Role:    aicore.USER,
					Content: []*Content{},
				}
				modelContent := &Content{}

				if msg.StructuredMessage != nil {
					for _, d := range msg.StructuredMessage {
						if d.Type == aicore.AIResponseFormatImageUrl.String() {
							if d.ImageUrl != nil {
								modelContent.Type = aicore.AIResponseFormatImageUrl.String()
								modelContent.ImageURL.URL = d.ImageUrl.Url
							} else {
								continue
							}
						} else {
							modelContent.Type = aicore.AIResponseFormatText.String()
							modelContent.Text = d.Text
						}
					}
				} else {
					fmt.Println("msg.Content is empty")
					modelContent.Type = "text"
					modelContent.Text = msg.Message
				}
				message.Content = append(message.Content, modelContent)
				messages = append(messages, message)
			} else {
				message := &Message{
					Role:    aicore.ASSISTANT,
					Content: []*Content{},
				}
				modelContent := &Content{}
				modelContent.Type = "text"
				modelContent.Text = msg.Message
				message.Content = append(message.Content, modelContent)
				messages = append(messages, message)
			}
		}
	}

	for _, c := range payload.Params.Messages {
		switch strings.ToLower(c.Role) {
		case aicore.USER:
			message := &Message{
				Role:    aicore.USER,
				Content: []*Content{},
			}
			modelContent := &Content{}
			// In case if we add the feature of image in chat in future
			if c.StructuredContent != nil {
				for _, d := range c.StructuredContent {
					if d.Type == aicore.AIResponseFormatImageUrl.String() {
						if d.ImageUrl != nil {
							modelContent.Type = aicore.AIResponseFormatImageUrl.String()
							modelContent.ImageURL.URL = d.ImageUrl.Url
						} else {
							continue
						}
					} else {
						fmt.Println("Content Type is text")
						modelContent.Type = "text"
						modelContent.Text = d.Text
					}
				}
			} else {
				fmt.Println("ContentPart is Empty")
				modelContent.Type = "text"
				modelContent.Text = c.Content
			}
			message.Content = append(message.Content, modelContent)
			messages = append(messages, message)
		case aicore.ASSISTANT:
			message := &Message{
				Role:    aicore.ASSISTANT,
				Content: []*Content{},
			}
			modelContent := &Content{}
			modelContent.Type = "text"
			modelContent.Text = c.Content
			message.Content = append(message.Content, modelContent)
			messages = append(messages, message)
		case aicore.SYSTEM:
			message := &Message{
				Role:    aicore.ASSISTANT,
				Content: []*Content{},
			}
			modelContent := &Content{}
			modelContent.Type = "text"
			modelContent.Text = c.Content
			message.Content = append(message.Content, modelContent)
			messages = append(messages, message)
		default:
			return nil, aicore.ErrInvalidAIModelRole
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
		Model:       model,
		Messages:    messages,
		Stream:      stream,
		Temperature: float64(payload.Params.Temperature),
		ToolChoice:  "auto",
		MaxTokens:   int(payload.Params.MaxCompletionTokens),
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
		x.log.Error().Err(err).Msg("Error marshalling request object")
		return nil, err
	}
	fmt.Println("Request built successfully")
	return reqBytes, err
}

func (x *Perplexity) buildResponse(resp *GenerativeResponse, payload *prompts.GenerativePrompterPayload) error {
	if len(resp.Choices) == 0 {
		x.log.Warn().Msg("No choices found in response")
		return aicore.ErrNoResponseFromAIModel
	}
	for _, choice := range resp.Choices {
		if choice.FinishReason == "tool_calls" {
			for _, tool := range choice.Message.ToolCalls {
				payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
					Id:   tool.Id,
					Type: tool.Type.String(),
					FunctionSchema: &mpb.FunctionCall{
						Name:       tool.Function.Name,
						Parameters: tool.Function.Arguments,
					},
				})
			}
			payload.ParseToolCallResponse()
			continue
		} else {
			payload.ParseOutput(&prompts.PayloadgenericResponse{
				Citations:    resp.Citations,
				FinishReason: string(resp.Choices[0].FinishReason),
				Data:         resp.Choices[0].Message.Content,
				Role:         resp.Choices[0].Message.Role,
			})
		}
	}
	return nil
}

func (x *Perplexity) buildStreamResponse(resp *http.Response, payload *prompts.GenerativePrompterPayload) (string, map[string]string) {
	// go func() {
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	content := ""
	toolCallsParams := map[string]string{}
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:    "",
					IsEnd:   true,
					IsError: false,
					Created: dmutils.GetEpochTime(),
					Id:      dmutils.GetUUID(),
				})
				break
			} else {
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:    "",
					IsEnd:   true,
					IsError: false,
					Created: dmutils.GetEpochTime(),
					Id:      dmutils.GetUUID(),
				})
				break
			}
		}

		if bytes.Equal(line, []byte("\n")) {
			continue
		}

		if bytes.HasPrefix(line, []byte("data: ")) {
			jsonLine := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data: ")))
			if bytes.Equal(jsonLine, []byte("[DONE]")) {
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:    "",
					IsEnd:   true,
					IsError: false,
					Created: dmutils.GetEpochTime(),
					Id:      dmutils.GetUUID(),
				})
				break
			}

			streamResponse := &GenerativeResponseStream{}
			if err := json.Unmarshal(jsonLine, streamResponse); err != nil {
				continue
			}
			if len(streamResponse.Choices) == 0 {
				continue
			}
			lcContent := ""
			for _, c := range streamResponse.Choices {
				if len(c.Delta.Content) > 0 {
					lcContent = lcContent + c.Delta.Content
				}
				if len(c.Delta.ToolCalls) > 0 {
					for _, t := range c.Delta.ToolCalls {
						val, ok := toolCallsParams[t.Function.Name]
						if ok {
							toolCallsParams[t.Function.Name] = val + t.Function.Arguments
						} else {
							toolCallsParams[t.Function.Name] = t.Function.Arguments
						}
					}
				}
			}
			payload.LogUsage(&prompts.UsageMetrics{
				TotalTokens:  int64(streamResponse.Usage.TotalTokens),
				InputTokens:  int64(streamResponse.Usage.PromptTokens),
				OutputTokens: int64(streamResponse.Usage.CompletionTokens),
			})
			content = content + lcContent
			err = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
				Data:    lcContent,
				IsEnd:   false,
				IsError: false,
				Created: dmutils.GetEpochTime(),
				Id:      dmutils.GetUUID(),
			})
			if err != nil {
				x.log.Err(err).Msg("error while sending stream response")
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:    "",
					IsEnd:   false,
					IsError: true,
					Created: dmutils.GetEpochTime(),
					Id:      dmutils.GetUUID(),
				})
				continue
			}
		}
	}
	// return
	// }()
	return content, toolCallsParams
}

func (x *Perplexity) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	fmt.Println("I am in Generate Content")
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	req, err := x.buildRequest(payload.Params.Model, payload, false)
	if err != nil {
		x.log.Err(err).Msg("Error building request object for Perplexity")
		return err
	}
	if req == nil {
		x.log.Error().Msg("Error building request object for Perplexity")
		return aicore.ErrInvalidAIModelRequest
	}
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	resp := &GenerativeResponse{}
	err = x.client.Post(ctx, generatePath, req, resp, jsonContentType)
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	// Yaha dikkat kaise aa raha hai???
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from Perplexity completion")
		return err
	}
	x.buildResponse(resp, payload)
	payload.LogUsage(&prompts.UsageMetrics{
		TotalTokens:  int64(resp.Usage.TotalTokens),
		InputTokens:  int64(resp.Usage.PromptTokens),
		OutputTokens: int64(resp.Usage.CompletionTokens),
	})
	return nil
}

func (x *Perplexity) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	req, err := x.buildRequest(payload.Params.Model, payload, true)
	if err != nil {
		x.log.Err(err).Msg("Error building request object for Perplexity")
		return err
	}
	if req == nil {
		x.log.Error().Msg("Error building request object for Perplexity")
		return aicore.ErrInvalidAIModelRequest
	}
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	response, err := x.client.StreamPost(ctx, generatePath, req, jsonContentType)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content stream from Perplexity completion")
		return err
	}
	if payload.StudioLog != nil {
		payload.StudioLog.TTFBAt = dmutils.GetMilliEpochTime()
	}
	if response == nil || response.Body == nil {
		x.log.Error().Msg("Error getting response from Perplexity")
		return aicore.ErrInvalidAIModelRequest
	}
	if response.StatusCode != http.StatusOK {
		x.log.Error().Msg("Error getting response from Perplexity")
		return aicore.ErrInvalidAIModelRequest
	}
	log.Println(dmutils.GetEpochTime(), "====================1")
	log.Println("Response from Perplexity stream", response)
	content, toolCallsParams := x.buildStreamResponse(response, payload)
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
		payload.StudioLog.Output = []*models.MessageLog{
			{
				Role:    aicore.ASSISTANT,
				Content: content,
			},
		}
		for k, v := range toolCallsParams {
			payload.StudioLog.ToolCallResponse = append(payload.StudioLog.ToolCallResponse, &mpb.ToolCall{
				Type: aicore.FUNCTION,
				FunctionSchema: &mpb.FunctionCall{
					Name:       k,
					Parameters: v,
				},
			})
		}
	}
	log.Println(dmutils.GetEpochTime(), "====================2")
	return nil
}

func (x *Perplexity) CrawlModels(ctx context.Context) ([]*models.AIModelBase, error) {
	result := []*models.AIModelBase{}

	resp := getHardcodedModels()
	for _, model := range resp.Data {
		result = append(result, &models.AIModelBase{
			ModelId:     model.ID,
			OwnedBy:     model.OwnedBy,
			ModelName:   model.Name,
			ModelType:   mpb.AIModelType_LLM.String(),
			ModelNature: []string{mpb.AIModelType_LLM.String()},
		})
	}
	return result, nil
}
