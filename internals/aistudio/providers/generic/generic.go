package generic

import (
	"context"
	"encoding/json"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

var (
	OpenAIEmbeddingMap = map[string]openai.EmbeddingModel{
		"text-embedding-3-large":  openai.EmbeddingModelTextEmbedding3Large,
		"text-embedding-3-small":  openai.EmbeddingModelTextEmbedding3Small,
		"ttext-embedding-3-small": openai.EmbeddingModelTextEmbeddingAda002,
	}
)

type OpenAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
	GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error
	GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error
}

type OpenAI struct {
	Client    *openai.Client
	Log       zerolog.Logger
	ModelNode *models.AIModelNode
	Params    map[string]any
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (OpenAIInterface, error) {
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}

	clientOptions := []option.RequestOption{
		option.WithAPIKey(token),
	}

	if node.NetworkParams.Url != "" {
		clientOptions = append(clientOptions, option.WithBaseURL(node.NetworkParams.Url))
	}
	if orgID, ok := node.NetworkParams.Params[aicore.OpenAIOrganizationID].(string); ok {
		clientOptions = append(clientOptions, option.WithOrganization(orgID))
	}

	if projectID, ok := node.NetworkParams.Params[aicore.OpenAIProjectID].(string); ok {
		clientOptions = append(clientOptions, option.WithProject(projectID))
	}

	client := openai.NewClient(clientOptions...)

	return &OpenAI{
		Client:    &client,
		Log:       dmlogger.GetSubDMLogger(logger, "ailogger", "openai"),
		ModelNode: node,
	}, nil
}

func (o *OpenAI) GenerateTranscription(ctx context.Context, payload *prompts.AudioParams) error {
	return nil
}

func (o *OpenAI) GenerateTranslation(ctx context.Context, payload *prompts.AudioParams) error {
	return nil
}

func (o *OpenAI) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload) error {
	o.Log.Info().Msgf("Generating embeddings from openai with model %s", payload.EmbeddingModel)
	var rModel openai.EmbeddingModel
	rModel, ok := OpenAIEmbeddingMap[payload.EmbeddingModel]
	if !ok {
		o.Log.Err(nil).Msg("invalid model name in request")
		return aicore.ErrInvalidAIModel
	}
	o.Log.Info().Msgf("Generating embeddings from openai with model %s", payload.EmbeddingModel)
	log.Println("INPUT------------------->>>>>>>>>>>>>>>>>>>>>|||||||||||||||||||||||||", payload.Input)
	reqObj := openai.EmbeddingNewParams{
		Model:      rModel,
		User:       param.NewOpt(string(aicore.USER)),
		Dimensions: param.NewOpt(int64(payload.Dimensions)),
	}
	if len(payload.InputArray) > 0 {
		reqObj.Input = openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: payload.InputArray,
		}
	} else {
		reqObj.Input = openai.EmbeddingNewParamsInputUnion{
			OfString: param.NewOpt(payload.Input),
		}
	}
	resp, err := o.Client.Embeddings.New(
		ctx,
		reqObj,
	)
	if err != nil {
		o.Log.Err(err).Msg("error while generating embeddings from openai")
		return err
	}
	float32Slice := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		float32Slice[i] = float32(v)
	}
	payload.Embeddings = &models.VectorEmbeddings{
		Vectors32: float32Slice,
	}
	return nil
}

