package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
	"github.com/rs/zerolog"

	//mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type AnthropicAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type AnthropicAI struct {
	Client    *anthropic.Client
	Log       zerolog.Logger
	ModelNode *models.AIModelNode
	Params    map[string]interface{}
}

const (
	defaultModel = "claude-3-5-sonnet-20241022"
)

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (AnthropicAIInterface, error) {
	// if node.NetworkParams.Url == "" {
	// node.NetworkParams.Url = "https://api.anthropic.com"	// No need of URL
	// }

	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	client := anthropic.NewClient(
		option.WithAPIKey(token),
	)
	return &AnthropicAI{
		Client:    &client,
		Log:       dmlogger.GetSubDMLogger(logger, "ailogger", "openai"),
		ModelNode: node,
	}, nil
}

func (o *AnthropicAI) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload) error {

	return nil
}

func (o *AnthropicAI) buildRequest(payload *prompts.GenerativePrompterPayload) anthropic.MessageNewParams {
	messages := []anthropic.MessageParam{}
	tools := []anthropic.ToolParam{}

	if payload.Mode == pb.AIInterfaceMode_CHAT_MODE {
		for _, c := range payload.SessionContext {
			if c.Role == aicore.USER {
				messages = append(messages, anthropic.NewUserMessage(BuildContentBlock(&pb.ChatMessageObject{
					Content:           c.Message,
					StructuredContent: c.StructuredMessage,
				})...))
			} else {
				messages = append(messages, anthropic.NewAssistantMessage(BuildContentBlock(&pb.ChatMessageObject{
					Content:           c.Message,
					StructuredContent: c.StructuredMessage,
				})...))
			}
		}
	}

	systemMessages := []anthropic.TextBlockParam{}
	for _, c := range payload.Params.Messages {
		switch c.Role {
		case aicore.USER:
			mess := anthropic.NewUserMessage(BuildContentBlock(c)...)
			messages = append(messages, mess)
		case aicore.SYSTEM:
			systemMessages = append(systemMessages, anthropic.TextBlockParam{Text: c.Content})
		case aicore.ASSISTANT:
			mess := anthropic.NewAssistantMessage(BuildContentBlock(c)...)
			messages = append(messages, mess)
		}
	}

	if len(payload.ToolCalls) > 0 {
		for _, tool := range payload.ToolCalls {
			if tool != nil && tool.FunctionSchema != nil {
				var paramObj map[string]interface{}
				err := json.Unmarshal([]byte(tool.FunctionSchema.Parameters), &paramObj)
				if err != nil {
					o.Log.Err(err).Msg("Error while unmarshalling tool call arguments")
					continue
				}
				ttype, ok := paramObj["type"].(string)
				if !ok {
					o.Log.Err(err).Msg("Error while getting type from tool call arguments")
					continue
				}
				properties, ok := paramObj["properties"].(map[string]interface{})
				if !ok {
					o.Log.Err(err).Msg("Error while getting properties from tool call arguments")
					continue
				}
				tools = append(tools, anthropic.ToolParam{
					Name:        tool.FunctionSchema.Name,
					Description: anthropic.String(tool.FunctionSchema.Description),
					InputSchema: anthropic.ToolInputSchemaParam{
						Type:       constant.Object(ttype),
						Properties: properties,
					},
				})
			}
		}
	}
	request := anthropic.MessageNewParams{
		Model:    payload.Params.Model,
		System:   systemMessages,
		Messages: messages,
	}
	if len(tools) > 0 {
		tools := make([]anthropic.ToolUnionParam, len(tools))
		for i, toolParam := range tools {
			tools[i] = anthropic.ToolUnionParam{OfTool: toolParam.OfTool}
		}
		request.Tools = tools
	}
	if payload.Params.Temperature >= 0.0 {
		request.Temperature = anthropic.Float(float64(payload.Params.Temperature))
	}
	if payload.Params.MaxCompletionTokens > 0 {
		request.MaxTokens = int64(payload.Params.MaxCompletionTokens)
	} else {
		request.MaxTokens = 6000
	}
	request.Messages = messages
	if payload.Params.StreamOptions != nil {
		request.TopK = anthropic.Int(int64(payload.Params.TopK))
	}
	if payload.Params.TopP > 0.0 {
		request.TopP = anthropic.Float(payload.Params.TopP)
	}
	if payload.Params.ToolChoiceParams != "" {
		switch payload.Params.ToolChoiceParams {
		case aicore.AnyToolChoice.String():
			request.ToolChoice = anthropic.ToolChoiceUnionParam{
				OfToolChoiceAny: &anthropic.ToolChoiceAnyParam{},
			}
		case aicore.NoToolChoice.String():
			request.ToolChoice = anthropic.ToolChoiceUnionParam{
				OfToolChoiceNone: &anthropic.ToolChoiceNoneParam{},
			}
		case aicore.AutoToolChoice.String():
			request.ToolChoice = anthropic.ToolChoiceUnionParam{
				OfToolChoiceAuto: &anthropic.ToolChoiceAutoParam{},
			}
		default:
			request.ToolChoice = anthropic.ToolChoiceUnionParam{
				OfToolChoiceAuto: &anthropic.ToolChoiceAutoParam{},
			}
		}
	} else if payload.Params.ToolChoice != nil && payload.Params.ToolChoice.Function != nil {
		request.ToolChoice = anthropic.ToolChoiceUnionParam{
			OfToolChoiceTool: &anthropic.ToolChoiceToolParam{
				Name: payload.Params.ResponseFormat.JsonSchema.Name,
			},
		}
	}
	return request
}

