package googlegenai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"iter"
	"log"

	"cloud.google.com/go/auth"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusai/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"google.golang.org/api/iterator"
	"google.golang.org/genai"
)

var defaultModel = "gemini-2.5-flash"
var defaultEmbeddingModel = "gemini-embedding-exp-03-07"

type GoogleGenAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type GoogleGenAI struct {
	client         *genai.Client
	log            zerolog.Logger
	modelNode      *models.AIModelNode
	maxRetries     int
	params         map[string]any
	IsVertexClient bool
}

type GoogleGenAIRequest struct {
	Model   string                       `json:"model"`
	Content []*genai.Content             `json:"parts"`
	Tools   []*genai.Tool                `json:"tools,omitempty"`
	Params  map[string]any               `json:"params,omitempty"`
	Config  *genai.GenerateContentConfig `json:"config,omitempty"`
}

func New(ctx context.Context, node *models.AIModelNode, retries int, isVertex bool, logger zerolog.Logger) (GoogleGenAIInterface, error) {
	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = "https://generativelanguage.googleapis.com"
	}
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	var genAiCl *genai.Client
	var err error
	if !isVertex {
		genAiCl, err = genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  token,
			Backend: genai.BackendGeminiAPI,
		})
	} else {
		decodeData, err := base64.StdEncoding.DecodeString(node.NetworkParams.Credentials.GcpCreds.ServiceAccountKey)
		if err != nil {
			logger.Err(err).Msg("Error decoding gcp service account key")
			return nil, err
		}
		log.Println("Decoded GCP Service Account Key: ", string(decodeData))
		genAiCl, err = genai.NewClient(ctx, &genai.ClientConfig{
			Project:  node.NetworkParams.Credentials.GcpCreds.ProjectId,
			Location: node.NetworkParams.Credentials.GcpCreds.Region,
			Backend:  genai.BackendGeminiAPI,
			Credentials: auth.NewCredentials(&auth.CredentialsOptions{
				JSON: decodeData,
			}),
		})
	}
	if err != nil || genAiCl == nil {
		logger.Error().Err(err).Msg("Error creating google gen ai client")
		return nil, err
	}
	return &GoogleGenAI{
		client:    genAiCl,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Google Gen AI"),
		modelNode: node,
	}, nil
}

func (x *GoogleGenAI) buildRequest(ctx context.Context, payload *prompts.GenerativePrompterPayload, stream bool) *GoogleGenAIRequest {
	request := &GoogleGenAIRequest{
		Config: &genai.GenerateContentConfig{
			Tools: []*genai.Tool{},
		},
	}
	log.Println("PAYLOAD for buildREQ -------------------------->>>>>>>>>>>>>>>>", payload.Context)
	if request.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		request.Model = defaultModel
	}
	// if payload.Params.MaxCompletionTokens == 0 {
	// 	payload.Params.MaxCompletionTokens = prompts.DefaultMaxOPTokenLength
	// }
	if payload.Params.Temperature > 0 {
		request.Config.Temperature = dmutils.ToPtr(float32(payload.Params.Temperature))
	}
	if payload.Params.TopP > 0 {
		request.Config.TopP = dmutils.ToPtr(float32(payload.Params.TopP))
	}
	if payload.Params.MaxCompletionTokens > 0 {
		request.Config.MaxOutputTokens = int32(payload.Params.MaxCompletionTokens)
	}
	tools := x.buildToolRequest(payload, false)
	if len(tools) > 0 {
		request.Config.Tools = append(request.Config.Tools, tools...)
	}
	req := []*genai.Part{}
	for _, msg := range payload.Params.Messages {
		switch msg.Role {
		case aicore.USER:
			req = append(req, BuildInputContent(ctx, x.client, msg)...)
		}
	}
	if len(req) > 0 {
		request.Content = []*genai.Content{
			{
				Parts: req,
			},
		}
	} else {
		x.log.Warn().Msg("No user messages found in the payload, using default message")
		defaultMsg := &pb.ChatMessageObject{
			Role:    aicore.USER,
			Content: "Hello, how can I assist you today?",
		}
		req = append(req, BuildInputContent(ctx, x.client, defaultMsg)...)
		request.Content = []*genai.Content{
			{
				Parts: req,
			},
		}
	}
	return request
}

func (x *GoogleGenAI) getSystemMessage(payload *prompts.GenerativePrompterPayload, stream bool) string {
	req := ""
	for _, msg := range payload.Params.Messages {
		if msg.Role == aicore.SYSTEM {
			req = req + msg.Content + " "
		}
	}
	return req
}

