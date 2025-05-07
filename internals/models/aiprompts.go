package models

import (
	fmt "fmt"

	guuid "github.com/google/uuid"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	types "github.com/vapusdata-ecosystem/vapusdata/core/types"
)

type AIPrompt struct {
	VapusBase       `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name            string      `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name"`
	PromptTypes     []string    `bun:"prompt_type,array" json:"promptType,omitempty" yaml:"promptType"`
	PreferredModels []string    `bun:"preferred_models,array" json:"preferredModels,omitempty" yaml:"preferredModels"`
	Editable        bool        `bun:"editable" json:"editable,omitempty" yaml:"editable" default:"true"`
	Spec            *PromptSpec `bun:"spec,type:jsonb" json:"spec,omitempty" yaml:"spec"`
	Labels          []string    `bun:"labels,array" json:"labels,omitempty" yaml:"labels"`
	SpecDigest      *DigestVal  `bun:"spec_digest,type:jsonb" json:"specDigest,omitempty" yaml:"specDigest"`
	UserTemplate    string      `bun:"user_template" json:"userTemplate,omitempty" yaml:"userTemplate"`
	Template        string      `bun:"template" json:"template,omitempty" yaml:"template"`
}

func (m *AIPrompt) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (d *AIPrompt) SetPromptId() {
	if d == nil {
		return
	}
	d.VapusID = fmt.Sprintf(types.PROMPT_ID, guuid.New())
}

func (d *AIPrompt) PreSaveCreate(authzClaim map[string]string) {
	if d == nil {
		return
	}
	d.PreSaveVapusBase(authzClaim)
}

func (dn *AIPrompt) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *AIPrompt) PreSaveDelete(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreDeleteVapusBase(authzClaim)
}

func (dn *AIPrompt) ConvertFromPb(pb *mpb.AIPrompt) *AIPrompt {
	if pb == nil {
		return nil
	}
	vv := &AIPrompt{
		Name:            pb.GetName(),
		PromptTypes:     pb.GetPromptTypes(),
		PreferredModels: pb.GetPreferredModels(),
		Editable:        pb.GetEditable(), // PromptScope:     mpb.ResourceScope(mpb.ResourceScope_value[pb.GetPromptScope()]),
		Spec:            (&PromptSpec{}).ConvertFromPb(pb.GetSpec()),
		Labels:          pb.GetLabels(),
	}
	vv.Scope = pb.GetPromptScope().String()
	return vv
}

func (dn *AIPrompt) ConvertToPb() *mpb.AIPrompt {
	if dn == nil {
		return nil
	}
	return &mpb.AIPrompt{
		Name:            dn.Name,
		PromptTypes:     dn.PromptTypes,
		PreferredModels: dn.PreferredModels,
		Editable:        dn.Editable,
		PromptScope:     mpb.ResourceScope(mpb.ResourceScope_value[dn.Scope]),
		Spec:            dn.Spec.ConvertToPb(),
		PromptId:        dn.VapusID,
		PromptOwner:     dn.CreatedBy,
		Labels:          dn.Labels,
		UserTemplate:    dn.UserTemplate,
		Template:        dn.Template,
		ResourceBase:    dn.ConvertToPbBase(),
	}
}

func (dn *AIPrompt) ConvertToListingPb() *mpb.AIPrompt {
	if dn == nil {
		return nil
	}
	return &mpb.AIPrompt{
		Name:         dn.Name,
		PromptTypes:  dn.PromptTypes,
		PromptId:     dn.VapusID,
		PromptOwner:  dn.CreatedBy,
		Labels:       dn.Labels,
		ResourceBase: dn.ConvertToPbBase(),
	}
}

type ToolPrompts struct {
	Schema       *FunctionCall `json:"schema,omitempty" yaml:"schema,omitempty"`
	AutoGenerate bool          `json:"autoGenerate,omitempty" yaml:"autoGenerate,omitempty"`
	RawJsonParam string        `json:"rawJsonParam,omitempty" yaml:"rawJsonParam,omitempty"`
	Type         string        `json:"type,omitempty" yaml:"type,omitempty"`
}

func (d *ToolPrompts) ConvertFromPb(pbTool *mpb.ToolPrompts) *ToolPrompts {
	if pbTool == nil {
		return nil
	}

	return &ToolPrompts{
		Schema:       (&FunctionCall{}).ConvertFromPb(pbTool.GetSchema()),
		AutoGenerate: pbTool.GetAutoGenerate(),
		RawJsonParam: pbTool.GetRawJsonParams(),
		Type:         pbTool.GetType(),
	}
}

