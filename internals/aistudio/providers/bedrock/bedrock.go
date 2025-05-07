package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	// "fmt"

	//"github.com/google/generative-ai-go/genai"

	"github.com/invopop/jsonschema"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"

	//pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"

	//"google.golang.org/api/iterator"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awscreds "github.com/aws/aws-sdk-go-v2/credentials"
	bedrockService "github.com/aws/aws-sdk-go-v2/service/bedrock"
	bedrockRuntime "github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	types "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

var defaultModel = "gemini-2.0-flash"
var defaultEmbeddingModel = "gpt-3.5-turbo"

type BedrockGenAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	ListInferenceProfiles(ctx context.Context) ([]*models.AIModelBase, error)
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
	GenerateImage(ctx context.Context, request *prompts.GenerativeImagePayload) error
}

type Bedrock struct {
	client         *bedrockRuntime.Client
	bedrockService *bedrockService.Client
	log            zerolog.Logger
	modelNode      *models.AIModelNode
	maxRetries     int
	params         map[string]any
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (BedrockGenAIInterface, error) {
	cfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion(node.NetworkParams.Credentials.AwsCreds.Region),
		awsConfig.WithCredentialsProvider(
			awscreds.NewStaticCredentialsProvider(node.GetCredentials("default").AwsCreds.AccessKeyId, node.GetCredentials("default").AwsCreds.SecretAccessKey, ""),
		),
	)
	genAiCl := bedrockRuntime.NewFromConfig(cfg)

	bedrockService := bedrockService.NewFromConfig(cfg)

	if err != nil {
		logger.Error().Err(err).Msg("Error creating aws bedrock gen ai client")
		return nil, err
	}
	return &Bedrock{
		client:         genAiCl,
		bedrockService: bedrockService,
		log:            dmlogger.GetSubDMLogger(logger, "ailogger", "Bedrock Gen AI"),
		modelNode:      node,
	}, nil
}

func (x *Bedrock) GenerateImage(ctx context.Context, input *prompts.GenerativeImagePayload) error {
	if input.ModelID == "" {
		input.ModelID = "stability.stable-diffusion-xl-v0"
	}
	if input.ImageCount == 0 {
		input.ImageCount = 1
	}
	if input.Size == "" {
		input.Size = "1024x1024"
	}
	if input.CfgScale == 0 {
		input.CfgScale = 7.0
	}

	jsonBody, err := json.Marshal(input.RequestBody)
	if err != nil {
		return err
	}

	invokeModelInput := &bedrockRuntime.InvokeModelInput{
		ModelId:     &input.ModelID,
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        jsonBody,
	}

	output, err := x.client.InvokeModel(ctx, invokeModelInput)

	if err != nil {
		x.log.Error().Err(err).Msg("Error generating image from Amazon Bedrock")
		return err
	}
	input.Response = output.Body

	return nil
}

func (x *Bedrock) buildRequest(payload *prompts.GenerativePrompterPayload, stream bool) (*bedrockRuntime.ConverseInput, *bedrockRuntime.ConverseStreamInput) {
	messages := []types.Message{}
	systemMessages := []types.SystemContentBlock{}

	for _, msg := range payload.Params.Messages {
		switch msg.Role {
		case aicore.USER:
			messages = append(messages, types.Message{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: msg.Content,
					},
				},
			})
		case aicore.ASSISTANT:
			messages = append(messages, types.Message{
				Role: types.ConversationRoleAssistant,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: msg.Content,
					},
				},
			})
		case aicore.SYSTEM:
			systemMessages = append(systemMessages, &types.SystemContentBlockMemberText{Value: msg.Content})
		}
	}
	for _, msg := range payload.SessionContext {
		switch msg.Role {
		case aicore.USER:
			messages = append(messages, types.Message{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: msg.Message,
					},
				},
			})
		case aicore.ASSISTANT:
			messages = append(messages, types.Message{
				Role: types.ConversationRoleAssistant,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: msg.Message,
					},
				},
			})
		case aicore.SYSTEM:
			systemMessages = append(systemMessages, &types.SystemContentBlockMemberText{Value: msg.Message})
		}
	}
	toolConfig := x.buildToolConfig(payload)
	topP := float32(payload.Params.TopP)
	if stream {
		converseInput := &bedrockRuntime.ConverseInput{
			ModelId: &payload.Params.Model,
			InferenceConfig: &types.InferenceConfiguration{
				Temperature: &payload.Params.Temperature,
				TopP:        &topP,
				MaxTokens:   &payload.Params.MaxCompletionTokens,
			},
			Messages: messages,
			System:   systemMessages,
		}
		if len(toolConfig.Tools) > 0 {
			converseInput.ToolConfig = toolConfig
		}
		return converseInput, nil
	} else {
		converseInput := &bedrockRuntime.ConverseInput{
			ModelId: &payload.Params.Model,
			InferenceConfig: &types.InferenceConfiguration{
				Temperature: &payload.Params.Temperature,
				TopP:        &topP,
				MaxTokens:   &payload.Params.MaxCompletionTokens,
			},
			Messages: messages,
			System:   systemMessages,
		}
		if len(toolConfig.Tools) > 0 {
			converseInput.ToolConfig = toolConfig
		}
		return converseInput, nil
	}
}