func (x *GoogleGenAI) buildToolRequest(payload *prompts.GenerativePrompterPayload, stream bool) []*genai.Tool {
	result := []*genai.Tool{}
	if len(payload.ToolCalls) > 0 {
		for _, toolCall := range payload.ToolCalls {
			tool := &genai.Tool{
				FunctionDeclarations: []*genai.FunctionDeclaration{},
			}
			log.Println("Tool Call", toolCall.FunctionSchema.Parameters)
			fMap := &map[string]any{}
			bbytes, err := dmutils.Marshall(toolCall.FunctionSchema.Parameters)
			if err != nil {
				x.log.Err(err).Msg("Error marshalling function schema")
				continue
			}
			log.Println("Tool Call Schema ++++++++++++++++++++++++++++++++++++2222222222", string(bbytes))
			err = dmutils.Unmarshall([]byte(toolCall.FunctionSchema.Parameters), fMap)
			if err != nil {
				x.log.Err(err).Msg("Error unmarshalling function schema")
				continue
			}
			hSchema, err := ConvertMapToSchema(*fMap, x.log)
			if err != nil {
				x.log.Err(err).Msg("Error marshalling function schema")
				continue
			}
			log.Println("toolCall.FunctionSchema.Name ++++++++++++++++++++++++++++++++++++1111", toolCall.FunctionSchema.Name)
			logToolCallSchema(hSchema)

			fc := &genai.FunctionDeclaration{
				Name:        toolCall.FunctionSchema.Name,
				Description: toolCall.FunctionSchema.Description,
				Parameters:  hSchema,
			}
			tool.FunctionDeclarations = append(tool.FunctionDeclarations, fc)
			result = append(result, tool)
		}
	}

	return result
}

func (x *GoogleGenAI) buildResponse(resp *genai.GenerateContentResponse, payload *prompts.GenerativePrompterPayload, parseOP bool) (string, map[string]string) {
	if resp == nil {
		return "", nil
	}
	var result string = ""
	toolCallMap := make(map[string]string)
	for _, cand := range resp.Candidates {
		log.Println("Response from Google Gen AI", cand.Content.Parts)
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				result = result + " " + fmt.Sprintf("%v", part)
			}
			// payload.ParseToolCallResponse()
		}
		if cand.Content.Parts != nil {
		funcLoop:
			for _, part := range cand.Content.Parts {
				if part.FunctionResponse == nil {
					continue funcLoop
				}
				log.Println("Function Call", part.FunctionResponse.Name)
				log.Println("Function Call Args", part.FunctionResponse.Response)
				bbytes, err := json.Marshal(part.FunctionResponse.Response)
				if err != nil {
					x.log.Err(err).Msg("Error marshalling function call args")
					continue funcLoop
				}

				toolCallMap[part.FunctionResponse.Name] = string(bbytes)
				payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
					Id:   dmutils.GetUUID(),
					Type: aicore.FUNCTION,
					FunctionSchema: &mpb.FunctionCall{
						Name:       part.FunctionResponse.Name,
						Parameters: string(bbytes),
					},
				})
			}
		}
		if parseOP {
			payload.ParseOutput(&prompts.PayloadgenericResponse{
				FinishReason: string(cand.FinishReason),
				Data:         result,
				Role:         aicore.ASSISTANT,
			})
		}
	}
	payload.ParseToolCallResponse()
	log.Println("payload.ToolCallResponse--------------------->>>>>>>>>>>++++++++++++++++++++++++", payload.ToolCallResponse)
	log.Println("Tool Call Map--------------------->>>>>>>>>>>++++++++++++++++++++++++", toolCallMap, "Result: ", result)
	return result, toolCallMap
}

