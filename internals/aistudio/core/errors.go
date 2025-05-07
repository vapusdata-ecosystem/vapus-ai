package aicore

import (
	"errors"
)

var (
	ErrUnknownServiceProvider = errors.New("unknown service provider for LLM models")
	Err404AIResult            = errors.New("result not found in the response body of the AI service")
	ErrInvalidAIModel         = errors.New("invalid AI model requested")
	ErrInvalidAIModelRequest  = errors.New("invalid AI model request")
	ErrEmbeddingNotSupported  = errors.New("embedding generation not supported for the AI model")
	ErrNoResponseFromAIModel  = errors.New("no response from the AI model")
	ErrInvalidAIModelRole     = errors.New("invalid role, Supported values are: 'system', 'assistant', 'user', 'function', 'tool', and 'developer'")
)
