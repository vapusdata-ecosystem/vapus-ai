package models

import (
	fmt "fmt"

	guuid "github.com/google/uuid"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type AIModelNode struct {
	VapusBase             `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name                  string                    `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name"`
	GenerativeModels      []*AIModelBase            `bun:"generative_models,type:jsonb" json:"generativeModels,omitempty" yaml:"generativeModels"`
	EmbeddingModels       []*AIModelBase            `bun:"embedding_models,type:jsonb" json:"embeddingModels,omitempty" yaml:"embeddingModels"`
	DiscoverModels        bool                      `bun:"discover_models" json:"discoverModels,omitempty" yaml:"discoverModels"`
	NetworkParams         *AIModelNodeNetworkParams `bun:"network_params,type:jsonb" json:"networkParams,omitempty" yaml:"networkParams"`
	ApprovedOrganizations []string                  `bun:"approved_Organizations,array" json:"approvedOrganizations,omitempty" yaml:"approvedOrganizations"`
	Hosting               string                    `bun:"hosting" json:"hosting,omitempty" yaml:"hosting"`
	ServiceProvider       string                    `bun:"service_provider" json:"serviceProvider,omitempty" yaml:"serviceProvider"`
	SecurityGuardrails    *SecurityGuardrails       `bun:"security_guardrails,type:jsonb" json:"securityGuardrails,omitempty" yaml:"securityGuardrails"`
}

func (m *AIModelNode) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (m *AIModelNode) ConvertToPb() *mpb.AIModelNode {
	if m != nil {
		obj := &mpb.AIModelNode{
			ModelNodeId:        m.VapusID,
			Name:               m.Name,
			NodeOwners:         m.Editors,
			SecurityGuardrails: m.SecurityGuardrails.ConvertToPb(),
			ResourceBase:       m.ConvertToPbBase(),
			Attributes: &mpb.AIModelNodeAttributes{
				DiscoverModels:        m.DiscoverModels,
				GenerativeModels:      make([]*mpb.AIModelBase, 0),
				EmbeddingModels:       make([]*mpb.AIModelBase, 0),
				NetworkParams:         m.NetworkParams.ConvertToPb(),
				Scope:                 mpb.ResourceScope(mpb.ResourceScope_value[m.Scope]),
				ApprovedOrganizations: m.ApprovedOrganizations,
				Hosting:               mpb.AIModelNodeHosting(mpb.AIModelNodeHosting_value[m.Hosting]),
				ServiceProvider:       mpb.ServiceProvider(mpb.ServiceProvider_value[m.ServiceProvider]),
			},
		}
		for _, gm := range m.GenerativeModels {
			obj.Attributes.GenerativeModels = append(obj.Attributes.GenerativeModels, gm.ConvertToPb())
		}
		for _, em := range m.EmbeddingModels {
			obj.Attributes.EmbeddingModels = append(obj.Attributes.EmbeddingModels, em.ConvertToPb())
		}
		return obj
	}
	return nil
}

func (m *AIModelNode) ConvertToListingPb() *mpb.AIModelNode {
	if m != nil {
		obj := &mpb.AIModelNode{
			ModelNodeId:        m.VapusID,
			Name:               m.Name,
			NodeOwners:         m.Editors,
			SecurityGuardrails: m.SecurityGuardrails.ConvertToPb(),
			ResourceBase:       m.ConvertToPbBase(),
			Attributes: &mpb.AIModelNodeAttributes{
				DiscoverModels:        m.DiscoverModels,
				GenerativeModels:      make([]*mpb.AIModelBase, 0),
				EmbeddingModels:       make([]*mpb.AIModelBase, 0),
				NetworkParams:         m.NetworkParams.ConvertToPb(),
				Scope:                 mpb.ResourceScope(mpb.ResourceScope_value[m.Scope]),
				ApprovedOrganizations: m.ApprovedOrganizations,
				Hosting:               mpb.AIModelNodeHosting(mpb.AIModelNodeHosting_value[m.Hosting]),
				ServiceProvider:       mpb.ServiceProvider(mpb.ServiceProvider_value[m.ServiceProvider]),
			},
		}
		for _, gm := range m.GenerativeModels {
			obj.Attributes.GenerativeModels = append(obj.Attributes.GenerativeModels, gm.ConvertToPb())
		}
		for _, em := range m.EmbeddingModels {
			obj.Attributes.EmbeddingModels = append(obj.Attributes.EmbeddingModels, em.ConvertToPb())
		}
		return obj
	}
	return nil
}

func (m *AIModelNode) ConvertFromPb(pb *mpb.AIModelNode) *AIModelNode {
	if pb == nil {
		return nil
	}
	obj := &AIModelNode{
		Name:                  pb.GetName(),
		DiscoverModels:        pb.GetAttributes().GetDiscoverModels(),
		GenerativeModels:      make([]*AIModelBase, 0),
		EmbeddingModels:       make([]*AIModelBase, 0),
		NetworkParams:         (&AIModelNodeNetworkParams{}).ConvertFromPb(pb.GetAttributes().GetNetworkParams()),
		ApprovedOrganizations: pb.GetAttributes().GetApprovedOrganizations(),
		Hosting:               pb.GetAttributes().GetHosting().String(),
		ServiceProvider:       pb.GetAttributes().GetServiceProvider().String(),
		SecurityGuardrails:    (&SecurityGuardrails{}).ConvertFromPb(pb.GetSecurityGuardrails()),
	}
	obj.Scope = pb.GetAttributes().GetScope().String()
	obj.Editors = pb.GetNodeOwners()
	for _, gm := range pb.GetAttributes().GetGenerativeModels() {
		obj.GenerativeModels = append(obj.GenerativeModels, (&AIModelBase{}).ConvertFromPb(gm))
	}
	for _, em := range pb.GetAttributes().GetEmbeddingModels() {
		obj.EmbeddingModels = append(obj.EmbeddingModels, (&AIModelBase{}).ConvertFromPb(em))
	}
	return obj
}

func (n *AIModelNode) GetModelNodeId() string {
	if n != nil {
		return n.VapusID
	}
	return ""
}

func (n *AIModelNode) GetName() string {
	if n != nil {
		return n.Name
	}
	return ""
}

func (n *AIModelNode) GetNodeOwners() []string {
	if n != nil {
		return n.Editors
	}
	return nil
}

func (n *AIModelNode) GetStatus() string {
	if n != nil {
		return n.Status
	}
	return ""
}

func (n *AIModelNode) GetGenerativeModels() []*AIModelBase {
	if n != nil {
		return n.GenerativeModels
	}
	return nil
}

func (n *AIModelNode) GetEmbeddingModels() []*AIModelBase {
	if n != nil {
		return n.EmbeddingModels
	}
	return nil
}

func (n *AIModelNode) GetDiscoverModels() bool {
	if n != nil {
		return n.DiscoverModels
	}
	return false
}

func (n *AIModelNode) GetNetworkParams() *AIModelNodeNetworkParams {
	if n != nil {
		return n.NetworkParams
	}
	return nil
}

func (n *AIModelNode) GetScope() string {
	if n != nil {
		return n.Scope
	}
	return ""
}

func (n *AIModelNode) GetApprovedOrganizations() []string {
	if n != nil {
		return n.ApprovedOrganizations
	}
	return nil
}

func (n *AIModelNode) GetHosting() string {
	if n != nil {
		return n.Hosting
	}
	return ""
}

func (n *AIModelNode) GetServiceProvider() string {
	if n != nil {
		return n.ServiceProvider
	}
	return ""
}

type AIModelNodeNetworkParams struct {
	Url                 string                  `json:"url,omitempty" yaml:"url"`
	ApiVersion          string                  `json:"apiVersion,omitempty" yaml:"apiVersion"`
	LocalPath           string                  `json:"localPath,omitempty" yaml:"localPath"`
	Credentials         *GenericCredentialModel `json:"credentials,omitempty" yaml:"credentials"`
	SecretName          string                  `json:"secretName,omitempty" yaml:"secretName"`
	IsAlreadyInSecretBs bool                    `json:"isAlreadyInSecretBS,omitempty" yaml:"isAlreadyInSecretBS"`
	Params              map[string]any
}

func (n *AIModelNodeNetworkParams) GetUrl() string {
	if n != nil {
		return n.Url
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetApiVersion() string {
	if n != nil {
		return n.ApiVersion
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetLocalPath() string {
	if n != nil {
		return n.LocalPath
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetCredentials() *GenericCredentialModel {
	if n != nil {
		return n.Credentials
	}
	return nil
}

func (n *AIModelNodeNetworkParams) GetSecretName() string {
	if n != nil {
		return n.SecretName
	}
	return ""
}

func (n *AIModelNodeNetworkParams) GetIsAlreadyInSecretBs() bool {
	if n != nil {
		return n.IsAlreadyInSecretBs
	}
	return false
}

func (m *AIModelNodeNetworkParams) ConvertToPb() *mpb.AIModelNodeNetworkParams {
	if m != nil {
		return &mpb.AIModelNodeNetworkParams{
			Url:                 m.Url,
			ApiVersion:          m.ApiVersion,
			LocalPath:           m.LocalPath,
			Credentials:         m.Credentials.ConvertToPb(),
			SecretName:          m.SecretName,
			IsAlreadyInSecretBs: m.IsAlreadyInSecretBs,
		}
	}
	return nil
}

func (m *AIModelNodeNetworkParams) ConvertFromPb(pb *mpb.AIModelNodeNetworkParams) *AIModelNodeNetworkParams {
	if pb == nil {
		return nil
	}
	return &AIModelNodeNetworkParams{
		Url:                 pb.GetUrl(),
		ApiVersion:          pb.GetApiVersion(),
		LocalPath:           pb.GetLocalPath(),
		Credentials:         (&GenericCredentialModel{}).ConvertFromPb(pb.GetCredentials()),
		SecretName:          pb.GetSecretName(),
		IsAlreadyInSecretBs: pb.GetIsAlreadyInSecretBs(),
	}
}

type AIModelBase struct {
	ModelName        string   `json:"modelName,omitempty" yaml:"modelName,omitempty"`
	ModelId          string   `json:"modelId,omitempty" yaml:"modelId,omitempty"`
	ModelType        string   `json:"modelType,omitempty" yaml:"modelType,omitempty"`
	OwnedBy          string   `json:"ownedBy,omitempty" yaml:"ownedBy,omitempty"`
	InputTokenLimit  int32    `json:"inputTokenLimit,omitempty" yaml:"inputTokenLimit,omitempty"`
	OutputTokenLimit int32    `json:"outputTokenLimit,omitempty" yaml:"outputTokenLimit,omitempty"`
	SupprtedOps      []string `json:"supportedOps,omitempty" yaml:"supportedOps,omitempty"`
	Version          string   `json:"version,omitempty" yaml:"version,omitempty"`
	ModelNature      []string `json:"modelNature,omitempty" yaml:"modelNature,omitempty"`
	ModelArn         string   `json:"modelArn,omitempty" yaml:"modelArn,omitempty"`
}

func (m *AIModelBase) ConvertToPb() *mpb.AIModelBase {
	if m != nil {
		obj := &mpb.AIModelBase{
			ModelName:        m.ModelName,
			ModelId:          m.ModelId,
			ModelType:        mpb.AIModelType(mpb.AIModelType_value[m.ModelType]),
			OwnedBy:          m.OwnedBy,
			InputTokenLimit:  m.InputTokenLimit,
			OutputTokenLimit: m.OutputTokenLimit,
			ModelNature:      make([]mpb.AIModelType, 0),
			ModelArn:         m.ModelArn,
		}
		for _, modelNature := range m.ModelNature {
			obj.ModelNature = append(obj.ModelNature, mpb.AIModelType(mpb.AIModelType_value[modelNature]))
		}
		return obj
	}
	return nil
}

func (m *AIModelBase) ConvertFromPb(pb *mpb.AIModelBase) *AIModelBase {
	if pb == nil {
		return nil
	}
	obj := &AIModelBase{
		ModelName:        pb.GetModelName(),
		ModelId:          pb.GetModelId(),
		ModelType:        pb.GetModelType().String(),
		OwnedBy:          pb.GetOwnedBy(),
		InputTokenLimit:  pb.GetInputTokenLimit(),
		OutputTokenLimit: pb.GetOutputTokenLimit(),
		ModelArn:         pb.GetModelArn(),
		ModelNature:      make([]string, 0),
	}
	for _, modelNature := range pb.GetModelNature() {
		obj.ModelNature = append(obj.ModelNature, modelNature.Enum().String())
	}
	return obj
}

func (dn *AIModelNode) SetAINodeId() {
	if dn == nil {
		return
	}
	if dn.VapusID == "" {
		dn.VapusID = fmt.Sprintf(types.VAPUS_AIMODEL_NODE_ID, guuid.New())
	}
}

func (dn *AIModelNode) PreSaveCreate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreSaveVapusBase(authzClaim)
}

func (dn *AIModelNode) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *AIModelNode) PreSaveDelete(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreDeleteVapusBase(authzClaim)
}

func (dn *AIModelNode) GetCredentials(param string) *GenericCredentialModel {
	if dn == nil {
		return nil
	}
	return dn.NetworkParams.Credentials
}

func (dn *AIModelNode) Delete(userId string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = userId
	dn.DeletedAt = dmutils.GetEpochTime()
}