func (x *GoogleGenAI) buildStreamResponse(response iter.Seq2[*genai.GenerateContentResponse, error], payload *prompts.GenerativePrompterPayload) (string, map[string]string, *prompts.UsageMetrics) {
	content := ""
	toolCallParams := make(map[string]string)
	usageMetrics := &prompts.UsageMetrics{
		TotalTokens:              0,
		InputTokens:              0,
		OutputTokens:             0,
		OutputCachedTokens:       0,
		InputCachedTokens:        0,
		InputAudioTokens:         0,
		ReasoningTokens:          0,
		OutputAudioTokens:        0,
		InputModalityMetrics:     make(map[string]*prompts.UsageModalityMetrics),
		OutputModalityMetrics:    make(map[string]*prompts.UsageModalityMetrics),
		ReasoningModalityMetrics: make(map[string]*prompts.UsageModalityMetrics),
	}
	func() {
		errCounter := 0
		for resp, err := range response {
			if err == iterator.Done {
				x.log.Info().Msg("Stream response done")
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:    "",
					IsEnd:   true,
					IsError: false,
					Created: dmutils.GetEpochTime(),
					Id:      dmutils.GetUUID(),
				})
				break
			}
			if err != nil {
				x.log.Error().Err(err).Msg("Error reading response from stream for google gen AI completion")
				_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
					Data:    "",
					IsEnd:   false,
					IsError: true,
					Created: dmutils.GetEpochTime(),
					Id:      dmutils.GetUUID(),
				})
				if errCounter > 3 {
					x.log.Err(err).Msg("error while reading stream response")
					_ = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
						Data:    "",
						IsEnd:   true,
						IsError: false,
						Created: dmutils.GetEpochTime(),
						Id:      dmutils.GetUUID(),
					})
					break
				} else {
					errCounter++
					continue
				}
			}
			result, funcMaps := x.buildResponse(resp, payload, false)
			for k, v := range funcMaps {
				val, ok := toolCallParams[k]
				if ok {
					toolCallParams[k] = val + v
				} else {
					toolCallParams[k] = v
				}
			}
			content = content + result
			//TO:DO , in future we can add more usage metrics based on the modality and other aspects as well
			usageMetrics.TotalTokens += int64(resp.UsageMetadata.TotalTokenCount)
			usageMetrics.InputTokens += int64(resp.UsageMetadata.PromptTokenCount)
			usageMetrics.OutputTokens += int64(resp.UsageMetadata.CandidatesTokenCount)
			usageMetrics.OutputCachedTokens += int64(resp.UsageMetadata.CachedContentTokenCount)
			usageMetrics.ReasoningTokens += int64(resp.UsageMetadata.ThoughtsTokenCount)
			CountTokenDetails(resp.UsageMetadata.PromptTokensDetails, usageMetrics.InputModalityMetrics)
			CountTokenDetails(resp.UsageMetadata.CandidatesTokensDetails, usageMetrics.OutputModalityMetrics)
			CountTokenDetails(resp.UsageMetadata.ToolUsePromptTokensDetails, usageMetrics.ReasoningModalityMetrics)
			err = payload.SendChatCompletionStreamData(&prompts.PayloadgenericResponse{
				Data:    result,
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
	}()
	return content, toolCallParams, usageMetrics
}

func (x *GoogleGenAI) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload) error {
	return aicore.ErrEmbeddingNotSupported
}

func (x *GoogleGenAI) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	requestObj := x.buildRequest(ctx, payload, false)
	log.Println("requestObj-----------------=========------------------->>>", requestObj.Content[0].Parts[0])
	switch payload.Mode {
	case pb.AIInterfaceMode_CHAT_MODE:
		return x.Chat(ctx, payload, requestObj)
	case pb.AIInterfaceMode_P2P:
		if payload.StudioLog != nil {
			payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
		}
		requestObj.Config.SystemInstruction = genai.Text(x.getSystemMessage(payload, false))[0]
		response, err := x.client.Models.GenerateContent(ctx, requestObj.Model, requestObj.Content, requestObj.Config)
		if payload.StudioLog != nil {
			payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
		}
		if err != nil {
			x.log.Error().Err(err).Msg("Error generating content from google gen ai")
			return err
		}
		x.buildResponse(response, payload, true)
		if response.UsageMetadata != nil {
			usageMetrics := &prompts.UsageMetrics{
				TotalTokens:              int64(response.UsageMetadata.TotalTokenCount),
				InputTokens:              int64(response.UsageMetadata.PromptTokenCount),
				OutputTokens:             int64(response.UsageMetadata.CandidatesTokenCount),
				OutputCachedTokens:       int64(response.UsageMetadata.CachedContentTokenCount),
				ReasoningTokens:          int64(response.UsageMetadata.ThoughtsTokenCount),
				InputModalityMetrics:     make(map[string]*prompts.UsageModalityMetrics),
				OutputModalityMetrics:    make(map[string]*prompts.UsageModalityMetrics),
				ReasoningModalityMetrics: make(map[string]*prompts.UsageModalityMetrics),
			}
			CountTokenDetails(response.UsageMetadata.PromptTokensDetails, usageMetrics.InputModalityMetrics)
			CountTokenDetails(response.UsageMetadata.CandidatesTokensDetails, usageMetrics.OutputModalityMetrics)
			CountTokenDetails(response.UsageMetadata.ToolUsePromptTokensDetails, usageMetrics.ReasoningModalityMetrics)

			payload.LogUsage(usageMetrics)
		}
	}
	return nil
}

