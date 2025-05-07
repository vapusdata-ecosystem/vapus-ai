package models

import dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"

type FileStoreLog struct {
	VapusBase  `bun:",embed" json:"base,omitempty" yaml:"base,omitempty" toml:"base,omitempty"`
	Name       string `bun:"name" json:"name" yaml:"name" toml:"name"`
	Path       string `bun:"path" json:"path" yaml:"path" toml:"path"`
	Format     string `bun:"format" json:"format" yaml:"format" toml:"format"`
	Size       int64  `bun:"size" json:"size" yaml:"size" toml:"size"`
	Checksum   string `bun:"checksum" json:"checksum" yaml:"checksum" toml:"checksum"`
	Resource   string `bun:"resource" json:"resource" yaml:"resource" toml:"resource"`
	ResourceId string `bun:"resource_id" json:"resource_id" yaml:"resource_id" toml:"resource_id"`
	DupCounter int64  `bun:"dup_counter" json:"dup_counter" yaml:"dup_counter" toml:"dup_counter"`
}

func (dn *FileStoreLog) PreSaveCreate(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.PreSaveVapusBase(authzClaim)
}

func (dn *FileStoreLog) Delete(userId string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = userId
	dn.DeletedAt = dmutils.GetEpochTime()
}

type FileStoreCache struct {
	Path      string   `json:"path" yaml:"path" toml:"path"`
	Name      string   `json:"name" yaml:"name" toml:"name"`
	Checksums []string `json:"checksums" yaml:"checksums" toml:"checksums"`
	Counter   int64    `json:"counter" yaml:"counter" toml:"counter"`
}
