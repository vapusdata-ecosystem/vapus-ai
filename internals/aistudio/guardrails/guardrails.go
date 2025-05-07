package guardrails

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	aipb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	aicore "github.com/vapusdata-ecosystem/vapusai/core/aistudio/core"
	"github.com/vapusdata-ecosystem/vapusai/core/aistudio/prompts"
	aimodels "github.com/vapusdata-ecosystem/vapusai/core/aistudio/providers"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
)

type GuardRailClient struct {
	Guardrail           *models.AIGuardrails
	modelPool           []*GuardModelNodePool
	dmstores            *apppkgs.VapusStore
	embeddingConnection aimodels.AIModelNodeInterface
}

type GuardModelNodePool struct {
	Connection aimodels.AIModelNodeInterface
	IsAccount  bool
	Model      string
}

type GuardrailScanner struct {
	ContentGuard    []string
	TopicGuard      []string
	WordGuard       []string
	SenstivityGuard []string
	client          *GuardRailClient
	ID              string
	Usage           *models.AIStudioUsages
}

type GuardrailsFunc func(*GuardRailClient)

func WithSpec(spec *models.AIGuardrails) GuardrailsFunc {
	return func(g *GuardRailClient) {
		g.Guardrail = spec
	}
}

func WithModelPool(pool []*GuardModelNodePool) GuardrailsFunc {
	return func(g *GuardRailClient) {
		g.modelPool = pool
	}
}

func WithDataStore(store *apppkgs.VapusStore) GuardrailsFunc {
	return func(g *GuardRailClient) {
		g.dmstores = store
	}
}

func WithEmbeddingConnection(conn aimodels.AIModelNodeInterface) GuardrailsFunc {
	return func(g *GuardRailClient) {
		g.embeddingConnection = conn
	}
}

func New(opts ...GuardrailsFunc) *GuardRailClient {
	cl := &GuardRailClient{}
	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

func (g *GuardRailClient) Scan(ctx context.Context, message string, logger zerolog.Logger, ctxClaim map[string]string) *GuardrailScanner {
	scanner := &GuardrailScanner{
		client: g,
	}

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		scanner.ScanContent(ctx, message, logger)
	}()
	go func() {
		defer wg.Done()
		scanner.ScanTopic(ctx, message, logger)
	}()
	go func() {
		defer wg.Done()
		scanner.ScanWords(ctx, message, logger)
	}()
	go func() {
		defer wg.Done()
		scanner.SensitivityDataAction(ctx, message, logger)
	}()
	wg.Wait()

	return scanner
}

func (g *GuardrailScanner) SetUsage(ctx context.Context, opts *models.AIStudioUsages) {
	if g.Usage == nil {
		g.Usage = opts
		return
	} else {
		g.Usage.InputTokens += opts.InputTokens
		g.Usage.InputCachedTokens += opts.InputCachedTokens
		g.Usage.OutputTokens += opts.OutputTokens
		g.Usage.OutputCachedTokens += opts.OutputCachedTokens
		g.Usage.InputAudioTokens += opts.InputAudioTokens
		g.Usage.OutputAudioTokens += opts.OutputAudioTokens
		g.Usage.TotalTokens += opts.TotalTokens
	}
}

func (g *GuardrailScanner) ScanContent(ctx context.Context, message string, logger zerolog.Logger) {}
func (g *GuardrailScanner) ScanTopic(ctx context.Context, message string, logger zerolog.Logger) {
	var contentHeatMap = map[string]string{
		"sexual":     g.client.Guardrail.Contents.Sexual,
		"hateSpeech": g.client.Guardrail.Contents.HateSpeech,
		"threats":    g.client.Guardrail.Contents.Threats,
		"insults":    g.client.Guardrail.Contents.Insults,
		"misconduct": g.client.Guardrail.Contents.Misconduct,
	}

	payload := prompts.NewPrompter(&aipb.ChatRequest{
		Model: g.client.modelPool[0].Model,
		Tools: func() []*mpb.ToolCall {
			var toolCalls []*mpb.ToolCall
			f := models.GetFunctionCallFromString(g.client.Guardrail.Schema)
			toolCalls = append(toolCalls, &mpb.ToolCall{
				Type: strings.ToLower(mpb.AIToolCallType_FUNCTION.String()),
				FunctionSchema: &mpb.FunctionCall{
					Name:        f.Name,
					Parameters:  f.GetStringParamSchema(),
					Description: f.Description,
				},
			})
			return toolCalls
		}(),
		Messages: []*aipb.ChatMessageObject{
			{
				Role:    aicore.SYSTEM,
				Content: "You are an AI guardrail inspector, please scan the user input based and provide the tool call response. Strictly, do not generate any other dataset.",
			},
			{
				Role:    aicore.USER,
				Content: message,
			},
		},
	}, nil, nil, nil, logger)

	payload.RenderPrompt()
	err := g.client.modelPool[0].Connection.GenerateContent(ctx, payload)
	if payload.Usage != nil {
		g.SetUsage(ctx, payload.Usage)
	}
	if err != nil {
		logger.Error().Err(err).Msg("error while generating content")
		return
	}
	for _, obj := range payload.ToolCallResponse {
		argsresult := obj.FunctionSchema.GetParameters()
		result := map[string]any{}
		err = json.Unmarshal([]byte(argsresult), &result)
		if err != nil {
			logger.Error().Err(err).Msg("error while unmarshalling function call arguments for topic")
			continue
		}
		ct, ok := result["content_guardrails"]
		if ok {
			cg, ok := ct.(map[string]any)
			if ok {
				for k, v := range cg {
					val := v.(string)
					acceptedVal, ok := mpb.GuardRailLevels_value[contentHeatMap[k]]
					if !ok {
						logger.Error().Msg("error while getting guardrail level")
						continue
					}
					foundLevel, ok := mpb.GuardRailLevels_value[strings.ToUpper(val)]
					if !ok {
						logger.Error().Msg("error while getting guardrail level")
						continue
					}
					if acceptedVal != 0 {
						if foundLevel >= acceptedVal {
							g.ContentGuard = append(g.ContentGuard, fmt.Sprintf("%s is %v for input provided by the user, guardrail check failed", k, v))
						}
					}
				}
			}
		}
		tg, ok := result["topic_guardrails"]
		if ok {
			topic, ok := tg.(map[string]any)
			if ok {
				for k, v := range topic {
					if v.(bool) == true {
						g.TopicGuard = append(g.TopicGuard, "user message contains topic "+k)
					}
				}
			}
		}
	}
}
func (g *GuardrailScanner) ScanWords(ctx context.Context, message string, logger zerolog.Logger) {
	for _, rule := range g.client.Guardrail.Words {
		for _, word := range rule.Words {
			if strings.Contains(message, word) {
				g.WordGuard = append(g.WordGuard, word)
			}
		}
	}
}
func (g *GuardrailScanner) SensitivityDataAction(ctx context.Context, message string, logger zerolog.Logger) {
}
