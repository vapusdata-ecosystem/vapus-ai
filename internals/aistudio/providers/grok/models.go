package grok

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

const (
	defaultModel      = "mistral-large-latest"
	defaultEmbedModel = "mistral-embed"
	jsonContentType   = "application/json"
	defaultEndpoint   = "https://api.x.ai"
	CodestralEndpoint = "https://codestral.mistral.ai"
	baseAPIPath       = "/v1"
	defaultAPIVersion = "2023-06-01"
	EOS               = "\x00"
	versionHeaderKey  = "anthropic-version"
	apiKeyHeader      = "Authorization"
	modelsPath        = "/models"
	embeddingsPath    = "/embeddings"
	generatePath      = "/chat/completions"
	fimPath           = "/fim/completions"
	defaultTopP       = 1.0
)

type contentType string

const (
	messageTypeText  contentType = "text"
	messageTypeImage contentType = "image"
)

type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbeddingObject struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	ID     string            `json:"id"`
	Object string            `json:"object"`
	Data   []EmbeddingObject `json:"data"`
	Model  string            `json:"model"`
	Usage  UsageInfo         `json:"usage"`
}

type GenerativeRequest struct {
	Model            string             `json:"model,omitempty"`
	Messages         []*Message         `json:"messages,omitempty"`
	SearchParameters *SearchParameters  `json:"search_parameters"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	RandomSeed       int                `json:"random_seed,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	SafePrompt       bool               `json:"safe_prompt,omitempty"`
	Tools            []Tool             `json:"tools,omitempty"`
	ToolChoice       string             `json:"tool_choice,omitempty"`
	ResponseFormat   *APIResponseFormat `json:"response_format,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
}

type APIResponseFormat struct {
	Type ResponseFormat `json:"type,omitempty"`
}

type GenerativeResponseMessage struct {
	Index        int          `json:"index,omitempty"`
	Message      Message      `json:"message,omitempty"`
	FinishReason FinishReason `json:"finish_reason,omitempty"`
}

type SearchParameters struct {
	Mode             string    `json:"mode"`
	ReturnCitation   bool      `json:"return_citation"`
	FromDate         string    `json:"from_date"`
	ToDate           string    `json:"to_date"`
	MaxSearchResults int32     `json:"max_search_results"`
	Sources          []*Source `json:sources`
}

type Source struct {
	Type           string   `json:"type"`
	Links          []string `json:"links"`
	ExcludeWebsite []string `json:"exclude_website"`
	XHandles       []string `json:"x_handles"`
	Country        string   `json:"country"`
	SafeSearch     bool     `json:"safe_search"`
}

type GenerativeResponseMessageStream struct {
	Index        int          `json:"index,omitempty"`
	Delta        DeltaMessage `json:"delta,omitempty"`
	FinishReason FinishReason `json:"finish_reason,omitempty"`
}

type GenerativeResponse struct {
	ID      string                       `json:"id,omitempty"`
	Object  string                       `json:"object,omitempty"`
	Created int                          `json:"created,omitempty"`
	Model   string                       `json:"model,omitempty"`
	Choices []*GenerativeResponseMessage `json:"choices,omitempty"`
	Usage   UsageInfo                    `json:"usage,omitempty"`
}

type GenerativeResponseStream struct {
	ID      string                             `json:"id,omitempty"`
	Model   string                             `json:"model,omitempty"`
	Choices []*GenerativeResponseMessageStream `json:"choices,omitempty"`
	Created int                                `json:"created,omitempty"`
	Object  string                             `json:"object,omitempty"`
	Usage   UsageInfo                          `json:"usage,omitempty"`
	Error   error                              `json:"error,omitempty"`
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
}

type FinishReason string

const (
	FinishReasonStop   FinishReason = "stop"
	FinishReasonLength FinishReason = "length"
	FinishReasonError  FinishReason = "error"
)

type ResponseFormat string

const (
	ResponseFormatText       ResponseFormat = "text"
	ResponseFormatJsonObject ResponseFormat = "json_object"
)

var ResponseFormatMap = map[string]ResponseFormat{
	mpb.AIResponseFormat_TEXT.String():        ResponseFormatText,
	mpb.AIResponseFormat_JSON_SCHEMA.String(): ResponseFormatJsonObject,
}

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

func (t ToolType) String() string {
	return string(t)
}

var ToolTypeMap = map[string]ToolType{
	mpb.AIToolCallType_FUNCTION.String(): ToolTypeFunction,
}

const (
	ToolChoiceAny  = "any"
	ToolChoiceAuto = "auto"
	ToolChoiceNone = "none"
)

type Tool struct {
	Type     ToolType `json:"type,omitempty"`
	Function Function `json:"function,omitempty"`
}

type Function struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type ToolCall struct {
	Id       string       `json:"id,omitempty"`
	Type     ToolType     `json:"type,omitempty"`
	Function FunctionCall `json:"function,omitempty"`
}

type DeltaMessage struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type Message struct {
	Role      string     `json:"role,omitempty"`
	Content   *string    `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type FIMRequests struct {
	Model       string   `json:"model,omitempty"`
	Prompt      string   `json:"prompt,omitempty"`
	Suffix      string   `json:"suffix,omitempty"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
	Stop        []string `json:"stop,omitempty"`
}

type FIMResponse struct {
	ID      string               `json:"id,omitempty"`
	Object  string               `json:"object,omitempty"`
	Created int                  `json:"created,omitempty"`
	Model   string               `json:"model,omitempty"`
	Choices []FIMResponseMessage `json:"choices,omitempty"`
	Usage   UsageInfo            `json:"usage,omitempty"`
}

type FIMResponseMessage struct {
	Index        int          `json:"index,omitempty"`
	Message      Message      `json:"message,omitempty"`
	FinishReason FinishReason `json:"finish_reason,omitempty"`
}

type ModelList struct {
	Object string   `json:"object,omitempty"`
	Data   []Models `json:"data,omitempty"`
}

type Models struct {
	ID           string `json:"id,omitempty"`
	Object       string `json:"object,omitempty"`
	Created      int    `json:"created,omitempty"`
	OwnedBy      string `json:"owned_by,omitempty"`
	Capabilities struct {
		CompletionChat  bool `json:"completion_chat,omitempty"`
		CompletionFim   bool `json:"completion_fim,omitempty"`
		FunctionCalling bool `json:"function_calling,omitempty"`
		FineTuning      bool `json:"fine_tuning,omitempty"`
		Vision          bool `json:"vision,omitempty"`
	} `json:"capabilities,omitempty"`
	Name                    string `json:"name,omitempty"`
	Description             string `json:"description,omitempty"`
	MaxContextLength        int    `json:"max_context_length,omitempty"`
	Aliases                 []any  `json:"aliases,omitempty"`
	Deprecation             any    `json:"deprecation,omitempty"`
	DefaultModelTemperature any    `json:"default_model_temperature,omitempty"`
	Type                    string `json:"type,omitempty"`
}
