package googlegenai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusdata/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusdata/core/aistudio/prompts"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var defaultModel = "gemini-2.0-flash"
var defaultEmbeddingModel = "gpt-3.5-turbo"

type GoogleGenAIInterface interface {
	GenerateEmbeddings(ctx context.Context, request *prompts.AIEmbeddingPayload) error
	GenerateContent(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	GenerateContentStream(ctx context.Context, request *prompts.GenerativePrompterPayload) error
	CrawlModels(ctx context.Context) ([]*models.AIModelBase, error)
}

type GoogleGenAI struct {
	client     *genai.Client
	log        zerolog.Logger
	modelNode  *models.AIModelNode
	maxRetries int
	params     map[string]any
}

func New(ctx context.Context, node *models.AIModelNode, retries int, logger zerolog.Logger) (GoogleGenAIInterface, error) {
	if node.NetworkParams.Url == "" {
		node.NetworkParams.Url = "https://generativelanguage.googleapis.com"
	}
	token := ""
	if node.GetCredentials("default") != nil {
		token = node.GetCredentials("default").ApiToken
	}
	genAiCl, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		logger.Error().Err(err).Msg("Error creating google gen ai client")
		return nil, err
	}
	return &GoogleGenAI{
		client:    genAiCl,
		log:       dmlogger.GetSubDMLogger(logger, "ailogger", "Google Gen AI"),
		modelNode: node,
	}, nil
}

func (x *GoogleGenAI) buildRequest(ctx context.Context, payload *prompts.GenerativePrompterPayload, stream bool) (*genai.GenerativeModel, []genai.Part) {
	if payload.Params.Model == "" {
		x.log.Warn().Msg("Model name is empty, using default model")
		payload.Params.Model = defaultModel
	}
	// if payload.Params.MaxCompletionTokens == 0 {
	// 	payload.Params.MaxCompletionTokens = prompts.DefaultMaxOPTokenLength
	// }
	modelCl := x.client.GenerativeModel(payload.Params.Model)
	if payload.Params.Temperature > 0 {
		modelCl.SetTemperature(payload.Params.Temperature)
	}
	if payload.Params.TopP > 0 {
		modelCl.SetTopP(float32(payload.Params.TopP))
	}
	if payload.Params.MaxCompletionTokens > 0 {
		modelCl.SetMaxOutputTokens(payload.Params.MaxCompletionTokens)
	}
	tools := x.buildToolRequest(payload, false)
	if len(tools) > 0 {
		modelCl.Tools = tools
	}
	req := []genai.Part{}
	for _, msg := range payload.Params.Messages {
		switch msg.Role {
		case aicore.USER:
			req = append(req, BuildInputContent(ctx, x.client, msg)...)
		}
	}
	log.Println("Request to Google Gen AI ==================>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", req)
	return modelCl, req
}

func (x *GoogleGenAI) getSystemMessage(payload *prompts.GenerativePrompterPayload, stream bool) string {
	req := ""
	for _, msg := range payload.Params.Messages {
		if msg.Role == aicore.SYSTEM {
			req = req + msg.Content + " "
		}
	}
	log.Println("Request to Google Gen AI ==================>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", req)
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
		log.Println("Response from Google Gen AI", cand)
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				result = result + " " + fmt.Sprintf("%v", part)
			}
			// payload.ParseToolCallResponse()
		}
		if cand.FunctionCalls() != nil {
		funcLoop:
			for _, funcCall := range cand.FunctionCalls() {
				log.Println("Function Call", funcCall.Name)
				log.Println("Function Call Args", funcCall.Args)
				bbytes, err := json.Marshal(funcCall.Args)
				if err != nil {
					x.log.Err(err).Msg("Error marshalling function call args")
					continue funcLoop
				}

				toolCallMap[funcCall.Name] = string(bbytes)
				payload.ToolCallResponse = append(payload.ToolCallResponse, &mpb.ToolCall{
					Id:   dmutils.GetUUID(),
					Type: aicore.FUNCTION,
					FunctionSchema: &mpb.FunctionCall{
						Name:       funcCall.Name,
						Parameters: string(bbytes),
					},
				})
			}
		}
		if parseOP {
			payload.ParseOutput(&prompts.PayloadgenericResponse{
				FinishReason: cand.FinishReason.String(),
				Data:         result,
				Role:         aicore.ASSISTANT,
			})
		}
	}
	payload.ParseToolCallResponse()
	log.Println("payload.ToolCallResponse--------------------->>>>>>>>>>>++++++++++++++++++++++++", payload.ToolCallResponse)
	return result, toolCallMap
}

