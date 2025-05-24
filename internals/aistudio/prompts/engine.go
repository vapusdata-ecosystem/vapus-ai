package prompts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	jinjarndr "github.com/vapusdata-ecosystem/vapusai/core/tools/jinja"
)

// var Baseregex = `{%s}\s*(.*?)\s*{/%s}`
var Baseregex = `(?s)\{%s\}(.*?)\{/%s\}`

type AudioParams struct {
	File                   io.Reader
	Model                  string
	Language               string
	Prompt                 string
	ResponseFormat         string
	Temperature            float64
	TimestampGranularities []string
	ResponseText           string
	ResponseJson           string
}

type AIEmbeddingPayload struct {
	Input          string
	Dimensions     int
	EmbeddingModel string
	Embeddings     *models.VectorEmbeddings
	SystemMessage  string
	UserMessage    string
	InputArray     []string
}

type GenerativeImagePayload struct {
	Prompt      string
	ModelID     string
	ImageCount  int32
	Size        string
	Seed        int64
	CfgScale    float32
	Images      []string
	RequestBody map[string]any
	Response    []byte
}

type GenerativePrompterPayload struct {
	Params           *pb.ChatRequest
	Prompt           *models.AIPrompt
	Response         *pb.ChatResponse
	Error            error
	ParsedOutput     string
	Context          []*mpb.Mapper
	ResultMetadata   map[string]any
	Suffix           string
	isRendered       bool
	SessionID        string
	SessionContext   []*SessionMessage
	SummaryOutput    string
	ToolCallResponse []*mpb.ToolCall
	opStart          string
	opEnd            string
	Stream           pb.AIStudio_ChatServer
	AgentStream      pb.AIStudio_BidiChatServer
	SSEChan          chan *aicore.GwChatCompletionChunk
	streamLog        string
	opStartReacted   bool
	opEndReacted     bool
	GuardrailsFailed bool
	GuardrailsResult map[string][]string
	ToolCalls        []*mpb.ToolCall
	logger           zerolog.Logger
	Usage            *models.AIStudioUsages
	StudioLog        *models.AIStudioLog
	Mode             pb.AIInterfaceMode
}

type UsageMetrics struct {
	InputTokens              int64 `json:"inputTokens,omitempty" yaml:"inputTokens,omitempty"`
	InputCachedTokens        int64 `json:"inputCachedTokens,omitempty" yaml:"inputCachedTokens,omitempty"`
	OutputTokens             int64 `json:"outputTokens,omitempty" yaml:"outputTokens,omitempty"`
	OutputCachedTokens       int64 `json:"outputCachedTokens,omitempty" yaml:"outputCachedTokens,omitempty"`
	InputAudioTokens         int64 `json:"inputAudioTokens,omitempty" yaml:"inputAudioTokens,omitempty"`
	OutputAudioTokens        int64 `json:"outputAudioTokens,omitempty" yaml:"outputAudioTokens,omitempty"`
	TotalTokens              int64 `json:"totalTokens,omitempty" yaml:"totalTokens,omitempty"`
	ReasoningTokens          int64 `json:"reasoningTokens,omitempty" yaml:"reasoningTokens,omitempty"`
	RejectedPredictionTokens int64 `json:"rejectedPredictionTokens,omitempty" yaml:"rejectedPredictionTokens,omitempty"`
	AcceptedPredictionTokens int64 `json:"acceptedPredictionTokens,omitempty" yaml:"acceptedPredictionTokens,omitempty"`
}

type PayloadgenericResponse struct {
	Data              string
	Role              string
	FinishReason      string
	Object            string
	Citations         []string
	ToolCall          *mpb.ToolCall
	Id                string
	Created           int64
	IsEnd             bool
	Event             string
	IsError           bool
	SystemFingerprint string
	Model             string
	StopSequence      string
}

type SessionMessage struct {
	Message           string
	Role              string
	StructuredMessage []*pb.RequestContentPart
	SearchParameters  *pb.SearchParameters
}

func NewPrompter(params *pb.ChatRequest, prompt *models.AIPrompt, stream pb.AIStudio_ChatServer, agentStream pb.AIStudio_BidiChatServer, logger zerolog.Logger) *GenerativePrompterPayload {
	return &GenerativePrompterPayload{
		Params:           params,
		Prompt:           prompt,
		ResultMetadata:   map[string]any{},
		Stream:           stream,
		GuardrailsFailed: false,
		GuardrailsResult: map[string][]string{},
		ToolCalls:        []*mpb.ToolCall{},
	}
}

