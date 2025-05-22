package mistral

const (
	defaultEndpoint = "https://api.mistral.ai"
	baseAPIPath     = "/v1"
	generatePath    = "/chat/moderations"
	jsonContentType = "application/json"
)

type MistralGuardrailResponse struct {
	ID      string     `json:"id"`
	Usage   *Usage     `json:"usage"`
	Model   string     `json:"model"`
	Results []*Results `json:"results"`
}

type Usage struct {
	PromptTokens       int `json:"prompt_tokens"`
	TotalTokens        int `json:"total_tokens"`
	CompletionTokens   int `json:"completion_tokens"`
	RequestCount       int `json:"request_count"`
	PromptTokenDetails any `json:"prompt_token_details"`
}

type Results struct {
	CategoryScores *CategoryScores `json:"category_scores"`
	Categories     *Categories     `json:"categories"`
}

type CategoryScores struct {
	Sexual                      float64 `json:"sexual"`
	HateAndDiscrimination       float64 `json:"hate_and_discrimination"`
	ViolenceAndThreats          float64 `json:"violence_and_threats"`
	DangerousAndCriminalContent float64 `json:"dangerous_and_criminal_content"`
	Selfharm                    float64 `json:"selfharm"`
	Health                      float64 `json:"health"`
	Financial                   float64 `json:"financial"`
	Law                         float64 `json:"law"`
	Pii                         float64 `json:"pii"`
}

type Categories struct {
	Sexual                      bool `json:"sexual"`
	HateAndDiscrimination       bool `json:"hate_and_discrimination"`
	ViolenceAndThreats          bool `json:"violence_and_threats"`
	DangerousAndCriminalContent bool `json:"dangerous_and_criminal_content"`
	Selfharm                    bool `json:"selfharm"`
	Health                      bool `json:"health"`
	Financial                   bool `json:"financial"`
	Law                         bool `json:"law"`
	Pii                         bool `json:"pii"`
}
