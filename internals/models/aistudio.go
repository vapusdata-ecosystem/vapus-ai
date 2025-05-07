package models

import (
	"context"

	"github.com/pgvector/pgvector-go"
	"github.com/uptrace/bun"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	types "github.com/vapusdata-ecosystem/vapusdata/core/types"
)

type VectorEmbeddings struct {
	Vectors32 []float32 `json:"vectors32,omitempty" yaml:"vectors32"`
	Vectors64 []float64 `json:"vectors64,omitempty" yaml:"vectors64"`
}

type MessageLog struct {
	Content           string                   `json:"content,omitempty" yaml:"content,omitempty"`
	Role              string                   `json:"role,omitempty" yaml:"role,omitempty"`
	StructuredContent []*pb.RequestContentPart `json:"structured,omitempty" yaml:"structured,omitempty"`
}

type AIStudioChat struct {
	VapusBase  `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	ChatId     string         `bun:"chat_id" json:"chatId,omitempty" yaml:"chatId"`
	Messages   []*AIStudioLog `bun:"rel:has-many,join:chat_id=chat_id" json:"messages,omitempty" yaml:"messages,omitempty"`
	EndedAt    int64          `bun:"ended_at" json:"endedAt,omitempty" yaml:"endedAt,omitempty"`
	StoppedAt  int64          `bun:"stopped_at" json:"stoppedAt,omitempty" yaml:"stoppedAt,omitempty"`
	MessageIds []string       `bun:"message_ids,array" json:"messageIds,omitempty" yaml:"messageIds,omitempty"`
	CurrentLog *Mapper        `bun:"current_log,type:jsonb" json:"currentLog,omitempty" yaml:"currentLog,omitempty"`
}

func (dm *AIStudioChat) ConvertToPb() *pb.AIStudioChat {
	if dm == nil {
		return nil
	}
	return &pb.AIStudioChat{
		ChatId: dm.ChatId,
		Messages: func() []*pb.ChatMessageObject {
			var messages []*pb.ChatMessageObject
			for _, msg := range dm.Messages {
				for _, input := range msg.Input {
					messages = append(messages, &pb.ChatMessageObject{
						Content: input.Content,
						Role:    input.Role,
					})
				}
				for _, output := range msg.Output {
					messages = append(messages, &pb.ChatMessageObject{
						Content: output.Content,
						Role:    output.Role,
					})
				}
			}
			return messages
		}(),
		EndedAt:    dm.EndedAt,
		StoppedAt:  dm.StoppedAt,
		MessageIds: dm.MessageIds,
		CurrentLog: dm.CurrentLog.ConvertToPb(),
		StartedAt:  dm.CreatedAt,
	}
}

func (dm *AIStudioChat) ConvertToListingPb() *pb.AIStudioChat {
	if dm == nil {
		return nil
	}
	return &pb.AIStudioChat{
		ChatId:     dm.ChatId,
		EndedAt:    dm.EndedAt,
		StoppedAt:  dm.StoppedAt,
		CurrentLog: dm.CurrentLog.ConvertToPb(),
		StartedAt:  dm.CreatedAt,
	}
}

func (dm *AIStudioChat) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.Status = mpb.CommonStatus_ACTIVE.String()
	dm.PreSaveVapusBase(authzClaim)
}

func (dm *AIStudioChat) PreSaveUpdate(status string) {
	if dm == nil {
		return
	}
	dm.Status = status
	dm.UpdatedAt = dmutils.GetEpochTime()
}

func (dm *AIStudioChat) PreSaveDelete(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.Status = mpb.CommonStatus_INACTIVE.String()
	dm.PreDeleteVapusBase(authzClaim)
}

type AIStudioLog struct {
	VapusBase        `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Input            []*MessageLog   `bun:"input,type:jsonb" json:"Input,omitempty" yaml:"Input"`
	ParsedInput      []*MessageLog   `bun:"parsed_input,type:jsonb" json:"parsedInput,omitempty" yaml:"parsedInput"`
	PromptId         string          `bun:"prompt_id" json:"promptId,omitempty" yaml:"promptId"`
	ResponseStatus   string          `bun:"response_status" json:"responseStatus,omitempty" yaml:"responseStatus"`
	Error            string          `bun:"error" json:"error,omitempty" yaml:"error"`
	Mode             string          `bun:"mode" json:"mode,omitempty" yaml:"mode,omitempty"`
	Output           []*MessageLog   `bun:"output,type:jsonb" json:"output,omitempty" yaml:"output"`
	ParsedOutput     []*MessageLog   `bun:"parsed_output,type:jsonb" json:"parsedOutput,omitempty" yaml:"parsedOutput"`
	AIModel          string          `bun:"ai_model" json:"aiModel,omitempty" yaml:"aiModel"`
	ModelNode        string          `bun:"model_node" json:"modelNode,omitempty" yaml:"modelNode"`
	ModelProvider    string          `bun:"model_provider" json:"modelProvider,omitempty" yaml:"modelProvider"`
	LogEmbeddings    pgvector.Vector `bun:"log_embeddings,type:vector(1536)"`
	StartedAt        int64           `bun:"started_at" json:"startedAt,omitempty" yaml:"startedAt,omitempty"`
	EndedAt          int64           `bun:"ended_at" json:"endedAt,omitempty" yaml:"endedAt,omitempty"`
	TTFBAt           int64           `bun:"ttfb_at" json:"ttfbAt,omitempty" yaml:"ttfbAt,omitempty"`
	InputFts         string          `bun:"input_fts,type:tsvector" json:"inputFts,omitempty" yaml:"inputFts,omitempty"`
	OutputFts        string          `bun:"output_fts,type:tsvector" json:"outputFts,omitempty" yaml:"outputFts,omitempty"`
	ToolCallSchema   []*mpb.ToolCall `bun:"tool_call_schema,type:jsonb" json:"toolCallSchema,omitempty" yaml:"toolCallSchema,omitempty"`
	ToolCallResponse []*mpb.ToolCall `bun:"tool_call_response,type:jsonb" json:"toolCallResponse,omitempty" yaml:"toolCallResponse,omitempty"`
	InputFiles       []*mpb.FileData `bun:"input_files,type:jsonb" json:"inputFiles,omitempty" yaml:"inputFiles,omitempty"`
	OutputFiles      []*mpb.FileData `bun:"output_files,type:jsonb" json:"outputFiles,omitempty" yaml:"outputFiles,omitempty"`
	ChatId           string          `bun:"chat_id" json:"chatId,omitempty" yaml:"chatId"`
	Summary          []string        `bun:"summary,array" json:"summary,omitempty" yaml:"summary"`
}

