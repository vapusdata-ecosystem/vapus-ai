package aicore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"google.golang.org/protobuf/types/known/structpb"
)

type AIGatewayChatRequest struct {
	Messages            []*ChatMessageObject `json:"messages,omitempty" yaml:"messages"`                 // @gotags: yaml:"messages"
	Contexts            []*mpb.Mapper        `json:"contexts,omitempty" yaml:"contexts"`                 // @gotags: yaml:"contexts"
	Temperature         float32              `json:"temperature,omitempty" yaml:"temperature"`           // @gotags: yaml:"temperature"
	ChatId              string               `json:"chat_id,omitempty" yaml:"chatId"`                    // @gotags: yaml:"chatId"
	Model               string               `json:"model,omitempty" yaml:"model"`                       // @gotags: yaml:"model"
	PromptId            string               `json:"prompt_id,omitempty" yaml:"promptId"`                // @gotags: yaml:"promptId"
	MaxCompletionTokens int32                `json:"max_output_tokens,omitempty" yaml:"maxOutputTokens"` // @gotags: yaml:"maxOutputTokens"
	TopP                float64              `json:"top_p,omitempty" yaml:"topP"`                        // @gotags: yaml:"topP"
	TopK                float64              `json:"top_k,omitempty" yaml:"topK"`                        // @gotags: yaml:"topK"
	Tools               []*mpb.ToolCall      `json:"tools,omitempty" yaml:"toolCalls"`                   // @gotags: yaml:"toolCalls"
	ToolChoice          ToolChoiceField      `json:"tool_choice,omitempty" yaml:"toolChoice"`            // @gotags: yaml:"toolChoice"`
	Stream              bool                 `json:"stream,omitempty" yaml:"stream"`                     // @gotags: yaml:"stream"
	PromptInput         *structpb.Struct     `json:"prompt_input,omitempty" yaml:"promptInput"`          // @gotags: yaml:"promptInput"
	StreamOptions       *pb.StreamOptions    `json:"stream_options,omitempty" yaml:"streamOptions"`      // @gotags: yaml:"streamOptions"`
	ResponseFormat      *mpb.ResponseFormat  `json:"response_format,omitempty" yaml:"responseFormat"`    // @gotags: yaml:"responseFormat"`
}

func (x *AIGatewayChatRequest) ConvertToPb() *pb.ChatRequest {
	req := &pb.ChatRequest{
		Contexts:            x.Contexts,
		Temperature:         x.Temperature,
		ChatId:              x.ChatId,
		Model:               x.Model,
		PromptId:            x.PromptId,
		MaxCompletionTokens: x.MaxCompletionTokens,
		TopP:                x.TopP,
		TopK:                x.TopK,
		Tools:               x.Tools,
		Stream:              x.Stream,
		PromptInput:         x.PromptInput,
		Messages:            make([]*pb.ChatMessageObject, 0),
		StreamOptions:       x.StreamOptions,
		ResponseFormat:      x.ResponseFormat,
	}
	log.Println("Started converting to pb")
	for _, msg := range x.Messages {
		mess := &pb.ChatMessageObject{
			Role:      msg.Role,
			ToolCalls: msg.ToolCalls,
		}
		if msg.Content.StringValue != nil {
			val := *msg.Content.StringValue
			if len(val) > 0 {
				mess.Content = *msg.Content.StringValue
			} else {
				mess.Content = ""
				continue
			}
		} else if msg.Content.StructValue != nil {
			mess.Content = ""
			mess.StructuredContent = msg.Content.StructValue
		} else {
			continue
		}
		req.Messages = append(req.Messages, mess)
	}
	if x.ToolChoice.StringValue != nil {
		tVal := *x.ToolChoice.StringValue
		if len(tVal) > 0 {
			req.ToolChoiceParams = tVal
			req.ToolChoice = nil
		} else {
			req.ToolChoiceParams = ""
		}
	}
	if x.ToolChoice.StructValue != nil {
		req.ToolChoiceParams = ""
		req.ToolChoice = x.ToolChoice.StructValue
	} else {
		req.ToolChoiceParams = ""
		req.ToolChoice = nil
	}

	return req
}

func (x *AIGatewayChatRequest) MarshalJSON() ([]byte, error) {
	if x == nil {
		return nil, errors.New("AIGatewayChatRequest is nil")
	}
	if x.Messages == nil {
		return nil, errors.New("Messages field is nil")
	}
	return dmutils.Marshall(x)
}

type ChatMessageObject struct {
	Role      string          `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty" yaml:"role"`                                 // @gotags: yaml:"role"
	Content   ContentField    `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty" yaml:"content"`                        // @gotags: yaml:"content"
	ToolCalls []*mpb.ToolCall `protobuf:"bytes,4,rep,name=tool_calls,json=toolCalls,proto3" json:"tool_calls,omitempty" yaml:"toolCalls"` // @gotags: yaml:"toolCalls"
}

type ContentField struct {
	StringValue *string
	StructValue []*pb.RequestContentPart
}

type ToolChoiceField struct {
	StringValue *string
	StructValue *mpb.ToolChoice
}