func (p *GenerativePrompterPayload) LogUsage(opts *UsageMetrics, err error) {
	if p.Usage == nil {
		p.Usage = &models.AIStudioUsages{}
	}
	p.Usage.InputTokens += opts.InputTokens
	p.Usage.InputCachedTokens += opts.InputCachedTokens
	p.Usage.OutputTokens += opts.OutputTokens
	p.Usage.OutputCachedTokens += opts.OutputCachedTokens
	p.Usage.InputAudioTokens += opts.InputAudioTokens
	p.Usage.OutputAudioTokens += opts.OutputAudioTokens
	p.Usage.TotalTokens += opts.TotalTokens
	if p.SSEChan == nil {
		p.Response.Usage = &pb.Usage{
			TotalTokens:  opts.TotalTokens,
			PromptTokens: opts.InputTokens,
			PromptTokensDetails: &pb.PromptTokensDetails{
				CachedTokens: opts.InputCachedTokens,
				AudioTokens:  opts.InputAudioTokens,
			},
			CompletionTokens: opts.OutputTokens,
			CompletionTokensDetails: &pb.CompletionTokensDetails{
				AudioTokens:              opts.OutputAudioTokens,
				ReasoningTokens:          opts.ReasoningTokens,
				AcceptedPredictionTokens: opts.AcceptedPredictionTokens,
				RejectedPredictionTokens: opts.RejectedPredictionTokens,
			},
		}
	}
}

func (p *GenerativePrompterPayload) ParseOutput(opts *PayloadgenericResponse) {
	if p.StudioLog != nil {
		p.StudioLog.Output = append(p.StudioLog.Output, &models.MessageLog{
			Content: opts.Data,
			Role:    aicore.ASSISTANT,
		})
	}
	outPut := opts.Data
	if p.Prompt != nil {
		// outPut = strings.Replace(outPut, "\n", "", -1)
		if p.Prompt.Spec.OutputTag == "" {
			p.ParsedOutput = outPut
		} else {
			outPut, _ = dmutils.ExtractBetweenDelimiters(outPut, p.opStart, p.opEnd)
			if outPut == "" {
				outPut = opts.Data
			}
			p.ParsedOutput = outPut
			if p.StudioLog != nil {
				p.StudioLog.ParsedOutput = append(p.StudioLog.ParsedOutput, &models.MessageLog{
					Content: p.ParsedOutput,
					Role:    aicore.ASSISTANT,
				})
			}
		}
	} else {
		p.ParsedOutput = outPut
	}
	opts.Data = p.ParsedOutput
	p.BuildResponseOP(aicore.StreamEventData.String(), opts, false)
	return
}

func (p *GenerativePrompterPayload) ParseToolCallResponse() {
	if len(p.ToolCallResponse) > 0 {
		log.Println("Tool call response", p.ToolCallResponse)
		p.Response.Choices = append(p.Response.Choices, &pb.ChatResponseChoice{
			Messages: &pb.ChatMessageObject{
				Role:      aicore.ASSISTANT,
				ToolCalls: p.ToolCallResponse,
			},
		})
	}
	if p.StudioLog != nil {
		p.StudioLog.ToolCallResponse = p.ToolCallResponse
	}
}

