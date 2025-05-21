package generic

import "time"

// Resquest
type GuardRailRequest struct {
	Messages []*Messages `json:"messages"`
}

type Messages struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// Response
type GuardrailResponse struct {
	RequestID    string    `json:"request_id"`
	RequestTime  time.Time `json:"request_time"`
	ResponseTime time.Time `json:"response_time"`
	Status       string    `json:"status"`
	Summary      string    `json:"summary"`
	Result       *Result   `json:"result"`
}

type Result struct {
	Recipe         string          `json:"recipe"`
	Blocked        bool            `json:"blocked"`
	PromptMessages *PromptMessages `json:"prompt_messages"`
	Detectors      *Detectors      `json:"detectors"`
}

type PromptMessages []struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Detectors struct {
	PromptInjection *PromptInjection `json:"prompt_injection"`
}

type PromptInjection struct {
	Detected bool `json:"detected"`
	Data     any  `json:"data"`
}
