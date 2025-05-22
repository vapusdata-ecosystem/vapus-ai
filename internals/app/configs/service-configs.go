package appconfigs

import (
	"log"
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

type PlatformServiceConfig struct {
	Path                 string
	VapusBESecretStorage *SecretMap     `yaml:"vapusBESecretStorage"`
	VapusBEDbStorage     *SecretMap     `yaml:"vapusBEDbStorage"`
	VapusBECacheStorage  *SecretMap     `yaml:"vapusBECacheStorage"`
	VapusFileStorage     *SecretMap     `yaml:"vapusFileStorage"`
	NetworkConfigFile    string         `yaml:"networkConfigFile"`
	JWTAuthnSecrets      *SecretMap     `yaml:"JWTAuthnSecrets"`
	AuthnSecrets         *SecretMap     `yaml:"authnSecrets"`
	ArtifactStore        *SecretMap     `yaml:"artifactStore"`
	ServerCerts          *TlsCertConfig `yaml:"serverCerts"`
	PbacConfig           struct {
		FilePath string `yaml:"filePath"`
	} `yaml:"pbacConfig"`
	AuthnMethod         string               `yaml:"authnMethod"`
	BaseOs              []*BaseOs            `yaml:"baseOs"`
	PlatformBaseAccount *PlatformBootConfig  `yaml:"platformBaseAccount"`
	TrinoSpecs          *TrinoDeploymentSpec `yaml:"trinoSpecs"`
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

// To get Secret Storage Path
func (sc *PlatformServiceConfig) GetSecretStoragePath() string {
	log.Println("Secret storage path: ", filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath))
	log.Println("Secret storage path: ", sc.VapusBESecretStorage.FilePath)
	log.Println("Secret storage path: ", sc.Path)
	return filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath)
}

func (sc *PlatformServiceConfig) GetDBStoragePath() string {
	return sc.VapusBEDbStorage.Secret
}

func (sc *PlatformServiceConfig) GetCachStoragePath() string {
	return sc.VapusBECacheStorage.Secret
}

func (sc *PlatformServiceConfig) GetJwtAuthSecretPath() string {
	return sc.JWTAuthnSecrets.Secret
}

func (sc *PlatformServiceConfig) GetPolicyConfPath() string {
	return sc.PbacConfig.FilePath
}

func (sc *PlatformServiceConfig) GetArtifactStorePath() string {
	return sc.ArtifactStore.Secret
}

func (sc *PlatformServiceConfig) GetMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *PlatformServiceConfig) GetPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *PlatformServiceConfig) GetCaCert() string {
	return sc.ServerCerts.CaCertFile
}

func (sc *PlatformServiceConfig) GetClientMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientKeyFile
}

func (sc *PlatformServiceConfig) GetClientPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientCertFile
}

func (sc *PlatformServiceConfig) GetAuthnSecrets() string {
	return sc.AuthnSecrets.Secret
}

func (sc *PlatformServiceConfig) GetNetworkConfigFile() string {
	return sc.NetworkConfigFile
}

func (sc *PlatformServiceConfig) GetTrinoConfig() *TrinoDeploymentSpec {
	return sc.TrinoSpecs
}

func (sc *PlatformServiceConfig) GetBaseOs() []*BaseOs {
	return sc.BaseOs
}

func (sc *PlatformServiceConfig) GetFileStorePath() string {
	return sc.VapusFileStorage.Secret
}

type NabhikServerConfig struct {
	Path                 string
	VapusBESecretStorage *SecretMap           `yaml:"vapusBESecretStorage"`
	VapusBEDbStorage     *SecretMap           `yaml:"vapusBEDbStorage"`
	VapusBECacheStorage  *SecretMap           `yaml:"vapusBECacheStorage"`
	VapusFileStorage     *SecretMap           `yaml:"vapusFileStorage"`
	NetworkConfigFile    string               `yaml:"networkConfigFile"`
	JWTAuthnSecrets      *SecretMap           `yaml:"JWTAuthnSecrets"`
	ServerCerts          *TlsCertConfig       `yaml:"serverCerts"`
	TrinoSpecs           *TrinoDeploymentSpec `yaml:"trinoSpecs"`
}

func (sc *NabhikServerConfig) GetFileStorePath() string {
	return sc.VapusFileStorage.Secret
}

func (sc *NabhikServerConfig) GetSecretStoragePath() string {
	return filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath)
}

func (sc *NabhikServerConfig) GetDBStoragePath() string {
	return sc.VapusBEDbStorage.Secret
}

func (sc *NabhikServerConfig) GetCachStoragePath() string {
	return sc.VapusBECacheStorage.Secret
}

func (sc *NabhikServerConfig) GetJwtAuthSecretPath() string {
	return sc.JWTAuthnSecrets.Secret
}

func (sc *NabhikServerConfig) GetMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *NabhikServerConfig) GetPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *NabhikServerConfig) GetCaCert() string {
	return sc.ServerCerts.CaCertFile
}

func (sc *NabhikServerConfig) GetClientMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientKeyFile
}

func (sc *NabhikServerConfig) GetClientPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientCertFile
}

func (sc *NabhikServerConfig) GetTrinoConfig() *TrinoDeploymentSpec {
	return sc.TrinoSpecs
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

type NabrunnerConfig struct {
	Path                 string
	VapusBESecretStorage *SecretMap           `yaml:"vapusBESecretStorage"`
	VapusBEDbStorage     *SecretMap           `yaml:"vapusBEDbStorage"`
	VapusBECacheStorage  *SecretMap           `yaml:"vapusBECacheStorage"`
	VapusFileStorage     *SecretMap           `yaml:"vapusFileStorage"`
	NetworkConfigFile    string               `yaml:"networkConfigFile"`
	ServerCerts          *TlsCertConfig       `yaml:"serverCerts"`
	ConcurrentRunners    int                  `yaml:"concurrentRunners"`
	TrinoSpecs           *TrinoDeploymentSpec `yaml:"trinoSpecs"`
}

func (sc *NabrunnerConfig) GetFileStorePath() string {
	return sc.VapusFileStorage.Secret
}

func (sc *NabrunnerConfig) GetSecretStoragePath() string {
	return filepath.Join(sc.Path, sc.VapusBESecretStorage.FilePath)
}

func (sc *NabrunnerConfig) GetDBStoragePath() string {
	return sc.VapusBEDbStorage.Secret
}

func (sc *NabrunnerConfig) GetCachStoragePath() string {
	log.Println("Cache storage path:>>>>>>>>>>>>>>>> ", sc.VapusBECacheStorage.Secret)
	return sc.VapusBECacheStorage.Secret
}

func (sc *NabrunnerConfig) GetMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *NabrunnerConfig) GetPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ServerCertFile,
		sc.ServerCerts.ServerKeyFile
}

func (sc *NabrunnerConfig) GetCaCert() string {
	return sc.ServerCerts.CaCertFile
}

func (sc *NabrunnerConfig) GetClientMtlsCerts() (string, string, string) {
	return sc.ServerCerts.CaCertFile,
		sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientKeyFile
}

func (sc *NabrunnerConfig) GetClientPlainTlsCerts() (string, string) {
	return sc.ServerCerts.ClientCertFile,
		sc.ServerCerts.ClientCertFile
}
