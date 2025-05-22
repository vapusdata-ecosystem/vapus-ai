package appconfigs

import (
	filepath "path/filepath"
)

type PlatformBootConfig struct {
	PlatformOwners  []string `yaml:"platformOwners" validate:"required"`
	PlatformAccount struct {
		Name    string `yaml:"name" validate:"required"`
		Creator string `yaml:"creator" validate:"required"`
	} `yaml:"platformAccount" validate:"required"`
	PlatformAccountOrganization struct {
		Name string `yaml:"name" `
	} `yaml:"platformAccountOrganization"`
}

type BaseOs struct {
	ArtifactType       string   `yaml:"artifactType"`
	URL                string   `yaml:"url"`
	Digest             string   `yaml:"digest"`
	Tag                string   `yaml:"tag"`
	OrganizationMounts []string `yaml:"OrganizationMounts"`
}

type LocalFSPaths struct {
	OrganizationFiles string `yaml:"OrganizationFiles"`
	DataSourceFiles   string `yaml:"dataSourceFiles"`
}

type VapusAISvcConfig struct {
	Path                 string
	VapusBESecretStorage *SecretMap           `yaml:"vapusBESecretStorage"`
	VapusBEDbStorage     *SecretMap           `yaml:"vapusBEDbStorage"`
	VapusBECacheStorage  *SecretMap           `yaml:"vapusBECacheStorage"`
	VapusFileStorage     *SecretMap           `yaml:"vapusFileStorage"`
	NetworkConfigFile    string               `yaml:"networkConfigFile"`
	JWTAuthnSecrets      *SecretMap           `yaml:"JWTAuthnSecrets"`
	AuthnSecrets         *SecretMap           `yaml:"authnSecrets"`
	ArtifactStore        *SecretMap           `yaml:"artifactStore"`
	ServerCerts          *TlsCertConfig       `yaml:"serverCerts"`
	TrinoSpecs           *TrinoDeploymentSpec `yaml:"trinoSpecs"`
	BaseOs               []*BaseOs            `yaml:"baseOs"`
	PlatformBaseAccount  *PlatformBootConfig  `yaml:"platformBaseAccount"`
}

func (sc *VapusAISvcConfig) GetFileStorePath() string {
	return sc.VapusFileStorage.Secret
}

func (sc *VapusAISvcConfig) GetSecretStoragePath() string {
	return filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath)
}

func (sc *VapusAISvcConfig) GetDBStoragePath() string {
	return sc.VapusBEDbStorage.Secret
}

func (sc *VapusAISvcConfig) GetCachStoragePath() string {
	return sc.VapusBECacheStorage.Secret
}

func (sc *VapusAISvcConfig) GetJwtAuthSecretPath() string {
	return sc.JWTAuthnSecrets.Secret
}

func (sc *VapusAISvcConfig) GetAuthnSecrets() string {
	return sc.AuthnSecrets.Secret
}

func (sc *VapusAISvcConfig) GetArtifactStorePath() string {
	return sc.ArtifactStore.Secret
}

func (sc *VapusAISvcConfig) GetMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *VapusAISvcConfig) GetPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *VapusAISvcConfig) GetCaCert() string {
	return sc.ServerCerts.CaCertFile
}

func (sc *VapusAISvcConfig) GetClientMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientKeyFile
}

func (sc *VapusAISvcConfig) GetClientPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientCertFile
}

func (sc *VapusAISvcConfig) GetBaseOs() []*BaseOs {
	return sc.BaseOs
}

func (sc *VapusAISvcConfig) GetPlatformBaseAccount() *PlatformBootConfig {
	return sc.PlatformBaseAccount
}