var _ bun.AfterCreateTableHook = (*AIStudioLog)(nil)

func (a *AIStudioLog) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	var err error
	_, err = query.DB().NewCreateIndex().IfNotExists().
		Model((*AIStudioLog)(nil)).TableExpr("ai_studio_logs").
		Index("ai_studio_logs_input_fts_idx").
		ColumnExpr("input_fts").
		Using("gin").
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = query.DB().NewCreateIndex().IfNotExists().
		Model((*AIStudioLog)(nil)).TableExpr("ai_studio_logs").
		Index("ai_studio_logs_output_fts_idx").
		ColumnExpr("output_fts").
		Using("gin").
		Exec(ctx)
	if err != nil {
		return err
	}

	// Add pgvector indexes for log_embeddings with l2, inner_product and cosine operators
	_, err = query.DB().NewCreateIndex().IfNotExists().
		Model((*AIStudioLog)(nil)).TableExpr("ai_studio_logs").
		Index("ai_studio_logs_log_embeddings_l2_idx").
		ColumnExpr("log_embeddings vector_l2_ops").
		Using("hnsw").
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = query.DB().NewCreateIndex().IfNotExists().
		Model((*AIStudioLog)(nil)).TableExpr("ai_studio_logs").
		Index("ai_studio_logs_log_embeddings_ip_idx").
		ColumnExpr("log_embeddings vector_ip_ops").
		Using("hnsw").
		Exec(ctx)
	if err != nil {
		return err
	}
	_, err = query.DB().NewCreateIndex().IfNotExists().
		Model((*AIStudioLog)(nil)).TableExpr("ai_studio_logs").
		Index("ai_studio_logs_log_embeddings_cosine_idx").
		ColumnExpr("log_embeddings vector_cosine_ops").
		Using("hnsw").
		Exec(ctx)
	if err != nil {
		return err
	}
	return err
}

