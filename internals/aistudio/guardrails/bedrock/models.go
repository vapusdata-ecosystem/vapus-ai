package bedrock

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/aws/smithy-go/middleware"
)

type BedrockGuardrailList struct {
	ARN  string
	Id   string
	Name string
}

type ApplyGuardrailOutput struct {
	Action            string                `json:"action" bson:"action"`
	Assessments       []GuardrailAssessment `json:"assessments" bson:"assessments"`
	Outputs           []string              `json:"outputs" bson:"outputs"`
	Usage             *GuardrailUsage       `json:"usage,omitempty" bson:"usage,omitempty"`
	GuardrailCoverage *GuardrailCoverage    `json:"guardrail_coverage,omitempty" bson:"guardrail_coverage,omitempty"`
	ResultMetadata    string                `json:"result_metadata" bson:"result_metadata"`
}

type GuardrailAssessment struct {
	ContentPolicy              *GuardrailContentFilter              `json:"content_policy" bson:"content_policy"`
	ContextualGroundingPolicy  string                               `json:"contextual_grounding_policy" bson:"contextual_grounding_policy"`
	InvocationMetrics          *GuardrailInvocationMetrics          `json:"invocation_metrics" bson:"invocation_metrics"`
	SensitiveInformationPolicy *GuardrailSensitiveInformationPolicy `json:"sensitive_information_policy" bson:"sensitive_information_policy"`
	TopicPolicy                []*GuardrailNamedTypedAction         `json:"topic_policy" bson:"topic_policy"`
	WordPolicy                 *GuardrailWordPolicyAssessment       `json:"word_policy" bson:"word_policy"`
}

type GuardrailInvocationMetrics struct {
	GuardrailCoverage          *GuardrailCoverage `json:"guardrail_coverage" bson:"guardrail_coverage"`
	GuardrailProcessingLatency *int64             `json:"guardrail_processing_latency" bson:"guardrail_processing_latency"`
	Usage                      *GuardrailUsage    `json:"usage" bson:"usage"`
}

type GuardrailSensitiveInformationPolicy struct {
	PiiEntities []*GuardrailNamedTypedAction `json:"pii_entities" bson:"pii_entities"` // Reused
	Regexes     []*GuardrailRegexFilter      `json:"regexes" bson:"regexes"`
}

type GuardrailRegexFilter struct {
	CommonActionMatch
	Name  string `json:"name" bson:"name"`
	Regex string `json:"regex" bson:"regex"`
}

type GuardrailWordPolicyAssessment struct {
	CustomWords      []*GuardrailActionMatch      `json:"custom_words" bson:"custom_words"`
	ManagedWordLists []*GuardrailNamedTypedAction `json:"managed_word_lists" bson:"managed_word_lists"`
}

type GuardrailContentFilter struct {
	Filters []GuardrailFilter `json:"filters" bson:"filters"`
}

type GuardrailFilter struct {
	CommonActionTyped
	Confidence     string `json:"confidence" bson:"confidence"`
	FilterStrength string `json:"filter_strength" bson:"filter_strength"`
}

type CommonActionMatch struct {
	Action string `json:"action" bson:"action"`
	Match  string `json:"match" bson:"match"`
}

type CommonActionTyped struct {
	Action string `json:"action" bson:"action"`
	Type   string `json:"type" bson:"type"`
}

type GuardrailActionMatch struct {
	CommonActionMatch
}

type GuardrailNamedTypedAction struct {
	CommonActionMatch
	Type string `json:"type" bson:"type"`
}

type GuardrailUsage struct {
	ContentPolicyUnits                  int `json:"content_policy_units" bson:"content_policy_units"`
	ContextualGroundingPolicyUnits      int `json:"contextual_grounding_policy_units" bson:"contextual_grounding_policy_units"`
	SensitiveInformationPolicyFreeUnits int `json:"sensitive_information_policy_free_units" bson:"sensitive_information_policy_free_units"`
	SensitiveInformationPolicyUnits     int `json:"sensitive_information_policy_units" bson:"sensitive_information_policy_units"`
	TopicPolicyUnits                    int `json:"topic_policy_units" bson:"topic_policy_units"`
	WordPolicyUnits                     int `json:"word_policy_units" bson:"word_policy_units"`
}

type GuardrailCoverage struct {
	Images         *GuardrailCoverageData `json:"images" bson:"images"`
	TextCharacters *GuardrailCoverageData `json:"text_characters" bson:"text_characters"`
}

type GuardrailCoverageData struct {
	Guarded int32 `json:"guarded" bson:"guarded"`
	Total   int32 `json:"total" bson:"total"`
}

