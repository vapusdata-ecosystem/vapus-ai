package pangea

import "github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ai_guard"

type PangeaResponse struct {
	RequestID    string  `json:"request_id"`
	RequestTime  string  `json:"request_time"`
	ResponseTime string  `json:"response_time"`
	Status       string  `json:"status"`
	Summary      string  `json:"summary"`
	Result       *Result `json:"result"`
}

type Result struct {
	Recipe         string            `json:"recipe"`
	Blocked        bool              `json:"blocked"`
	PromptMessages []*PromptMessages `json:"prompt_messages"`
	Detectors      *Detectors        `json:"detectors"`
}

type PromptMessages struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Detectors struct {
	PromptInjection      *DetectorsDetails   `json:"prompt_injection"`
	Gibberish            *DetectorsDetails   `json:"gibberish"`
	Sentiment            *DetectorsDetails   `json:"sentiment"`
	SelfHarm             *DetectorsDetails   `json:"selfHarm"`
	PiiEntity            []*RedactEntity     `json:"pii_entity"`
	MaliciousEntity      []*MaliciousEntity  `json:"malicious_entity,omitempty"`
	CustomEntity         []*RedactEntity     `json:"custom_entity"`
	SecretsDetection     []*SecretsEntity    `json:"secrets_detection"`
	Competitors          *SingleEntityResult `json:"competitors"`
	ProfanityAndToxicity *DetectorsDetails   `json:"profanity_and_toxicity"`
	LanguageDetection    *DetectorsDetails   `json:"language_detection"`
	TopicDetection       *DetectorsDetails   `json:"topic_detection"`
	CodeDetection        *DetectorsDetails   `json:"code_detection"`
}

type DetectorsDetails struct {
	Detected bool  `json:"detected"`
	Data     *Data `json:"data"`
}

type Data struct {
	Action            string               `json:"action"`
	Language          string               `json:"language"`
	AnalyzerResponses []*AnalyzerResponses `json:"analyzer_responses"`
	Classifications   []*Classifications   `json:"classifications"`
	Topics            []*Topic             `json:"topics"`
}

type AnalyzerResponses struct {
	Analyzer   string  `json:"analyzer"`
	Confidence float64 `json:"confidence"`
}
type Classifications struct {
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
}
type Topic struct {
	Topic      string  `json:"topic`
	Confidence float64 `json:"confidence"`
}

type RedactEntity struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Action   string `json:"action"`
	StartPos *int   `json:"start_pos,omitempty"`
}

type MaliciousEntity struct {
	RedactEntity *RedactEntity          `json:"redact_entity"`
	Raw          map[string]interface{} `json:"raw,omitempty"`
}

type SecretsEntity struct {
	RedactEntity  *RedactEntity `json:"redact_entity"`
	RedactedValue string        `json:"redacted_value,omitempty"`
}

type SingleEntityResult struct {
	Action   string   `json:"action"`
	Entities []string `json:"entities"`
}

// extractorFn is a generic alias for “take a *T and return *Data”
type extractorFn[T any] func(*T) *Data

// convertDetector takes any ai_guard.TextGuardDetector[T] and an extractor,
// and populates your DetectorsDetails.
func convertDetector[T any](tg *ai_guard.TextGuardDetector[T], ex extractorFn[T]) *DetectorsDetails {
	if tg == nil {
		return nil
	}
	out := &DetectorsDetails{Detected: tg.Detected}
	if tg.Data != nil {
		out.Data = ex(tg.Data)
	}
	return out
}

// classify pulls ai_guard.ClassificationResult → your Data.Classifications
func classify(cr *ai_guard.ClassificationResult) *Data {
	if cr == nil {
		return nil
	}
	d := &Data{Action: cr.Action}
	for _, c := range cr.Classifications {
		d.Classifications = append(d.Classifications, &Classifications{
			Category:   c.Category,
			Confidence: c.Confidence,
		})
	}
	return d
}

// analyze pulls ai_guard.PromptInjectionResult → your Data.AnalyzerResponses
func analyze(pr *ai_guard.PromptInjectionResult) *Data {
	if pr == nil {
		return nil
	}
	d := &Data{Action: pr.Action}
	for _, a := range pr.AnalyzerResponses {
		d.AnalyzerResponses = append(d.AnalyzerResponses, &AnalyzerResponses{
			Analyzer:   a.Analyzer,
			Confidence: a.Confidence,
		})
	}
	return d
}

// lang pulls ai_guard.LanguageDetectionResult → your Data.Action only
func lang(ld *ai_guard.LanguageDetectionResult) *Data {
	if ld == nil {
		return nil
	}
	return &Data{Action: ld.Action}
}

func topic(tp *ai_guard.TopicDetectionResult) *Data {
	if tp == nil {
		return nil
	}
	d := &Data{Action: tp.Action}
	for _, a := range tp.Topics {
		d.Topics = append(d.Topics, &Topic{
			Topic:      a.Topic,
			Confidence: a.Confidence,
		})
	}
	return d
}

func code(cd *ai_guard.CodeDetectionResult) *Data {
	if cd == nil {
		return nil
	}
	return &Data{
		Action:   cd.Action,
		Language: cd.Language,
	}
}

func convertRedactEntities(tg *ai_guard.TextGuardDetector[ai_guard.RedactEntityResult]) []*RedactEntity {
	if tg == nil || tg.Data == nil {
		return nil
	}
	var result []*RedactEntity
	for _, e := range tg.Data.Entities {
		result = append(result, &RedactEntity{
			Type:     e.Type,
			Value:    e.Value,
			Action:   e.Action,
			StartPos: e.StartPos,
		})
	}
	return result
}

func convertPiiEntities(tg *ai_guard.TextGuardDetector[ai_guard.PiiEntityResult]) []*RedactEntity {
	if tg == nil || tg.Data == nil {
		return nil
	}
	var result []*RedactEntity
	for _, e := range tg.Data.Entities {
		result = append(result, &RedactEntity{
			Type:     e.Type,
			Value:    e.Value,
			Action:   e.Action,
			StartPos: e.StartPos,
		})
	}
	return result
}

func convertMaliciousEntities(tg *ai_guard.TextGuardDetector[ai_guard.MaliciousEntityResult]) []*MaliciousEntity {
	if tg == nil || tg.Data == nil {
		return nil
	}
	var result []*MaliciousEntity
	for _, e := range tg.Data.Entities {
		result = append(result, &MaliciousEntity{
			RedactEntity: &RedactEntity{
				Type:     e.Type,
				Value:    e.Value,
				Action:   e.Action,
				StartPos: e.StartPos,
			},
			Raw: e.Raw,
		})
	}
	return result
}

func convertSecretsEntities(tg *ai_guard.TextGuardDetector[ai_guard.SecretsEntityResult]) []*SecretsEntity {
	if tg == nil || tg.Data == nil {
		return nil
	}
	var result []*SecretsEntity
	for _, e := range tg.Data.Entities {
		result = append(result, &SecretsEntity{
			RedactEntity: &RedactEntity{
				Type:     e.Type,
				Value:    e.Value,
				Action:   e.Action,
				StartPos: e.StartPos,
			},
			RedactedValue: e.RedactedValue,
		})
	}
	return result
}

func convertSingleEntity(tg *ai_guard.TextGuardDetector[ai_guard.SingleEntityResult]) *SingleEntityResult {
	if tg == nil || tg.Data == nil {
		return nil
	}
	return &SingleEntityResult{
		Action:   tg.Data.Action,
		Entities: tg.Data.Entities,
	}
}
