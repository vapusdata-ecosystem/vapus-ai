package models

import (
	"fmt"
	"reflect"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
)

type GuardModels struct {
	ModelID     string `json:"modelId" yaml:"modelId"`
	ModelNodeID string `json:"modelNodeId" yaml:"modelNodeId"`
}

type AIGuardrails struct {
	VapusBase          `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name               string                     `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	MinConfidence      float64                    `bun:"min_confidence" json:"minConfidence,omitempty" yaml:"minConfidence,omitempty" toml:"minConfidence,omitempty"`
	MaxConfidence      float64                    `bun:"max_confidence" json:"maxConfidence,omitempty" yaml:"maxConfidence,omitempty" toml:"maxConfidence,omitempty"`
	Contents           *ContentGuardrailLevel     `bun:"contents,type:jsonb" json:"contents,omitempty" yaml:"contents,omitempty" toml:"contents,omitempty"`
	Topics             []*TopicGuardrails         `bun:"topics,type:jsonb" json:"topics,omitempty" yaml:"topics,omitempty" toml:"topics,omitempty"`
	Words              []*WordGuardRails          `bun:"words,type:jsonb" json:"words,omitempty" yaml:"words,omitempty" toml:"words,omitempty"`
	SensitiveDataset   []*SensitiveDataGuardrails `bun:"sensitive_dataset,type:jsonb" json:"sensitiveDataset,omitempty" yaml:"sensitiveDataset,omitempty" toml:"sensitiveDataset,omitempty"`
	Description        string                     `bun:"description" json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	FailureMessage     string                     `bun:"failure_message" json:"failureMessage,omitempty" yaml:"failureMessage,omitempty" toml:"failureMessage,omitempty"`
	DisplayName        string                     `bun:"display_name" json:"displayName,omitempty" yaml:"displayName,omitempty" toml:"displayName,omitempty"`
	Schema             string                     `bun:"schema" json:"schema,omitempty" yaml:"schema,omitempty" toml:"schema,omitempty"`
	ScanMode           string                     `bun:"scan_mode" json:"scanMode,omitempty" yaml:"scanMode,omitempty" toml:"scanMode,omitempty"`
	GuardModel         *GuardModels               `bun:"guard_model,type:jsonb" json:"guardModel,omitempty" yaml:"guardModel,omitempty" toml:"guardModel,omitempty"`
	EligibleModelNodes []string                   `bun:"eligible_model_nodes,array" json:"eligibleModelNodes,omitempty" yaml:"eligibleModelNodes,omitempty" toml:"eligibleModelNodes,omitempty"`
	Partner            []*ThirdParty              `bun:"partner,type:jsonb" json:"partner,omitempty" yaml:"partner,omitempty" toml:"partner,omitempty"`
}

func (dm *AIGuardrails) PreSaveCreate(authzClaim map[string]string) {
	if dm == nil {
		return
	}
	dm.PreSaveVapusBase(authzClaim)
}

func (dn *AIGuardrails) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *AIGuardrails) PreSaveDelete(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreDeleteVapusBase(authzClaim)
}

func (a *AIGuardrails) ConvertToPb() *mpb.AIGuardrails {
	obj := &mpb.AIGuardrails{
		Name:               a.Name,
		MinConfidence:      a.MinConfidence,
		MaxConfidence:      a.MaxConfidence,
		Description:        a.Description,
		FailureMessage:     a.FailureMessage,
		DisplayName:        a.DisplayName,
		Contents:           a.Contents.ConvertToPb(),
		Topics:             make([]*mpb.TopicGuardrails, len(a.Topics)),
		Words:              make([]*mpb.WordGuardRails, len(a.Words)),
		SensitiveDataset:   make([]*mpb.SensitiveDataGuardrails, len(a.SensitiveDataset)),
		GuardrailId:        a.VapusID,
		Schema:             a.Schema,
		ScanMode:           mpb.AIGuardrailScanMode(mpb.AIGuardrailScanMode_value[a.ScanMode]),
		EligibleModelNodes: a.EligibleModelNodes,
		Partner:            make([]*mpb.ThirdParty, len(a.Partner)),

		ResourceBase: a.ConvertToPbBase(),
	}
	if a.GuardModel != nil {
		obj.GuardModel = &mpb.GuardModels{
			ModelId:     a.GuardModel.ModelID,
			ModelNodeId: a.GuardModel.ModelNodeID,
		}
	}
	for i := range a.Topics {
		obj.Topics[i] = a.Topics[i].ConvertToPb()
	}
	for i := range a.Words {
		obj.Words[i] = a.Words[i].ConvertToPb()
	}
	for i := range a.SensitiveDataset {
		obj.SensitiveDataset[i] = a.SensitiveDataset[i].ConvertToPb()
	}
	for i := range a.Partner {
		obj.Partner[i] = a.Partner[i].ConvertToPb()
	}
	return obj
}

func (a *AIGuardrails) ConvertToListingPb() *mpb.AIGuardrails {
	obj := &mpb.AIGuardrails{
		Name:         a.Name,
		Description:  a.Description,
		DisplayName:  a.DisplayName,
		ScanMode:     mpb.AIGuardrailScanMode(mpb.AIGuardrailScanMode_value[a.ScanMode]),
		GuardrailId:  a.VapusID,
		ResourceBase: a.ConvertToPbBase(),
	}
	return obj
}

func (a *AIGuardrails) ConvertFromPb(obpb *mpb.AIGuardrails) *AIGuardrails {
	a.Name = obpb.Name
	a.MinConfidence = obpb.MinConfidence
	a.MaxConfidence = obpb.MaxConfidence
	a.Description = obpb.Description
	a.FailureMessage = obpb.FailureMessage
	a.DisplayName = obpb.DisplayName
	a.Contents = new(ContentGuardrailLevel).ConvertFromPb(obpb.Contents)
	a.Topics = make([]*TopicGuardrails, len(obpb.Topics))
	a.Words = make([]*WordGuardRails, len(obpb.Words))
	a.VapusBase = VapusBase{}
	//a.Editors = obpb.ResourceBase.Owners
	a.Scope = obpb.ResourceBase.Scope.String()
	a.SensitiveDataset = make([]*SensitiveDataGuardrails, len(obpb.SensitiveDataset))
	a.ScanMode = obpb.ScanMode.String()
	if obpb.GuardModel != nil {
		a.GuardModel = &GuardModels{
			ModelID:     obpb.GuardModel.ModelId,
			ModelNodeID: obpb.GuardModel.ModelNodeId,
		}
	}
	a.EligibleModelNodes = obpb.EligibleModelNodes
	a.Partner = make([]*ThirdParty, len(obpb.Partner))
	for i := range obpb.Topics {
		a.Topics[i] = new(TopicGuardrails).ConvertFromPb(obpb.Topics[i])
	}
	for i := range obpb.Words {
		a.Words[i] = new(WordGuardRails).ConvertFromPb(obpb.Words[i])
	}
	for i := range obpb.SensitiveDataset {
		a.SensitiveDataset[i] = new(SensitiveDataGuardrails).ConvertFromPb(obpb.SensitiveDataset[i])
	}
	fmt.Println("obpb.Partner: ", reflect.ValueOf(obpb.Partner))
	for i := range obpb.Partner {
		a.Partner[i] = new(ThirdParty).ConvertFromPb(obpb.Partner[i])
	}
	return a
}

type ContentGuardrailLevel struct {
	HateSpeech string `json:"hateSpeech,omitempty" yaml:"hateSpeech,omitempty" toml:"hateSpeech,omitempty"`
	Insults    string `json:"insults,omitempty" yaml:"insults,omitempty" toml:"insults,omitempty"`
	Sexual     string `json:"sexual,omitempty" yaml:"sexual,omitempty" toml:"sexual,omitempty"`
	Threats    string `json:"threats,omitempty" yaml:"threats,omitempty" toml:"threats,omitempty"`
	Misconduct string `json:"misconduct,omitempty" yaml:"misconduct,omitempty" toml:"misconduct,omitempty"`
}

func (c *ContentGuardrailLevel) ConvertToPb() *mpb.ContentGuardrailLevel {
	return &mpb.ContentGuardrailLevel{
		HateSpeech: mpb.GuardRailLevels(mpb.GuardRailLevels_value[c.HateSpeech]),
		Insults:    mpb.GuardRailLevels(mpb.GuardRailLevels_value[c.Insults]),
		Sexual:     mpb.GuardRailLevels(mpb.GuardRailLevels_value[c.Sexual]),
		Threats:    mpb.GuardRailLevels(mpb.GuardRailLevels_value[c.Threats]),
		Misconduct: mpb.GuardRailLevels(mpb.GuardRailLevels_value[c.Misconduct]),
	}
}

func (c *ContentGuardrailLevel) ConvertFromPb(pb *mpb.ContentGuardrailLevel) *ContentGuardrailLevel {
	c.HateSpeech = pb.HateSpeech.String()
	c.Insults = pb.Insults.String()
	c.Sexual = pb.Sexual.String()
	c.Threats = pb.Threats.String()
	c.Misconduct = pb.Misconduct.String()
	return c
}

type TopicGuardrails struct {
	Topic       string   `json:"topic,omitempty" yaml:"topic,omitempty" toml:"topic,omitempty"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Samples     []string `json:"samples,omitempty" yaml:"samples,omitempty" toml:"samples,omitempty"`
}