func (x *GoogleGenAI) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	requestObj := x.buildRequest(ctx, payload, true)
	switch payload.Mode {
	case pb.AIInterfaceMode_CHAT_MODE:
		return x.ChatStream(ctx, payload, requestObj)
	case pb.AIInterfaceMode_P2P:
		if payload.StudioLog != nil {
			payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
		}
		requestObj.Config.SystemInstruction = genai.Text(x.getSystemMessage(payload, false))[0]
		response := x.client.Models.GenerateContentStream(ctx, requestObj.Model, requestObj.Content, requestObj.Config)
		if payload.StudioLog != nil {
			payload.StudioLog.TTFBAt = dmutils.GetMilliEpochTime()
		}
		op, toolCallsParams, usageMetrics := x.buildStreamResponse(response, payload)
		payload.ParsedOutput = op
		if payload.StudioLog != nil {
			payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
			payload.StudioLog.Output = []*models.MessageLog{
				{
					Role:    aicore.ASSISTANT,
					Content: op,
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

		payload.LogUsage(usageMetrics)
	}
	return nil
}

func (x *GoogleGenAI) Chat(ctx context.Context, payload *prompts.GenerativePrompterPayload, request *GoogleGenAIRequest) error {
	request.Config.SystemInstruction = genai.Text(x.getSystemMessage(payload, false))[0]
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	sessionContent := []*genai.Content{}
	for _, msg := range payload.SessionContext {
		sessionParts := BuildInputContent(ctx, x.client, msg)
		sessionContent = append(sessionContent, &genai.Content{
			Parts: sessionParts,
			Role:  msg.Role,
		})
	}
	session, err := x.client.Chats.Create(ctx, request.Model, request.Config, sessionContent)
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	if err != nil {
		x.log.Error().Err(err).Msg("Error creating chat session in google gen ai")
		return err
	}
	reqParts := []genai.Part{}
	for _, msg := range request.Content {
		for _, part := range msg.Parts {
			reqParts = append(reqParts, *part)
		}
	}
	if len(request.Content) < 1 {
		x.log.Warn().Msg("No content provided for chat session, using default message")
		reqParts = append(reqParts, genai.Part{
			Text: "Hello, how can I assist you today?",
		})
	}
	response, err := session.SendMessage(ctx, reqParts...)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from google gen ai")
		return err
	}
	x.buildResponse(response, payload, true)
	if response.UsageMetadata != nil {
		usageMetrics := &prompts.UsageMetrics{
			TotalTokens:              int64(response.UsageMetadata.TotalTokenCount),
			InputTokens:              int64(response.UsageMetadata.PromptTokenCount),
			OutputTokens:             int64(response.UsageMetadata.CandidatesTokenCount),
			OutputCachedTokens:       int64(response.UsageMetadata.CachedContentTokenCount),
			ReasoningTokens:          int64(response.UsageMetadata.ThoughtsTokenCount),
			InputModalityMetrics:     make(map[string]*prompts.UsageModalityMetrics),
			OutputModalityMetrics:    make(map[string]*prompts.UsageModalityMetrics),
			ReasoningModalityMetrics: make(map[string]*prompts.UsageModalityMetrics),
		}
		CountTokenDetails(response.UsageMetadata.PromptTokensDetails, usageMetrics.InputModalityMetrics)
		CountTokenDetails(response.UsageMetadata.CandidatesTokensDetails, usageMetrics.OutputModalityMetrics)
		CountTokenDetails(response.UsageMetadata.ToolUsePromptTokensDetails, usageMetrics.ReasoningModalityMetrics)

		payload.LogUsage(usageMetrics)
	}
	return nil
}

func (x *GoogleGenAI) ChatStream(ctx context.Context, payload *prompts.GenerativePrompterPayload, request *GoogleGenAIRequest) error {
	request.Config.SystemInstruction = genai.Text(x.getSystemMessage(payload, false))[0]
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	sessionContent := []*genai.Content{}
	for _, msg := range payload.SessionContext {
		sessionParts := BuildInputContent(ctx, x.client, msg)
		sessionContent = append(sessionContent, &genai.Content{
			Parts: sessionParts,
			Role:  msg.Role,
		})
	}
	session, err := x.client.Chats.Create(ctx, request.Model, request.Config, sessionContent)
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	if err != nil {
		x.log.Error().Err(err).Msg("Error creating chat session in google gen ai")
		return err
	}
	reqParts := []genai.Part{}
	for _, msg := range request.Content {
		for _, part := range msg.Parts {
			reqParts = append(reqParts, *part)
		}
	}
	if len(request.Content) < 1 {
		x.log.Warn().Msg("No content provided for chat session, using default message")
		reqParts = append(reqParts, genai.Part{
			Text: "Hello, how can I assist you today?",
		})
	}
	response := session.SendMessageStream(ctx, reqParts...)
	op, toolCallsParams, usageMetrics := x.buildStreamResponse(response, payload)
	payload.ParsedOutput = op
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
		payload.StudioLog.Output = []*models.MessageLog{
			{
				Role:    aicore.ASSISTANT,
				Content: op,
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
	fmt.Println("I am in Gemini beforeee UsageMetadata ===================>>>>")

	payload.LogUsage(usageMetrics)
	return nil
}