func (p *GenerativePrompterPayload) RenderPrompt() error {
	if p.isRendered {
		return nil
	}
	contextsMap := make(map[string]any)
	contexts := ""
	if len(p.Context) > 0 {
		for _, context := range p.Context {
			contextsMap[context.Key] = context.Value
		}
		contextBytes, err := json.Marshal(contextsMap)
		if err != nil {
			p.logger.Err(err).Msg("error while marshalling context")
			contexts = ""
		} else {
			contexts = string(contextBytes)
		}
	}
	if p.Prompt != nil {
		// if p.Params.MaxCompletionTokens == 0 {
		// 	p.Params.MaxCompletionTokens = DefaultMaxOPTokenLength
		// }
		p.isRendered = true
		p.opStart = fmt.Sprintf("{%s}", p.Prompt.Spec.OutputTag)
		p.opEnd = fmt.Sprintf("{/%s}", p.Prompt.Spec.OutputTag)
		p.opEndReacted = false
		p.opStartReacted = false
		systemMessRendered := false

		for _, mess := range p.Params.Messages {
			if mess.Role == aicore.SYSTEM {
				mess.Content = p.Prompt.Spec.SystemMessage + "\n" + mess.Content
				systemMessRendered = true
			}
			if mess.Role == aicore.USER {
				mes := mess.Content
				if p.Prompt.UserTemplate != "" {
					if strings.Contains(p.Prompt.UserTemplate, "["+p.Prompt.Spec.InputTag+"]") {
						mess.Content = strings.Replace(p.Prompt.UserTemplate, "["+p.Prompt.Spec.InputTag+"]", mes, -1)
					}
				}
				if len(contexts) > 0 && contexts != "{}" && contexts != "" {
					mess.Content = strings.Replace(mess.Content, "["+p.Prompt.Spec.ContextTag+"]", contexts, -1)
				}
			}

		}
		if len(p.Prompt.Spec.Variables) > 0 {
			var uMess string
			var err error
			if p.Params.PromptInput != nil {
				varMap := p.Params.PromptInput.AsMap()
				uMess, err = jinjarndr.RenderJinjaTemplate(p.Prompt.Spec.UserMessage, varMap)
				if err != nil {
					p.logger.Err(err).Msg("error while rendering user message template from prompt")
					return err
				}
			} else {
				uMess = p.Prompt.Spec.UserMessage
			}
			p.Params.Messages = append(p.Params.Messages, &pb.ChatMessageObject{
				Role:    aicore.USER,
				Content: uMess,
			})
		}
		if len(p.Prompt.Spec.Tools) > 0 {
			for _, tool := range p.Prompt.Spec.Tools {
				argBytes, err := json.MarshalIndent(tool.Schema.Parameters, "", "  ")
				if err != nil {
					p.logger.Err(err).Msg("error while marshalling tool arguments")
					return err
				}
				p.ToolCalls = append(p.ToolCalls, &mpb.ToolCall{
					Type: tool.Type,
					FunctionSchema: &mpb.FunctionCall{
						Name:        tool.Schema.Name,
						Description: tool.Schema.Description,
						Parameters:  string(argBytes),
					},
				})
			}
		}
		if !systemMessRendered {
			p.Params.Messages = append(p.Params.Messages, &pb.ChatMessageObject{
				Role:    aicore.SYSTEM,
				Content: p.Prompt.Spec.SystemMessage,
			})
		}
	} else {
		if len(contexts) > 0 && contexts != "{}" && contexts != "" {
			p.Params.Messages = append(p.Params.Messages, &pb.ChatMessageObject{
				Role:    aicore.USER,
				Content: "Context: " + contexts,
			})
		}
	}
	if p.StudioLog != nil {
		for _, mess := range p.Params.Messages {
			if mess == nil {
				continue
			}
			log.Println("Message", mess)
			if mess.StructuredContent != nil {
				p.StudioLog.Input = append(p.StudioLog.Input, &models.MessageLog{
					StructuredContent: mess.StructuredContent,
					Role:              mess.Role,
				})
			} else if mess.Content != "" {
				p.StudioLog.Input = append(p.StudioLog.Input, &models.MessageLog{
					Content: mess.Content,
					Role:    mess.Role,
				})
			}
		}
		if len(p.ToolCalls) > 0 {
			p.StudioLog.ToolCallSchema = p.ToolCalls
		}
	}
	// if p.Params.MaxCompletionTokens == 0 {
	// 	p.Params.MaxCompletionTokens = DefaultMaxOPTokenLength
	// }
	if len(p.Params.Tools) > 0 {
		p.ToolCalls = append(p.ToolCalls, p.Params.Tools...)
	}
	p.Response = &pb.ChatResponse{}
	return nil

}

func (p *GenerativePrompterPayload) GetUserMessage() string {
	return p.Prompt.Spec.UserMessage
}
func (p *GenerativePrompterPayload) FilterTag(text string) string {
	if !p.opStartReacted {
		if strings.Contains(text, p.opStart) {
			p.opStartReacted = true
			return strings.ReplaceAll(text, p.opStart, "")
		} else {
			return text
		}
	} else if !p.opEndReacted {
		if strings.Contains(text, p.opEnd) {
			p.opEndReacted = true
			return strings.ReplaceAll(text, p.opEnd, "")
		} else {
			return text
		}
	} else {
		return text
	}
}