func (cf ToolChoiceField) MarshalJSON() ([]byte, error) {
	if cf.StringValue != nil && cf.StructValue != nil {
		return nil, errors.New("ToolChoiceField cannot have both StringValue and StructValue set")
	}
	if cf.StringValue != nil {
		return dmutils.Marshall(cf.StringValue)
	}
	if cf.StructValue != nil {
		return dmutils.Marshall(cf.StructValue)
	}
	return dmutils.Marshall(nil)
}
func (cf *ToolChoiceField) UnmarshalJSON(data []byte) error {
	// 1. Try to unmarshal as a string first
	var strVal string
	// Use a decoder to check for extraneous data after the string
	decoderStr := json.NewDecoder(bytes.NewReader(data))
	decoderStr.DisallowUnknownFields() // Optional: Be strict
	errStr := decoderStr.Decode(&strVal)
	if errStr == nil && decoderStr.More() == false { // Successfully decoded a pure string
		cf.StringValue = &strVal
		cf.StructValue = nil
		log.Println("UnmarshalJSON: Detected as string")
		return nil
	}
	// 2. If string fails (or isn't just a string), try to unmarshal as the struct
	var structVal mpb.ToolChoice
	// Use a decoder to check for extraneous data after the struct
	decoderStruct := json.NewDecoder(bytes.NewReader(data))
	decoderStruct.DisallowUnknownFields() // Optional: Be strict
	errStruct := decoderStruct.Decode(&structVal)
	if errStruct == nil && decoderStruct.More() == false { // Successfully decoded the struct
		// Check if it's an empty struct that could also be null/empty string
		// This check prevents ambiguous cases like unmarshalling "null" or ""
		// into an empty struct. Add more checks if needed.
		if bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte(`""`)) {
			// It looked like the target struct but was actually null or empty string
			// Let the string path (or null path below) handle it.
			// Fall through or explicitly handle null/empty string if required.
			log.Println("UnmarshalJSON: Detected as struct, but was null/empty string")
			// Treat as string if errStr was nil earlier? Or handle as null?
			// Resetting cf just in case:
			cf.StringValue = nil
			cf.StructValue = nil
			// Re-evaluate based on errStr if you want "" to be string
			if errStr == nil {
				cf.StringValue = &strVal // It *was* a string originally
				return nil
			}
			// Otherwise, fall through to treat as null or error
		} else {
			cf.StructValue = &structVal
			cf.StringValue = nil
			log.Println("UnmarshalJSON: Detected as struct ToolChoice")
			return nil
		}
	}
	// 3. Handle JSON null explicitly
	if bytes.Equal(data, []byte("null")) {
		log.Println("UnmarshalJSON: Detected as null")
		cf.StringValue = nil
		cf.StructValue = nil
		return nil // Successfully unmarshalled null
	}
	// 4. If both fail, return an error
	// Provide a more informative error message
	log.Printf("UnmarshalJSON: Failed to unmarshal as string (err: %v) or struct (err: %v)", errStr, errStruct)
	return fmt.Errorf("data '%s' could not be unmarshaled as string or ToolChoice struct", string(data))
}

func (cf ContentField) MarshalJSON() ([]byte, error) {
	if cf.StringValue != nil && cf.StructValue != nil {
		return nil, errors.New("ContentField cannot have both StringValue and StructValue set")
	}
	if cf.StringValue != nil {
		return dmutils.Marshall(cf.StringValue)
	}
	if cf.StructValue != nil {
		return dmutils.Marshall(cf.StructValue)
	}
	return dmutils.Marshall(nil)
}

func (cf *ContentField) UnmarshalJSON(data []byte) error {
	// 1. Try to unmarshal as a string first
	var strVal string
	// Use a decoder to check for extraneous data after the string
	decoderStr := json.NewDecoder(bytes.NewReader(data))
	decoderStr.DisallowUnknownFields() // Optional: Be strict
	errStr := decoderStr.Decode(&strVal)

	if errStr == nil && decoderStr.More() == false { // Successfully decoded a pure string
		cf.StringValue = &strVal
		cf.StructValue = nil
		log.Println("UnmarshalJSON: Detected as string")
		return nil
	}

	// 2. If string fails (or isn't just a string), try to unmarshal as the struct
	var structVal []*pb.RequestContentPart
	// Use a decoder to check for extraneous data after the struct
	decoderStruct := json.NewDecoder(bytes.NewReader(data))
	decoderStruct.DisallowUnknownFields() // Optional: Be strict
	errStruct := decoderStruct.Decode(&structVal)

	if errStruct == nil && decoderStruct.More() == false { // Successfully decoded the struct
		// Check if it's an empty struct that could also be null/empty string
		// This check prevents ambiguous cases like unmarshalling "null" or ""
		// into an empty struct. Add more checks if needed.
		if bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte(`""`)) {
			// It looked like the target struct but was actually null or empty string
			// Let the string path (or null path below) handle it.
			// Fall through or explicitly handle null/empty string if required.
			log.Println("UnmarshalJSON: Detected as struct, but was null/empty string")
			// Treat as string if errStr was nil earlier? Or handle as null?
			// Resetting cf just in case:
			cf.StringValue = nil
			cf.StructValue = nil
			// Re-evaluate based on errStr if you want "" to be string
			if errStr == nil {
				cf.StringValue = &strVal // It *was* a string originally
				return nil
			}
			// Otherwise, fall through to treat as null or error

		} else {
			cf.StructValue = structVal
			cf.StringValue = nil
			log.Println("UnmarshalJSON: Detected as struct RequestContentPart")
			return nil
		}

	}

	// 3. Handle JSON null explicitly
	if bytes.Equal(data, []byte("null")) {
		log.Println("UnmarshalJSON: Detected as null")
		cf.StringValue = nil
		cf.StructValue = nil
		return nil // Successfully unmarshalled null
	}

	// 4. If both fail, return an error
	// Provide a more informative error message
	log.Printf("UnmarshalJSON: Failed to unmarshal as string (err: %v) or struct (err: %v)", errStr, errStruct)
	return fmt.Errorf("data '%s' could not be unmarshaled as string or RequestContentPart struct", string(data))
}
