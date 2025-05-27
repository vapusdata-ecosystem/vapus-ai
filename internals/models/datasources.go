package models

import (
	"encoding/json"
	fmt "fmt"
	"strings"

	guuid "github.com/google/uuid"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"gopkg.in/yaml.v3"
)

type DataSource struct {
	VapusBase `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name      string               `bun:"name" json:"name"`
	NetParams *DataSourceNetParams `bun:"net_params,type:jsonb" json:"netParams"`
	// Attributes                   *mpb.DataSourceAttributes `bun:"attributes,type:jsonb" json:"attributes"`
	Sharable        bool                     `bun:"sharable" json:"sharable"`
	DataSourceType  string                   `bun:"data_source_type" json:"dataSourceType"`
	Tags            []*Mapper                `bun:"tags,type:jsonb" json:"tags"`
	SharingParams   *DataSourceSharingParams `bun:"sharing_params,type:jsonb" json:"sharingParams"`
	StorageEngine   string                   `bun:"storage_engine" json:"storageEngine"`
	ServiceName     string                   `bun:"service_name" json:"serviceName"`
	ServiceProvider string                   `bun:"service_provider" json:"serviceProvider"`
	Goal            string                   `bun:"goal" json:"goal"`
}

func (x *DataSource) GetFormatedObject(fileFormat string) string {
	if x != nil {
		switch strings.ToLower(fileFormat) {
		case strings.ToLower(mpb.ContentFormats_YAML.String()):
			yamlData, err := yaml.Marshal(x)
			if err != nil {
				return ""
			}
			return string(yamlData)
		case strings.ToLower(mpb.ContentFormats_JSON.String()):
			jsonData, err := json.Marshal(x)
			if err != nil {
				return ""
			}
			return string(jsonData)
		default:
			yamlData, err := yaml.Marshal(x)
			if err != nil {
				return ""
			}
			return string(yamlData)
		}
	}
	return ""
}

func (m *DataSource) SetAccountId(accountId string) {
	if m != nil {
		m.OwnerAccount = accountId
	}
}

func (ds *DataSource) GetName() string {
	if ds != nil {
		return ds.Name
	}
	return ""
}

func (ds *DataSource) GetNetParams() *DataSourceNetParams {
	if ds != nil {
		return ds.NetParams
	}
	return nil
}

// func (ds *DataSource) GetAttributes() *mpb.DataSourceAttributes {
// 	if ds != nil {
// 		return ds.Attributes
// 	}
// 	return nil
// }

func (ds *DataSource) GetDataSourceId() string {
	if ds != nil {
		return ds.VapusID
	}
	return ""
}

func (ds *DataSource) GetOwnerOrganization() string {
	if ds != nil {
		return ds.Organization
	}
	return ""
}

func (ds *DataSource) GetOwners() []string {
	if ds != nil {
		return ds.Editors
	}
	return nil
}

func (ds *DataSource) GetSharable() bool {
	if ds != nil {
		return ds.Sharable
	}
	return false
}

func (ds *DataSource) GetDataSourceType() string {
	if ds != nil {
		return ds.DataSourceType
	}
	return ""
}

func (ds *DataSource) GetStatus() string {
	if ds != nil {
		return ds.Status
	}
	return ""
}

func (ds *DataSource) GetTags() []*Mapper {
	if ds != nil {
		return ds.Tags
	}
	return nil
}

func (ds *DataSource) GetSharingParams() *DataSourceSharingParams {
	if ds != nil {
		return ds.SharingParams
	}
	return nil
}

func (ds *DataSource) GetStorageEngine() string {
	if ds != nil {
		return ds.StorageEngine
	}
	return ""
}

func (ds *DataSource) GetServiceName() string {
	if ds != nil {
		return ds.ServiceName
	}
	return ""
}

func (ds *DataSource) GetServiceProvider() string {
	if ds != nil {
		return ds.ServiceProvider
	}
	return ""
}

func (ds *DataSource) GetGoal() string {
	if ds != nil {
		return ds.Goal
	}
	return ""
}

func (j *DataSource) ConvertToPb() *mpb.DataSource {
	if j != nil {
		strEngine, ok := types.StorageEngineMap[mpb.DataSourceServices(mpb.DataSourceServices_value[j.ServiceName])]
		if !ok {
			strEngine = ""
		}
		obj := &mpb.DataSource{
			Name:      j.Name,
			NetParams: j.NetParams.ConvertToPb(),
			Attributes: &mpb.DataSourceAttributes{
				StorageEngine:   strEngine.String(),
				ServiceName:     mpb.DataSourceServices(mpb.DataSourceServices_value[j.ServiceName]),
				ServiceProvider: mpb.ServiceProvider(mpb.ServiceProvider_value[j.ServiceProvider]),
				// Goal:            mpb.DataSourceGoal(mpb.DataSourceGoal_value[j.Goal]),
				// SupportedArtifactTypes: j.Attributes.SupportedArtifactTypes,
			},
			DataSourceId:   j.VapusID,
			Owners:         j.Editors,
			Sharable:       j.Sharable,
			DataSourceType: mpb.DataSourceType(mpb.DataSourceType_value[j.DataSourceType]),
			Status:         j.Status,
			Tags:           make([]*mpb.Mapper, 0),
			SharingParams:  j.SharingParams.ConvertToPb(),
			ResourceBase:   j.ConvertToPbBase(),
		}
		for _, c := range j.Tags {
			obj.Tags = append(obj.Tags, c.ConvertToPb())
		}
		return obj
	}
	return nil
}

func (j *DataSource) ConvertToListingPb() *mpb.DataSource {
	if j != nil {
		obj := &mpb.DataSource{
			Name: j.Name,
			Attributes: &mpb.DataSourceAttributes{
				StorageEngine:   j.StorageEngine,
				ServiceName:     mpb.DataSourceServices(mpb.DataSourceServices_value[j.ServiceName]),
				ServiceProvider: mpb.ServiceProvider(mpb.ServiceProvider_value[j.ServiceProvider]),
				// Goal:            mpb.DataSourceGoal(mpb.DataSourceGoal_value[j.Goal]),
			},
			DataSourceId: j.VapusID,
			Owners:       j.Editors,
			// Sharable:                     j.Sharable,
			DataSourceType: mpb.DataSourceType(mpb.DataSourceType_value[j.DataSourceType]),
			ResourceBase:   j.ConvertToPbBase(),
		}
		return obj
	}
	return nil
}

func (j *DataSource) ConvertFromPb(pb *mpb.DataSource) *DataSource {
	if pb == nil {
		return nil
	}
	va := &DataSource{
		Name:      pb.GetName(),
		NetParams: (&DataSourceNetParams{}).ConvertFromPb(pb.GetNetParams()),
		// Attributes:                   pb.GetAttributes(),
		Sharable:       pb.GetSharable(),
		DataSourceType: mpb.DataSourceType_name[int32(pb.GetDataSourceType())],
		Tags:           make([]*Mapper, 0),
		SharingParams:  (&DataSourceSharingParams{}).ConvertFromPb(pb.GetSharingParams()),
		ServiceName:    pb.Attributes.ServiceName.String(),
		StorageEngine:  pb.Attributes.StorageEngine,
	}
	va.Editors = pb.GetOwners()
	return va
}

type DataSourceNetParams struct {
	Address          string             `json:"address,omitempty" yaml:"address"`
	Port             int32              `json:"port,omitempty" yaml:"port"`
	Databases        []string           `json:"databases,omitempty" yaml:"databases"`
	DsCreds          []*DataSourceCreds `json:"dsCreds,omitempty" yaml:"dsCreds"`
	DatabasePrefixes []string           `json:"databasePrefixes,omitempty" yaml:"databasePrefixes"`
	Version          string             `json:"version,omitempty" yaml:"version"`
}

func (j *DataSourceNetParams) ConvertToPb() *mpb.DataSourceNetParams {
	if j != nil {
		obj := &mpb.DataSourceNetParams{
			Address:          j.Address,
			Port:             j.Port,
			Databases:        j.Databases,
			DsCreds:          make([]*mpb.DataSourceCreds, 0),
			DatabasePrefixes: j.DatabasePrefixes,
			Version:          j.Version,
		}
		for _, c := range j.DsCreds {
			obj.DsCreds = append(obj.DsCreds, c.ConvertToPb())
		}
		return obj
	}
	return nil
}

func (j *DataSourceNetParams) ConvertFromPb(pb *mpb.DataSourceNetParams) *DataSourceNetParams {
	if pb == nil {
		return nil
	}
	return &DataSourceNetParams{
		Address:          pb.GetAddress(),
		Port:             pb.GetPort(),
		Databases:        pb.GetDatabases(),
		DatabasePrefixes: pb.GetDatabasePrefixes(),
		Version:          pb.GetVersion(),
		DsCreds: func() []*DataSourceCreds {
			var creds []*DataSourceCreds
			for _, c := range pb.GetDsCreds() {
				creds = append(creds, (&DataSourceCreds{}).ConvertFromPb(c))
			}
			return creds
		}(),
	}
}

type DataSourceCreds struct {
	Name                string                  `json:"name,omitempty" yaml:"name"`
	IsAlreadyInSecretBs bool                    `json:"isAlreadyInSecretBS,omitempty" yaml:"isAlreadyInSecretBS"`
	Credentials         *GenericCredentialModel `json:"credentials,omitempty" yaml:"credentials"`
	Priority            int32                   `json:"priority,omitempty" yaml:"priority"`
	AccessScope         string                  `json:"accessScope,omitempty" yaml:"accessScope"`
	DB                  string                  `json:"db,omitempty" yaml:"db"`
	SecretName          string                  `json:"secretName,omitempty" yaml:"secretName"`
}

func (j *DataSourceCreds) ConvertToPb() *mpb.DataSourceCreds {
	if j != nil {
		return &mpb.DataSourceCreds{
			Name:                j.Name,
			IsAlreadyInSecretBs: j.IsAlreadyInSecretBs,
			Credentials:         j.Credentials.ConvertToPb(),
			Priority:            j.Priority,
			AccessScope:         mpb.DataSourceAccessScope(mpb.DataSourceAccessScope_value[j.AccessScope]),
			Db:                  j.DB,
			SecretName:          j.SecretName,
		}
	}
	return nil
}

func (j *DataSourceCreds) ConvertFromPb(pb *mpb.DataSourceCreds) *DataSourceCreds {
	if pb == nil {
		return nil
	}
	return &DataSourceCreds{
		Name:                pb.GetName(),
		IsAlreadyInSecretBs: pb.GetIsAlreadyInSecretBs(),
		Credentials:         (&GenericCredentialModel{}).ConvertFromPb(pb.GetCredentials()),
		Priority:            pb.GetPriority(),
		AccessScope:         mpb.DataSourceAccessScope_name[int32(pb.GetAccessScope())],
		DB:                  pb.GetDb(),
		SecretName:          pb.GetSecretName(),
	}
}

type DataSourceSharingParams struct {
	OrganizationId string `json:"organizationId" yaml:"organizationId"`
	AccessScope    string `json:"accessScope" yaml:"accessScope"`
	ValidFrom      int32  `json:"validFrom" yaml:"validFrom"`
	ValidTill      int32  `json:"validTill" yaml:"validTill"`
}

func (j *DataSourceSharingParams) ConvertToPb() *mpb.DataSourceSharingParams {
	if j != nil {
		return &mpb.DataSourceSharingParams{
			OrganizationId: j.OrganizationId,
			AccessScope:    mpb.DataSourceAccessScope(mpb.DataSourceAccessScope_value[j.AccessScope]),
			ValidFrom:      j.ValidFrom,
			ValidTill:      j.ValidTill,
		}
	}
	return nil
}

func (j *DataSourceSharingParams) ConvertFromPb(pb *mpb.DataSourceSharingParams) *DataSourceSharingParams {
	if pb == nil {
		return nil
	}
	return &DataSourceSharingParams{
		OrganizationId: pb.GetOrganizationId(),
		AccessScope:    mpb.DataSourceAccessScope_name[int32(pb.GetAccessScope())],
		ValidFrom:      pb.GetValidFrom(),
		ValidTill:      pb.GetValidTill(),
	}
}

func (dn *DataSource) SetDataSourceUuid() {
	if dn == nil {
		return
	}
	if dn.VapusID == "" {
		dn.VapusID = fmt.Sprintf(types.DATA_SOURCEID_TEMPLATE, guuid.New())
	}
}

func (dn *DataSource) PreSaveCreate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreSaveVapusBase(authzClaim)
}

func (dn *DataSource) PreSaveUpdate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.UpdatedBy = authzClaim[encryption.ClaimUserIdKey]
	dn.UpdatedAt = dmutils.GetEpochTime()
}

func (dn *DataSource) PreSaveDelete(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.VapusBase.PreDeleteVapusBase(authzClaim)
}

func (dn *DataSource) GetCredentials(as string, hp bool, database string) (*DataSourceCreds, error) {
	if dn == nil {
		return nil, fmt.Errorf("nil value in the pointer refrences")
	}
	if as == "" && len(dn.NetParams.DsCreds) > 0 {
		if database != "" {
			for _, c := range dn.NetParams.DsCreds {
				if c.DB == database {
					return c, nil
				}
			}
		}
		return dn.NetParams.DsCreds[0], nil
	}
	for _, c := range dn.NetParams.DsCreds {
		if database != "" && c.DB == database {
			return c, nil
		}
		if c.AccessScope == as {
			return c, nil
		}
	}
	return nil, dmerrors.DMError(fmt.Errorf("no credentials found for the access scope %s", as), ErrDataSourceCredsNotFound)
}