func (x *GoogleGenAI) buildStreamResponse(resp *genai.GenerateContentResponseIterator, payload *prompts.GenerativePrompterPayload) (string, map[string]string) {
	content := ""
	toolCallParams := make(map[string]string)
	func() {
		errCounter := 0
		for {
			resp, err := resp.Next()
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
	return content, toolCallParams
}

func (x *GoogleGenAI) GenerateEmbeddings(ctx context.Context, payload *prompts.AIEmbeddingPayload) error {
	return aicore.ErrEmbeddingNotSupported
}

func (x *GoogleGenAI) GenerateContent(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	modelCl, parts := x.buildRequest(ctx, payload, false)
	switch payload.Mode {
	case pb.AIInterfaceMode_CHAT_MODE:
		return x.Chat(ctx, modelCl, payload, parts)
	case pb.AIInterfaceMode_P2P:
		if payload.StudioLog != nil {
			payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
		}
		modelCl.SystemInstruction = genai.NewUserContent(genai.Text(x.getSystemMessage(payload, false)))
		response, err := modelCl.GenerateContent(ctx, parts...)
		if payload.StudioLog != nil {
			payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
		}
		if err != nil {
			x.log.Error().Err(err).Msg("Error generating content from google gen ai")
			return err
		}
		x.buildResponse(response, payload, true)
		if response.UsageMetadata != nil {
			payload.LogUsage(&prompts.UsageMetrics{
				TotalTokens:        int64(response.UsageMetadata.TotalTokenCount),
				InputTokens:        int64(response.UsageMetadata.PromptTokenCount),
				OutputTokens:       int64(response.UsageMetadata.CandidatesTokenCount),
				OutputCachedTokens: int64(response.UsageMetadata.CachedContentTokenCount),
			}, nil)
		}
	}
	return nil
}

func (x *GoogleGenAI) GenerateContentStream(ctx context.Context, payload *prompts.GenerativePrompterPayload) error {
	modelCl, parts := x.buildRequest(ctx, payload, true)
	switch payload.Mode {
	case pb.AIInterfaceMode_CHAT_MODE:
		return x.ChatStream(ctx, modelCl, payload, parts)
	case pb.AIInterfaceMode_P2P:
		if payload.StudioLog != nil {
			payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
		}
		modelCl.SystemInstruction = genai.NewUserContent(genai.Text(x.getSystemMessage(payload, false)))
		response := modelCl.GenerateContentStream(ctx, parts...)
		if payload.StudioLog != nil {
			payload.StudioLog.TTFBAt = dmutils.GetMilliEpochTime()
		}
		op, toolCallsParams := x.buildStreamResponse(response, payload)
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
		if response.MergedResponse != nil && response.MergedResponse().UsageMetadata != nil {
			payload.LogUsage(&prompts.UsageMetrics{
				TotalTokens:        int64(response.MergedResponse().UsageMetadata.TotalTokenCount),
				InputTokens:        int64(response.MergedResponse().UsageMetadata.PromptTokenCount),
				OutputTokens:       int64(response.MergedResponse().UsageMetadata.CandidatesTokenCount),
				OutputCachedTokens: int64(response.MergedResponse().UsageMetadata.CachedContentTokenCount),
			}, nil)
		}
	}
	return nil
}

func (x *GoogleGenAI) Chat(ctx context.Context, modelCl *genai.GenerativeModel, payload *prompts.GenerativePrompterPayload, parts []genai.Part) error {
	modelCl.ResponseMIMEType = "text/plain"
	modelCl.SystemInstruction = genai.NewUserContent(genai.Text(x.getSystemMessage(payload, false)))
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	session := modelCl.StartChat()
	if payload.StudioLog != nil {
		payload.StudioLog.EndedAt = dmutils.GetMilliEpochTime()
	}
	for _, msg := range payload.SessionContext {
		parts := BuildInputContent(ctx, x.client, msg)
		session.History = append(session.History, &genai.Content{
			Parts: parts,
			Role:  msg.Role,
		})
	}
	response, err := session.SendMessage(ctx, parts...)
	if err != nil {
		x.log.Error().Err(err).Msg("Error generating content from google gen ai")
		return err
	}
	x.buildResponse(response, payload, true)
	if response.UsageMetadata != nil {
		payload.LogUsage(&prompts.UsageMetrics{
			TotalTokens:        int64(response.UsageMetadata.TotalTokenCount),
			InputTokens:        int64(response.UsageMetadata.PromptTokenCount),
			OutputTokens:       int64(response.UsageMetadata.CandidatesTokenCount),
			OutputCachedTokens: int64(response.UsageMetadata.CachedContentTokenCount),
		}, nil)
	}
	return nil
}

func (x *GoogleGenAI) ChatStream(ctx context.Context, modelCl *genai.GenerativeModel, payload *prompts.GenerativePrompterPayload, parts []genai.Part) error {
	modelCl.ResponseMIMEType = "text/plain"
	modelCl.SystemInstruction = genai.NewUserContent(genai.Text(x.getSystemMessage(payload, false)))
	if payload.StudioLog != nil {
		payload.StudioLog.StartedAt = dmutils.GetMilliEpochTime()
	}
	session := modelCl.StartChat()
	if payload.StudioLog != nil {
		payload.StudioLog.TTFBAt = dmutils.GetMilliEpochTime()
	}
	for _, msg := range payload.SessionContext {
		parts := BuildInputContent(ctx, x.client, msg)
		session.History = append(session.History, &genai.Content{
			Parts: parts,
			Role:  msg.Role,
		})
	}
	response := session.SendMessageStream(ctx, parts...)
	op, toolCallsParams := x.buildStreamResponse(response, payload)
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
	if response.MergedResponse() != nil && response.MergedResponse().UsageMetadata != nil {
		fmt.Println("I am in Gemini after UsageMetadata")
		fmt.Println(response.MergedResponse().UsageMetadata.CandidatesTokenCount) // why I am getting zero here????
		payload.LogUsage(&prompts.UsageMetrics{
			TotalTokens:        int64(response.MergedResponse().UsageMetadata.TotalTokenCount),
			InputTokens:        int64(response.MergedResponse().UsageMetadata.PromptTokenCount),
			OutputTokens:       int64(response.MergedResponse().UsageMetadata.CandidatesTokenCount),
			OutputCachedTokens: int64(response.MergedResponse().UsageMetadata.CachedContentTokenCount),
		}, nil)
	}
	return nil
}