func (o *OpenAI) buildRequest(payload *prompts.GenerativePrompterPayload) openai.ChatCompletionNewParams {
	messages := []openai.ChatCompletionMessageParamUnion{}
	tools := []openai.ChatCompletionToolParam{}

	if payload.Mode == pb.AIInterfaceMode_CHAT_MODE {
		for _, c := range payload.SessionContext {
			if c.Role == aicore.USER {
				messages = append(messages, openai.UserMessage(c.Message))
			} else {
				messages = append(messages, openai.UserMessage(c.Message))
			}
		}
	}
	for _, c := range payload.Params.Messages {
		if c.Role == aicore.SYSTEM {
			if c.Content != "" {
				messages = append(messages, openai.SystemMessage(ConvertToCompletionSystemPart[string](c)))
				// } else {
				// 	request = append(request, openai.SystemMessage(ConvertToCompletionAssistantPart[[]openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion](c)))
			}
		}
		if c.Role == aicore.ASSISTANT {
			if c.Content != "" {
				messages = append(messages, openai.AssistantMessage(ConvertToCompletionAssistantPart[string](c)))
			} else {
				messages = append(messages, openai.AssistantMessage(ConvertToCompletionAssistantPart[[]openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion](c)))
			}
		}
		if c.Role == aicore.USER {
			if c.Content != "" {
				messages = append(messages, openai.UserMessage(ConvertToCompletionUserPart[string](c)))
			} else {
				messages = append(messages, openai.UserMessage(ConvertToCompletionUserPart[[]openai.ChatCompletionContentPartUnionParam](c)))
			}
		}
	}
	if len(payload.ToolCalls) > 0 {
		for _, tool := range payload.ToolCalls {
			if tool != nil && tool.FunctionSchema != nil {
				var paramObj shared.FunctionParameters
				paramObj = map[string]interface{}{}
				err := json.Unmarshal([]byte(tool.FunctionSchema.Parameters), &paramObj)
				if err != nil {
					o.Log.Err(err).Msg("error while unmarshalling tool call arguments")
					continue
				}
				tools = append(tools, openai.ChatCompletionToolParam{
					Type: "function",
					Function: openai.FunctionDefinitionParam{
						Name:        tool.FunctionSchema.Name,
						Description: openai.String(tool.FunctionSchema.Description),
						Parameters:  paramObj,
					},
				})
			}
		}
	}
	request := openai.ChatCompletionNewParams{
		Model: payload.Params.Model,
	}
	if payload.Params.Temperature >= 0.0 {
		request.Temperature = param.NewOpt(float64(payload.Params.Temperature))
	}
	if payload.Params.MaxCompletionTokens > 0 {
		request.MaxCompletionTokens = param.NewOpt(int64(payload.Params.MaxCompletionTokens))
	}
	if len(tools) > 0 {
		request.Tools = tools
	}
	request.Messages = messages
	if payload.Params.Stream {
		if payload.Params.StreamOptions != nil {
			request.StreamOptions = openai.ChatCompletionStreamOptionsParam{
				IncludeUsage: param.NewOpt(payload.Params.StreamOptions.IncludeUsage),
			}
		}
	}
	if payload.Params.TopP > 0.0 {
		request.TopP = param.NewOpt(float64(payload.Params.TopP))
	}
	if payload.Params.ToolChoiceParams != "" {
		request.ToolChoice = openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: param.NewOpt(payload.Params.ToolChoiceParams),
		}
	} else if payload.Params.ToolChoice != nil && payload.Params.ToolChoice.Function != nil {
		request.ToolChoice = openai.ChatCompletionToolChoiceOptionUnionParam{
			OfChatCompletionNamedToolChoice: &openai.ChatCompletionNamedToolChoiceParam{
				Type: aicore.FUNCTION,
				Function: openai.ChatCompletionNamedToolChoiceFunctionParam{
					Name: payload.Params.ToolChoice.Function.GetName(),
				},
			},
		}
	}

	if payload.Params.ResponseFormat != nil {
		switch payload.Params.ResponseFormat.Type {
		case aicore.AIResponseFormatJSONObject.String():
			request.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{},
			}
		case aicore.AIResponseFormatText.String():
			request.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfText: &shared.ResponseFormatTextParam{},
			}
		case aicore.AIResponseFormatJSONSchema.String():
			if payload.Params.ResponseFormat.JsonSchema != nil {
				request.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
					OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
						JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
							Name:   payload.Params.ResponseFormat.JsonSchema.Name,
							Strict: param.NewOpt(payload.Params.ResponseFormat.JsonSchema.Strict),
							Schema: payload.Params.ResponseFormat.JsonSchema.Schema,
						},
					},
				}
			}
		default:
			request.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{},
			}
		}
	}
	// // To get the used token
	// // if payload.Params.StreamOptions.IncludeUsage {
	// request.StreamOptions.IncludeUsage = param.Opt[bool]{
	// 	Value: true,
	// }
	// // }
	return request
}