func (p *GenerativePrompterPayload) SendChatCompletionStreamData(opts *PayloadgenericResponse) error {
	var err error
	p.streamLog = p.FilterTag(p.streamLog)
	p.streamLog = p.streamLog + opts.Data
	p.ParsedOutput += opts.Data
	if len(p.streamLog) < 25 {
		if opts.IsEnd {
			p.EndStream(aicore.StreamEventEnd.String(), opts)
		} else if opts.IsError {
			p.EndStream(aicore.StreamEventError.String(), opts)
		} else {
			return nil
		}
	} else {
		if opts.IsEnd {
			p.EndStream(aicore.StreamEventEnd.String(), opts)
		} else if opts.IsError {
			p.EndStream(aicore.StreamEventError.String(), opts)
		} else {
			opts.Data = p.streamLog
			opts.Role = aicore.ASSISTANT
			err = p.HandleStream(p.BuildStreamResponseOP(aicore.StreamEventData.String(), opts))
			p.streamLog = ""
			return err
		}
		return err
	}
	return err
}

func (p *GenerativePrompterPayload) EndStream(event string, opts *PayloadgenericResponse) error {
	if len(p.streamLog) > 0 {
		opts.Data = p.streamLog
		err := p.HandleStream(p.BuildStreamResponseOP(aicore.StreamEventData.String(), opts))
		if err != nil {
			return err
		}
		p.streamLog = ""
	}
	return p.HandleStream(&pb.ChatStreamResponse{
		Event:   event,
		Model:   p.Params.Model,
		Id:      dmutils.GetUUID(),
		Created: dmutils.GetEpochTime(),
	})
}

func (p *GenerativePrompterPayload) BuildResponseOP(event string, opts *PayloadgenericResponse, isStream bool) *pb.ChatResponse {
	if p.Response == nil {
		p.Response = &pb.ChatResponse{
			Choices: make([]*pb.ChatResponseChoice, 0),
			Model:   p.Params.Model,
			Created: dmutils.GetEpochTime(),
		}
	}
	ch := &pb.ChatResponseChoice{
		Messages: &pb.ChatMessageObject{
			Role:    opts.Role,
			Content: opts.Data,
		},
		FinishReason: opts.FinishReason,
	}
	p.Response.Choices = append(p.Response.Choices, ch)
	p.Response.Id = opts.Id
	p.Response.Created = opts.Created
	p.Response.Object = opts.Object
	return p.Response
}

func (p *GenerativePrompterPayload) BuildStreamResponseOP(event string, opts *PayloadgenericResponse) *pb.ChatStreamResponse {
	return &pb.ChatStreamResponse{
		Choices: []*pb.ChatResponseStreamChoice{
			{
				Messages: &pb.ChatMessageObject{
					Role:    opts.Role,
					Content: opts.Data,
				},
				FinishReason: opts.FinishReason,
			},
		},
		Model:   p.Params.Model,
		Created: opts.Created,
		Event:   event,
		Id:      opts.Id,
		Object:  opts.Object,
	}
}

func (x *GenerativePrompterPayload) HandleStream(obj *pb.ChatStreamResponse) error {
	if x.Stream != nil {
		select {
		case <-x.Stream.Context().Done():
			return x.Stream.Context().Err()
		default:
			if err := x.Stream.Send(obj); err != nil {
				return err
			}
		}
	} else if x.AgentStream != nil {
		select {
		case <-x.AgentStream.Context().Done():
			return x.AgentStream.Context().Err()
		default:
			if err := x.AgentStream.Send(obj); err != nil {
				return err
			}
		}
	} else if x.SSEChan != nil {
		c := 0
		ev := &aicore.GwChatCompletionChunk{
			Object:  obj.Object,
			Created: obj.Created,
			ID:      obj.Id,
			Model:   x.Params.Model,
			Choices: make([]aicore.GwChatChoice, 0),
		}
		for _, tool := range obj.Choices {
			ev.Choices = append(ev.Choices, aicore.GwChatChoice{
				Delta: &aicore.GwChatDelta{
					Content: tool.Messages.Content,
					Role:    tool.Messages.Role,
				},
				Index:        0,
				FinishReason: &tool.FinishReason,
			})
		}

		for {
			select {
			case x.SSEChan <- ev:
				return nil
			case <-time.After(50 * time.Millisecond):
				x.logger.Info().Msg("waiting for channel to be ready")
			}
			c++
		}
	}
	return nil
}
