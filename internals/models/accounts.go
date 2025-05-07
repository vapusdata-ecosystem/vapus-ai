package models

import (
	fmt "fmt"

	guuid "github.com/google/uuid"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	types "github.com/vapusdata-ecosystem/vapusai/core/types"
)

type Account struct {
	VapusBase            `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name                 string               `bun:"name,unique" json:"name"`
	Status               string               `bun:"status" json:"status"`
	AuthnMethod          string               `bun:"authn_method" json:"authnMethod"`
	DmAccessJwtKeys      *JWTParams           `bun:"dm_access_jwt_keys,type:jsonb" json:"dmAccessJwtKeys"`
	BackendSecretStorage *BackendStorages     `bun:"backend_secret_storage,type:jsonb" json:"backendSecretStorage"`
	BackendDataStorage   *BackendStorages     `bun:"backend_data_storage,type:jsonb" json:"backendDataStorage"`
	ArtifactStorage      *BackendStorages     `bun:"artifact_storage,type:jsonb" json:"artifactStorage"`
	AuthnParams          *AuthnOIDC           `bun:"authn_params,type:jsonb" json:"authnParams"`
	Profile              *AccountProfile      `bun:"profile,type:jsonb" json:"profile"`
	AIAttributes         *AccountAIAttributes `bun:"ai_attributes,type:jsonb" json:"aiAttributes"`
	Settings             *AccountSettings     `bun:"settings,type:jsonb" json:"settings"`
}

func (a *Account) ConvertToPb() *mpb.Account {
	if a == nil {
		return nil
	}
	obj := &mpb.Account{
		Name:                 a.Name,
		Status:               a.Status,
		AuthnMethod:          mpb.AuthnMethod(mpb.AuthnMethod_value[a.AuthnMethod]),
		DmAccessJwtKeys:      a.DmAccessJwtKeys.ConvertToPb(),
		BackendSecretStorage: a.BackendSecretStorage.ConvertToPb(),
		BackendDataStorage:   a.BackendDataStorage.ConvertToPb(),
		ArtifactStorage:      a.ArtifactStorage.ConvertToPb(),
		OidcParams:           a.AuthnParams.ConvertToPb(),
		AccountId:            a.VapusID,
		AiAttributes:         a.AIAttributes.ConvertToPb(),
		Profile:              a.Profile.ConvertToPb(),
		Settings:             a.Settings.ConvertToPb(),
	}
	return obj
}

func (a *Account) ConvertFromPb(pb *mpb.Account) *Account {
	if a == nil {
		return nil
	}
	obj := &Account{
		Name:                 pb.GetName(),
		Status:               pb.GetStatus(),
		AuthnMethod:          pb.GetAuthnMethod().String(),
		DmAccessJwtKeys:      (&JWTParams{}).ConvertFromPb(pb.DmAccessJwtKeys),
		BackendSecretStorage: (&BackendStorages{}).ConvertFromPb(pb.BackendSecretStorage),
		BackendDataStorage:   (&BackendStorages{}).ConvertFromPb(pb.BackendDataStorage),
		ArtifactStorage:      (&BackendStorages{}).ConvertFromPb(pb.ArtifactStorage),
		AuthnParams:          (&AuthnOIDC{}).ConvertFromPb(pb.GetOidcParams()),
		AIAttributes:         (&AccountAIAttributes{}).ConvertFromPb(pb.GetAiAttributes()),
		Profile:              (&AccountProfile{}).ConvertFromPb(pb.Profile),
		Settings:             (&AccountSettings{}).ConvertFromPb(pb.Settings),
	}
	return obj
}

func (dm *Account) GetAiAttributes() *AccountAIAttributes {
	if dm == nil {
		return nil
	}
	return dm.AIAttributes
}

func (a *Account) GetName() string {
	if a == nil {
		return ""
	}
	return a.Name
}

func (a *Account) GetStatus() string {
	if a == nil {
		return ""
	}
	return a.Status
}

func (a *Account) GetAuthnMethod() string {
	if a == nil {
		return ""
	}
	return a.AuthnMethod
}

func (a *Account) GetDmAccessJwtKeys() *JWTParams {
	if a == nil {
		return &JWTParams{}
	}
	return a.DmAccessJwtKeys
}

func (a *Account) GetBackendSecretStorage() *BackendStorages {
	if a == nil {
		return &BackendStorages{}
	}
	return a.BackendSecretStorage
}

func (a *Account) GetBackendDataStorage() *BackendStorages {
	if a == nil {
		return &BackendStorages{}
	}
	return a.BackendDataStorage
}

func (a *Account) GetArtifactStorage() *BackendStorages {
	if a == nil {
		return &BackendStorages{}
	}
	return a.ArtifactStorage
}

func (a *Account) GetAuthnParams() *AuthnOIDC {
	if a == nil {
		return &AuthnOIDC{}
	}
	return a.AuthnParams
}

func (dm *Account) SetAccountId() {
	if dm == nil {
		return
	}
	if dm.VapusID == types.EMPTYSTR {
		dm.VapusID = fmt.Sprintf(types.ACCOUNTIDTEMPLATE, guuid.New())
	}
}

func (dm *Account) PreSaveCreate(userId string) {
	if dm == nil {
		return
	}
	if dm.CreatedBy == types.EMPTYSTR {
		dm.CreatedBy = userId
	}
	if dm.CreatedAt == 0 {
		dm.CreatedAt = dmutils.GetEpochTime()
	}

	if dm.Status == mpb.AccountStatus_VAS_INVALID.String() {
		dm.Status = mpb.AccountStatus_VAS_ACTIVE.String()
	}
}

func (dn *Account) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *Account) Delete(userId string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = userId
	dn.DeletedAt = dmutils.GetEpochTime()
}

type AccountProfile struct {
	Addresses   []*Address `json:"addresses"`
	Logo        string     `json:"logo"`
	Description string     `json:"description"`
	Moto        string     `json:"moto"`
	Favicon     string     `json:"favicon"`
}

func (a *AccountProfile) GetAddress() []*Address {
	if a == nil {
		return nil
	}
	return a.Addresses
}

func (a *AccountProfile) GetLogo() string {
	if a == nil {
		return ""
	}
	return a.Logo
}

func (a *AccountProfile) GetDescription() string {
	if a == nil {
		return ""
	}
	return a.Description
}

func (a *AccountProfile) GetMoto() string {
	if a == nil {
		return ""
	}
	return a.Moto
}

func (x *AccountProfile) ConvertFromPb(a *mpb.AccountProfile) *AccountProfile {
	if x == nil {
		return nil
	}
	obj := &AccountProfile{
		Logo:        a.GetLogo(),
		Description: a.GetDescription(),
		Moto:        a.GetMoto(),
		Favicon:     a.GetFavicon(),
		Addresses: func() []*Address {
			var address []*Address
			for _, d := range a.GetAddresses() {
				address = append(address, (&Address{}).ConvertFromPb(d))
			}
			return address
		}(),
	}
	return obj
}

func (a *AccountProfile) ConvertToPb() *mpb.AccountProfile {
	if a == nil {
		return nil
	}
	obj := &mpb.AccountProfile{
		Logo:        a.Logo,
		Description: a.Description,
		Moto:        a.Moto,
		Favicon:     a.Favicon,
	}
	for _, d := range a.Addresses {
		obj.Addresses = append(obj.Addresses, d.ConvertToPb())
	}
	return obj
}

type AccountAIAttributes struct {
	EmbeddingModelNode  string `json:"embeddingModelNode,omitempty" yaml:"embeddingModelNode"`
	EmbeddingModel      string `json:"embeddingModel,omitempty" yaml:"embeddingModel"`
	GenerativeModelNode string `json:"generativeModelNode,omitempty" yaml:"generativeModelNode"`
	GenerativeModel     string `json:"generativeModel,omitempty" yaml:"generativeModel"`
	GuardrailModelNode  string `json:"guardrailModelNode,omitempty" yaml:"guardrailModelNode"`
	GuardrailModel      string `json:"guardrailModel,omitempty" yaml:"guardrailModel"`
}

func (a *AccountAIAttributes) ConvertFromPb(pb *mpb.AIAttributes) *AccountAIAttributes {
	if a == nil {
		return nil
	}
	obj := &AccountAIAttributes{
		EmbeddingModelNode:  pb.GetEmbeddingModelNode(),
		EmbeddingModel:      pb.GetEmbeddingModel(),
		GenerativeModelNode: pb.GetGenerativeModelNode(),
		GenerativeModel:     pb.GetGenerativeModel(),
		GuardrailModelNode:  pb.GetGuardrailModelNode(),
		GuardrailModel:      pb.GetGuardrailModel(),
	}
	return obj
}

func (a *AccountAIAttributes) ConvertToPb() *mpb.AIAttributes {
	if a == nil {
		return nil
	}
	obj := &mpb.AIAttributes{
		EmbeddingModelNode:  a.EmbeddingModelNode,
		EmbeddingModel:      a.EmbeddingModel,
		GenerativeModelNode: a.GenerativeModelNode,
		GenerativeModel:     a.GenerativeModel,
		GuardrailModelNode:  a.GuardrailModelNode,
		GuardrailModel:      a.GuardrailModel,
	}
	return obj
}

type AccountSettings struct {
	GoogleAnalyticsTagId string `json:"googleAnalyticsTagId" yaml:"googleAnalyticsTagId"`
}

func (a *AccountSettings) ConvertFromPb(pb *mpb.AccountSettings) *AccountSettings {
	if a == nil {
		return nil
	}
	obj := &AccountSettings{
		GoogleAnalyticsTagId: pb.GetGoogleAnalyticsTagId(),
	}
	return obj
}

func (a *AccountSettings) ConvertToPb() *mpb.AccountSettings {
	if a == nil {
		return nil
	}
	obj := &mpb.AccountSettings{
		GoogleAnalyticsTagId: a.GoogleAnalyticsTagId,
	}
	return obj
}