func (t *TopicGuardrails) ConvertToPb() *mpb.TopicGuardrails {
	return &mpb.TopicGuardrails{
		Topic:       t.Topic,
		Description: t.Description,
		Samples:     t.Samples,
	}
}

func (t *TopicGuardrails) ConvertFromPb(pb *mpb.TopicGuardrails) *TopicGuardrails {
	t.Topic = pb.Topic
	t.Description = pb.Description
	t.Samples = pb.Samples
	return t
}

type WordGuardRails struct {
	Words        []string `json:"words,omitempty" yaml:"words,omitempty" toml:"words,omitempty"`
	FileLocation string   `json:"fileLocation,omitempty" yaml:"fileLocation,omitempty" toml:"fileLocation,omitempty"`
}

func (w *WordGuardRails) ConvertToPb() *mpb.WordGuardRails {
	return &mpb.WordGuardRails{
		Words:        w.Words,
		FileLocation: w.FileLocation,
	}
}

func (w *WordGuardRails) ConvertFromPb(pb *mpb.WordGuardRails) *WordGuardRails {
	w.Words = pb.Words
	w.FileLocation = pb.FileLocation
	return w
}

type SensitiveDataGuardrails struct {
	PIIType string `json:"piiType,omitempty" yaml:"piiType,omitempty" toml:"piiType,omitempty"`
	Action  string `json:"action,omitempty" yaml:"action,omitempty" toml:"action,omitempty"`
	Regex   string `json:"regex,omitempty" yaml:"regex,omitempty" toml:"regex,omitempty"`
}

