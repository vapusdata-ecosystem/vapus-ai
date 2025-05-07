package models

// import (
// 	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
// 	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
// 	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
// )

// type AIAgent struct {
// 	VapusBase    `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
// 	Name         string           `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
// 	AgentType    string           `bun:"agent_type" json:"agentType,omitempty" yaml:"agentType,omitempty" toml:"agentType,omitempty"`
// 	Description  string           `bun:"description" json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
// 	AgentVersion string           `bun:"agent_version" json:"agentVersion,omitempty" yaml:"agentVersion,omitempty" toml:"agentVersion,omitempty"`
// 	AIModelMap   []*AIModelMap    `bun:"ai_model_map,type:jsonb" son:"aiModelMap,omitempty" yaml:"aiModelMap,omitempty" toml:"aiModelMap,omitempty"`
// 	Editable     bool             `bun:"editable" json:"editable,omitempty" yaml:"editable,omitempty" toml:"editable,omitempty" default:"true"`
// 	Steps        []*Steps         `bun:"steps,type:jsonb" json:"steps,omitempty" yaml:"steps,omitempty" toml:"steps,omitempty"`
// 	Labels       []string         `bun:"labels,array" json:"labels,omitempty" yaml:"labels,omitempty" toml:"labels,omitempty"`
// 	Settings     *AIAgentSettings `bun:"settings,type:jsonb" json:"settings,omitempty" yaml:"settings,omitempty" toml:"settings,omitempty"`
// }

// type AIModelMap struct {
// 	ModelNodeId string `json:"modelNodeId,omitempty" yaml:"modelNodeId,omitempty" toml:"modelNodeId,omitempty"`
// 	ModelName   string `json:"modelName,omitempty" yaml:"modelName,omitempty" toml:"modelName,omitempty"`
// }

// type AIAgentSettings struct {
// 	ToolCallSchema string `json:"toolCallSchema,omitempty" yaml:"toolCallSchema,omitempty" toml:"toolCallSchema,omitempty"`
// }

// func (dm *AIAgentSettings) ConvertToPb() *mpb.AIAgentSettings {
// 	if dm == nil {
// 		return nil
// 	}
// 	return &mpb.AIAgentSettings{
// 		ToolCallSchema: dm.ToolCallSchema,
// 	}
// }

// func (dm *AIAgentSettings) ConvertFromPb(pb *mpb.AIAgentSettings) *AIAgentSettings {
// 	if pb == nil {
// 		return nil
// 	}
// 	dm.ToolCallSchema = pb.ToolCallSchema
// 	return dm
// }

// func (dm *AIAgentSettings) MarshalToString(format mpb.ContentFormats) string {
// 	if dm == nil {
// 		return ""
// 	}
// 	switch format {
// 	case mpb.ContentFormats_JSON:
// 		bbytes, err := filetools.GenericMarshaler(dm, format.String())
// 		if err != nil {
// 			return ""
// 		}
// 		return string(bbytes)
// 	default:
// 		bbytes, err := filetools.GenericMarshaler(dm, format.String())
// 		if err != nil {
// 			return ""
// 		}
// 		return string(bbytes)
// 	}
// }

// func (dm *AIModelMap) ConvertToPb() *mpb.AIModelMap {
// 	if dm == nil {
// 		return nil
// 	}
// 	return &mpb.AIModelMap{
// 		ModelNodeId: dm.ModelNodeId,
// 		ModelName:   dm.ModelName,
// 	}
// }

// func (dm *AIModelMap) ConvertFromPb(pb *mpb.AIModelMap) *AIModelMap {
// 	if pb == nil {
// 		return nil
// 	}
// 	dm.ModelNodeId = pb.ModelNodeId
// 	dm.ModelName = pb.ModelName
// 	return dm
// }

