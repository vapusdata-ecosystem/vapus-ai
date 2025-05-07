package aicore

type GwChatDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
	Refusal string `json:"refusal,omitempty"`
}

type GwChatChoice struct {
	Delta        *GwChatDelta `json:"delta"`
	Index        int          `json:"index"`
	FinishReason *string      `json:"finish_reason,omitempty"`
}

type GwChatCompletionChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []GwChatChoice `json:"choices"`
}

type (
	AiGatewayError struct {
		Error    AiGatewayErrorDetail `json:"error"`              // The error object.
		Status   int                  `json:"status,omitempty"`   // The HTTP status code.
		Provider string               `json:"provider,omitempty"` // The provider of the error.
	}
	AiGatewayErrorDetail struct {
		Message string `json:"message"`         // A human-readable message providing more details about the error.
		Type    string `json:"type,omitempty"`  // The type of error, e.g., "invalid_request_error".
		Param   string `json:"param,omitempty"` // The parameter associated with the error (if any).
		Code    string `json:"code,omitempty"`  // An error code for programmatic handling (if any).
	}
)