func (s *SensitiveDataGuardrails) ConvertToPb() *mpb.SensitiveDataGuardrails {
	return &mpb.SensitiveDataGuardrails{
		PiiType: s.PIIType,
		Action:  s.Action,
		Regex:   s.Regex,
	}
}

func (s *SensitiveDataGuardrails) ConvertFromPb(pb *mpb.SensitiveDataGuardrails) *SensitiveDataGuardrails {
	s.PIIType = pb.PiiType
	s.Action = pb.Action
	s.Regex = pb.Regex
	return s
}

type ThirdParty struct {
	Bedrock []*BedrockGuardrailModel    `json:"bedrock,omitempty" yaml:"bedrock,omitempty" toml:"bedrock,omitempty"`
	Mistral []*ThirdPartyGuardrailModel `json:"mistral,omitempty" yaml:"mistral,omitempty" toml:"mistral,omitempty"`
	Pangea  []*ThirdPartyGuardrailModel `json:"pangea,omitempty" yaml:"pangea,omitempty" toml:"pangea,omitempty"`
	// Nemo    []*FileData                 `json:"nemo,omitempty" yaml:"nemo,omitempty" toml:"nemo,omitempty"`
}

func (a *ThirdParty) ConvertToPb() *mpb.ThirdParty {
	obj := &mpb.ThirdParty{}

	for i := range a.Bedrock {
		obj.Bedrock[i] = a.Bedrock[i].ConvertToPb()
	}
	for i := range a.Mistral {
		obj.Mistral[i] = a.Mistral[i].ConvertToPb()
	}
	for i := range a.Pangea {
		obj.Pangea[i] = a.Pangea[i].ConvertToPb()
	}
	return obj
}

