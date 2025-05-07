package aicore

import (
	"time"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
)

const (
	USER       = "user"
	SYSTEM     = "system"
	ASSISTANT  = "assistant"
	FUNCTION   = "function"
	TOOL       = "tool"
	DEVELOPER  = "developer"
	VAPUSGUARD = "vapusguard"
)

type StreamEvent string

const (
	StreamEventStart      StreamEvent = "start"
	StreamEventEnd        StreamEvent = "end"
	StreamEventData       StreamEvent = "data"
	StreamEventError      StreamEvent = "error"
	StreamStop            StreamEvent = "stop"
	StreamGuardrailFailed StreamEvent = "guardrail_failed"
)

func (s StreamEvent) String() string {
	return string(s)
}

var retryStatusCodes = map[int]bool{
	429: true,
	500: true,
	502: true,
	503: true,
	504: true,
}

var defaultRetryWaitTime = 2 * time.Second

var StartTagTemplate = `{TAG}`
var EndTagTemplate = `{/TAG}`

const (
	OpenAIOrganizationID = "organization_id"
	OpenAIProjectID      = "project_id"
)

type AIResponseFormats string

func (f AIResponseFormats) String() string {
	return string(f)
}

const (
	AIResponseFormatJSON        AIResponseFormats = "json"
	AIResponseFormatText        AIResponseFormats = "text"
	AIResponseFormatSRT         AIResponseFormats = "srt"
	AIResponseFormatVerboseJSON AIResponseFormats = "verbose_json"
	AIResponseFormatVTT         AIResponseFormats = "vtt"
	AIResponseFormatImageUrl    AIResponseFormats = "image_url"
	AIResponseFormatInputAudio  AIResponseFormats = "input_audio"
	AIResponseFormatInputFile   AIResponseFormats = "file"
	AIResponseFormatRefusal     AIResponseFormats = "refusal"
	AIResponseFormatJSONObject  AIResponseFormats = "json_object"
	AIResponseFormatJSONSchema  AIResponseFormats = "json_schema"
)

var AudioResponseMapper = map[pb.AudioResponseFormat]AIResponseFormats{
	pb.AudioResponseFormat_ARF_JSON:         AIResponseFormatJSON,
	pb.AudioResponseFormat_ARF_TEXT:         AIResponseFormatText,
	pb.AudioResponseFormat_ARF_SRT:          AIResponseFormatSRT,
	pb.AudioResponseFormat_ARF_VERBOSE_JSON: AIResponseFormatVerboseJSON,
	pb.AudioResponseFormat_ARF_VTT:          AIResponseFormatVTT,
}

type ToolChoice string

func (t ToolChoice) String() string {
	return string(t)
}

const (
	AutoToolChoice      ToolChoice = "auto"
	AnyToolChoice       ToolChoice = "any"
	NoToolChoice        ToolChoice = "none"
	ToolChoiceTool      ToolChoice = "tool"
	ToolChopiceRequired ToolChoice = "required"
)