func (o *AnthropicAI) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		o.Log.Err(nil).Msg("invalid model name in request")
		payload.Params.Model = defaultModel
	}
	o.Log.Info().Msgf("Generating content from anthropic with model %s", payload.Params.Model)
	request := o.buildRequest(payload)

	vbytes, _ := json.MarshalIndent(request, "", "  ")
	log.Println("request------------------->>>>>>>>>>>>>>>>>>>>>|||||||||||||||||||||||||", string(vbytes))
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	resp, err := o.Client.Messages.New(
		ctx,
		request,
	)
	if payload.StudioLog != nil {
		payload.StudioLog.TTFBAt = dmutils.GetMilliEpochTime()
	}
	if err != nil {
		o.Log.Err(err).Msg("error while generating content from openai")
		return err
	}
	for _, c := range resp.Content {
		payload.ParseOutput(&prompts.PayloadgenericResponse{
			FinishReason: string(resp.StopReason),
			Data:         c.Text,
			Role:         string(resp.Role),
			Id:           resp.ID,
			Model:        resp.Model,
			StopSequence: string(resp.StopSequence),
		})
	}
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	payload.LogUsage(&prompts.UsageMetrics{
		InputTokens:        int64(resp.Usage.InputTokens),
		OutputTokens:       int64(resp.Usage.OutputTokens),
		InputCachedTokens:  int64(resp.Usage.CacheCreationInputTokens),
		OutputCachedTokens: int64(resp.Usage.CacheReadInputTokens),
	}, nil)
	return nil
}

func (o *AnthropicAI) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		o.Log.Err(nil).Msg("invalid model name in request")
		payload.Params.Model = defaultModel
	}
	o.Log.Info().Msgf("Generating content from anthropic stream with model %s", payload.Params.Model)
	request := o.buildRequest(payload)

	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	stream := o.Client.Messages.NewStreaming(
		ctx,
		request,
	)

	if stream == nil {
		errMsg := "failed to create anthropic stream"
		o.Log.Info().Msg("Error creating anthropic stream")
		return fmt.Errorf(errMsg)
	}

	o.Log.Info().Msg("Stream created successfully, reading response")

	content := ""

	func() {
		defer stream.Close()

		for stream.Next() {
			resp := stream.Current()
			err := stream.Err()

			if err != nil {
				o.Log.Err(err).Msg("error streaming response")
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:              "",
					IsEnd:             false,
					IsError:           true,
					Created:           dmutils.GetEpochTime(),
					Id:                dmutils.AnyToStr(resp.Index),
					SystemFingerprint: resp.Delta.Signature,
				})
				break
			}

			lcContent := ""
			if len(resp.Delta.Text) > 0 {
				lcContent = resp.Delta.Text
			}

			content = content + lcContent

			err = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
				Data:              lcContent,
				IsEnd:             false,
				IsError:           false,
				Created:           dmutils.GetEpochTime(),
				Id:                dmutils.AnyToStr(resp.Index),
				SystemFingerprint: resp.Delta.Signature,
			})
			if err != nil {
				o.Log.Err(err).Msg("error sending stream data")
				break
			}

			{
				usageMetrics := &prompts.UsageMetrics{
					OutputTokens: int64(resp.Usage.OutputTokens),
				}
				payload.LogUsage(usageMetrics, nil)
			}
		}
	}()

	err := payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
		Data:    "",
		IsEnd:   true,
		IsError: false,
		Created: dmutils.GetEpochTime(),
		Id:      dmutils.GetUUID(),
	})

	if err != nil {
		o.Log.Err(err).Msg("error sending final stream message")
	}

	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}

	return nil
}