func (dm *AIStudioLog) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.PreSaveVapusBase(authzClaim)
}

func (dm *AIStudioLog) PreSaveUpdate(userId string) {
	if dm == nil {
		return
	}
	dm.UpdatedBy = userId
	dm.UpdatedAt = dmutils.GetEpochTime()
}

type SecurityGuardrails struct {
	Guardrails []string `json:"guardrails,omitempty" yaml:"guardrails,omitempty"`
}

func (s *SecurityGuardrails) ConvertFromPb(pb *mpb.SecurityGuardrails) *SecurityGuardrails {
	if pb == nil {
		return nil
	}
	return &SecurityGuardrails{
		Guardrails: pb.GetGuardrails(),
	}
}

func (s *SecurityGuardrails) ConvertToPb() *mpb.SecurityGuardrails {
	if s != nil {
		return &mpb.SecurityGuardrails{
			Guardrails: s.Guardrails,
		}
	}
	return nil
}

type AIStudioUsages struct {
	VapusBase                `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	AIStudioLogId            string           `bun:"ai_studio_log_id,notnull" json:"aiStudioLogId,omitempty" yaml:"aiStudioLogId,omitempty"`
	AIStudioLog              *AIStudioLog     `bun:"rel:belongs-to,join:ai_studio_log_id=vapus_id" json:"-" yaml:"-"`
	InputTokens              int64            `bun:"input_tokens" json:"inputTokens,omitempty" yaml:"inputTokens,omitempty"`
	InputCachedTokens        int64            `bun:"input_cached_tokens" json:"inputCachedTokens,omitempty" yaml:"inputCachedTokens,omitempty"`
	OutputTokens             int64            `bun:"output_tokens" json:"outputTokens,omitempty" yaml:"outputTokens,omitempty"`
	OutputCachedTokens       int64            `bun:"output_cached_tokens" json:"outputCachedTokens,omitempty" yaml:"outputCachedTokens,omitempty"`
	InputAudioTokens         int64            `bun:"input_audio_tokens" json:"inputAudioTokens,omitempty" yaml:"inputAudioTokens,omitempty"`
	OutputAudioTokens        int64            `bun:"output_audio_tokens" json:"outputAudioTokens,omitempty" yaml:"outputAudioTokens,omitempty"`
	TotalTokens              int64            `bun:"total_tokens" json:"totalTokens,omitempty" yaml:"totalTokens,omitempty"`
	GuardrailLog             *AIGuardrailsLog `bun:"rel:belongs-to,join:ai_guardrails_logs_id=vapus_id" json:"-" yaml:"-"`
	GuardrailLogId           string           `bun:"ai_guardrails_logs_id" json:"guardrailLogId,omitempty" yaml:"guardrailLogId,omitempty"`
	Charges                  float64          `bun:"charges" json:"charges,omitempty" yaml:"charges,omitempty"`
	ReasoningTokens          int64            `bun:"reasoning_tokens" json:"reasoningTokens,omitempty" yaml:"reasoningTokens,omitempty"`
	AcceptedPredictionTokens int64            `bun:"accepted_prediction_tokens" json:"acceptedPredictionTokens,omitempty" yaml:"acceptedPredictionTokens,omitempty"`
	RejectedPredictionTokens int64            `bun:"rejected_prediction_tokens" json:"rejectedPredictionTokens,omitempty" yaml:"rejectedPredictionTokens,omitempty"`
}

func (dm *AIStudioUsages) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.PreSaveVapusBase(authzClaim)
}

type AIGuardrailsLog struct {
	VapusBase       `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Input           []*MessageLog   `bun:"input,type:jsonb" json:"Input,omitempty" yaml:"Input"`
	InputEmbeddings pgvector.Vector `bun:"input_embeddings,type:vector(1536)"`
	Output          []string        `bun:"output,type:jsonb" json:"output,omitempty" yaml:"output"`
	Failed          bool            `bun:"failed" json:"failed,omitempty" yaml:"failed"`
	StartedAt       int64           `bun:"started_at" json:"startedAt,omitempty" yaml:"startedAt,omitempty"`
	EndedAt         int64           `bun:"ended_at" json:"endedAt,omitempty" yaml:"endedAt,omitempty"`
	TTFBAt          int64           `bun:"ttfb_at" json:"ttfbAt,omitempty" yaml:"ttfbAt,omitempty"`
	FailedMessage   string          `bun:"failed_message" json:"failedMessage,omitempty" yaml:"failedMessage"`
}