func (d *ToolPrompts) ConvertToPb() *mpb.ToolPrompts {
	if d == nil {
		return nil
	}
	return &mpb.ToolPrompts{
		Schema:        d.Schema.ConvertToPb(),
		AutoGenerate:  d.AutoGenerate,
		RawJsonParams: d.RawJsonParam,
		Type:          d.Type,
	}
}

type PromptSpec struct {
	SystemMessage  string                    `json:"systemMessage,omitempty" yaml:"systemMessage,omitempty"`
	UserMessage    string                    `json:"userMessage,omitempty" yaml:"userMessage,omitempty"`
	Tools          []*ToolPrompts            `json:"tools,omitempty" yaml:"tools,omitempty"`
	InputTag       string                    `json:"inputTag,omitempty" yaml:"inputTag,omitempty"`
	OutputTag      string                    `json:"outputTag,omitempty" yaml:"outputTag,omitempty"`
	ContextTag     string                    `json:"contextTag,omitempty" yaml:"contextTag,omitempty"`
	Sample         *Sample                   `json:"sample,omitempty" yaml:"sample,omitempty"`
	ResponseFormat *StructuredResponseFormat `json:"responseFormat,omitempty" yaml:"responseFormat,omitempty"`
	Variables      []string                  `json:"variables,omitempty" yaml:"variables,omitempty"`
}

func (d *PromptSpec) ConvertFromPb(pbPrompt *mpb.PromptSpec) *PromptSpec {
	if pbPrompt == nil {
		return nil
	}
	return &PromptSpec{
		SystemMessage: pbPrompt.GetSystemMessage(),
		UserMessage:   pbPrompt.GetUserMessage(),
		Tools: func() []*ToolPrompts {
			var tools []*ToolPrompts
			for _, tool := range pbPrompt.GetTools() {
				tools = append(tools, (&ToolPrompts{}).ConvertFromPb(tool))
			}
			return tools
		}(),
		InputTag:       pbPrompt.GetInputTag(),
		OutputTag:      pbPrompt.GetOutputTag(),
		ContextTag:     pbPrompt.GetContextTag(),
		Sample:         (&Sample{}).ConvertFromPb(pbPrompt.GetSample()),
		ResponseFormat: (&StructuredResponseFormat{}).ConvertFromPb(pbPrompt.GetResponseFormat()),
		Variables:      pbPrompt.GetVariables(),
	}
}

func (d *PromptSpec) ConvertToPb() *mpb.PromptSpec {
	if d == nil {
		return nil
	}
	return &mpb.PromptSpec{
		SystemMessage: d.SystemMessage,
		UserMessage:   d.UserMessage,
		Tools: func() []*mpb.ToolPrompts {
			var tools []*mpb.ToolPrompts
			for _, tool := range d.Tools {
				tools = append(tools, tool.ConvertToPb())
			}
			return tools
		}(),
		InputTag:       d.InputTag,
		OutputTag:      d.OutputTag,
		ContextTag:     d.ContextTag,
		Sample:         d.Sample.ConvertToPb(),
		ResponseFormat: d.ResponseFormat.ConvertToPb(),
		Variables:      d.Variables,
	}
}

type Sample struct {
	InputText string `json:"inputText,omitempty" yaml:"inputText,omitempty"`
	Response  string `json:"response,omitempty" yaml:"response,omitempty"`
}

func (d *Sample) ConvertFromPb(pbSample *mpb.Sample) *Sample {
	if pbSample == nil {
		return nil
	}
	return &Sample{
		InputText: pbSample.GetInputText(),
		Response:  pbSample.GetResponse(),
	}
}

func (d *Sample) ConvertToPb() *mpb.Sample {
	if d == nil {
		return nil
	}
	return &mpb.Sample{
		InputText: d.InputText,
		Response:  d.Response,
	}
}

type PromptTag struct {
	Start string `json:"start,omitempty" yaml:"start,omitempty"`
	End   string `json:"end,omitempty" yaml:"end,omitempty"`
}

func (d *PromptTag) ConvertFromPb(pbTag *mpb.PromptTag) *PromptTag {
	if pbTag == nil {
		return nil
	}
	return &PromptTag{
		Start: pbTag.GetStart(),
		End:   pbTag.GetEnd(),
	}
}

func (d *PromptTag) ConvertToPb() *mpb.PromptTag {
	if d == nil {
		return nil
	}
	return &mpb.PromptTag{
		Start: d.Start,
		End:   d.End,
	}
}