func (x *Bedrock) buildToolConfig(payload *prompts.GenerativePrompterPayload) *types.ToolConfiguration {
	if len(payload.ToolCalls) == 0 {
		return &types.ToolConfiguration{Tools: []types.Tool{}}
	}

	var tools []types.Tool
	reflector := jsonschema.Reflector{}

	for _, toolCall := range payload.ToolCalls {
		schema := reflector.Reflect(toolCall.FunctionSchema.Parameters)
		doc := document.NewLazyDocument(schema)
		toolSpec := &types.ToolMemberToolSpec{
			Value: types.ToolSpecification{
				InputSchema: &types.ToolInputSchemaMemberJson{
					Value: doc,
				},
				Description: &toolCall.FunctionSchema.Description,
				Name:        &toolCall.FunctionSchema.Name,
			},
		}
		tools = append(tools, toolSpec)
	}

	return &types.ToolConfiguration{Tools: tools}
}

func (x *Bedrock) buildResponse(resp *bedrockRuntime.ConverseOutput, payload *prompts.GenerativePrompterPayload, parseOP bool) (string, map[string]string) {
	if resp == nil {
		return "", nil
	}

	var resultBuilder strings.Builder
	toolCallMap := make(map[string]string)

	if messageOutput, ok := resp.Output.(*types.ConverseOutputMemberMessage); ok {
		// Handle text content
		for _, content := range messageOutput.Value.Content {
			if textBlock, ok := content.(*types.ContentBlockMemberText); ok {
				resultBuilder.WriteString(textBlock.Value)
				resultBuilder.WriteString(" ")
			}
		}

		// Handle tool calls
	}

	if parseOP {
		payload.ParseOutput(&prompts.PayloadgenericResponse{
			FinishReason: string(resp.StopReason),
			Data:         strings.TrimSpace(resultBuilder.String()),
			Role:         aicore.ASSISTANT,
		})
	}

	return strings.TrimSpace(resultBuilder.String()), toolCallMap
}

func (x *Bedrock) buildStreamResponse(eventStream *bedrockRuntime.ConverseStreamEventStream, payload *prompts.GenerativePrompterPayload) (string, map[string]string) {
	var fullResponse strings.Builder
	toolCallParams := make(map[string]string)

	defer eventStream.Reader.Close()

	for event := range eventStream.Reader.Events() {
		var deltaText string

		switch v := event.(type) {

		case *types.ConverseStreamOutputMemberContentBlockDelta:
			if delta, ok := v.Value.Delta.(*types.ContentBlockDeltaMemberText); ok {
				deltaText = delta.Value
				fullResponse.WriteString(deltaText)
			} else {
				x.log.Warn().Msg("Received unknown content block delta type")
				continue
			}

			err := payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
				Data:    deltaText,
				IsEnd:   false,
				IsError: false,
				Created: dmutils.GetEpochTime(),
				Id:      dmutils.GetUUID(),
				Object:  "chat.completion.chunk",
			})
			if err != nil {
				x.log.Err(err).Msg("Error while sending stream response")
			}

		case *types.ConverseStreamOutputMemberMessageStart:
			x.log.Info().Msg("Message streaming started")
		case *types.ConverseStreamOutputMemberMessageStop:
			dd := v.Value.StopReason.Values()
			for _, ll := range dd {
				reason := ""
				for _, rr := range ll.Values() {
					if rr != "" {
						reason += string(rr) + "\n"
					}
				}
				err := payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:         deltaText,
					IsEnd:        true,
					IsError:      false,
					FinishReason: reason,
					Created:      dmutils.GetEpochTime(),
					Id:           dmutils.GetUUID(),
					Object:       "chat.completion.chunk",
				})
				if err != nil {
					x.log.Err(err).Msg("Error while sending stream response")
				}
			}

			x.log.Info().Msg("Message streaming completed")

		case *types.ConverseStreamOutputMemberMetadata:
			x.log.Info().Msgf("Received metadata: %+v", v.Value)
			payload.LogUsage(&prompts.UsageMetrics{
				TotalTokens:  int64(*v.Value.Usage.InputTokens) + int64(*v.Value.Usage.OutputTokens),
				InputTokens:  int64(*v.Value.Usage.InputTokens),
				OutputTokens: int64(*v.Value.Usage.OutputTokens),
			}, nil)

		default:
			x.log.Warn().Msg("Received unknown event type in Amazon Bedrock stream")
		}
	}

	return fullResponse.String(), toolCallParams
}

func (x *Bedrock) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload) error {
	return aicore.ErrEmbeddingNotSupported
}

func (x *Bedrock) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}

	x.log.Info().Msg(fmt.Sprintf("++++++ ModelId %s", payload.Params.Model))
	converseInput, _ := x.buildRequest(payload, false)

	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}

	response, err := x.client.Converse(ctx, converseInput)
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from bedrock gen ai")
		return err
	}
	x.buildResponse(response, payload, true)
	if response.Usage != nil {
		payload.LogUsage(&prompts.UsageMetrics{
			TotalTokens:  int64(*response.Usage.TotalTokens),
			InputTokens:  int64(*response.Usage.InputTokens),
			OutputTokens: int64(*response.Usage.OutputTokens),
		}, nil)
	}

	return nil
}

func (x *Bedrock) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	// if payload.Params.MaxCompletionTokens == 0 {
	// 	payload.Params.MaxCompletionTokens = prompts.DefaultMaxOPTokenLength
	// }

	_, converseInput := x.buildRequest(payload, true)

	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}

	response, err := x.client.ConverseStream(ctx, converseInput)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from Amazon Bedrock")
		return err
	}

	if payload.StudioLog != nil {
		payload.StudioLog.TTFBAt = dmutils.GetMilliEpochTime()
	}

	op, toolCallsParams := x.buildStreamResponse(response.GetStream(), payload)
	payload.ParsedOutput = op

	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
		payload.StudioLog.Output = []*models.MessageLog{
			{Role: aicore.ASSISTANT, Content: op},
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