// type Steps struct {
// 	Id             string    `json:"id,omitempty" yaml:"id,omitempty" toml:"id,omitempty"`
// 	Prompt         string    `json:"prompt,omitempty" yaml:"prompt,omitempty" toml:"prompt,omitempty"`
// 	Required       bool      `json:"required,omitempty" yaml:"required,omitempty" toml:"required,omitempty" default:"true"`
// 	AutoGenerate   bool      `json:"autoGenerate,omitempty" yaml:"autoGenerate,omitempty" toml:"autoGenerate,omitempty" default:"true"`
// 	InputTemplates []*Mapper `json:"input_templates,omitempty" yaml:"inputTemplates,omitempty" toml:"inputTemplates,omitempty"`
// 	PromptId       string    `json:"prompt_id,omitempty" yaml:"promptId,omitempty" toml:"promptId,omitempty"`
// 	ValueType      string    `json:"value_type,omitempty" yaml:"valueType,omitempty" toml:"valueType,omitempty"`
// }

// func (dm *Steps) ConvertToPb() *mpb.Steps {
// 	if dm == nil {
// 		return nil
// 	}
// 	return &mpb.Steps{
// 		Id:     mpb.AgentStepEnum(mpb.AgentStepEnum_value[dm.Id]),
// 		Prompt: dm.Prompt,
// 		// InputTemplates: MapperSliceToPb(dm.InputTemplates),
// 		AutoGenerate: dm.AutoGenerate,
// 		Required:     dm.Required,
// 		PromptId:     dm.PromptId,
// 		ValueType:    mpb.AgentStepValueType(mpb.AgentStepValueType_value[dm.ValueType]),
// 	}
// }

// func (dm *Steps) ConvertFromPb(cpb *mpb.Steps) *Steps {
// 	if cpb == nil {
// 		return nil
// 	}
// 	dm.Id = cpb.Id.String()
// 	dm.Prompt = cpb.Prompt
// 	// dm.InputTemplates = MapperSliceFromPb(cpb.InputTemplates)
// 	dm.AutoGenerate = cpb.AutoGenerate
// 	dm.Required = cpb.Required
// 	dm.PromptId = cpb.PromptId
// 	dm.ValueType = cpb.ValueType.String()
// 	return dm
// }

// func (dm *AIAgent) ConvertToPb() *mpb.AIAgent {
// 	if dm == nil {
// 		return nil
// 	}
// 	return &mpb.AIAgent{
// 		AgentId:      dm.VapusID,
// 		Name:         dm.Name,
// 		AgentType:    mpb.VapusAiAgentTypes(mpb.VapusAiAgentTypes_value[dm.AgentType]),
// 		Description:  dm.Description,
// 		AgentVersion: dm.AgentVersion,
// 		Owners:       dm.Editors,
// 		Labels:       dm.Labels,
// 		Status:       dm.Status,
// 		Organization:       dm.Organization,
// 		ResourceBase: dm.ConvertToPbBase(),
// 		AiModelMap: func() []*mpb.AIModelMap {
// 			var pb []*mpb.AIModelMap
// 			for _, v := range dm.AIModelMap {
// 				pb = append(pb, v.ConvertToPb())
// 			}
// 			return pb

// 		}(),
// 		Editable: dm.Editable,
// 		Settings: dm.Settings.ConvertToPb(),
// 		Steps: func() []*mpb.Steps {
// 			var pb []*mpb.Steps
// 			for _, v := range dm.Steps {
// 				pb = append(pb, v.ConvertToPb())
// 			}
// 			return pb
// 		}(),
// 	}
// }

// func (dm *AIAgent) ConvertToListingPb() *mpb.AIAgent {
// 	if dm == nil {
// 		return nil
// 	}
// 	return &mpb.AIAgent{
// 		AgentId:      dm.VapusID,
// 		Name:         dm.Name,
// 		AgentType:    mpb.VapusAiAgentTypes(mpb.VapusAiAgentTypes_value[dm.AgentType]),
// 		AgentVersion: dm.AgentVersion,
// 		Owners:       dm.Editors,
// 		Labels:       dm.Labels,
// 		Status:       dm.Status,
// 		Organization:       dm.Organization,
// 	}
// }