func (a *ThirdParty) ConvertFromPb(obpb *mpb.ThirdParty) *ThirdParty {
	for i := range obpb.Bedrock {
		a.Bedrock[i] = new(BedrockGuardrailModel).ConvertFromPb(obpb.Bedrock[i])
	}
	for i := range obpb.Mistral {
		a.Mistral[i] = new(ThirdPartyGuardrailModel).ConvertFromPb(obpb.Mistral[i])
	}
	for i := range obpb.Bedrock {
		a.Pangea[i] = new(ThirdPartyGuardrailModel).ConvertFromPb(obpb.Pangea[i])
	}
	return a
}

type BedrockGuardrailModel struct {
	Arn  string `json:"arn,omitempty" yaml:"arn,omitempty" toml:"arn,omitempty"`
	Id   string `json:"id,omitempty" yaml:"id,omitempty" toml:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
}

func (a *BedrockGuardrailModel) ConvertToPb() *mpb.BedrockGuardrailModel {
	return &mpb.BedrockGuardrailModel{
		Arn:  a.Arn,
		Id:   a.Id,
		Name: a.Name,
	}
}

func (a *BedrockGuardrailModel) ConvertFromPb(obpb *mpb.BedrockGuardrailModel) *BedrockGuardrailModel {
	a.Arn = obpb.Arn
	a.Id = obpb.Id
	a.Name = obpb.Name
	return a
}

type ThirdPartyGuardrailModel struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Id   string `json:"id,omitempty" yaml:"id,omitempty" toml:"id,omitempty"`
}

func (a *ThirdPartyGuardrailModel) ConvertToPb() *mpb.ThirdPartyGuardrailModel {
	return &mpb.ThirdPartyGuardrailModel{
		Id:   a.Id,
		Name: a.Name,
	}
}

func (a *ThirdPartyGuardrailModel) ConvertFromPb(obpb *mpb.ThirdPartyGuardrailModel) *ThirdPartyGuardrailModel {
	a.Id = obpb.Id
	a.Name = obpb.Name
	return a
}

// type FileData struct {
// 	Name          string            `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
// 	Data          []byte            `json:"data,omitempty" yaml:"data,omitempty" toml:"data,omitempty"`
// 	ContentFormat string            `json:"contentFormat,omitempty" yaml:"contentFormat,omitempty" toml:"contentFormat,omitempty"`
// 	Path          string            `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
// 	Eof           bool              `json:"eof,omitempty" yaml:"eof,omitempty" toml:"eof,omitempty"`
// 	Params        map[string]string `json:"params,omitempty" yaml:"params,omitempty" toml:"params,omitempty"`
// 	Description   string            `json:"nadescriptionme,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
// 	RedirectUrl   string            `json:"redirectUrl,omitempty" yaml:"redirectUrl,omitempty" toml:"redirectUrl,omitempty"`
// }