func (dm *AIGuardrailsLog) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.PreSaveVapusBase(authzClaim)
}

type (
	AIGatewayError struct {
		Error    AIGatewayErrorDetail `json:"error"`              // The error object.
		Status   int                  `json:"status,omitempty"`   // The HTTP status code.
		Provider string               `json:"provider,omitempty"` // The provider of the error.
	}
	AIGatewayErrorDetail struct {
		Message string `json:"message"`         // A human-readable message providing more details about the error.
		Type    string `json:"type,omitempty"`  // The type of error, e.g., "invalid_request_error".
		Param   string `json:"param,omitempty"` // The parameter associated with the error (if any).
		Code    string `json:"code,omitempty"`  // An error code for programmatic handling (if any).
	}
)

type AIToolCallLog struct {
	VapusBase  `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	PlainInput string `bun:"plain_input" json:"plainInput,omitempty" yaml:"plainInput,omitempty" toml:"plainInput,omitempty"`
	Input      string `bun:"input,type:tsvector" json:"input,omitempty" yaml:"input,omitempty" toml:"input,omitempty"`
	// InputEmbeddings pgvector.Vector `bun:"text_query_embeddings,type:vector(1536)"`
	OutputSchema   string `bun:"output_schema" json:"outputSchema,omitempty" yaml:"outputSchema,omitempty" toml:"outputSchema,omitempty"`
	ActionAnalyzer bool   `bun:"action_analyzer" json:"actionAnalyzer,omitempty" yaml:"actionAnalyzer,omitempty" toml:"actionAnalyzer,omitempty"`
	ParamAnalyzer  bool   `bun:"param_analyzer" json:"paramAnalyzer,omitempty" yaml:"paramAnalyzer,omitempty" toml:"paramAnalyzer,omitempty"`
}

func (dmn *AIToolCallLog) PreSaveCreate(authzClaim map[string]string) {
	if dmn == nil {
		return
	}
	if dmn.CreatedBy == types.EMPTYSTR {
		dmn.CreatedBy = authzClaim[encryption.ClaimUserIdKey]
	}
	if dmn.CreatedAt == 0 {
		dmn.CreatedAt = dmutils.GetEpochTime()
	}
	dmn.OwnerAccount = authzClaim[encryption.ClaimAccountKey]
}

var _ bun.AfterCreateTableHook = (*AIToolCallLog)(nil)

func (*AIToolCallLog) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	var err error
	_, err = query.DB().NewCreateIndex().IfNotExists().
		Model((*AIToolCallLog)(nil)).TableExpr("fabric_chat_tool_logs").
		Index("input_idx").
		ColumnExpr("input").
		Using("gin").
		Exec(ctx)
	return err
}