func (o *OpenAI) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		o.Log.Err(nil).Msg("invalid model name in request")
		payload.Params.Model = openai.ChatModelGPT4o
	}
	o.Log.Info().Msgf("Generating content from openai with model %s", payload.Params.Model)
	request := o.buildRequest(payload)
	vbytes, _ := json.MarshalIndent(request, "", "  ")
	log.Println("request------------------->>>>>>>>>>>>>>>>>>>>>|||||||||||||||||||||||||", string(vbytes))
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	resp, err := o.Client.Chat.Completions.New(
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
	for _, c := range resp.Choices {
		if len(c.Message.ToolCalls) > 0 {
			for _, t := range c.Message.ToolCalls {
				payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
					Id:   t.ID,
					Type: string(t.Type),
					FunctionSchema: &mpb.FunctionCall{
						Name:       t.Function.Name,
						Parameters: t.Function.Arguments,
					},
				})
			}
			payload.ParseToolCallResponse()
		}
		if len(c.Message.Content) > 0 {
			payload.ParseOutput(&prompts.PayloadgenericResponse{
				FinishReason: string(resp.Choices[0].FinishReason),
				Data:         c.Message.Content,
				Role:         string(c.Message.Role),
			})
		}

	}
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	payload.Response.Object = string(resp.Object)
	payload.Response.Created = resp.Created
	payload.Response.Id = resp.ID
	payload.Response.Model = payload.Params.Model

	payload.LogUsage(&prompts.UsageMetrics{
		TotalTokens:              int64(resp.Usage.TotalTokens),
		InputTokens:              int64(resp.Usage.PromptTokens),
		OutputTokens:             int64(resp.Usage.CompletionTokens),
		InputCachedTokens:        int64(resp.Usage.PromptTokensDetails.CachedTokens),
		InputAudioTokens:         int64(resp.Usage.PromptTokensDetails.AudioTokens),
		OutputAudioTokens:        int64(resp.Usage.CompletionTokensDetails.AudioTokens),
		ReasoningTokens:          int64(resp.Usage.CompletionTokensDetails.ReasoningTokens),
		RejectedPredictionTokens: int64(resp.Usage.CompletionTokensDetails.RejectedPredictionTokens),
		AcceptedPredictionTokens: int64(resp.Usage.CompletionTokensDetails.AcceptedPredictionTokens),
	}, nil)
	return nil
}

func (o *OpenAI) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		o.Log.Err(nil).Msg("invalid model name in request")
		payload.Params.Model = openai.ChatModelGPT4o
	}
	o.Log.Info().Msgf("Generating content from openai stream with model %s", payload.Params.Model)
	request := o.buildRequest(payload)

	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	stream := o.Client.Chat.Completions.NewStreaming(
		ctx,
		request,
	)

	content := ""
	toolCallsParams := map[string]string{}
	resp := openai.ChatCompletionChunk{}
	func() {
		defer stream.Close()
		o.Log.Info().Msg("Stream created successfully, reading response")
		log.Println("Stream created successfully, reading response", stream)
		for stream.Next() {
			resp = stream.Current()
			err := stream.Err()
			// if err != nil {
			// 	o.Log.Err(err).Msg("error streaming response")
			// }
			// if len(resp.Choices) == 0 {
			// 	continue
			// }
			if err != nil {
				o.Log.Err(err).Msg("error while reading stream")
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:              "",
					IsEnd:             false,
					IsError:           true,
					Created:           resp.Created,
					Id:                resp.ID,
					Object:            string(resp.Object),
					SystemFingerprint: resp.SystemFingerprint,
					Model:             payload.Params.Model,
				})
				break
			}
			if len(resp.Choices) == 0 {
				continue
			}
			lcContent := ""
			for _, c := range resp.Choices {
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
			content = content + lcContent
			err = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
				Data:              lcContent,
				IsEnd:             false,
				IsError:           false,
				Created:           resp.Created,
				Id:                resp.ID,
				Object:            string(resp.Object),
				SystemFingerprint: resp.SystemFingerprint,
				Model:             payload.Params.Model,
			})
			if err != nil {
				break
			}
		}
	}()

	payload.LogUsage(&prompts.UsageMetrics{
		TotalTokens:              int64(resp.Usage.TotalTokens),
		InputTokens:              int64(resp.Usage.PromptTokens),
		OutputTokens:             int64(resp.Usage.CompletionTokens),
		InputCachedTokens:        int64(resp.Usage.PromptTokensDetails.CachedTokens),
		InputAudioTokens:         int64(resp.Usage.PromptTokensDetails.AudioTokens),
		OutputAudioTokens:        int64(resp.Usage.CompletionTokensDetails.AudioTokens),
		ReasoningTokens:          int64(resp.Usage.CompletionTokensDetails.ReasoningTokens),
		RejectedPredictionTokens: int64(resp.Usage.CompletionTokensDetails.RejectedPredictionTokens),
		AcceptedPredictionTokens: int64(resp.Usage.CompletionTokensDetails.AcceptedPredictionTokens),
	}, nil)

	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()

		// To save the Output
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
	return nil
}
