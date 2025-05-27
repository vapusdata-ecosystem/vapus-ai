package models

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type Organization struct {
	VapusBase                `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name                     string            `bun:"name,notnull,unique" json:"name,omitempty" yaml:"name"`
	DisplayName              string            `bun:"display_name" json:"displayName,omitempty" yaml:"displayName"`
	Users                    []string          `bun:"users,array" json:"users,omitempty" yaml:"users"`
	SecretPasscode           string            `bun:"salt_val" json:"saltVal,omitempty" yaml:"saltVal"`
	AuthnJwtParams           *JWTParams        `bun:"authn_jwt_params,type:jsonb" json:"authnJwtParams,omitempty" yaml:"authnJwtParams"`
	OrganizationType         string            `bun:"organization_type" json:"organizationType,omitempty" yaml:"organizationType"`
	BackendSecretStorage     *BackendStorages  `bun:"backend_secret_storage,type:jsonb" json:"backendSecretStorage,omitempty" yaml:"backendSecretStorage"`
	ArtifactStorage          *BackendStorages  `bun:"artifact_storage,type:jsonb" json:"artifactStorage,omitempty" yaml:"artifactStorage"`
	DataProductInfraPlatform []*K8SInfraParams `bun:"data_product_infra_platform,type:jsonb" json:"dataProductInfraPlatform,omitempty" yaml:"dataProductInfraPlatform"`
}

func (m *Organization) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (d *Organization) GetName() string {
	return d.Name
}

func (d *Organization) GetDisplayName() string {
	return d.DisplayName
}

func (d *Organization) GetOrganizationId() string {
	return d.VapusID
}

func (d *Organization) GetUsers() []string {
	return d.Users
}

func (d *Organization) GetSecretPasscode() string {
	return d.SecretPasscode
}

func (d *Organization) GetStatus() string {
	return d.Status
}

func (d *Organization) GetOrganizationType() string {
	return d.OrganizationType
}

func (d *Organization) GetBackendSecretStorage() *BackendStorages {
	return d.BackendSecretStorage
}

func (d *Organization) GetArtifactStorage() *BackendStorages {
	return d.ArtifactStorage
}

func (d *Organization) GetDataProductInfraPlatform() []*K8SInfraParams {
	return d.DataProductInfraPlatform
}

func (d *Organization) GetAuthnJwtParams() *JWTParams {
	return d.AuthnJwtParams
}

func (d *Organization) HasArtifactStore() bool {
	if d != nil {
		if d.ArtifactStorage != nil && d.ArtifactStorage.NetParams != nil && d.ArtifactStorage.NetParams.Address != "" && len(d.ArtifactStorage.NetParams.DsCreds) > 0 {
			return true
		}
	}
	return false
}

func (dmn *Organization) ConvertToPb() *mpb.Organization {
	if dmn != nil {
		obj := &mpb.Organization{
			Name:                     dmn.Name,
			DisplayName:              dmn.DisplayName,
			OrganizationId:           dmn.VapusID,
			Users:                    dmn.Users,
			SecretPasscode:           &mpb.CredentialSalt{SaltVal: dmn.SecretPasscode},
			Status:                   dmn.Status,
			OrganizationType:         mpb.OrganizationType(mpb.OrganizationType_value[dmn.OrganizationType]),
			BackendSecretStorage:     dmn.BackendSecretStorage.ConvertToPb(),
			ArtifactStorage:          dmn.ArtifactStorage.ConvertToPb(),
			DataProductInfraPlatform: make([]*mpb.K8SInfraParams, 0),
			ResourceBase:             dmn.ConvertToPbBase(),
			Attributes: &mpb.OrganizationAttributes{
				AuthnJwtParams: dmn.AuthnJwtParams.ConvertToPb(),
			},
		}
		for _, ds := range dmn.DataProductInfraPlatform {
			obj.DataProductInfraPlatform = append(obj.DataProductInfraPlatform, ds.ConvertToPb())
		}
		return obj
	}
	return nil
}

func (dmn *Organization) ConvertToListingPb() *mpb.Organization {
	if dmn != nil {
		obj := &mpb.Organization{
			Name:             dmn.Name,
			DisplayName:      dmn.DisplayName,
			OrganizationId:   dmn.VapusID,
			Status:           dmn.Status,
			OrganizationType: mpb.OrganizationType(mpb.OrganizationType_value[dmn.OrganizationType]),
			ResourceBase:     dmn.ConvertToPbBase(),
		}
		return obj
	}
	return nil
}

func (dmn *Organization) ConvertFromPb(pb *mpb.Organization) *Organization {
	if pb == nil {
		return nil
	}
	obj := &Organization{
		Name:                     pb.GetName(),
		DisplayName:              pb.GetDisplayName(),
		Users:                    pb.GetUsers(),
		SecretPasscode:           pb.GetSecretPasscode().GetSaltVal(),
		OrganizationType:         mpb.OrganizationType_name[int32(pb.GetOrganizationType())],
		BackendSecretStorage:     (&BackendStorages{}).ConvertFromPb(pb.GetBackendSecretStorage()),
		ArtifactStorage:          (&BackendStorages{}).ConvertFromPb(pb.GetArtifactStorage()),
		DataProductInfraPlatform: make([]*K8SInfraParams, 0),
		AuthnJwtParams:           (&JWTParams{}).ConvertFromPb(pb.GetAttributes().GetAuthnJwtParams()),
	}
	for _, ds := range pb.GetDataProductInfraPlatform() {
		obj.DataProductInfraPlatform = append(obj.DataProductInfraPlatform, (&K8SInfraParams{}).ConvertFromPb(ds))
	}
	return obj
}

func (dmn *Organization) SetOrganizationId() {
	if dmn == nil {
		return
	}
	if dmn.VapusID == "" {
		dmn.VapusID = dmutils.GetUUID()
	}
}

func (dmn *Organization) PreSaveCreate(authzClaim map[string]string) {
	if dmn == nil {
		return
	}
	if dmn.CreatedBy == types.EMPTYSTR {
		dmn.CreatedBy = authzClaim[encryption.ClaimUserIdKey]
	}
	if dmn.CreatedAt == 0 {
		dmn.CreatedAt = dmutils.GetEpochTime()
	}
	if dmn.DisplayName == "" {
		dmn.DisplayName = dmn.Name
	}
	dmn.OwnerAccount = authzClaim[encryption.ClaimAccountKey]
}

func (dn *Organization) PreSaveUpdate(userId string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = userId
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *Organization) Delete(userId string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = userId
	dn.DeletedAt = dmutils.GetEpochTime()
}

func (dn *Organization) GetK8sInfra(id string) *K8SInfraParams {
	if dn != nil {
		for _, infra := range dn.DataProductInfraPlatform {
			if id != "" {
				if infra.InfraId == id {
					return infra
				}
			} else {
				return infra
			}
		}
	}
	return nil
}