func ConvertFromBedrockOutput(resp *bedrockruntime.ApplyGuardrailOutput) *ApplyGuardrailOutput {
	if resp == nil {
		return nil
	}

	// Convert Assessments
	var assessments []GuardrailAssessment
	for _, a := range resp.Assessments {
		assessments = append(assessments, GuardrailAssessment{
			ContentPolicy:              convertContentPolicy(a.ContentPolicy),
			ContextualGroundingPolicy:  toStr(a.ContextualGroundingPolicy),
			InvocationMetrics:          convertInvocationMetrics(a.InvocationMetrics),
			SensitiveInformationPolicy: convertSensitiveInfoPolicy(a.SensitiveInformationPolicy),
			TopicPolicy:                convertTopicPolicy(a.TopicPolicy),
			WordPolicy:                 convertWordPolicy(a.WordPolicy),
		})
	}

	// Convert Outputs
	var outputs []string
	for _, o := range resp.Outputs {
		if o.Text != nil {
			outputs = append(outputs, *o.Text)
		}
	}

	// Convert Usage
	var usage *GuardrailUsage
	if resp.Usage != nil {
		usage = &GuardrailUsage{
			ContentPolicyUnits:                  toInt(resp.Usage.ContentPolicyUnits),
			ContextualGroundingPolicyUnits:      toInt(resp.Usage.ContextualGroundingPolicyUnits),
			SensitiveInformationPolicyFreeUnits: toInt(resp.Usage.SensitiveInformationPolicyFreeUnits),
			SensitiveInformationPolicyUnits:     toInt(resp.Usage.SensitiveInformationPolicyUnits),
			TopicPolicyUnits:                    toInt(resp.Usage.TopicPolicyUnits),
			WordPolicyUnits:                     toInt(resp.Usage.WordPolicyUnits),
		}
	}

	// Convert GuardrailCoverage
	var coverage *GuardrailCoverage
	if resp.GuardrailCoverage != nil {
		coverage = &GuardrailCoverage{
			Images:         convertImageCoverage(resp.GuardrailCoverage.Images),
			TextCharacters: convertTextCharCoverage(resp.GuardrailCoverage.TextCharacters),
		}
	}

	// Convert Metadata
	metadataStr := marshalMetadata(resp.ResultMetadata)

	return &ApplyGuardrailOutput{
		Action:            string(resp.Action),
		Assessments:       assessments,
		Outputs:           outputs,
		Usage:             usage,
		GuardrailCoverage: coverage,
		ResultMetadata:    metadataStr,
	}
}

func toStr(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}

	return string(bytes)
}

func toInt(ptr *int32) int {
	if ptr == nil {
		return 0
	}
	return int(*ptr)
}

func convertContentPolicy(v interface{}) *GuardrailContentFilter {
	if v == nil {
		return nil
	}

	// Attempt type assertion if SDK struct type is available
	if cp, ok := v.(types.GuardrailContentPolicyAssessment); ok {
		var filters []GuardrailFilter
		for _, f := range cp.Filters {
			filters = append(filters, GuardrailFilter{
				CommonActionTyped: CommonActionTyped{
					Action: string(f.Action),
					Type:   string(f.Type),
				},
				Confidence:     string(f.Confidence),
				FilterStrength: string(f.FilterStrength),
			})
		}
		return &GuardrailContentFilter{
			Filters: filters,
		}
	}

	// Fallback: try to parse from marshaled JSON if type info is missing
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	var fallback struct {
		Filters []GuardrailFilter `json:"Filters"`
	}
	_ = json.Unmarshal(bytes, &fallback)

	return &GuardrailContentFilter{
		Filters: fallback.Filters,
	}
}

func convertImageCoverage(v *types.GuardrailImageCoverage) *GuardrailCoverageData {
	if v == nil {
		return nil
	}
	return &GuardrailCoverageData{
		Guarded: *v.Guarded,
		Total:   *v.Total,
	}
}

func convertTextCharCoverage(v *types.GuardrailTextCharactersCoverage) *GuardrailCoverageData {
	if v == nil {
		return nil
	}
	return &GuardrailCoverageData{
		Guarded: *v.Guarded,
		Total:   *v.Total,
	}
}

func marshalMetadata(meta middleware.Metadata) string {
	bytes, err := json.Marshal(meta)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func convertInvocationMetrics(v interface{}) *GuardrailInvocationMetrics {
	if v == nil {
		return nil
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var metrics GuardrailInvocationMetrics
	if err := json.Unmarshal(bytes, &metrics); err != nil {
		return nil
	}
	return &metrics
}

func convertSensitiveInfoPolicy(v interface{}) *GuardrailSensitiveInformationPolicy {
	if v == nil {
		return nil
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var policy GuardrailSensitiveInformationPolicy
	if err := json.Unmarshal(bytes, &policy); err != nil {
		return nil
	}
	return &policy
}

func convertTopicPolicy(v interface{}) []*GuardrailNamedTypedAction {
	if v == nil {
		return nil
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var policy []*GuardrailNamedTypedAction
	if err := json.Unmarshal(bytes, &policy); err != nil {
		return nil
	}
	return policy
}

func convertWordPolicy(v interface{}) *GuardrailWordPolicyAssessment {
	if v == nil {
		return nil
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var policy GuardrailWordPolicyAssessment
	if err := json.Unmarshal(bytes, &policy); err != nil {
		return nil
	}
	return &policy
}