// func (dm *AIAgent) ConvertFromPb(cpb *mpb.AIAgent) *AIAgent {
// 	if cpb == nil {
// 		return nil
// 	}
// 	dm.Name = cpb.Name
// 	dm.AgentType = cpb.AgentType.String()
// 	dm.Description = cpb.Description
// 	dm.AgentVersion = cpb.AgentVersion
// 	dm.Editors = cpb.Owners
// 	dm.Editable = cpb.Editable
// 	dm.Labels = cpb.Labels
// 	dm.AIModelMap = func() []*AIModelMap {
// 		var pb []*AIModelMap
// 		for _, v := range cpb.AiModelMap {
// 			pb = append(pb, (&AIModelMap{}).ConvertFromPb(v))
// 		}
// 		return pb
// 	}()
// 	dm.Steps = func() []*Steps {
// 		var pb []*Steps
// 		for _, v := range cpb.Steps {
// 			pb = append(pb, (&Steps{}).ConvertFromPb(v))
// 		}
// 		return pb
// 	}()
// 	return dm
// }

// func (dm *AIAgent) PreSaveCreate(authzClaim map[string]string) {
// 	if dm == nil {
// 		return
// 	}
// 	dm.PreSaveVapusBase(authzClaim)
// }

// func (dn *AIAgent) PreSaveUpdate(userId string) {
// 	if dn == nil {
// 		return
// 	}
// 	dn.UpdatedBy = userId
// 	dn.UpdatedAt = dmutils.GetEpochTime()
// }

// func (dn *AIAgent) PreSaveDelete(authzClaim map[string]string) {
// 	if dn == nil {
// 		return
// 	}
// 	dn.PreDeleteVapusBase(authzClaim)
// }

// type AIAgentRunMD struct {
// 	StartedAt int64  `json:"startedAt,omitempty" yaml:"startedAt,omitempty" toml:"startedAt,omitempty"`
// 	EndedAt   int64  `json:"endedAt,omitempty" yaml:"endedAt,omitempty" toml:"endedAt,omitempty"`
// 	RunId     string `json:"runId,omitempty" yaml:"runId,omitempty" toml:"runId,omitempty"`
// }

// type AIAgentThread struct {
// 	VapusBase `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
// 	// AgentId   string            `bun:"agent_id" json:"agentId,omitempty" yaml:"agentId,omitempty" toml:"agentId,omitempty"`
// 	ThreadId string `bun:"thread_id" json:"threadId,omitempty" yaml:"threadId,omitempty" toml:"threadId,omitempty"`
// 	// AgentType string            `bun:"agent_type" json:"agentType,omitempty" yaml:"agentType,omitempty" toml:"agentType,omitempty"`
// 	Error    string              `bun:"error" json:"error,omitempty" yaml:"error,omitempty" toml:"error,omitempty"`
// 	Status   string              `bun:"status" json:"status,omitempty" yaml:"status,omitempty" toml:"status,omitempty"`
// 	Log      []*AIAgentThreadLog `bun:"log,type:jsonb" json:"log,omitempty" yaml:"log,omitempty" toml:"log,omitempty"`
// 	Schedule *VapusSchedule      `bun:"schedule,type:jsonb" json:"schedule,omitempty" yaml:"schedule,omitempty" toml:"schedule,omitempty"`
// }

// type AIAgentThreadLog struct {
// 	Steps      string        `json:"steps,omitempty" yaml:"steps,omitempty" toml:"steps,omitempty"`
// 	EOLReason  string        `json:"eolReason,omitempty" yaml:"eolReason,omitempty" toml:"eolReason,omitempty"`
// 	FailedStep string        `json:"failedStep,omitempty" yaml:"failedStep,omitempty" toml:"failedStep,omitempty"`
// 	AgentId    string        `json:"agentId,omitempty" yaml:"agentId,omitempty" toml:"agentId,omitempty"`
// 	AgentType  string        `json:"agentType,omitempty" yaml:"agentType,omitempty" toml:"agentType,omitempty"`
// 	Metadata   *AIAgentRunMD `json:"metadata,omitempty" yaml:"metadata,omitempty" toml:"metadata,omitempty"`
// }

// func (dm *AIAgentThread) PreSaveCreate(authzClaim map[string]string) {
// 	dm.PreSaveVapusBase(authzClaim)
// }

// func (dm *AIAgentThread) PreSaveUpdate(userId string) {
// 	if dm == nil {
// 		return
// 	}
// 	dm.UpdatedBy = userId
// 	dm.UpdatedAt = dmutils.GetEpochTime()
// }
