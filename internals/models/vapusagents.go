package models

import (
	"fmt"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type VapusAgents struct {
	VapusBase       `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name            string               `bun:"name,unique" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Description     string               `bun:"description" json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	AgentType       string               `bun:"agent_type" json:"agent_type,omitempty" yaml:"agent_type,omitempty" toml:"agent_type,omitempty"`
	Specs           []*AgentSpec         `bun:"specs,type:jsonb" json:"specs,omitempty" yaml:"specs,omitempty" toml:"specs,omitempty"`
	CurrentVersion  string               `bun:"current_version" json:"current_version,omitempty" yaml:"current_version,omitempty" toml:"current_version,omitempty"`
	CurrentRunLogId string               `bun:"current_run_log_id" json:"current_run_log_id,omitempty" yaml:"current_run_log_id,omitempty" toml:"current_run_log_id,omitempty"`
	Attributes      VapusAgentAttributes `bun:"attributes,type:jsonb" json:"attributes,omitempty" yaml:"attributes,omitempty" toml:"attributes,omitempty"`
	Logs            []*VapusAgentLog     `bun:"rel:has-many,join:vapus_id=agent_id" json:"logs,omitempty" yaml:"logs,omitempty" toml:"logs,omitempty"`
	AssetStore      string               `bun:"asset_store" json:"asset_store,omitempty" yaml:"asset_store,omitempty" toml:"asset_store,omitempty"`
	LastRunStartAt  int64                `bun:"last_start_at" json:"last_start_at,omitempty" yaml:"last_start_at,omitempty" toml:"last_start_at,omitempty"`
	LastRunEndAt    int64                `bun:"last_end_at" json:"last_end_at,omitempty" yaml:"last_end_at,omitempty" toml:"last_end_at,omitempty"`
	LastRunStatus   string               `bun:"last_run_status" json:"last_run_status,omitempty" yaml:"last_run_status,omitempty" toml:"last_run_status,omitempty"`
	ModelNode       string               `bun:"model_node" json:"model_node,omitempty" yaml:"model_node,omitempty" toml:"model_node,omitempty"`
	Model           string               `bun:"model" json:"model,omitempty" yaml:"model,omitempty" toml:"model,omitempty"`
}

func (dn *VapusAgents) PreSaveCreate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreSaveVapusBase(authzClaim)
}

func (dn *VapusAgents) PreSaveDelete(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.DeletedAt = dmutils.GetEpochTime()
	dn.DeletedBy = dn.CreatedBy
	dn.Status = mpb.CommonStatus_EXPIRED.String()
	dn.PreDeleteVapusBase(authzClaim)
}

func (dn *VapusAgents) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *VapusAgents) GetVersionSpec(version string) (*AgentSpec, error) {
	if dn == nil {
		return nil, fmt.Errorf("VapusAgent is nil")
	}
	if version == "" {
		version = dn.CurrentVersion
	}
	for _, v := range dn.Specs {
		if v.Version == version {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Version not found")
}

type VapusAgentAttributes struct {
	LogRetentionCount      int            `bun:"log_retention_count" json:"log_retention_count,omitempty" yaml:"log_retention_count,omitempty" toml:"log_retention_count,omitempty"`
	LogRetentionDays       int            `bun:"log_retention_days" json:"log_retention_days,omitempty" yaml:"log_retention_days,omitempty" toml:"log_retention_days,omitempty"`
	ErrorLogRetentionCount int            `bun:"error_log_retention_count" json:"error_log_retention_count,omitempty" yaml:"error_log_retention_count,omitempty" toml:"error_log_retention_count,omitempty"`
	RetryCount             int            `bun:"retry_count" json:"retry_count,omitempty" yaml:"retry_count,omitempty" toml:"retry_count,omitempty"`
	RetryFailureExit       bool           `bun:"retry_failure_exit" json:"retry_failure_exit,omitempty" yaml:"retry_failure_exit,omitempty" toml:"retry_failure_exit,omitempty"`
	Schedule               *VapusSchedule `bun:"schedule" json:"schedule,omitempty" yaml:"schedule,omitempty" toml:"schedule,omitempty"`
	FailChannelPlugin      string         `bun:"fail_channel_plugin" json:"fail_channel_plugin,omitempty" yaml:"fail_channel_plugin,omitempty" toml:"fail_channel_plugin,omitempty"`
}

type AgentSpec struct {
	Instructions   string          `json:"instructions,omitempty" yaml:"instructions,omitempty" toml:"instructions,omitempty"`
	Dataproducts   []string        `json:"dataproducts,omitempty" yaml:"dataproducts,omitempty" toml:"dataproducts,omitempty"`
	Goal           string          `json:"goal,omitempty" yaml:"goal,omitempty" toml:"goal,omitempty"`
	Version        string          `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`
	InputToolCalls []*FunctionCall `json:"tool_call,omitempty" yaml:"tool_call,omitempty" toml:"tool_call,omitempty"`
	Files          []*mpb.FileData `json:"files,omitempty" yaml:"files,omitempty" toml:"files,omitempty"`
	SchemaArgs     string          `json:"schema_args,omitempty" yaml:"schema_args,omitempty" toml:"schema_args,omitempty"`
	VersionStatus  string          `json:"version_status,omitempty" yaml:"version_status,omitempty" toml:"version_status,omitempty"`
	ReleaseNotes   string          `json:"releaseNotes,omitempty" yaml:"releaseNotes,omitempty" toml:"releaseNotes,omitempty"`
}

type VapusAgentLog struct {
	VapusBase     `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	AgentID       string                   `bun:"agent_id" json:"agent_id,omitempty" yaml:"agent_id,omitempty" toml:"agent_id,omitempty"`
	Log           []string                 `bun:"log" json:"log,omitempty" yaml:"log,omitempty" toml:"log,omitempty"`
	Message       *VapusAgentMessage       `bun:"message" json:"message,omitempty" yaml:"message,omitempty" toml:"message,omitempty"`
	ErrorLog      []string                 `bun:"error_log" json:"error_log,omitempty" yaml:"error_log,omitempty" toml:"error_log,omitempty"`
	Timestamp     int64                    `bun:"timestamp" json:"timestamp,omitempty" yaml:"timestamp,omitempty" toml:"timestamp,omitempty"`
	Observability *VapusAgentObservability `bun:"observability" json:"observability,omitempty" yaml:"observability,omitempty" toml:"observability,omitempty"`
}

type VapusAgentObservability struct {
	CurrentCrewAgent string              `json:"current_crew_agent,omitempty" yaml:"current_crew_agent,omitempty" toml:"current_crew_agent,omitempty"`
	Reasonings       map[string][]string `json:"reasonings,omitempty" yaml:"reasonings,omitempty" toml:"reasonings,omitempty"`
	CrewRunningLogs  []*CrewRunningLog   `json:"crew_running_logs,omitempty" yaml:"crew_running_logs,omitempty" toml:"crew_running_logs,omitempty"`
}

type VapusAgentMessage struct {
	MessageId   string          `json:"messageId,omitempty" yaml:"messageId,omitempty" toml:"messageId,omitempty"`
	Input       string          `json:"input,omitempty" yaml:"input,omitempty" toml:"input,omitempty"`
	OutputFiles []*mpb.FileData `json:"outputFiles,omitempty" yaml:"outputFiles,omitempty" toml:"outputFiles,omitempty"`
	Error       string          `json:"error,omitempty" yaml:"error,omitempty" toml:"error,omitempty"`
	// OutputAssets      []*NabhikAsset          `json:"assets,omitempty" yaml:"assets,omitempty" toml:"assets,omitempty"`
	DataResultSummary *FabricDataQuerySummary `json:"dataResultSummary,omitempty" yaml:"dataResultSummary,omitempty" toml:"dataResultSummary,omitempty"`
}

type CrewRunningLog struct {
	CrewAgent string   `json:"crewAgent,omitempty" yaml:"crewAgent,omitempty" toml:"crewAgent,omitempty"`
	RunId     string   `json:"runId,omitempty" yaml:"runId,omitempty" toml:"runId,omitempty"`
	Logs      []string `json:"logs,omitempty" yaml:"logs,omitempty" toml:"logs,omitempty"`
	Errors    []string `json:"errors,omitempty" yaml:"errors,omitempty" toml:"errors,omitempty"`
	// Assets             []*NabhikAsset `json:"assets,omitempty" yaml:"assets,omitempty" toml:"assets,omitempty"`
	StartedAt          int64          `json:"startedAt,omitempty" yaml:"startedAt,omitempty" toml:"startedAt,omitempty"`
	EndedAt            int64          `json:"endedAt,omitempty" yaml:"endedAt,omitempty" toml:"endedAt,omitempty"`
	CallerCrewAgent    string         `json:"callerCrewAgent,omitempty" yaml:"callerCrewAgent,omitempty" toml:"callerCrewAgent,omitempty"`
	CrewTools          []string       `json:"crewTools,omitempty" yaml:"crewTools,omitempty" toml:"crewTools,omitempty"`
	CrewToolSchemaArgs map[string]any `json:"crewToolSchemaArgs,omitempty" yaml:"crewToolSchemaArgs,omitempty" toml:"crewToolSchemaArgs,omitempty"`
}

func (f *VapusAgents) ConvertToPb() *mpb.VapusAgent {
	return &mpb.VapusAgent{
		AgentId:     f.VapusID,
		Name:        f.Name,
		Description: f.Description,
		Specs: func() []*mpb.AgentSpec {
			var specs []*mpb.AgentSpec
			for _, v := range f.Specs {
				specs = append(specs, v.ConvertToPb())
			}
			return specs
		}(),
		CurrentVersion: f.CurrentVersion,
		CurrentMessage: f.CurrentRunLogId,
		Attributes:     f.Attributes.ConvertToPb(),
		ResourceBase:   f.ConvertToPbBase(),
		ModelNode:      f.ModelNode,
		Model:          f.Model,
	}
}

func (f *VapusAgents) ConvertFromPb(pb *mpb.VapusAgent) {
	f.VapusID = pb.AgentId
	f.Name = pb.Name
	f.Description = pb.Description
	f.Specs = make([]*AgentSpec, 0)
	for _, v := range pb.Specs {
		spec := &AgentSpec{}
		spec.ConvertFromPb(v)
		f.Specs = append(f.Specs, spec)
	}
	f.CurrentVersion = pb.CurrentVersion
	f.CurrentRunLogId = pb.CurrentMessage
	f.Attributes = VapusAgentAttributes{}
	f.Attributes.ConvertFromPb(pb.Attributes)
	f.ModelNode = pb.ModelNode
	f.Model = pb.Model
}

func (f *VapusAgents) IsOwner(authz map[string]string) bool {
	return f.CreatedBy == authz[encryption.ClaimUserIdKey] && f.Organization == authz[encryption.ClaimOrganizationKey]
}

func (f *AgentSpec) ConvertToPb() *mpb.AgentSpec {
	var toolCalls []*mpb.FunctionCall
	for _, v := range f.InputToolCalls {
		toolCalls = append(toolCalls, v.ConvertToPb())
	}
	return &mpb.AgentSpec{
		Instructions:   f.Instructions,
		Dataproducts:   f.Dataproducts,
		Goal:           f.Goal,
		ReleaseNotes:   f.ReleaseNotes,
		Version:        f.Version,
		InputToolCalls: toolCalls,
	}
}

func (f *AgentSpec) ConvertFromPb(pb *mpb.AgentSpec) {
	var toolCalls []*FunctionCall
	for _, v := range pb.InputToolCalls {
		toolCall := (&FunctionCall{}).ConvertFromPb(v)
		toolCalls = append(toolCalls, toolCall)
	}
	f.Instructions = pb.Instructions
	f.Dataproducts = pb.Dataproducts
	f.Goal = pb.Goal
	f.ReleaseNotes = pb.ReleaseNotes
	f.Version = pb.Version
	f.InputToolCalls = toolCalls
}

func (f *VapusAgentAttributes) ConvertToPb() *mpb.AgentAttributes {
	return &mpb.AgentAttributes{
		LogRetentionCount:      int32(f.LogRetentionCount),
		LogRetentionDays:       int32(f.LogRetentionDays),
		ErrorLogRetentionCount: int32(f.ErrorLogRetentionCount),
		RetryCount:             int32(f.RetryCount),
		RetryFailureExit:       f.RetryFailureExit,
		Schedule:               f.Schedule.ConvertToPb(),
		FailChannelPlugin:      f.FailChannelPlugin,
	}
}

func (f *VapusAgentAttributes) ConvertFromPb(pb *mpb.AgentAttributes) {
	f.LogRetentionCount = int(pb.LogRetentionCount)
	f.LogRetentionDays = int(pb.LogRetentionDays)
	f.ErrorLogRetentionCount = int(pb.ErrorLogRetentionCount)
	f.RetryCount = int(pb.RetryCount)
	f.RetryFailureExit = pb.RetryFailureExit
	f.Schedule = (&VapusSchedule{}).ConvertFromPb(pb.Schedule)
	f.FailChannelPlugin = pb.FailChannelPlugin

}
